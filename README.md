# 🚀 FAD - Fast Advanced Downloader

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS%20%7C%20ARM64-lightgrey)]()

A powerful, multi-threaded downloader with support for **HTTP/HTTPS**, **FTP/FTPS**, **SFTP**, **Metalink**, **MITM Capture Proxy**, **proxy connections**, **smart resume capabilities**, **JSON-RPC interface**, and **complete website backup tool**. Built for speed and reliability.

## 🎉 New in v3.3.5

FAD now includes a **MITM Capture Proxy** that intercepts HTTPS traffic to automatically detect and capture downloadable files, **batch download from captured JSON**, and **automatic CA certificate installation** for seamless HTTPS interception.

<a href="https://github.com/Mr-Spect3r/fad/releases/tag/3.3.5">Release v3.3.5</a>

## ✨ Features

### Core Downloader
- 🧵 **Multi-threaded downloads** - Maximize bandwidth utilization
- 📡 **Multiple protocols** - HTTP, HTTPS, FTP, FTPS, SFTP
- 🔄 **Resume support** - Interrupt and resume downloads seamlessly
- 🕸️ **Proxy support** - SOCKS4, SOCKS5, and HTTP proxies
- 📁 **Batch downloading** - Download from file lists
- 🕷️ **Web scraping** - Extract and download links from web pages
- 🔍 **Extension filtering** - Filter downloads by file extensions
- ⚡ **Adaptive buffering** - Automatically optimizes buffer sizes (16KB-1MB)
- 💾 **Session saving** - Save progress and resume later
- 🎨 **Beautiful progress bars** - Real-time visual feedback
- 🌍 **Cross-platform** - Windows, Linux, macOS, ARM64
- 🔐 **Integrity checks** - SHA256, SHA1, MD5 verification

### MITM Capture Proxy (NEW in v3.3.5)
- 🔒 **HTTPS interception** - Full Man-in-the-Middle proxy capabilities
- 🎯 **Auto-detection** - Automatically detects videos, music, images, documents, archives
- 🔍 **Hidden URL extraction** - Scrapes HTML, JSON, and JavaScript for buried links
- 📊 **Confidence scoring** - 0-100% confidence system to filter false positives
- 📝 **Multi-format output** - Saves captured links to both TXT and JSON
- 🔧 **Custom extensions** - Add your own file extensions to capture
- 🎯 **Domain filtering** - Focus on specific domains only
- 📐 **Size filtering** - Min/max file size constraints
- 🔄 **Auto-download** - Option to automatically download captured files
- 🍪 **Cookie support** - Pass cookies through the proxy
- 🔐 **Auto-certificate** - Automatic CA certificate installation

### Download from Captured JSON (NEW in v3.3.5)
- 📦 **Batch download** - Download all captured files with one command
- ⚡ **Smart threading** - Auto-adjusts threads based on file size
- 🔄 **Concurrent downloads** - Download multiple files simultaneously
- 📊 **Unified progress** - Track all downloads in a single beautiful interface
- 💾 **Resume support** - Interrupted downloads can be resumed
- 🧹 **Filename sanitization** - Automatically cleans titles for safe filenames

### Website Downloader
- 🌐 **Full site mirroring** - Crawl and backup entire websites
- 📄 **Single page backup** - Save page with all dependencies
- 🎯 **SPA support** - Handle hash-based routing (#!/ and #! paths)
- 🖼️ **Asset rewriting** - Automatic URL rewriting for offline browsing
- 🚀 **Concurrent crawling** - Configurable page and asset workers
- 💾 **Resumable crawls** - Save and resume interrupted backups
- 🎨 **CSS/JS processing** - Rewrite URLs in stylesheets and scripts
- 🖼️ **Iframe support** - Download iframe content recursively
- ⚡ **Rate limiting** - Per-domain request throttling

### Metalink Support (RFC 5854)
- 📦 **Version 3 & 4 support** - Full Metalink specification compliance
- 🔄 **Mirror selection** - Automatic priority-based mirror selection
- ✅ **Integrity verification** - Built-in checksum validation
- 📊 **Piece information** - File segmentation with hash verification

### RPC Interface
- 🔌 **JSON-RPC API** - Full remote control capabilities
- 🌐 **REST endpoints** - HTTP endpoints for status and control
- 📊 **Real-time monitoring** - Download progress and statistics
- 🎮 **Dynamic control** - Pause, resume, speed limiting via API

### SFTP Support
- 🔑 **SSH key authentication** - RSA, ECDSA, Ed25519 support
- 🔐 **Password authentication** - With fallback to keys
- 📁 **Full resume capability** - Interrupt and resume SFTP transfers
- ⏱️ **Configurable timeouts** - Connection and operation timeouts

## 📸 Screenshots

| Feature | Preview |
|---------|---------|
| **Normal Download** | ![Normal](https://github.com/user-attachments/assets/e1cd29f0-2d04-4af2-9e71-0d1492fc3aba) |
| **Web Scraping** | ![Scraping](https://github.com/user-attachments/assets/48a0409b-8bea-48bf-93ee-fa4e2f66f32d) |
| **Website Backup** | ![Website Backup](https://github.com/user-attachments/assets/826d5701-1f31-4992-97c0-89545386b083) |
| **Progress Bars** | ![Progress](https://github.com/user-attachments/assets/3c39c0be-3e9e-4004-8d75-58962bd90158) |

## 📦 Installation

### Binary Downloads

**Windows:**
- `fad-windows-amd64.exe` - 64-bit
- `fad-windows-386.exe` - 32-bit

**Linux:**
- `fad-linux-amd64` - 64-bit
- `fad-linux-arm64` - ARM64 (Raspberry Pi, ARM servers)
- `fad-linux-386` - 32-bit

**macOS:**
- `fad-darwin-amd64` - Intel Mac
- `fad-darwin-arm64` - Apple Silicon (M1/M2/M3)

### Go Installation

```bash
go install github.com/Mr-Spect3r/fad@3.3.5
```

### From Source

```bash
git clone https://github.com/Mr-Spect3r/fad.git
cd fad
go build -o fad .
```

## 🚀 Quick Start

### MITM Capture Proxy (NEW)

Capture download links while browsing:

```bash
# Start capture proxy on port 8085
./fad -capture-proxy :8085 -capture-types video,music

# Capture with custom extensions
./fad -capture-proxy :9090 -capture-types video,archive -capture-exts .webm,.mka

# Auto-download captured files
./fad -capture-proxy :8085 -capture-auto -capture-output ./downloads

# Filter by domain and confidence
./fad -capture-proxy :8085 -filter-domain example.com -capture-confidence 50

# Install certificate only
./fad -install-cert
```

### Download from Captured JSON (NEW)

After capturing links, download everything with one command:

```bash
# Download all captured files
./fad -download-json captured_links.json

# Custom output with concurrent downloads
./fad -download-json captured_links.json -o ./videos -u 5

# High performance (8 threads per file, 4 concurrent)
./fad -download-json captured_links.json -t 8 -u 4
```

### Basic Downloads

```bash
# Download a single file
./fad https://example.com/file.zip

# Download with 16 threads
./fad -t 16 https://example.com/large-file.zip

# Download multiple files
./fad https://example.com/file1.zip https://example.com/file2.zip

# Download from file list
./fad -f urls.txt

# Download with speed limit (1MB/s)
./fad -max-speed 1048576 https://example.com/file.zip

# Download with checksum verification
./fad -checksum-sha256 abc123... https://example.com/file.zip
```

### Website Backup

```bash
# Basic full site backup
./fad web -url https://example.com -mode full

# Backup entire website to specific directory
./fad web -url https://example.com -mode full -output ./backup

# Backup single page with all assets
./fad web -url https://example.com/about -mode single -download-external

# Backup with external CDN assets
./fad web -url https://example.com -mode full -download-external -external-domains cdn.example.com,images.example.com

# High performance crawling (10 concurrent workers)
./fad web -url https://example.com -mode full -concurrency 10 -max-pages 500

# Resume interrupted backup
./fad web -url https://example.com -mode full -resume -output ./backup

# SPA with hash routing support
./fad web -url https://app.example.com/#!/home -mode full -crawl-hash-routes

# Limit asset size and rate
./fad web -url https://example.com -mode full -max-asset-size 20 -rate-limit 5
```

### Metalink Downloads

```bash
# Download from Metalink URL
./fad -metalink https://example.com/file.metalink

# Download from local Metalink file
./fad -metalink ./downloads/ubuntu.metalink4

# Metalink with custom output directory
./fad -metalink https://example.com/file.metalink -o ./downloads
```

### RPC Server Mode

```bash
# Start RPC server on default port
./fad -rpc

# Start RPC server on custom address
./fad -rpc -rpc-addr 0.0.0.0:6800

# Enable WebSocket RPC (experimental)
./fad -rpc -rpc-websocket -rpc-addr :6800

# RPC with downloads directory
./fad -rpc -rpc-addr localhost:6800 -o /downloads
```

**RPC Example Requests:**

```bash
# Get version info
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"fad.version","id":1}'

# Get global statistics
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"fad.getGlobalStat","id":2}'

# Add download URL
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"fad.addUri","params":{"uris":["https://example.com/file.zip"]},"id":3}'

# Get all files status
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"fad.tellAllStatus","id":4}'

# List all available methods
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"system.listMethods","id":5}'

# Pause all downloads
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"fad.pauseAll","id":6}'

# Set speed limit to 5MB/s
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"fad.setSpeedLimit","params":{"speed":5242880},"id":7}'

# Shutdown FAD
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"fad.shutdown","id":8}'
```

**REST API Endpoints:**

```bash
# Get global status
curl http://localhost:6800/api/status

# Get all files
curl http://localhost:6800/api/files

# Pause all downloads
curl http://localhost:6800/api/pause

# Resume all downloads
curl http://localhost:6800/api/resume

# Version info
curl http://localhost:6800/api/version
```

### 🔐 Proxy Downloads

```bash
# SOCKS5 proxy
./fad -proxy socks5://127.0.0.1:1080 https://example.com/file.zip

# SOCKS5 with authentication
./fad -proxy socks5://user:pass@127.0.0.1:1080 https://example.com/file.zip

# SOCKS4 proxy with custom threads
./fad -proxy socks4://192.168.1.1:9050 -t 16 https://example.com/file.zip

# HTTP proxy
./fad -proxy http://proxy.company.com:8080 https://example.com/file.zip

# HTTPS proxy with auth
./fad -proxy https://user:pass@proxy.company.com:8080 https://example.com/file.zip
```

### 📡 FTP/SFTP Downloads

```bash
# Standard FTP
./fad -protocol ftp ftp://example.com/file.zip

# FTP with custom credentials
./fad -protocol ftp -ftp-user myuser -ftp-pass mypass ftp://example.com/file.zip

# FTPS (FTP over TLS)
./fad -protocol ftps ftps://example.com/secure-file.zip

# SFTP with password
./fad -protocol sftp -sftp-user myuser -sftp-pass mypass sftp://example.com/file.zip

# SFTP with SSH key
./fad -protocol sftp -ssh-key ~/.ssh/id_rsa sftp://example.com/file.zip

# SFTP with encrypted SSH key
./fad -protocol sftp -ssh-key ~/.ssh/id_rsa -ssh-key-pass mypassphrase sftp://example.com/file.zip
```

### 🕷️ Web Scraping

```bash
# Extract and download all files from a page
./fad -scrape https://example.com/downloads/

# Filter by extensions
./fad -scrape https://example.com/downloads/ -ex .mp4,.mp3,.zip

# Scrape with custom threads
./fad -scrape https://example.com/files/ -t 16 -ex .pdf,.doc,.xls

# Scrape and download with progress
./fad -scrape https://example.com/media/ -ex .jpg,.png,.gif -v
```

### 🔄 Parameterized URLs

```bash
# Simple numeric placeholder
./fad -parameterized-url 'https://example.com/file{}.zip' -start 1 -end 50

# Zero-padded placeholders
./fad -parameterized-url 'https://example.com/image{0}.jpg' -start 1 -end 100

# Triple zero-padded
./fad -parameterized-url 'https://example.com/page{00}.html' -start 1 -end 500 -step 2

# Custom step size
./fad -parameterized-url 'https://example.com/chunk{}.bin' -start 0 -end 200 -step 10
```

### 🔄 Resume Downloads

```bash
# Resume from saved session
./fad session_20231215_143022.json

# Session auto-saves on interrupt (Ctrl+C)
# Progress saves every 10 seconds automatically
```

### 🍪 Cookie Support

```bash
# Load cookies from Netscape format file (Firefox/Chrome export)
./fad -load-cookies cookies.txt https://example.com/private-file.zip

# Save cookies after download
./fad -save-cookies output.txt https://example.com/file.zip

# Direct cookie string
./fad -c "sessionid=abc123; user=test" https://example.com/file.zip

# Load and save cookies
./fad -load-cookies cookies.txt -save-cookies newcookies.txt https://example.com/file.zip
```

### 🔐 NetRC Authentication

```bash
# Use .netrc file for authentication
./fad -netrc ~/.netrc https://example.com/private/file.zip

# .netrc file format:
# machine example.com login myuser password mypass
# default login anonymous password user@example.com
```

## ⚙️ Command Line Options

### Core Options

| Option | Default | Description |
|--------|---------|-------------|
| `-t` | CPU cores | Number of parallel download threads per file |
| `-o` | `.` | Destination directory for downloads |
| `-u` | `2` | Maximum simultaneous file downloads |
| `-r` | `5` | Retries per segment |
| `-timeout` | `30` | Network timeout in seconds |
| `-v` | `false` | Verbose mode with per-thread progress |
| `-save-session` | `true` | Save session to JSON if interrupted |
| `-f` | `""` | File containing download URLs (one per line) |

### Capture Proxy Options (NEW in v3.3.5)

| Option | Default | Description |
|--------|---------|-------------|
| `-capture-proxy` | `""` | Start MITM proxy (e.g., :8085) |
| `-capture-types` | `video,music` | File types: video,music,image,document,archive,all |
| `-capture-exts` | `""` | Custom extensions (comma-separated) |
| `-capture-auto` | `false` | Auto-download captured files |
| `-capture-output` | `captured` | Output directory for auto-downloads |
| `-capture-confidence` | `30` | Confidence threshold (0-100) |
| `-capture-min-size` | `1024` | Minimum file size in bytes |
| `-capture-max-size` | `0` | Maximum file size (0=unlimited) |
| `-capture-save` | `captured_links.txt` | File to save captured links |
| `-capture-header` | `""` | Custom HTTP headers (can be repeated) |
| `-capture-cookie` | `""` | Cookie for requests |
| `-filter-domain` | `""` | Filter specific domain |
| `-filter-pattern` | `""` | URL pattern filter |
| `-capture-body` | `false` | Capture request/response bodies |
| `-install-cert` | `true` | Auto-install CA certificate |

### Download from JSON Options (NEW in v3.3.5)

| Option | Default | Description |
|--------|---------|-------------|
| `-download-json` | `""` | Download all files from captured JSON file |
| `-u` | `3` | Max concurrent downloads (when using -download-json) |
| `-o` | `captured_downloads` | Output directory |
| `-t` | Auto | Threads per file (auto-adjusted by size) |

### Network Options

| Option | Description |
|--------|-------------|
| `-proxy` | Proxy address (socks4://, socks5://, http://) |
| `-protocol` | Force protocol: auto, http, https, ftp, ftps, sftp |
| `-H` | Custom HTTP header (can be repeated) |
| `-c` | Cookie header value |
| `-load-cookies` | Load cookies from Netscape format file |
| `-save-cookies` | Save cookies to file in Netscape format |
| `-netrc` | Path to .netrc file for authentication |
| `-gzip` | Enable gzip/deflate encoding (default: true) |

### Speed & Cache Options

| Option | Default | Description |
|--------|---------|-------------|
| `-max-speed` | `0` | Maximum download speed in bytes/second (0 = unlimited) |
| `-disk-cache` | `16MB` | Disk cache size in bytes |

### Integrity Check Options

| Option | Description |
|--------|-------------|
| `-check-integrity` | Verify file integrity after download |
| `-checksum-sha256` | Expected SHA256 hash for integrity check |
| `-checksum-md5` | Expected MD5 hash for integrity check |
| `-checksum-sha1` | Expected SHA1 hash for integrity check |

### FTP/SFTP Options

| Option | Default | Description |
|--------|---------|-------------|
| `-ftp-user` | `anonymous` | FTP/SFTP username |
| `-ftp-pass` | `anonymous@example.com` | FTP/SFTP password |
| `-ssh-user` | `""` | SSH username for SFTP |
| `-ssh-pass` | `""` | SSH password for SFTP |
| `-ssh-key` | `""` | SSH private key file for SFTP |
| `-ssh-key-pass` | `""` | SSH private key passphrase |
| `-ftp-multipart` | `true` | Enable FTP multi-part download |
| `-ftp-parts` | `0` | Number of FTP parts (0 = auto) |

### Website Downloader Options

| Option | Default | Description |
|--------|---------|-------------|
| `-url` | Required | Target URL to backup |
| `-output` | domain name | Output directory |
| `-mode` | `single` | Crawl mode: 'single' or 'full' |
| `-max-pages` | `100` | Maximum pages for full-site mode |
| `-concurrency` | `5` | Number of concurrent workers |
| `-download-external` | `false` | Download external assets |
| `-external-domains` | `""` | Comma-separated external domains to include |
| `-cookies` | `""` | Cookies (format: name1=value1; name2=value2) |
| `-user-agent` | `Mozilla/5.0...` | User-Agent header |
| `-timeout` | `30` | Request timeout in seconds |
| `-retries` | `3` | Number of retries on failure |
| `-minify` | `false` | Minify HTML output |
| `-resume` | `false` | Resume interrupted crawl |
| `-rate-limit` | `10` | Requests per second per domain |
| `-max-asset-size` | `50` | Maximum asset size in MB |
| `-crawl-iframes` | `true` | Download iframe content |
| `-crawl-hash-routes` | `true` | Handle hash-based routing for SPAs |

### Metalink Options

| Option | Description |
|--------|-------------|
| `-metalink` | Metalink URL or file path (RFC 5854) |

### RPC Options

| Option | Default | Description |
|--------|---------|-------------|
| `-rpc` | `false` | Enable JSON-RPC interface |
| `-rpc-addr` | `localhost:6800` | RPC server address |
| `-rpc-websocket` | `false` | Enable WebSocket RPC (experimental) |

### Parameterized URL Options

| Option | Default | Description |
|--------|---------|-------------|
| `-parameterized-url` | `""` | URL pattern with {} as placeholder |
| `-start` | `1` | Start index for parameterized URLs |
| `-end` | `100` | End index for parameterized URLs |
| `-step` | `1` | Step for parameterized URLs |

### Scraping Options

| Option | Description |
|--------|-------------|
| `-scrape` | URL to scrape for downloadable links |
| `-ex` | Filter extensions (e.g., .mp4,.mp3,.zip) |

### Daemon Options

| Option | Default | Description |
|--------|---------|-------------|
| `-daemon` | `false` | Run as daemon process in background |
| `-pid-file` | `/tmp/fad.pid` | PID file path for daemon mode |

## 📝 File Format Examples

### URLs File (urls.txt)

```text
# This is a comment
https://example.com/file1.zip
https://example.com/file2.zip
ftp://ftp.example.com/large-file.iso
https://example.com/document.pdf
sftp://sftp.example.com/backup.tar.gz
```

### Captured Links JSON (captured_links.json)

```json
[
  {
    "URL": "https://example.com/video.mp4",
    "FileType": "video",
    "Extension": ".mp4",
    "Size": 104857600,
    "Title": "sample video",
    "Confidence": 85,
    "Timestamp": "2024-01-15T10:30:00Z"
  },
  {
    "URL": "https://example.com/music.mp3",
    "FileType": "music",
    "Extension": ".mp3",
    "Size": 5242880,
    "Title": "sample song",
    "Confidence": 90,
    "Timestamp": "2024-01-15T10:31:00Z"
  }
]
```

### Session Files

Session files are auto-saved as `{filename}.json` and `{filename}.progress`. To resume:

```bash
./fad file.zip.json
```

### .netrc File Format

```text
machine example.com
login myusername
password mysecretpass

machine github.com
login mytoken
password ghp_xxxxxxxxxxxx

default
login anonymous
password user@example.com
```

### Netscape Cookies File Format

```text
# Netscape HTTP Cookie File
.example.com	TRUE	/	FALSE	1735689600	SESSION	abc123def456
.example.com	TRUE	/	TRUE	1735689600	SECURE	token789
```

## 🎯 Workflow Example

**Step 1: Start the capture proxy**

```bash
./fad -capture-proxy :8085 -capture-types video -capture-body -verbose
```

**Step 2: Configure your browser**

Set FoxyProxy or system proxy to `localhost:8085`

**Step 3: Browse normally**

The proxy captures all video links automatically and saves them to `captured_links.txt` and `captured_links.json`

**Step 4: Download all captured files**

```bash
./fad -download-json captured_links.json -o ./videos -u 5 -t 8
```

## 🎨 Output Preview

### Normal Download Progress

```text
══════════════════════════════════════════════════════════════════
                      DOWNLOAD STATUS
══════════════════════════════════════════════════════════════════
⬇️ 1. large-file.zip ████████████████████░░░░░░░░░░░░ 65.2%  1.2GB/1.8GB  ⚡ 12.5 MB/s  ETA: 45s
⬇️ 2. document.pdf   ████████████████████████████░░░░ 82.1%  8.2MB/10.0MB  ⚡ 2.3 MB/s  ETA: 2s
⏳ 3. image.jpg      ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 0.0%  0B/5.2MB       ⏳ Waiting
──────────────────────────────────────────────────────────────────
Avg Speed: 14.8 MB/s  Instant: 12.5 MB/s  Active: 2
Files: 1/3  Downloaded: 1.2GB / 1.8GB (65.2%)  Elapsed: 45.2s
Remaining: 45s  Left: 0.6GB
```

### Capture Proxy Output

```text
╔════════════════════════════════════════════════════════════════╗
║              CAPTURE PROXY - ADVANCED MODE v3.3.5             ║
╚════════════════════════════════════════════════════════════════╝

✓ Proxy: :8085
✓ Capturing: video, music
✓ Save file: captured_links.txt
⚠ Body capture: enabled (may slow down proxy)

Configure FoxyProxy:
  • HTTP Proxy: localhost:8085

🎯 Waiting for traffic...

[HIDDEN] Found in body: https://cdn.example.com/video.mp4
🎬 [VIDEO]  | 85% | unknown | sample video | GET
  https://cdn.example.com/video.mp4
```

### Verbose Mode Output (Per-thread Progress)

```text
     └─ Threads Progress [6/8 completed]:
        T1: ✅ [██████████] Complete
        T2: ✅ [██████████] Complete
        T3: ⬇️ [████████░░] Downloading (85.3%)
        T4: ⬇️ [██████░░░░] Downloading (62.1%)
        T5: ✅ [██████████] Complete
        T6: ✅ [██████████] Complete
        T7: ✅ [██████████] Complete
        T8: ⏳ [░░░░░░░░░░] Waiting
```

### Website Downloader Output

```text
[*] Mode: Full Site (max 100 pages)
[*] Full-site mode: crawling example.com
[1] https://example.com/
[2] https://example.com/about
[3] https://example.com/contact
  [↓] style.css
  [↓] script.js
  [↓] logo.png
  [↓] font.woff2

[✓] Backup completed in 12.5s
    Pages: 15 | Assets: 87 | Size: 24.32 MB
```

### RPC Server Output

```text
[INFO] Starting JSON-RPC server on localhost:6800
[INFO] WebSocket RPC: disabled
[INFO] Available methods:
  - fad.addUri
  - fad.remove
  - fad.tellStatus
  - fad.tellAllStatus
  - fad.getGlobalStat
  - fad.pause
  - fad.pauseAll
  - fad.resume
  - fad.resumeAll
  - fad.setSpeedLimit
  - fad.shutdown
  - system.listMethods
```

## 🛠️ Building from Source

### Prerequisites

- Go 1.26 or higher
- GCC (for Windows builds)

### Build Commands

```bash
# Linux/macOS
go build -o fad main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o fad.exe main.go

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o fad-linux-arm64 main.go

# With optimizations (smaller binary)
go build -ldflags="-s -w" -o fad main.go

# Cross-compile all platforms
GOOS=linux GOARCH=amd64 go build -o fad-linux-amd64 main.go
GOOS=linux GOARCH=386 go build -o fad-linux-386 main.go
GOOS=linux GOARCH=arm64 go build -o fad-linux-arm64 main.go
GOOS=windows GOARCH=amd64 go build -o fad-windows-amd64.exe main.go
GOOS=windows GOARCH=386 go build -o fad-windows-386.exe main.go
GOOS=darwin GOARCH=amd64 go build -o fad-darwin-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o fad-darwin-arm64 main.go
```

## 🐛 Troubleshooting

### Common Issues

| Question | Solution |
|----------|----------|
| **Q: Capture proxy not capturing anything?** | • Install CA certificate: `./fad -install-cert`<br>• Check browser proxy settings (localhost:port)<br>• Ensure HTTPS sites show certificate warning (accept it)<br>• Use `-capture-body` to see hidden links |
| **Q: Download from JSON not working?** | • Verify captured_links.json exists<br>• Check URLs are valid in the JSON file<br>• Try with `-v` flag for debug output<br>• Ensure output directory is writable |
| **Q: Certificate installation fails?** | • Run as administrator/root<br>• Follow manual instructions displayed<br>• Check if Fad.crt file was created<br>• Disable antivirus temporarily |
| **Q: Slow download speeds?** | • Increase thread count: `-t 16`<br>• Check if server supports range requests<br>• Try FTP multi-part for FTP files<br>• Reduce proxy latency if using proxy |
| **Q: Proxy not working?** | • Verify proxy format: `socks5://host:port`<br>• Ensure proxy is reachable<br>• Try without authentication first<br>• Check firewall settings |
| **Q: Resume not working?** | • Ensure server supports `Accept-Ranges: bytes`<br>• Check if session files exist in download directory<br>• Manual resume: `./fad session_file.json`<br>• Some servers don't support resume |
| **Q: Website crawl hangs?** | • Use `-resume` flag to continue interrupted crawl<br>• Reduce concurrency: `-concurrency 3`<br>• Increase timeout: `-timeout 60`<br>• Add rate limiting: `-rate-limit 5` |
| **Q: Assets not downloading?** | • Enable `-download-external` for external resources<br>• Add external domains with `-external-domains`<br>• Check `-max-asset-size` limit<br>• Verify network connectivity to CDN |
| **Q: SFTP connection fails?** | • Verify SSH key permissions (600)<br>• Check if server supports SFTP (not just SSH)<br>• Try password authentication first<br>• Increase timeout: `-timeout 60` |
| **Q: Metalink not working?** | • Verify file is valid XML/RFC 5854 compliant<br>• Check if URLs in metalink are accessible<br>• Try downloading individual URLs directly<br>• Validate checksums if provided |
| **Q: RPC server not responding?** | • Check if port is open: `netstat -an | grep 6800`<br>• Verify firewall allows the port<br>• Use `localhost` instead of `0.0.0.0` for testing<br>• Check logs for binding errors |
| **Q: Daemon mode not starting?** | • Verify PID file directory is writable<br>• Check if another instance is running<br>• Review system logs for errors<br>• Try running without daemon mode first |

## 📊 Performance Tips

1. **For large files (1GB+):** Use `-t` flag with 8-16 threads
2. **For many small files:** Use `-u` flag to increase parallel downloads (10-20)
3. **For slow connections:** Reduce threads to 2-4 and increase timeout to 60s
4. **For capture proxy:** Use `-capture-confidence 30` to balance accuracy and speed
5. **For batch downloads from JSON:** Use `-u 5` and `-t auto` (auto-adjusts by size)
6. **For website backup:** Start with `-mode single` to test, then use `-mode full`
7. **For rate-limited sites:** Use `-rate-limit` to avoid being blocked (5-10 req/s)
8. **For high bandwidth connections:** Increase threads to 32-64 for maximum speed
9. **For limited disk space:** Reduce `-disk-cache` to 8MB or lower
10. **For many concurrent downloads:** Use `-u` with 5-10 to avoid connection limits
11. **For SFTP transfers:** Use SSH keys instead of passwords for better performance
12. **For Metalink downloads:** Let FAD auto-select best mirrors

## 🔧 Advanced Usage

### Running as Daemon (Linux)

```bash
# Start daemon
./fad -daemon -o /downloads https://example.com/bigfile.zip

# Check status
cat /tmp/fad.pid

# Stop daemon
kill $(cat /tmp/fad.pid)

# Systemd service
sudo cat > /etc/systemd/system/fad.service << EOF
[Unit]
Description=FAD Downloader Service
After=network.target

[Service]
Type=simple
User=downloader
ExecStart=/usr/local/bin/fad -daemon -o /downloads -rpc
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable fad
sudo systemctl start fad
```

### Using Environment Variables

```bash
# Set proxy via environment
export HTTP_PROXY=http://proxy:8080
export HTTPS_PROXY=http://proxy:8080
export NO_PROXY=localhost,127.0.0.1

# Run with environment
./fad https://example.com/file.zip
```

### Combining with Other Tools

```bash
# Download and extract
./fad https://example.com/archive.zip && unzip archive.zip

# Download and verify signature
./fad https://example.com/file.iso && sha256sum -c file.iso.sha256

# Download list from another program output
cat urls.txt | xargs -n1 ./fad

# Scheduled downloads with cron
0 2 * * * /usr/local/bin/fad https://example.com/daily-backup.zip -o /backups
```

### Custom Headers Example

```bash
# Multiple custom headers
./fad -H "Authorization: Bearer token123" -H "X-Custom-Header: value" https://api.example.com/file.zip

# User-Agent override
./fad -H "User-Agent: MyCustomBot/1.0" https://example.com/file.zip

# Referer header
./fad -H "Referer: https://google.com" https://example.com/file.zip
```

### Capture Proxy with Custom Headers

```bash
# Add authentication headers to capture proxy
./fad -capture-proxy :8085 -capture-header "Authorization: Bearer token123" -capture-header "X-API-Key: abc123"

# Pass cookies through proxy
./fad -capture-proxy :8085 -capture-cookie "sessionid=abc123; user=test"
```

## 🔒 Security Note

The capture proxy uses a self-signed CA certificate to intercept HTTPS traffic. This certificate must be trusted by your system for HTTPS interception to work. The `-install-cert` flag attempts to install it automatically, but you may need to do it manually depending on your system permissions.

**Important:** Only use this tool on networks and websites you own or have permission to test. Intercepting HTTPS traffic without authorization may violate terms of service or laws in your jurisdiction.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Development Setup

```bash
git clone https://github.com/Mr-Spect3r/fad.git
cd fad
go mod download
go build -o fad main.go
./fad -v https://example.com/test.zip
```

## 🙏 Acknowledgments

- Thanks to all contributors and users who reported issues
- Built with ❤️ using Go
- Special thanks to the goproxy library for MITM proxy capabilities

---

**⭐ Star this repository if you find it useful!**

[Report Bug](https://github.com/Mr-Spect3r/fad/issues) · [Request Feature](https://github.com/Mr-Spect3r/fad/issues) · [View Releases](https://github.com/Mr-Spect3r/fad/releases)
