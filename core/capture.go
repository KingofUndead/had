package core

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/mattn/go-colorable"
)

func init() {
	_ = colorable.NewColorable(os.Stdout)
}

type FileType string

const (
	TypeVideo    FileType = "video"
	TypeMusic    FileType = "music"
	TypeImage    FileType = "image"
	TypeDocument FileType = "document"
	TypeArchive  FileType = "archive"
	TypeAll      FileType = "all"
)

var ExtensionGroups = map[FileType][]string{
	TypeVideo: {
		".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm", ".m4v",
		".mpg", ".mpeg", ".m2ts", ".mts", ".ts", ".m3u8", ".mpd", ".iso",
		".vob", ".3gp", ".ogv", ".ogg", ".qt", ".rm", ".rmvb", ".asf",
		".divx", ".xvid", ".264", ".265", ".hevc",
	},
	TypeMusic: {
		".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a", ".opus", ".wma",
		".alac", ".ape", ".dsd", ".dff", ".dsf", ".mid", ".midi", ".ra",
		".voc", ".vox", ".aiff", ".au", ".snd", ".amr", ".awb", ".weba",
	},
	TypeImage: {
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp", ".ico",
		".tiff", ".tif", ".raw", ".cr2", ".nef", ".arw", ".dng", ".heif",
		".heic", ".jfif", ".pjpeg", ".pjp", ".avif", ".apng", ".cur",
	},
	TypeDocument: {
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt",
		".rtf", ".odt", ".ods", ".odp", ".odg", ".csv", ".json", ".xml",
		".md", ".epub", ".mobi", ".azw", ".azw3", ".cbr", ".cbz", ".ps",
		".tex", ".log", ".ini", ".cfg", ".conf", ".yaml", ".yml",
	},
	TypeArchive: {
		".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz", ".tgz",
		".tbz2", ".txz", ".zst", ".lzma", ".lz", ".lzh", ".cab", ".arj",
		".deb", ".rpm", ".pkg", ".msi", ".apk", ".ipa", ".jar", ".war",
	},
}

type CaptureConfig struct {
	Port             string
	FileTypes        []FileType
	CustomExtensions []string
	Headers          map[string]string
	Cookie           string
	AutoDownload     bool
	OutputDir        string
	MinFileSize      int64
	MaxFileSize      int64
	ConfidenceLevel  int
	SaveToFile       string
	Verbose          bool
	CaptureBody      bool
	FilterDomain     string
	FilterPattern    string
}

type CapturedItem struct {
	URL          string
	FileType     FileType
	Extension    string
	Size         int64
	Title        string
	SourceURL    string
	Timestamp    time.Time
	Confidence   int
	Method       string
	StatusCode   int
	Headers      map[string]string
	RequestBody  []byte
	ResponseBody []byte
}

type RequestLog struct {
	Method      string
	URL         string
	Host        string
	Path        string
	QueryParams map[string]string
	Headers     map[string]string
	Body        string
	Timestamp   time.Time
}


type CaptureProxy struct {
	config      *CaptureConfig
	proxy       *goproxy.ProxyHttpServer
	captured    []CapturedItem
	requestLogs []RequestLog
	mu          sync.RWMutex
	urlFilters  []*regexp.Regexp
}

func NewCaptureProxy(config *CaptureConfig) *CaptureProxy {
	if config == nil {
		config = &CaptureConfig{
			Port:             ":8085",
			FileTypes:        []FileType{TypeVideo, TypeMusic},
			CustomExtensions: []string{},
			Headers:          make(map[string]string),
			AutoDownload:     false,
			OutputDir:        "captured",
			MinFileSize:      1024,
			MaxFileSize:      0,
			ConfidenceLevel:  30,
			SaveToFile:       "captured_links.txt",
			Verbose:          true,
			CaptureBody:      false,
		}
	}

	os.MkdirAll(config.OutputDir, 0755)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.Tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cp := &CaptureProxy{
		config:      config,
		proxy:       proxy,
		captured:    make([]CapturedItem, 0),
		requestLogs: make([]RequestLog, 0),
		urlFilters:  make([]*regexp.Regexp, 0),
	}

	if config.FilterPattern != "" {
		if re, err := regexp.Compile(config.FilterPattern); err == nil {
			cp.urlFilters = append(cp.urlFilters, re)
		}
	}

	cp.setupHandlers()
	return cp
}

func (cp *CaptureProxy) setupHandlers() {

	cp.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	cp.proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

		for k, v := range cp.config.Headers {
			r.Header.Set(k, v)
		}
		if cp.config.Cookie != "" {
			r.Header.Set("Cookie", cp.config.Cookie)
		}

		cp.logRequest(r)

		if cp.config.FilterDomain != "" && !strings.Contains(r.Host, cp.config.FilterDomain) {
			return r, nil
		}

		if len(cp.urlFilters) > 0 {
			matched := false
			for _, filter := range cp.urlFilters {
				if filter.MatchString(r.URL.String()) {
					matched = true
					break
				}
			}
			if !matched {
				return r, nil
			}
		}

		url := r.URL.String()

		if cp.config.Verbose {
			cp.logAdvanced(r)
		}

		if cp.config.CaptureBody && (r.Method == "POST" || r.Method == "PUT") {
			body, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(strings.NewReader(string(body)))

			cp.analyzeBodyForURLs(string(body), r)
		}

		cp.capture(url, r)

		return r, nil
	})

	cp.proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp != nil && resp.Request != nil {
			url := resp.Request.URL.String()

			cp.logResponse(resp)

			if cp.config.CaptureBody {
				body, _ := io.ReadAll(resp.Body)
				resp.Body = io.NopCloser(strings.NewReader(string(body)))
				cp.analyzeBodyForURLs(string(body), resp.Request)
			}

			cp.capture(url, resp.Request)

			cp.captureInterestingHeaders(resp.Header, url)
		}
		return resp
	})

	cp.proxy.OnRequest().HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		if cp.config.Verbose {
			log.Printf("[CONNECT] %s", host)
		}
		return goproxy.OkConnect, host
	})
}

func (cp *CaptureProxy) logRequest(r *http.Request) {
	logEntry := RequestLog{
		Method:      r.Method,
		URL:         r.URL.String(),
		Host:        r.Host,
		Path:        r.URL.Path,
		QueryParams: make(map[string]string),
		Headers:     make(map[string]string),
		Timestamp:   time.Now(),
	}

	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			logEntry.QueryParams[k] = v[0]
		}
	}

	for k, v := range r.Header {
		if len(v) > 0 {
			logEntry.Headers[k] = v[0]
		}
	}

	cp.mu.Lock()
	cp.requestLogs = append(cp.requestLogs, logEntry)

	if len(cp.requestLogs) > 1000 {
		cp.requestLogs = cp.requestLogs[1:]
	}
	cp.mu.Unlock()
}


func (cp *CaptureProxy) logAdvanced(r *http.Request) {
	methodColor := "\033[36m"
	switch r.Method {
	case "GET":
		methodColor = "\033[32m"
	case "POST":
		methodColor = "\033[33m"
	case "PUT":
		methodColor = "\033[34m" 
	case "DELETE":
		methodColor = "\033[31m"
	}

	fmt.Printf("\n%s[%s]\033[0m \033[1m%s\033[0m",
		methodColor, r.Method, r.Host)

	if r.URL.Path != "/" && r.URL.Path != "" {
		fmt.Printf("\033[90m%s\033[0m", r.URL.Path)
	}

	if len(r.URL.Query()) > 0 {
		fmt.Printf("\n  \033[90mParams:\033[0m ")
		first := true
		for k, v := range r.URL.Query() {
			if !first {
				fmt.Printf(", ")
			}
			fmt.Printf("%s=%s", k, strings.Join(v, ","))
			first = false
		}
	}

	importantHeaders := []string{"Referer", "User-Agent", "Content-Type", "Origin"}
	hasHeaders := false
	for _, h := range importantHeaders {
		if val := r.Header.Get(h); val != "" {
			if !hasHeaders {
				fmt.Printf("\n  \033[90mHeaders:\033[0m ")
				hasHeaders = true
			} else {
				fmt.Printf(", ")
			}
			fmt.Printf("%s: %s", h, truncate(val, 50))
		}
	}

	fmt.Println()
}

func (cp *CaptureProxy) logResponse(resp *http.Response) {
	if cp.config.Verbose {
		statusColor := "\033[32m" 
		if resp.StatusCode >= 400 {
			statusColor = "\033[31m" 
		} else if resp.StatusCode >= 300 {
			statusColor = "\033[33m"
		}

		fmt.Printf("  \033[90m→\033[0m [%s%d\033[0m] \033[90m%s\033[0m",
			statusColor, resp.StatusCode, http.StatusText(resp.StatusCode))

		if cl := resp.Header.Get("Content-Length"); cl != "" {
			fmt.Printf(" \033[90m(%s bytes)\033[0m", cl)
		}
		fmt.Println()
	}
}

func (cp *CaptureProxy) analyzeBodyForURLs(body string, r *http.Request) {
	patterns := []string{
		`https?://[^\s"']+\.(mp4|mkv|m3u8|ts|mp3|flac|jpg|png|pdf|zip)`,
		`"url"\s*:\s*"([^"]+)"`,
		`"src"\s*:\s*"([^"]+)"`,
		`"file"\s*:\s*"([^"]+)"`,
		`"video_url"\s*:\s*"([^"]+)"`,
		`<source[^>]+src=["']([^"']+)["']`,
		`<video[^>]+src=["']([^"']+)["']`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(body, -1)
		for _, match := range matches {
			if len(match) > 1 {
				foundURL := match[1]
				if !strings.HasPrefix(foundURL, "http") {
					if strings.HasPrefix(foundURL, "//") {
						foundURL = "https:" + foundURL
					} else if strings.HasPrefix(foundURL, "/") && r != nil {
						foundURL = r.URL.Scheme + "://" + r.Host + foundURL
					}
				}
				log.Printf("\033[35m[HIDDEN]\033[0m Found in body: %s", foundURL)
				cp.capture(foundURL, r)
			}
		}
	}
}

func (cp *CaptureProxy) captureInterestingHeaders(headers http.Header, url string) {
	interestingHeaders := []string{
		"Location", "Content-Disposition", "X-Accel-Redirect",
		"X-Sendfile", "Content-Range", "Accept-Ranges",
	}

	for _, h := range interestingHeaders {
		if val := headers.Get(h); val != "" {
			log.Printf("\033[36m[HEADER]\033[0m %s: %s (from %s)", h, val, url)

			if h == "Location" || h == "X-Accel-Redirect" {
				cp.capture(val, nil)
			}
		}
	}
}

func (cp *CaptureProxy) getFileSizeWithTimeout(urlStr string, timeout time.Duration) int64 {
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Timeout:   timeout,
	}

	req, _ := http.NewRequest("HEAD", urlStr, nil)
	for k, v := range cp.config.Headers {
		req.Header.Set(k, v)
	}
	if cp.config.Cookie != "" {
		req.Header.Set("Cookie", cp.config.Cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	return resp.ContentLength
}

func (cp *CaptureProxy) capture(urlStr string, r *http.Request) {
	if urlStr == "" {
		return
	}

	ext := strings.ToLower(filepath.Ext(strings.Split(urlStr, "?")[0]))

	cleanURL := strings.Split(urlStr, "?")[0]

	fileType := cp.detectFileType(ext, cleanURL)
	if fileType == "" {
		return
	}

	if !cp.isFileTypeAllowed(fileType) {
		return
	}

	cp.mu.RLock()
	for _, item := range cp.captured {
		if item.URL == urlStr {
			cp.mu.RUnlock()
			return
		}
	}
	cp.mu.RUnlock()

	confidence := cp.calculateConfidence(ext, urlStr)
	if confidence < cp.config.ConfidenceLevel {
		return
	}

	var size int64 = -1
	if cp.config.AutoDownload {
		size = cp.getFileSizeWithTimeout(urlStr, 3*time.Second)
		if size < cp.config.MinFileSize {
			return
		}
		if cp.config.MaxFileSize > 0 && size > cp.config.MaxFileSize {
			return
		}
	}

	title := ""
	if r != nil {
		title = cp.extractTitle(urlStr, r.Header.Get("Referer"))
	} else {
		title = cp.extractTitle(urlStr, "")
	}

	item := CapturedItem{
		URL:        urlStr,
		FileType:   fileType,
		Extension:  ext,
		Size:       size,
		Title:      title,
		SourceURL:  "",
		Timestamp:  time.Now(),
		Confidence: confidence,
		Method:     "",
		StatusCode: 0,
	}

	if r != nil {
		item.Method = r.Method
		item.SourceURL = r.Header.Get("Referer")
	}

	cp.mu.Lock()
	cp.captured = append(cp.captured, item)
	cp.mu.Unlock()

	cp.displayCapturedItem(item)

	cp.saveToFile(item)

	if cp.config.AutoDownload && confidence >= cp.config.ConfidenceLevel {
		go cp.downloadItem(item)
	}
}


func (cp *CaptureProxy) displayCapturedItem(item CapturedItem) {
	color := "\033[32m" 
	if item.Confidence < 60 {
		color = "\033[33m" 
	}
	if item.Confidence < 40 {
		color = "\033[31m" 
	}

	icon := "📁"
	switch item.FileType {
	case TypeVideo:
		icon = "🎬"
	case TypeMusic:
		icon = "🎵"
	case TypeImage:
		icon = "🖼️"
	case TypeDocument:
		icon = "📄"
	case TypeArchive:
		icon = "🗜️"
	}

	fmt.Printf("\n%s \033[36m[%s]\033[0m %s | %d%% | %s",
		icon, strings.ToUpper(string(item.FileType)),
		color, item.Confidence,
		cp.formatSize(item.Size))

	if item.Title != "" && item.Title != "unknown" {
		fmt.Printf(" | \033[33m%s\033[0m", item.Title)
	}

	if item.Method != "" {
		fmt.Printf(" | \033[90m%s\033[0m", item.Method)
	}

	fmt.Printf("\n  \033[90m%s\033[0m\n", item.URL)
}

func (cp *CaptureProxy) detectFileType(ext, url string) FileType {
	urlLower := strings.ToLower(url)

	for fileType, exts := range ExtensionGroups {
		for _, e := range exts {
			if ext == e {
				return fileType
			}
		}
	}

	keywords := map[FileType][]string{
		TypeVideo:    {"video", "movie", "film", "stream", "play", "watch", "vod", "hls", "dash"},
		TypeMusic:    {"music", "audio", "song", "track", "album", "listen", "radio"},
		TypeImage:    {"image", "photo", "picture", "img", "gallery"},
		TypeDocument: {"document", "download", "file", "get", "dl"},
		TypeArchive:  {"archive", "package", "bundle"},
	}

	for fileType, kwList := range keywords {
		for _, kw := range kwList {
			if strings.Contains(urlLower, "/"+kw+"/") ||
				strings.Contains(urlLower, "/"+kw+"?") ||
				strings.HasSuffix(urlLower, "/"+kw) {
				return fileType
			}
		}
	}

	for _, customExt := range cp.config.CustomExtensions {
		if ext == customExt {
			return TypeAll
		}
	}

	return ""
}

func (cp *CaptureProxy) isFileTypeAllowed(ft FileType) bool {
	for _, allowed := range cp.config.FileTypes {
		if allowed == ft || allowed == TypeAll {
			return true
		}
	}
	return false
}

func (cp *CaptureProxy) calculateConfidence(ext, urlStr string) int {
	conf := 50

	for _, exts := range ExtensionGroups {
		for _, e := range exts {
			if ext == e {
				conf += 35
				break
			}
		}
	}

	urlLower := strings.ToLower(urlStr)
	keywords := []string{"download", "video", "music", "stream", "get", "file", "media", "content"}
	for _, kw := range keywords {
		if strings.Contains(urlLower, kw) {
			conf += 5
		}
	}

	if strings.Contains(urlLower, "/download/") ||
		strings.Contains(urlLower, "/get/") ||
		strings.Contains(urlLower, "/file/") {
		conf += 10
	}

	if conf > 100 {
		conf = 100
	}
	return conf
}

func (cp *CaptureProxy) getFileSize(urlStr string) int64 {
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Timeout:   5 * time.Second,
	}

	req, _ := http.NewRequest("HEAD", urlStr, nil)
	for k, v := range cp.config.Headers {
		req.Header.Set(k, v)
	}
	if cp.config.Cookie != "" {
		req.Header.Set("Cookie", cp.config.Cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	return resp.ContentLength
}

func (cp *CaptureProxy) extractTitle(urlStr, referer string) string {
	base := filepath.Base(strings.Split(urlStr, "?")[0])
	if idx := strings.LastIndex(base, "."); idx > 0 {
		base = base[:idx]
	}
	base = strings.ReplaceAll(base, "-", " ")
	base = strings.ReplaceAll(base, "_", " ")
	base = strings.ReplaceAll(base, "%20", " ")
	base = strings.TrimSpace(base)

	if len(base) > 3 && len(base) < 100 && !isGarbageTitle(base) {
		return base
	}

	if referer != "" {
		parts := strings.Split(referer, "/")
		for i := len(parts) - 1; i >= 0; i-- {
			part := strings.TrimSpace(strings.ReplaceAll(parts[i], "-", " "))
			if len(part) > 3 && len(part) < 100 && !strings.Contains(part, ".") && !isGarbageTitle(part) {
				return part
			}
		}
	}

	return "unknown"
}

func isGarbageTitle(title string) bool {
	garbage := []string{"index", "default", "main", "home", "api", "v1", "v2", "static", "assets"}
	titleLower := strings.ToLower(title)
	for _, g := range garbage {
		if titleLower == g || strings.HasPrefix(titleLower, g+" ") {
			return true
		}
	}
	return false
}

func (cp *CaptureProxy) saveToFile(item CapturedItem) {
	f, err := os.OpenFile(cp.config.SaveToFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	var line string
	if item.Title != "unknown" {
		line = fmt.Sprintf("[%s] [%s] %s | %d%% | %s | %s | %s\n",
			item.Timestamp.Format("2006-01-02 15:04:05"),
			strings.ToUpper(string(item.FileType)),
			item.Title,
			item.Confidence,
			cp.formatSize(item.Size),
			item.Method,
			item.URL)
	} else {
		line = fmt.Sprintf("[%s] [%s] %d%% | %s | %s | %s\n",
			item.Timestamp.Format("2006-01-02 15:04:05"),
			strings.ToUpper(string(item.FileType)),
			item.Confidence,
			cp.formatSize(item.Size),
			item.Method,
			item.URL)
	}

	f.WriteString(line)

	cp.saveToJSON(item)
}

func (cp *CaptureProxy) saveToJSON(item CapturedItem) {
	jsonFile := strings.TrimSuffix(cp.config.SaveToFile, ".txt") + ".json"

	var items []CapturedItem
	data, err := os.ReadFile(jsonFile)
	if err == nil {
		json.Unmarshal(data, &items)
	}

	items = append(items, item)

	newData, _ := json.MarshalIndent(items, "", "  ")
	os.WriteFile(jsonFile, newData, 0644)
}

func (cp *CaptureProxy) GetStats() map[string]interface{} {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["total_captured"] = len(cp.captured)
	stats["total_requests"] = len(cp.requestLogs)

	byType := make(map[FileType]int)
	for _, item := range cp.captured {
		byType[item.FileType]++
	}
	stats["by_type"] = byType

	return stats
}

func (cp *CaptureProxy) GetCapturedItems() []CapturedItem {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return cp.captured
}

func (cp *CaptureProxy) GetRequestLogs() []RequestLog {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	return cp.requestLogs
}

func (cp *CaptureProxy) ExportLogs(format string, filename string) error {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	switch format {
	case "json":
		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "  ")
		return encoder.Encode(cp.requestLogs)
	case "csv":
		f.WriteString("Timestamp,Method,URL,Host,Path\n")
		for _, log := range cp.requestLogs {
			f.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s\n",
				log.Timestamp.Format("2006-01-02 15:04:05"),
				log.Method,
				log.URL,
				log.Host,
				log.Path))
		}
	}
	return nil
}

func (cp *CaptureProxy) downloadItem(item CapturedItem) {
	log.Printf("\033[36m[DOWNLOAD]\033[0m Starting: %s\n", item.Title)

	global := NewGlobalStatus()

	oldOutDir := outDir
	oldNumThreads := numThreads
	oldMaxParallel := maxParallel
	oldVerbose := verbose
	oldRetries := retries
	oldTimeoutSec := timeoutSec
	oldEnableGzip := enableGzip

	outDir = cp.config.OutputDir
	numThreads = 4 
	maxParallel = 1 
	verbose = cp.config.Verbose
	retries = 3
	timeoutSec = 30
	enableGzip = true

	if item.Size > 100*1024*1024 {
		numThreads = 8
	} else if item.Size > 50*1024*1024 {
		numThreads = 6
	} else if item.Size > 10*1024*1024 {
		numThreads = 4
	} else if item.Size > 1024*1024 { 
		numThreads = 2
	} else {
		numThreads = 1
	}

	log.Printf("\033[33m[INFO]\033[0m Using %d threads for %s (size: %s)\n",
		numThreads, item.Title, cp.formatSize(item.Size))

	fileName := filepath.Base(strings.Split(item.URL, "?")[0])
	if fileName == "" || fileName == "/" || fileName == "." {
		fileName = fmt.Sprintf("captured_%d%s", time.Now().Unix(), item.Extension)
	}

	global.addFile(fileName, item.Size)

	go func() {
		client := createHTTPClient()

		if len(cp.config.Headers) > 0 {
		}

		downloadSingle(item.URL, client, global)

		outDir = oldOutDir
		numThreads = oldNumThreads
		maxParallel = oldMaxParallel
		verbose = oldVerbose
		retries = oldRetries
		timeoutSec = oldTimeoutSec
		enableGzip = oldEnableGzip
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			global.mu.RLock()
			for _, f := range global.files {
				if f.Name == fileName {
					if f.Status == "downloaded" {
						log.Printf("\033[32m[COMPLETE]\033[0m Downloaded: %s (%s)\n",
							item.Title, cp.formatSize(f.Done))
						return
					}
					if f.Done > 0 && f.Total > 0 {
						pct := float64(f.Done) * 100 / float64(f.Total)
						log.Printf("\033[36m[PROGRESS]\033[0m %s: %.1f%% (%s/%s)\n",
							item.Title, pct, cp.formatSize(f.Done), cp.formatSize(f.Total))
					}
					break
				}
			}
			global.mu.RUnlock()
		}
	}()
}

func (cp *CaptureProxy) formatSize(size int64) string {
	if size <= 0 {
		return "unknown"
	}
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(size)/1024)
	}
	if size < 1024*1024*1024 {
		return fmt.Sprintf("%.1fMB", float64(size)/(1024*1024))
	}
	return fmt.Sprintf("%.2fGB", float64(size)/(1024*1024*1024))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func (cp *CaptureProxy) Start() error {
	port := cp.config.Port
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	fmt.Printf("\033[36m╔════════════════════════════════════════════════════════════════╗\033[0m\n")
	fmt.Printf("\033[36m║              CAPTURE PROXY - ADVANCED MODE v2.0               ║\033[0m\n")
	fmt.Printf("\033[36m╚════════════════════════════════════════════════════════════════╝\033[0m\n\n")

	displayPort := port
	fmt.Printf("\033[32m✓\033[0m Proxy: \033[33m%s\033[0m\n", displayPort)
	fmt.Printf("\033[32m✓\033[0m Capturing: ")
	for i, ft := range cp.config.FileTypes {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("\033[33m%s\033[0m", ft)
	}
	fmt.Printf("\n")

	if cp.config.AutoDownload {
		fmt.Printf("\033[32m✓\033[0m Auto-download: \033[33menabled\033[0m (≥%d%%)\n", cp.config.ConfidenceLevel)
		fmt.Printf("\033[32m✓\033[0m Output: \033[33m%s\033[0m\n", cp.config.OutputDir)
	}

	fmt.Printf("\033[32m✓\033[0m Save file: \033[33m%s\033[0m\n", cp.config.SaveToFile)

	if cp.config.FilterDomain != "" {
		fmt.Printf("\033[32m✓\033[0m Domain filter: \033[33m%s\033[0m\n", cp.config.FilterDomain)
	}

	if cp.config.FilterPattern != "" {
		fmt.Printf("\033[32m✓\033[0m URL pattern: \033[33m%s\033[0m\n", cp.config.FilterPattern)
	}

	if cp.config.CaptureBody {
		fmt.Printf("\033[33m⚠\033[0m Body capture: \033[33menabled\033[0m (may slow down proxy)\n")
	}

	fmt.Printf("\n\033[90mConfigure FoxyProxy:\033[0m\n")

	cleanPort := strings.TrimPrefix(port, ":")
	fmt.Printf("  • HTTP Proxy: localhost:%s\n", cleanPort)
	fmt.Printf("\n\033[33m🎯 Waiting for traffic...\033[0m\n\n")

	return http.ListenAndServe(port, cp.proxy)
}
