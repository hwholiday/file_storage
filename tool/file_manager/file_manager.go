package file_manager

import "sync"

var fileManager *FileManager

type FileManager struct {
	fileItems *sync.Map
	clearItem chan int64
}

type FileUploadItem struct {
	Fid   int64
	Part  int
	Data  []byte
	Md5   string
	IsEnd bool
}

func NewFileManager() {
	fileManager = &FileManager{
		fileItems: new(sync.Map),
		clearItem: make(chan int64, 10),
	}
	go fileManager.run()
	return
}

func GetFileManager() *FileManager {
	return fileManager
}

func (f *FileManager) SendFidToChan(fid int64) {
	f.clearItem <- fid
}

func (f *FileManager) AddItem(item *FileItem) {
	_, ok := f.fileItems.Load(item.Fid)
	if ok {
		return
	}
	f.fileItems.Store(item.Fid, item)
}

func (f *FileManager) run() {
	for {
		select {
		case fid := <-f.clearItem:
			f.fileItems.Delete(fid)
		}
	}
}
