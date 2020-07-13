package file_manager

import "sync"

type FileManager struct {
	mu    sync.Mutex
	items map[int][]byte
}

type FileUploadItem struct {
	Fid   int64
	Part  int
	Data  []byte
	Md5   string
	IsEnd bool
}
