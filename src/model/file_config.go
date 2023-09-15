package model

const DEFAULT_CONFIG_NAME = "files.json"

type FileConfig struct {
	RawFile         RawFileInfo       `json:"raw_file"`
	ChunkedFileList []ChunkedFileInfo `json:"chunked_file_list"`
	EndPoint        []string          `json:"end_point"`
}

type RawFileInfo struct {
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
}

type ChunkedFileInfo struct {
	FileName string `json:"file_name"`
	Md5      string `json:"md5"`
	Size     int64  `json:"size"`
	Offset   int64  `json:"offset"`
}
