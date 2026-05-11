package lib

import (
    "bytes"
    "compress/gzip"
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "net/http"
    neturl "net/url"
    "os"
    "os/signal"
    "path/filepath"
    "regexp"
    "strings"
    "sync"
    "sync/atomic"
    "syscall"
    "time"

    "golang.org/x/net/html"
    "golang.org/x/net/html/atom"
    "golang.org/x/sync/errgroup"
    "golang.org/x/sync/semaphore"
    "golang.org/x/time/rate"
)

type CrawlMode int

const (
    ModeFullSite CrawlMode = iota
    ModeSinglePage
)

const (
    maxPathLen     = 240
    maxSegLen      = 120
    chunkSize      = 8192
    defaultTimeout = 30 * time.Second
    maxRetries     = 3
    stateVersion   = 1
)

var (
    assetExtensions = map[string]bool{
        ".css": true, ".js": true, ".mjs": true, ".map": true,
        ".json": true, ".wasm": true, ".webmanifest": true,
        ".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
        ".webp": true, ".avif": true, ".svg": true, ".ico": true,
        ".woff": true, ".woff2": true, ".ttf": true, ".eot": true,
        ".mp4": true, ".webm": true, ".mp3": true, ".pdf": true,
    }

    skipSchemes = map[string]bool{
        "mailto": true, "tel": true, "sms": true, "javascript": true,
        "data": true, "geo": true, "blob": true, "about": true,
    }

    cssURLRegex    = regexp.MustCompile(`url\((?:['"]?)([^'"()]+)(?:['"]?)\)`)
    cssImportRegex = regexp.MustCompile(`@import\s+(?:url\(['"]?|['"])([^'"\);]+)`)
    jsStringRegex  = regexp.MustCompile(`['"](/[^'"?#]+\.(?:png|jpg|jpeg|gif|svg|webp|avif|ico|css|js|mjs|map|woff2?|ttf|eot|json|wasm|webmanifest)(?:\?[^'"]*)?)['"]`)
    hashRouteRegex = regexp.MustCompile(`^#!/?|#!?/`)
)

type Config struct {
    TargetURL          string
    OutputDir          string
    Mode               CrawlMode
    MaxPages           int
    Concurrency        int
    DownloadExternal   bool
    ExternalDomains    []string
    Cookies            map[string]string
    UserAgent          string
    Timeout            time.Duration
    Retries            int
    PreserveStructure  bool
    MinifyOutput       bool
    Resume             bool
    RateLimit          float64
    MaxAssetSize       int64
    CrawlIframes       bool
    CrawlHashRoutes    bool
}

type Dependency struct {
    URL       string
    Type      string
    LocalPath string
    IsIframe  bool
}

type CrawlState struct {
    Version          int       `json:"version"`
    StartURL         string    `json:"start_url"`
    OutputDir        string    `json:"output_dir"`
    VisitedPages     []string  `json:"visited_pages"`
    DownloadedAssets []string  `json:"downloaded_assets"`
    LastUpdate       time.Time `json:"last_update"`
    PagesCount       int       `json:"pages_count"`
    AssetsCount      int       `json:"assets_count"`
    TotalBytes       int64     `json:"total_bytes"`
}

type RateLimiter struct {
    domainLimiters sync.Map
    defaultRate    rate.Limit
    burst          int
}

func NewRateLimiter(defaultRate float64, burst int) *RateLimiter {
    return &RateLimiter{
        defaultRate: rate.Limit(defaultRate),
        burst:       burst,
    }
}

func (rl *RateLimiter) Wait(ctx context.Context, domain string) error {
    limiterI, _ := rl.domainLimiters.LoadOrStore(domain, rate.NewLimiter(rl.defaultRate, rl.burst))
    limiter := limiterI.(*rate.Limiter)
    return limiter.Wait(ctx)
}

type ProgressBar struct {
    total     int64
    current   int64
    mu        sync.Mutex
    lastPrint time.Time
    stopped   atomic.Bool
}

func NewProgressBar(total int64) *ProgressBar {
    pb := &ProgressBar{
        total:     total,
        lastPrint: time.Now(),
    }
    go pb.startReporter()
    return pb
}

func (pb *ProgressBar) Add(n int64) {
    atomic.AddInt64(&pb.current, n)
}

func (pb *ProgressBar) startReporter() {
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    for range ticker.C {
        if pb.stopped.Load() {
            return
        }
        current := atomic.LoadInt64(&pb.current)
        total := atomic.LoadInt64(&pb.total)
        if total > 0 {
            percent := float64(current) / float64(total) * 100
            fmt.Printf("\r[*] Progress: %.1f%% (%d/%d bytes)", percent, current, total)
        }
        if current >= total && total > 0 {
            fmt.Println()
            pb.stopped.Store(true)
            return
        }
    }
}

func (pb *ProgressBar) Stop() {
    pb.stopped.Store(true)
}

type Crawler struct {
    cfg              *Config
    client           *http.Client
    baseURL          *neturl.URL
    outputRoot       string
    visitedPages     sync.Map
    queuedAssets     sync.Map
    pageSem          *semaphore.Weighted
    assetSem         *semaphore.Weighted
    downloadedPages  atomic.Int32
    downloadedAssets atomic.Int32
    totalBytes       atomic.Int64
    errGroup         *errgroup.Group
    ctx              context.Context
    cancel           context.CancelFunc
    dependencies     chan *Dependency
    pageQueue        chan string
    rateLimiter      *RateLimiter
    progressBar      *ProgressBar
    stateFile        string
    iframeQueue      chan string
    pageWorkersDone  sync.WaitGroup
    assetWg          sync.WaitGroup
}

func NewCrawler(cfg *Config) (*Crawler, error) {
    parsedURL, err := neturl.Parse(cfg.TargetURL)
    if err != nil {
        return nil, fmt.Errorf("invalid target URL: %w", err)
    }

    if parsedURL.Scheme == "" {
        parsedURL.Scheme = "https"
    }

    transport := &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 20,
        IdleConnTimeout:     90 * time.Second,
        TLSHandshakeTimeout: 10 * time.Second,
        DisableCompression:  false,
        ForceAttemptHTTP2:   true,
        ResponseHeaderTimeout: cfg.Timeout,
    }

    client := &http.Client{
        Transport: transport,
        Timeout:   cfg.Timeout,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if len(via) >= 10 {
                return fmt.Errorf("too many redirects")
            }
            return nil
        },
    }

    outputRoot := cfg.OutputDir
    if outputRoot == "" {
        outputRoot = sanitizeFilename(parsedURL.Hostname())
    }

    if err := os.MkdirAll(outputRoot, 0755); err != nil {
        return nil, fmt.Errorf("failed to create output directory: %w", err)
    }

    ctx, cancel := context.WithCancel(context.Background())
    eg, ctx := errgroup.WithContext(ctx)

    pageConcurrency := cfg.Concurrency
    if cfg.Mode == ModeSinglePage {
        pageConcurrency = 1
    }

    pageQueueSize := cfg.MaxPages
    if pageQueueSize < 100 {
        pageQueueSize = 100
    }

    stateFile := filepath.Join(outputRoot, ".crawl_state.json")
    
    if cfg.Resume {
        if err := loadState(stateFile, cfg, parsedURL); err != nil {
            fmt.Printf("[*] No previous state found, starting fresh\n")
        }
    }

    rateLimiter := NewRateLimiter(cfg.RateLimit, cfg.Concurrency)

    return &Crawler{
        cfg:          cfg,
        client:       client,
        baseURL:      parsedURL,
        outputRoot:   outputRoot,
        pageSem:      semaphore.NewWeighted(int64(pageConcurrency)),
        assetSem:     semaphore.NewWeighted(int64(cfg.Concurrency * 2)),
        errGroup:     eg,
        ctx:          ctx,
        cancel:       cancel,
        dependencies: make(chan *Dependency, 1000),
        pageQueue:    make(chan string, pageQueueSize),
        iframeQueue:  make(chan string, 100),
        rateLimiter:  rateLimiter,
        stateFile:    stateFile,
    }, nil
}

func loadState(stateFile string, cfg *Config, parsedURL *neturl.URL) error {
    data, err := os.ReadFile(stateFile)
    if err != nil {
        return err
    }
    
    var state CrawlState
    if err := json.Unmarshal(data, &state); err != nil {
        return err
    }
    
    if state.Version != stateVersion {
        return fmt.Errorf("state version mismatch")
    }
    
    fmt.Printf("[✓] Resuming from previous crawl: %d pages, %d assets\n", 
        state.PagesCount, state.AssetsCount)
    
    return nil
}

func (c *Crawler) saveState() error {
    if !c.cfg.Resume {
        return nil
    }
    
    var visitedPages []string
    c.visitedPages.Range(func(key, value interface{}) bool {
        visitedPages = append(visitedPages, key.(string))
        return true
    })
    
    var downloadedAssets []string
    c.queuedAssets.Range(func(key, value interface{}) bool {
        downloadedAssets = append(downloadedAssets, key.(string))
        return true
    })
    
    state := CrawlState{
        Version:          stateVersion,
        StartURL:         c.cfg.TargetURL,
        OutputDir:        c.outputRoot,
        VisitedPages:     visitedPages,
        DownloadedAssets: downloadedAssets,
        LastUpdate:       time.Now(),
        PagesCount:       int(c.downloadedPages.Load()),
        AssetsCount:      int(c.downloadedAssets.Load()),
        TotalBytes:       c.totalBytes.Load(),
    }
    
    data, err := json.MarshalIndent(state, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(c.stateFile, data, 0644)
}

func (c *Crawler) Run() error {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        fmt.Println("\n[!] Received interrupt signal, saving state...")
        c.cancel()
        c.saveState()
    }()

    c.assetWg.Add(1)
    go func() {
        defer c.assetWg.Done()
        c.processDependencies()
    }()
    
    if c.cfg.CrawlIframes {
        c.assetWg.Add(1)
        go func() {
            defer c.assetWg.Done()
            c.processIframes()
        }()
    }

    startTime := time.Now()
    c.progressBar = NewProgressBar(0)

    if c.cfg.Mode == ModeSinglePage {
        fmt.Printf("[*] Single-page mode: backing up %s\n", c.cfg.TargetURL)
        if err := c.crawlSinglePage(); err != nil {
            c.cancel()
            return err
        }
    } else {
        fmt.Printf("[*] Full-site mode: crawling %s\n", c.baseURL.String())
        if err := c.crawlFullSite(); err != nil {
            c.cancel()
            return err
        }
    }

    c.pageWorkersDone.Wait()
    close(c.pageQueue)
    close(c.dependencies)
    close(c.iframeQueue)
    c.assetWg.Wait()

    c.saveState()
    if c.progressBar != nil {
        c.progressBar.Stop()
    }

    elapsed := time.Since(startTime)
    fmt.Printf("\n[✓] Backup completed in %s\n", elapsed.Round(time.Millisecond))
    fmt.Printf("    Pages: %d | Assets: %d | Size: %.2f MB\n",
        c.downloadedPages.Load(),
        c.downloadedAssets.Load(),
        float64(c.totalBytes.Load())/(1024*1024))

    return nil
}

func (c *Crawler) crawlSinglePage() error {
    mainURL := c.cfg.TargetURL
    
    if c.cfg.CrawlHashRoutes && hashRouteRegex.MatchString(mainURL) {
        mainURL = hashRouteRegex.ReplaceAllString(mainURL, "")
        fmt.Printf("[*] Normalized hash route: %s\n", mainURL)
    }
    
    c.visitedPages.Store(mainURL, true)

    doc, htmlBytes, err := c.fetchHTML(mainURL)
    if err != nil {
        return fmt.Errorf("failed to fetch main page: %w", err)
    }

    if _, err := c.savePage(doc, mainURL, htmlBytes); err != nil {
        return fmt.Errorf("failed to save main page: %w", err)
    }
    c.downloadedPages.Add(1)

    deps := c.extractDependencies(doc, mainURL)
    for _, dep := range deps {
        select {
        case c.dependencies <- dep:
        case <-c.ctx.Done():
            return c.ctx.Err()
        }
    }

    return nil
}

func (c *Crawler) crawlFullSite() error {
    select {
    case <-c.ctx.Done():
        return c.ctx.Err()
    default:
    }
    
    c.pageQueue <- c.cfg.TargetURL
    c.visitedPages.Store(c.cfg.TargetURL, true)

    for i := 0; i < c.cfg.Concurrency; i++ {
        c.pageWorkersDone.Add(1)
        go func() {
            defer c.pageWorkersDone.Done()
            c.pageWorker()
        }()
    }

    return nil
}

func (c *Crawler) pageWorker() {
    for {
        select {
        case <-c.ctx.Done():
            return
        case url, ok := <-c.pageQueue:
            if !ok {
                return
            }
            if int(c.downloadedPages.Load()) >= c.cfg.MaxPages {
                return
            }
            if err := c.processPage(url); err != nil {
                continue
            }
        }
    }
}

func (c *Crawler) processIframes() {
    for {
        select {
        case <-c.ctx.Done():
            return
        case iframeURL, ok := <-c.iframeQueue:
            if !ok {
                return
            }
            if err := c.processIframe(iframeURL); err != nil {
                fmt.Printf("  [!] Iframe failed: %s -> %v\n", iframeURL, err)
            }
        }
    }
}

func (c *Crawler) processIframe(iframeURL string) error {
    if err := c.rateLimiter.Wait(c.ctx, c.baseURL.Host); err != nil {
        return err
    }
    
    doc, htmlBytes, err := c.fetchHTML(iframeURL)
    if err != nil {
        return err
    }
    
    parsedURL, _ := neturl.Parse(iframeURL)
    if parsedURL == nil {
        return fmt.Errorf("invalid iframe URL")
    }
    
    localPath := c.urlToLocalPath(parsedURL, true)
    if _, err := c.savePage(doc, iframeURL, htmlBytes); err != nil {
        return err
    }
    
    deps := c.extractDependencies(doc, iframeURL)
    for _, dep := range deps {
        select {
        case c.dependencies <- dep:
        case <-c.ctx.Done():
            return c.ctx.Err()
        }
    }
    
    _ = localPath
    return nil
}

func (c *Crawler) fetchHTML(rawURL string) (*html.Node, []byte, error) {
    body, err := c.fetchWithRetry(rawURL)
    if err != nil {
        return nil, nil, err
    }

    if len(body) >= 3 && body[0] == 0xEF && body[1] == 0xBB && body[2] == 0xBF {
        body = body[3:]
    }

    doc, err := html.Parse(bytes.NewReader(body))
    if err != nil {
        return nil, nil, fmt.Errorf("failed to parse HTML: %w", err)
    }

    return doc, body, nil
}

func (c *Crawler) fetchWithRetry(rawURL string) ([]byte, error) {
    var lastErr error

    for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
        if attempt > 0 {
            select {
            case <-c.ctx.Done():
                return nil, c.ctx.Err()
            case <-time.After(time.Duration(attempt) * time.Second):
            }
        }

        parsedURL, _ := neturl.Parse(rawURL)
        if parsedURL != nil {
            if err := c.rateLimiter.Wait(c.ctx, parsedURL.Host); err != nil {
                return nil, err
            }
        }

        req, err := http.NewRequestWithContext(c.ctx, "GET", rawURL, nil)
        if err != nil {
            lastErr = err
            continue
        }

        req.Header.Set("User-Agent", c.cfg.UserAgent)
        req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
        req.Header.Set("Accept-Language", "en-US,en;q=0.9")
        req.Header.Set("Accept-Encoding", "gzip, deflate, br")
        req.Header.Set("Connection", "keep-alive")
        req.Header.Set("Upgrade-Insecure-Requests", "1")

        for name, value := range c.cfg.Cookies {
            req.AddCookie(&http.Cookie{Name: name, Value: value})
        }

        resp, err := c.client.Do(req)
        if err != nil {
            lastErr = err
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
            lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
            continue
        }

        if resp.StatusCode != http.StatusOK {
            return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
        }

        if c.cfg.MaxAssetSize > 0 && resp.ContentLength > c.cfg.MaxAssetSize {
            return nil, fmt.Errorf("asset too large: %d bytes", resp.ContentLength)
        }

        var reader io.ReadCloser = resp.Body
        switch resp.Header.Get("Content-Encoding") {
        case "gzip":
            reader, err = gzip.NewReader(resp.Body)
            if err != nil {
                return nil, err
            }
            defer reader.Close()
        }

        body, err := io.ReadAll(reader)
        if err != nil {
            lastErr = err
            continue
        }

        c.totalBytes.Add(int64(len(body)))
        if c.progressBar != nil {
            c.progressBar.Add(int64(len(body)))
        }
        return body, nil
    }

    return nil, fmt.Errorf("failed after %d retries: %w", c.cfg.Retries, lastErr)
}

func (c *Crawler) savePage(doc *html.Node, pageURL string, htmlBytes []byte) (string, error) {
    parsedURL, err := neturl.Parse(pageURL)
    if err != nil {
        return "", err
    }

    localPath := c.urlToLocalPath(parsedURL, true)

    if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
        return "", err
    }

    c.rewriteHTML(doc, pageURL, localPath)

    var output []byte
    if c.cfg.MinifyOutput {
        output = minifyHTML(htmlBytes)
    } else {
        var buf bytes.Buffer
        if err := html.Render(&buf, doc); err != nil {
            return "", err
        }
        output = buf.Bytes()
    }

    if err := os.WriteFile(localPath, output, 0644); err != nil {
        return "", err
    }

    return localPath, nil
}

func (c *Crawler) extractDependencies(n *html.Node, baseURL string) []*Dependency {
    var deps []*Dependency

    var traverse func(*html.Node)
    traverse = func(node *html.Node) {
        if node.Type == html.ElementNode {
            switch node.DataAtom {
            case atom.Img, atom.Script, atom.Link:
                for _, attr := range node.Attr {
                    if attr.Key == "src" || (attr.Key == "href" && c.isResourceLink(node)) {
                        if dep := c.createDependency(attr.Val, baseURL, node.Data); dep != nil {
                            deps = append(deps, dep)
                        }
                    }
                }

            case atom.Video, atom.Audio, atom.Source:
                for _, attr := range node.Attr {
                    if attr.Key == "src" || attr.Key == "poster" {
                        if dep := c.createDependency(attr.Val, baseURL, node.Data); dep != nil {
                            deps = append(deps, dep)
                        }
                    }
                }
                
            case atom.Iframe:
                if c.cfg.CrawlIframes {
                    for _, attr := range node.Attr {
                        if attr.Key == "src" {
                            absURL := c.resolveURL(attr.Val, baseURL)
                            if absURL != "" {
                                select {
                                case c.iframeQueue <- absURL:
                                default:
                                }
                            }
                            break
                        }
                    }
                }
            }

            if node.DataAtom == atom.Img {
                for _, attr := range node.Attr {
                    if attr.Key == "srcset" {
                        urls := parseSrcSet(attr.Val)
                        for _, u := range urls {
                            if dep := c.createDependency(u, baseURL, "srcset"); dep != nil {
                                deps = append(deps, dep)
                            }
                        }
                    }
                }
            }
        }

        for child := node.FirstChild; child != nil; child = child.NextSibling {
            traverse(child)
        }
    }

    traverse(n)
    return deps
}

func (c *Crawler) createDependency(rawURL, baseURL, tagName string) *Dependency {
    if rawURL == "" || strings.HasPrefix(rawURL, "#") || strings.HasPrefix(rawURL, "data:") {
        return nil
    }

    absURL := c.resolveURL(rawURL, baseURL)
    if absURL == "" {
        return nil
    }

    parsedURL, err := neturl.Parse(absURL)
    if err != nil {
        return nil
    }

    ext := strings.ToLower(filepath.Ext(parsedURL.Path))
    if ext != "" && !assetExtensions[ext] && tagName != "srcset" {
        return nil
    }

    isExternal := !c.isSameDomain(parsedURL)
    if isExternal && !c.cfg.DownloadExternal {
        return nil
    }

    if isExternal && len(c.cfg.ExternalDomains) > 0 {
        allowed := false
        for _, domain := range c.cfg.ExternalDomains {
            if strings.Contains(parsedURL.Hostname(), domain) {
                allowed = true
                break
            }
        }
        if !allowed {
            return nil
        }
    }

    localPath := c.urlToLocalPath(parsedURL, false)

    return &Dependency{
        URL:       absURL,
        Type:      ext,
        LocalPath: localPath,
    }
}

func (c *Crawler) processDependencies() {
    for {
        select {
        case <-c.ctx.Done():
            return
        case dep, ok := <-c.dependencies:
            if !ok {
                return
            }
            if _, loaded := c.queuedAssets.LoadOrStore(dep.URL, true); loaded {
                continue
            }

            if err := c.assetSem.Acquire(c.ctx, 1); err != nil {
                return
            }

            go func(d *Dependency) {
                defer c.assetSem.Release(1)
                
                select {
                case <-c.ctx.Done():
                    return
                default:
                    if err := c.downloadAsset(d.URL, d.LocalPath); err != nil {
                        fmt.Printf("  [!] Failed: %s -> %v\n", d.URL, err)
                        return
                    }
                    c.downloadedAssets.Add(1)
                    fmt.Printf("  [↓] %s\n", filepath.Base(d.LocalPath))
                }
            }(dep)
        }
    }
}

func (c *Crawler) downloadAsset(assetURL, localPath string) error {
    if _, err := os.Stat(localPath); err == nil {
        return nil
    }

    body, err := c.fetchWithRetry(assetURL)
    if err != nil {
        return err
    }

    if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
        return err
    }

    ext := filepath.Ext(localPath)
    if ext == ".css" {
        body = c.rewriteCSSContent(body, assetURL, filepath.Dir(localPath))
    } else if ext == ".js" || ext == ".mjs" {
        body = c.rewriteJSContent(body, assetURL, filepath.Dir(localPath))
    }

    if err := os.WriteFile(localPath, body, 0644); err != nil {
        return err
    }

    c.totalBytes.Add(int64(len(body)))
    return nil
}

func (c *Crawler) rewriteCSSContent(data []byte, baseURL, baseDir string) []byte {
    content := string(data)

    content = cssURLRegex.ReplaceAllStringFunc(content, func(match string) string {
        parts := cssURLRegex.FindStringSubmatch(match)
        if len(parts) < 2 {
            return match
        }

        urlStr := strings.TrimSpace(parts[1])
        urlStr = strings.Trim(urlStr, "'\"")

        if strings.HasPrefix(urlStr, "data:") {
            return match
        }
        
        if (strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://")) && !c.cfg.DownloadExternal {
            return match
        }

        resolved := c.resolveURL(urlStr, baseURL)
        if resolved == "" {
            return match
        }

        parsed, _ := neturl.Parse(resolved)
        if parsed == nil {
            return match
        }

        localPath := c.urlToLocalPath(parsed, false)
        relPath, err := filepath.Rel(baseDir, localPath)
        if err != nil {
            return match
        }

        relPath = filepath.ToSlash(relPath)
        return fmt.Sprintf("url('%s')", relPath)
    })

    content = cssImportRegex.ReplaceAllStringFunc(content, func(match string) string {
        parts := cssImportRegex.FindStringSubmatch(match)
        if len(parts) < 2 {
            return match
        }

        urlStr := strings.TrimSpace(parts[1])
        urlStr = strings.Trim(urlStr, "'\"")

        resolved := c.resolveURL(urlStr, baseURL)
        if resolved == "" {
            return match
        }

        parsed, _ := neturl.Parse(resolved)
        if parsed == nil {
            return match
        }

        localPath := c.urlToLocalPath(parsed, false)
        relPath, err := filepath.Rel(baseDir, localPath)
        if err != nil {
            return match
        }

        return fmt.Sprintf("@import '%s'", filepath.ToSlash(relPath))
    })

    return []byte(content)
}

func (c *Crawler) rewriteJSContent(data []byte, baseURL, baseDir string) []byte {
    content := string(data)

    content = jsStringRegex.ReplaceAllStringFunc(content, func(match string) string {
        if len(match) < 2 {
            return match
        }

        quote := match[0]
        urlStr := match[1 : len(match)-1]

        resolved := c.resolveURL(urlStr, baseURL)
        if resolved == "" {
            return match
        }

        parsed, _ := neturl.Parse(resolved)
        if parsed == nil {
            return match
        }

        localPath := c.urlToLocalPath(parsed, false)
        relPath, err := filepath.Rel(baseDir, localPath)
        if err != nil {
            return match
        }

        return string(quote) + filepath.ToSlash(relPath) + string(quote)
    })

    return []byte(content)
}

func (c *Crawler) rewriteHTML(n *html.Node, pageURL, pagePath string) {
    if n.Type == html.ElementNode {
        switch n.DataAtom {
        case atom.Base:
            n.Parent.RemoveChild(n)
            return

        case atom.A:
            for i, attr := range n.Attr {
                if attr.Key == "href" {
                    newVal := c.rewriteURL(attr.Val, pageURL, pagePath, true)
                    if c.cfg.CrawlHashRoutes && hashRouteRegex.MatchString(newVal) {
                        newVal = hashRouteRegex.ReplaceAllString(newVal, "")
                    }
                    n.Attr[i].Val = newVal
                    break
                }
            }

        case atom.Img, atom.Script, atom.Link, atom.Video, atom.Audio, atom.Source:
            for i, attr := range n.Attr {
                if attr.Key == "src" || attr.Key == "href" || attr.Key == "poster" {
                    if n.DataAtom == atom.Link && !c.isResourceLink(n) {
                        continue
                    }
                    n.Attr[i].Val = c.rewriteURL(attr.Val, pageURL, pagePath, false)
                    break
                }
            }

            for i, attr := range n.Attr {
                if attr.Key == "srcset" {
                    urls := parseSrcSet(attr.Val)
                    var newURLs []string
                    for _, u := range urls {
                        newURL := c.rewriteURL(u, pageURL, pagePath, false)
                        newURLs = append(newURLs, newURL)
                    }
                    n.Attr[i].Val = strings.Join(newURLs, ", ")
                    break
                }
            }

        case atom.Style:
            if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
                rewritten := c.rewriteCSSContent([]byte(n.FirstChild.Data), pageURL, filepath.Dir(pagePath))
                n.FirstChild.Data = string(rewritten)
            }
            
        case atom.Iframe:
            for i, attr := range n.Attr {
                if attr.Key == "src" && c.cfg.CrawlIframes {
                    n.Attr[i].Val = c.rewriteURL(attr.Val, pageURL, pagePath, true)
                    break
                }
            }
        }

        for i, attr := range n.Attr {
            if attr.Key == "style" {
                rewritten := c.rewriteCSSContent([]byte(attr.Val), pageURL, filepath.Dir(pagePath))
                n.Attr[i].Val = string(rewritten)
                break
            }
        }
    }

    for child := n.FirstChild; child != nil; child = child.NextSibling {
        c.rewriteHTML(child, pageURL, pagePath)
    }
}

func (c *Crawler) rewriteURL(rawURL, pageURL, pagePath string, isPageLink bool) string {
    if rawURL == "" || strings.HasPrefix(rawURL, "#") || strings.HasPrefix(rawURL, "javascript:") || strings.HasPrefix(rawURL, "data:") {
        return rawURL
    }

    resolved := c.resolveURL(rawURL, pageURL)
    if resolved == "" {
        return rawURL
    }

    parsed, err := neturl.Parse(resolved)
    if err != nil {
        return rawURL
    }

    isExternal := !c.isSameDomain(parsed)

    if isPageLink && isExternal {
        return rawURL
    }

    if isExternal && !c.cfg.DownloadExternal {
        return rawURL
    }

    localPath := c.urlToLocalPath(parsed, isPageLink)
    relPath, err := filepath.Rel(filepath.Dir(pagePath), localPath)
    if err != nil {
        return rawURL
    }

    relPath = filepath.ToSlash(relPath)

    if parsed.Fragment != "" {
        relPath += "#" + parsed.Fragment
    }

    if parsed.RawQuery != "" && !strings.Contains(relPath, "?") {
        relPath += "?" + parsed.RawQuery
    }

    return relPath
}

func (c *Crawler) urlToLocalPath(parsedURL *neturl.URL, isPage bool) string {
    var relPath string

    if isPage {
        relPath = strings.TrimPrefix(parsedURL.Path, "/")
        if relPath == "" {
            relPath = "index.html"
        } else if strings.HasSuffix(relPath, "/") {
            relPath += "index.html"
        } else if filepath.Ext(relPath) == "" {
            relPath += ".html"
        }
    } else {
        relPath = strings.TrimPrefix(parsedURL.Path, "/")
        if relPath == "" {
            relPath = "index"
        }
    }

    if parsedURL.RawQuery != "" {
        hash := sha256.Sum256([]byte(parsedURL.RawQuery))
        hashStr := hex.EncodeToString(hash[:])[:10]
        ext := filepath.Ext(relPath)
        base := strings.TrimSuffix(relPath, ext)
        relPath = base + "-" + hashStr + ext
    }

    parts := strings.Split(relPath, "/")
    for i, part := range parts {
        parts[i] = sanitizeSegment(part)
        if len(parts[i]) > maxSegLen {
            parts[i] = shortenSegment(parts[i], maxSegLen)
        }
    }

    localPath := filepath.Join(c.outputRoot, filepath.Join(parts...))

    if len(localPath) > maxPathLen {
        hash := sha256.Sum256([]byte(parsedURL.String()))
        hashStr := hex.EncodeToString(hash[:])[:16]
        dir := filepath.Dir(localPath)
        base := filepath.Base(localPath)
        ext := filepath.Ext(base)
        stem := strings.TrimSuffix(base, ext)
        newBase := fmt.Sprintf("%s-%s%s", stem[:min(50, len(stem))], hashStr, ext)
        localPath = filepath.Join(dir, newBase)
    }

    return localPath
}

func (c *Crawler) resolveURL(rawURL, baseURL string) string {
    if strings.HasPrefix(rawURL, "//") {
        rawURL = c.baseURL.Scheme + ":" + rawURL
    }

    base, err := neturl.Parse(baseURL)
    if err != nil {
        return ""
    }

    ref, err := neturl.Parse(rawURL)
    if err != nil {
        return ""
    }

    resolved := base.ResolveReference(ref)
    resolved.Fragment = ""

    return resolved.String()
}

func (c *Crawler) isSameDomain(parsedURL *neturl.URL) bool {
    return parsedURL.Hostname() == c.baseURL.Hostname()
}

func (c *Crawler) isResourceLink(node *html.Node) bool {
    for _, attr := range node.Attr {
        if attr.Key == "rel" {
            rel := strings.ToLower(attr.Val)
            return rel == "stylesheet" || rel == "icon" || rel == "preload" || rel == "modulepreload" || rel == "manifest"
        }
    }
    return false
}

func (c *Crawler) processPage(url string) error {
    if err := c.pageSem.Acquire(c.ctx, 1); err != nil {
        return err
    }
    defer c.pageSem.Release(1)

    doc, htmlBytes, err := c.fetchHTML(url)
    if err != nil {
        return nil
    }

    if _, err := c.savePage(doc, url, htmlBytes); err != nil {
        return nil
    }

    c.downloadedPages.Add(1)
    fmt.Printf("[%d] %s\n", c.downloadedPages.Load(), url)

    deps := c.extractDependencies(doc, url)
    for _, dep := range deps {
        select {
        case c.dependencies <- dep:
        case <-c.ctx.Done():
            return c.ctx.Err()
        }
    }

    if c.cfg.Mode == ModeFullSite {
        c.extractPageLinks(doc, url)
    }

    return nil
}

func (c *Crawler) extractPageLinks(n *html.Node, baseURL string) {
    var traverse func(*html.Node)
    traverse = func(node *html.Node) {
        select {
        case <-c.ctx.Done():
            return
        default:
        }
        
        if node.Type == html.ElementNode && node.DataAtom == atom.A {
            for _, attr := range node.Attr {
                if attr.Key == "href" {
                    absURL := c.resolveURL(attr.Val, baseURL)
                    if absURL != "" && c.isSameDomainOrSubdomain(absURL) && absURL != baseURL {
                        if c.cfg.CrawlHashRoutes {
                            absURL = hashRouteRegex.ReplaceAllString(absURL, "")
                        }
                        if _, loaded := c.visitedPages.LoadOrStore(absURL, true); !loaded {
                            if int(c.downloadedPages.Load()) < c.cfg.MaxPages {
                                select {
                                case c.pageQueue <- absURL:
                                case <-c.ctx.Done():
                                    return
                                default:
                                }
                            }
                        }
                    }
                    break
                }
            }
        }
        for child := node.FirstChild; child != nil; child = child.NextSibling {
            traverse(child)
        }
    }
    traverse(n)
}

func (c *Crawler) isSameDomainOrSubdomain(rawURL string) bool {
    parsed, err := neturl.Parse(rawURL)
    if err != nil {
        return false
    }
    return parsed.Hostname() == c.baseURL.Hostname() || 
           strings.HasSuffix(parsed.Hostname(), "."+c.baseURL.Hostname())
}

func parseSrcSet(srcset string) []string {
    var urls []string
    entries := strings.Split(srcset, ",")
    for _, entry := range entries {
        parts := strings.Fields(strings.TrimSpace(entry))
        if len(parts) > 0 {
            urls = append(urls, parts[0])
        }
    }
    return urls
}

func sanitizeSegment(segment string) string {
    segment = strings.TrimSpace(segment)
    segment = strings.Trim(segment, ". ")

    replacer := strings.NewReplacer(
        "<", "_", ">", "_", ":", "_", `"`, "_", `|`, "_",
        "?", "_", "*", "_", "\\", "_", "/", "_",
    )
    segment = replacer.Replace(segment)

    if segment == "" || segment == "." || segment == ".." {
        segment = "_"
    }

    return segment
}

func shortenSegment(segment string, limit int) string {
    if len(segment) <= limit {
        return segment
    }

    ext := filepath.Ext(segment)
    stem := strings.TrimSuffix(segment, ext)

    hash := sha256.Sum256([]byte(segment))
    hashStr := hex.EncodeToString(hash[:])[:8]

    keep := limit - len(ext) - 9
    if keep < 8 {
        keep = 8
    }

    if keep > len(stem) {
        keep = len(stem)
    }

    return stem[:keep] + "-" + hashStr + ext
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func sanitizeFilename(name string) string {
    return strings.NewReplacer(".", "_", ":", "_", "/", "_").Replace(name)
}

func minifyHTML(htmlBytes []byte) []byte {
    s := string(htmlBytes)
    s = regexp.MustCompile(`>\s+<`).ReplaceAllString(s, "><")
    s = regexp.MustCompile(`\s{2,}`).ReplaceAllString(s, " ")
    return []byte(s)
}

func RunWebDownloader() {
    var (
        targetURL       string
        outputDir       string
        mode            string
        maxPages        int
        concurrency     int
        downloadExt     bool
        externalDomStr  string
        cookieStr       string
        userAgent       string
        timeoutSec      int
        retries         int
        minify          bool
        resume          bool
        rateLimit       float64
        maxAssetSizeMB  int
        crawlIframes    bool
        crawlHashRoutes bool
    )

    webFlags := flag.NewFlagSet("web", flag.ExitOnError)
    
    webFlags.StringVar(&targetURL, "url", "", "Target URL to backup")
    webFlags.StringVar(&outputDir, "output", "", "Output directory (default: domain name)")
    webFlags.StringVar(&mode, "mode", "single", "Crawl mode: 'single' or 'full'")
    webFlags.IntVar(&maxPages, "max-pages", 100, "Maximum pages for full-site mode")
    webFlags.IntVar(&concurrency, "concurrency", 5, "Number of concurrent workers")
    webFlags.BoolVar(&downloadExt, "download-external", false, "Download external assets")
    webFlags.StringVar(&externalDomStr, "external-domains", "", "Comma-separated external domains to include")
    webFlags.StringVar(&cookieStr, "cookies", "", "Cookies (format: name1=value1; name2=value2)")
    webFlags.StringVar(&userAgent, "user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36", "User-Agent")
    webFlags.IntVar(&timeoutSec, "timeout", 30, "Request timeout in seconds")
    webFlags.IntVar(&retries, "retries", 3, "Number of retries on failure")
    webFlags.BoolVar(&minify, "minify", false, "Minify HTML output")
    webFlags.BoolVar(&resume, "resume", false, "Resume interrupted crawl")
    webFlags.Float64Var(&rateLimit, "rate-limit", 10, "Requests per second per domain")
    webFlags.IntVar(&maxAssetSizeMB, "max-asset-size", 50, "Maximum asset size in MB")
    webFlags.BoolVar(&crawlIframes, "crawl-iframes", true, "Download iframe content")
    webFlags.BoolVar(&crawlHashRoutes, "crawl-hash-routes", true, "Handle hash-based routing (SPA)")

    webFlags.Parse(os.Args[1:])

    if targetURL == "" {
        fmt.Fprintf(os.Stderr, "Error: --url is required\n")
        webFlags.Usage()
        os.Exit(1)
    }

    var crawlMode CrawlMode
    switch strings.ToLower(mode) {
    case "single":
        crawlMode = ModeSinglePage
        fmt.Println("[*] Mode: Single Page (only target URL + dependencies)")
    case "full":
        crawlMode = ModeFullSite
        fmt.Printf("[*] Mode: Full Site (max %d pages)\n", maxPages)
    default:
        fmt.Fprintf(os.Stderr, "Invalid mode: %s (use 'single' or 'full')\n", mode)
        os.Exit(1)
    }

    var externalDomains []string
    if externalDomStr != "" {
        externalDomains = strings.Split(externalDomStr, ",")
        for i := range externalDomains {
            externalDomains[i] = strings.TrimSpace(externalDomains[i])
        }
    }

    cookies := make(map[string]string)
    if cookieStr != "" {
        parts := strings.Split(cookieStr, ";")
        for _, part := range parts {
            kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
            if len(kv) == 2 {
                cookies[kv[0]] = kv[1]
            }
        }
    }

    config := &Config{
        TargetURL:         targetURL,
        OutputDir:         outputDir,
        Mode:              crawlMode,
        MaxPages:          maxPages,
        Concurrency:       concurrency,
        DownloadExternal:  downloadExt,
        ExternalDomains:   externalDomains,
        Cookies:           cookies,
        UserAgent:         userAgent,
        Timeout:           time.Duration(timeoutSec) * time.Second,
        Retries:           retries,
        PreserveStructure: true,
        MinifyOutput:      minify,
        Resume:            resume,
        RateLimit:         rateLimit,
        MaxAssetSize:      int64(maxAssetSizeMB) * 1024 * 1024,
        CrawlIframes:      crawlIframes,
        CrawlHashRoutes:   crawlHashRoutes,
    }

    crawler, err := NewCrawler(config)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    if err := crawler.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
