package file_manager

import (
	"bytes"
	"sort"
	"strings"
	"sync"
	"time"
)

type FileItem struct {
	mu          *sync.Mutex
	Fid         int64
	Name        string
	Size        int64 //文件总大小
	UploadSize  int64 //已经上传大小
	Md5         string
	ExName      string
	IsSuccess   bool  //上传完成
	ExpiredTime int64 //到这个点没上传完成,自动删除
	Items       map[int][]byte
	autoTime    *time.Timer
}

var imageExName = []string{"JPG", "JPEG", "PNG"}

func NewFileItem(s *FileItem) *FileItem {
	s.IsSuccess = false
	s.ExpiredTime = 60 * 30
	s.Items = make(map[int][]byte)
	s.mu = new(sync.Mutex)
	s.AutoClear()
	return s
}

func (s *FileItem) AutoClear() {
	s.autoTime = time.AfterFunc(time.Second*time.Duration(s.ExpiredTime), func() {
		if s == nil {
			return
		}
		if s.IsSuccess {
			return
		}
		GetFileManager().SendFidToChan(s.Fid)
	})
}

func (f *FileItem) AddItem(part int, data []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.IsSuccess {
		return
	}
	if _, ok := f.Items[part]; ok {
		return
	}
	f.Items[part] = data
	f.UploadSize += int64(len(data))
	if f.UploadSize >= f.Size {
		f.IsSuccess = true
		go f.MergeUp()
	}
}

func (f *FileItem) MergeUp() {
	var (
		exName        = strings.ToUpper(f.ExName)
		needThumbnail bool
		sortItems     []int
		data          = make([]byte, 0, f.Size)
		buffer        = bytes.NewBuffer(data)
	)
	for _, v := range imageExName {
		if v == exName {
			needThumbnail = true
		}
	}
	for k, _ := range f.Items {
		sortItems = append(sortItems, k)
	}
	sort.Ints(sortItems)
	for _, v := range sortItems {
		buffer.Write(f.Items[v])
	}
	f.autoTime.Stop()
	defer func() {
		GetFileManager().SendFidToChan(f.Fid)
	}()
	//上传文件
	//生成缩略图
	if needThumbnail {

	}
	//处理完成删除该信息
}
