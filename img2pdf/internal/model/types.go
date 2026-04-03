package model

type FileInfo struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Thumbnail string `json:"thumbnail"` // base64
}

type PDFOptions struct {
	PageSize string `json:"pageSize"` // "A4" or "original"
	Output   string `json:"output"`
}

type ImageOptions struct {
	Format    string `json:"format"`    // "jpg" or "png"
	DPI       int    `json:"dpi"`       // 72, 150, 300
	PageRange string `json:"pageRange"` // "all" or "1-5"
	OutputDir string `json:"outputDir"`
}
