package core

import (
	"encoding/xml"
	"io"
	"net/http"
)

type Metalink struct {
	XMLName xml.Name `xml:"metalink"`
	Files   []File   `xml:"file"`
}

type File struct {
	Name   string  `xml:"name,attr"`
	Size   int64   `xml:"size"`
	Hashes []Hash  `xml:"hash"`
	URLs   []URL   `xml:"url"`
	Pieces []Piece `xml:"pieces"`
}

type Hash struct {
	Type string `xml:"type,attr"`
	Data string `xml:",chardata"`
}

type URL struct {
	Location string `xml:"location,attr"`
	URL      string `xml:",chardata"`
	Priority int    `xml:"priority,attr"`
}

type Piece struct {
	Length int    `xml:"length,attr"`
	Type   string `xml:"type,attr"`
	Hash   string `xml:",chardata"`
}

type Metalink4 struct {
	XMLName xml.Name `xml:"metalink"`
	Files   []File4  `xml:"file"`
}

type File4 struct {
	Name string `xml:"name,attr"`
	Size int64  `xml:"size"`
	URLs []URL4 `xml:"url"`
}

type URL4 struct {
	URL string `xml:",chardata"`
}

func downloadMetalink(metalinkURL string, global *GlobalStatus) {
	resp, err := http.Get(metalinkURL)
	if err != nil {
		logError("Failed to download metalink: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logError("Failed to read metalink: %v", err)
		return
	}

	var metalink Metalink
	if err := xml.Unmarshal(body, &metalink); err != nil {
		var metalink4 Metalink4
		if err := xml.Unmarshal(body, &metalink4); err != nil {
			logError("Failed to parse metalink: %v", err)
			return
		}
		processMetalink4(metalink4, global)
		return
	}

	processMetalink(metalink, global)
}

func processMetalink(metalink Metalink, global *GlobalStatus) {
	for _, file := range metalink.Files {
		var bestURL string
		for _, url := range file.URLs {
			if bestURL == "" || url.Priority > 0 {
				bestURL = url.URL
			}
		}

		if bestURL == "" && len(file.URLs) > 0 {
			bestURL = file.URLs[0].URL
		}

		if bestURL != "" {
			logInfo("Downloading from metalink: %s", file.Name)
			downloadSingle(bestURL, createHTTPClient(), global)
		}
	}
}

func processMetalink4(metalink Metalink4, global *GlobalStatus) {
	for _, file := range metalink.Files {
		if len(file.URLs) > 0 {
			logInfo("Downloading from metalink4: %s", file.Name)
			downloadSingle(file.URLs[0].URL, createHTTPClient(), global)
		}
	}
}
