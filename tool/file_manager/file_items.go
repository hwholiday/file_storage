package file_manager

import (
	"bytes"
	"file_storage/library/utils"
	"sort"
	"strings"
	"sync"
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
}

var imageExName = []string{"JPG", "JPEG", "PNG"}

func NewFileItem(s *FileItem) *FileItem {
	s.IsSuccess = false
	s.ExpiredTime = utils.GetTimeUnix() + 60*30
	s.Items = make(map[int][]byte)
	s.mu = new(sync.Mutex)
	return s
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
	//上传文件
	//生成缩略图
	if needThumbnail {

	}
}
