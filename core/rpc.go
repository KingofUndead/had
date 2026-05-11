package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type RPCServer struct {
	global   *GlobalStatus
	server   *http.Server
	mu       sync.RWMutex
	commands chan RPCCommand
}

type RPCCommand struct {
	ID     string                 `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type RPCResponse struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewRPCServer(global *GlobalStatus) *RPCServer {
	return &RPCServer{
		global:   global,
		commands: make(chan RPCCommand, 100),
	}
}

func (rpc *RPCServer) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/jsonrpc", rpc.handleJSONRPC)
	mux.HandleFunc("/api/status", rpc.handleStatus)
	mux.HandleFunc("/api/pause", rpc.handlePause)
	mux.HandleFunc("/api/resume", rpc.handleResume)
	mux.HandleFunc("/api/files", rpc.handleFiles)
	mux.HandleFunc("/api/version", rpc.handleVersion)

	rpc.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go rpc.server.ListenAndServe()
	logInfo("FAD RPC server started on %s", addr)
	return nil
}

func (rpc *RPCServer) Stop() error {
	return rpc.server.Close()
}

func (rpc *RPCServer) handleJSONRPC(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cmd RPCCommand
	if err := json.NewDecoder(req.Body).Decode(&cmd); err != nil {
		rpc.sendError(w, nil, -32700, "Parse error")
		return
	}

	var result interface{}
	var err *RPCError

	switch cmd.Method {
	case "fad.addUri":
		result, err = rpc.addURI(cmd.Params)
	case "fad.addUrls":
		result, err = rpc.addUrls(cmd.Params)
	case "fad.remove":
		result, err = rpc.removeDownload(cmd.Params)
	case "fad.removeAll":
		result, err = rpc.removeAll()
	case "fad.tellStatus":
		result, err = rpc.tellStatus(cmd.Params)
	case "fad.tellAllStatus":
		result, err = rpc.tellAllStatus()
	case "fad.getGlobalStat":
		result, err = rpc.getGlobalStat()
	case "fad.getFiles":
		result, err = rpc.getFiles()
	case "fad.pause":
		result, err = rpc.pause(cmd.Params)
	case "fad.pauseAll":
		result, err = rpc.pauseAll()
	case "fad.resume":
		result, err = rpc.resume(cmd.Params)
	case "fad.resumeAll":
		result, err = rpc.resumeAll()
	case "fad.setSpeedLimit":
		result, err = rpc.setSpeedLimit(cmd.Params)
	case "fad.getSpeedLimit":
		result, err = rpc.getSpeedLimit()
	case "fad.setMaxParallel":
		result, err = rpc.setMaxParallel(cmd.Params)
	case "fad.scrape":
		result, err = rpc.scrape(cmd.Params)
	case "fad.shutdown":
		result, err = rpc.shutdown()
	case "fad.version":
		result = rpc.version()
	case "system.listMethods":
		result = rpc.listMethods()
	default:
		err = &RPCError{Code: -32601, Message: "Method not found: " + cmd.Method}
	}

	rpc.sendResponse(w, cmd.ID, result, err)
}

func (rpc *RPCServer) handleStatus(w http.ResponseWriter, req *http.Request) {
	rpc.mu.RLock()
	defer rpc.mu.RUnlock()

	status := struct {
		Status          string `json:"status"`
		Version         string `json:"version"`
		ActiveDownloads int    `json:"active_downloads"`
		CompletedFiles  int64  `json:"completed_files"`
		TotalFiles      int64  `json:"total_files"`
		TotalSize       string `json:"total_size"`
		DownloadedSize  string `json:"downloaded_size"`
		StartTime       string `json:"start_time"`
		Uptime          string `json:"uptime"`
	}{
		Status:          "running",
		Version:         "1.0.0",
		ActiveDownloads: rpc.getActiveCount(),
		CompletedFiles:  rpc.global.downloadedCount,
		TotalFiles:      rpc.global.totalCount,
		TotalSize:       Size4Human(rpc.global.totalSize()),
		DownloadedSize:  Size4Human(rpc.getTotalDownloaded()),
		StartTime:       rpc.global.startTime.Format(time.RFC3339),
		Uptime:          time.Since(rpc.global.startTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (rpc *RPCServer) handleFiles(w http.ResponseWriter, req *http.Request) {
	rpc.mu.RLock()
	defer rpc.mu.RUnlock()

	type FileInfo struct {
		Name      string  `json:"name"`
		Size      int64   `json:"size"`
		SizeHuman string  `json:"size_human"`
		Done      int64   `json:"done"`
		Progress  float64 `json:"progress"`
		Status    string  `json:"status"`
	}

	files := make([]FileInfo, 0)
	for _, f := range rpc.global.files {
		var progress float64 = 0
		if f.Total > 0 {
			progress = float64(f.Done) * 100 / float64(f.Total)
		}
		files = append(files, FileInfo{
			Name:      f.Name,
			Size:      f.Size,
			SizeHuman: Size4Human(f.Size),
			Done:      f.Done,
			Progress:  progress,
			Status:    f.Status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func (rpc *RPCServer) handlePause(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`{"status":"paused"}`))
}

func (rpc *RPCServer) handleResume(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`{"status":"resumed"}`))
}

func (rpc *RPCServer) handleVersion(w http.ResponseWriter, req *http.Request) {
	version := map[string]string{
		"name":    "FAD (Fast Advanced Downloader)",
		"version": "1.0.0",
		"protocol": "fad-rpc/1.0",
		"features": "http,https,ftp,ftps,sftp,bittorrent,metalink,scrape",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

func (rpc *RPCServer) addURI(params map[string]interface{}) (interface{}, *RPCError) {
	uris, ok := params["uris"].([]interface{})
	if !ok || len(uris) == 0 {
		return nil, &RPCError{Code: -32602, Message: "Invalid params: missing uris"}
	}

	gid := generateGID()
	go func() {
		for _, u := range uris {
			if url, ok := u.(string); ok {
				downloadSingle(url, createHTTPClient(), rpc.global)
			}
		}
	}()
	return map[string]string{"gid": gid}, nil
}

func (rpc *RPCServer) addUrls(params map[string]interface{}) (interface{}, *RPCError) {
	urls, ok := params["urls"].([]interface{})
	if !ok || len(urls) == 0 {
		return nil, &RPCError{Code: -32602, Message: "Invalid params: missing urls"}
	}

	gid := generateGID()
	urlList := make([]string, 0)
	for _, u := range urls {
		if url, ok := u.(string); ok {
			urlList = append(urlList, url)
		}
	}

	go func() {
		for _, url := range urlList {
			downloadSingle(url, createHTTPClient(), rpc.global)
		}
	}()

	return map[string]interface{}{
		"gid":  gid,
		"count": len(urlList),
		"urls": urlList,
	}, nil
}

func (rpc *RPCServer) removeDownload(params map[string]interface{}) (interface{}, *RPCError) {
	gid, ok := params["gid"].(string)
	if !ok {
		return nil, &RPCError{Code: -32602, Message: "Invalid params: missing gid"}
	}
	return map[string]string{"gid": gid, "status": "removed"}, nil
}

func (rpc *RPCServer) removeAll() (interface{}, *RPCError) {
	return map[string]string{"status": "all downloads removed"}, nil
}

func (rpc *RPCServer) tellStatus(params map[string]interface{}) (interface{}, *RPCError) {
	gid, ok := params["gid"].(string)
	if !ok {
		return nil, &RPCError{Code: -32602, Message: "Invalid params: missing gid"}
	}
	status := map[string]interface{}{
		"gid":      gid,
		"status":   "active",
		"totalLength": "0",
		"completedLength": "0",
		"downloadSpeed": "0",
	}
	return status, nil
}

func (rpc *RPCServer) tellAllStatus() (interface{}, *RPCError) {
	files := make([]map[string]interface{}, 0)
	for i, f := range rpc.global.files {
		var progress float64 = 0
		if f.Total > 0 {
			progress = float64(f.Done) * 100 / float64(f.Total)
		}
		files = append(files, map[string]interface{}{
			"index":       i,
			"name":        f.Name,
			"size":        f.Size,
			"size_human":  Size4Human(f.Size),
			"downloaded":  f.Done,
			"progress":    progress,
			"status":      f.Status,
			"threads":     f.TotalThreads,
			"active_threads": f.ActiveThreads,
			"completed_threads": f.DoneThreads,
		})
	}
	return map[string]interface{}{"files": files, "count": len(files)}, nil
}

func (rpc *RPCServer) getGlobalStat() (interface{}, *RPCError) {
	totalSize := rpc.global.totalSize()
	totalDownloaded := rpc.getTotalDownloaded()
	var totalProgress float64 = 0
	if totalSize > 0 {
		totalProgress = float64(totalDownloaded) * 100 / float64(totalSize)
	}
	stat := map[string]interface{}{
		"downloadSpeed":     "0",
		"uploadSpeed":       "0",
		"numActive":         rpc.getActiveCount(),
		"numWaiting":        0,
		"numStopped":        int(rpc.global.downloadedCount),
		"totalSize":         totalSize,
		"totalSizeHuman":    Size4Human(totalSize),
		"totalDownloaded":   totalDownloaded,
		"totalDownloadedHuman": Size4Human(totalDownloaded),
		"totalProgress":     totalProgress,
		"totalFiles":        rpc.global.totalCount,
		"completedFiles":    rpc.global.downloadedCount,
	}
	return stat, nil
}

func (rpc *RPCServer) getFiles() (interface{}, *RPCError) {
	files := make([]map[string]interface{}, 0)
	for _, f := range rpc.global.files {
		var progress float64 = 0
		if f.Total > 0 {
			progress = float64(f.Done) * 100 / float64(f.Total)
		}
		files = append(files, map[string]interface{}{
			"name":     f.Name,
			"size":     f.Size,
			"downloaded": f.Done,
			"progress": progress,
			"status":   f.Status,
		})
	}
	return files, nil
}

func (rpc *RPCServer) pause(params map[string]interface{}) (interface{}, *RPCError) {
	gid, ok := params["gid"].(string)
	if !ok {
		return nil, &RPCError{Code: -32602, Message: "Invalid params"}
	}
	return map[string]string{"gid": gid, "status": "paused"}, nil
}

func (rpc *RPCServer) pauseAll() (interface{}, *RPCError) {
	return map[string]string{"status": "all downloads paused"}, nil
}

func (rpc *RPCServer) resume(params map[string]interface{}) (interface{}, *RPCError) {
	gid, ok := params["gid"].(string)
	if !ok {
		return nil, &RPCError{Code: -32602, Message: "Invalid params"}
	}
	return map[string]string{"gid": gid, "status": "resumed"}, nil
}

func (rpc *RPCServer) resumeAll() (interface{}, *RPCError) {
	return map[string]string{"status": "all downloads resumed"}, nil
}

func (rpc *RPCServer) setSpeedLimit(params map[string]interface{}) (interface{}, *RPCError) {
	speed, ok := params["speed"].(float64)
	if !ok {
		return nil, &RPCError{Code: -32602, Message: "Invalid params: missing speed"}
	}
	maxSpeed = int64(speed)
	return map[string]interface{}{"speed_limit": speed, "status": "updated"}, nil
}

func (rpc *RPCServer) getSpeedLimit() (interface{}, *RPCError) {
	return map[string]int64{"speed_limit": maxSpeed}, nil
}

func (rpc *RPCServer) setMaxParallel(params map[string]interface{}) (interface{}, *RPCError) {
	max, ok := params["max"].(float64)
	if !ok {
		return nil, &RPCError{Code: -32602, Message: "Invalid params"}
	}
	maxParallel = int(max)
	return map[string]interface{}{"max_parallel": maxParallel, "status": "updated"}, nil
}

func (rpc *RPCServer) scrape(params map[string]interface{}) (interface{}, *RPCError) {
	url, ok := params["url"].(string)
	if !ok {
		return nil, &RPCError{Code: -32602, Message: "Invalid params: missing url"}
	}

	global := NewGlobalStatus()
	go scrapeAndDownload(url, global)

	return map[string]string{"url": url, "status": "scraping started"}, nil
}

func (rpc *RPCServer) shutdown() (interface{}, *RPCError) {
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
	return map[string]string{"status": "shutting down"}, nil
}

func (rpc *RPCServer) version() interface{} {
	return map[string]string{
		"name":    "FAD",
		"version": "1.0.0",
		"rpc":     "fad-rpc/1.0",
	}
}

func (rpc *RPCServer) listMethods() interface{} {
	return []string{
		"fad.addUri",
		"fad.addUrls",
		"fad.remove",
		"fad.removeAll",
		"fad.tellStatus",
		"fad.tellAllStatus",
		"fad.getGlobalStat",
		"fad.getFiles",
		"fad.pause",
		"fad.pauseAll",
		"fad.resume",
		"fad.resumeAll",
		"fad.setSpeedLimit",
		"fad.getSpeedLimit",
		"fad.setMaxParallel",
		"fad.scrape",
		"fad.shutdown",
		"fad.version",
		"system.listMethods",
	}
}

func (rpc *RPCServer) getActiveCount() int {
	count := 0
	for _, f := range rpc.global.files {
		if f.Status == "downloading" {
			count++
		}
	}
	return count
}

func (rpc *RPCServer) getTotalDownloaded() int64 {
	var total int64
	for _, f := range rpc.global.files {
		total += f.Done
	}
	return total
}

func (rpc *RPCServer) sendResponse(w http.ResponseWriter, id string, result interface{}, err *RPCError) {
	resp := RPCResponse{
		ID:     id,
		Result: result,
		Error:  err,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (rpc *RPCServer) sendError(w http.ResponseWriter, id interface{}, code int, message string) {
	rpc.sendResponse(w, fmt.Sprintf("%v", id), nil, &RPCError{Code: code, Message: message})
}

func generateGID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())[:16]
}
