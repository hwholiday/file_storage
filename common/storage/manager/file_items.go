package manager

import (
	"bytes"
	"filesrv/conf"
	"filesrv/library/utils"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type FileItem struct {
	mu            *sync.Mutex
	Fid           int64  //文件ID
	Size          int64  //文件总大小
	UploadSize    int64  //已经上传大小
	Md5           string //文件MD5
	ExName        string //文件扩展名
	SliceTotal    int    // 1   为不分片文件  (1~3000)
	SliceSize     int    //上传除开最后一片的大小,用来判断最后一片外的每片大小是否相等
	IsSuccess     bool   //上传完成
	AutoClearTime int64  //到这个点没上传完成,自动删除
	Items         map[int][]byte
	autoTime      *time.Timer
}

var imageExName = []string{"JPG", "JPEG", "PNG"}

func NewFileItem(s *FileItem) *FileItem {
	s.IsSuccess = false
	s.AutoClearTime = 60 * 30
	s.Items = make(map[int][]byte)
	s.mu = new(sync.Mutex)
	s.AutoClear()
	return s
}

func (s *FileItem) AutoClear() {
	s.autoTime = time.AfterFunc(time.Second*time.Duration(s.AutoClearTime), func() {
		if s == nil {
			return
		}
		if s.IsSuccess {
			return
		}
		GetFileManager().SendFidToChan(s.Fid)
	})
}

func (f *FileItem) AddItem(upItem *FileUploadItem) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.IsSuccess {
		return conf.ErrFileUploadCompleted
	}
	if upItem.Part < 1 && upItem.Part > 3000 {
		return conf.ErrFilePartsInvalid
	}
	dataLen := len(upItem.Data)
	if dataLen <= 0 {
		return conf.ErrFilePartEmpty
	}
	if dataLen > 524288 {
		return conf.ErrFilePartTooBig
	}
	if upItem.Part != f.SliceTotal { //不是最后一片
		if dataLen%1024 != 0 {
			return conf.ErrFilePartSizeInvalid1KB
		}
		if 524288%dataLen != 0 {
			return conf.ErrFilePartSizeInvalid512KB
		}
		if f.SliceSize == 0 {
			f.SliceSize = dataLen
		} else {
			if f.SliceSize != dataLen {
				return conf.ErrFilePartSizeChanged
			}
		}
	}
	if _, ok := f.Items[upItem.Part]; ok {
		return conf.ErrFilePartUploadCompleted
	}
	if upItem.Md5 != utils.Md5(upItem.Data) {
		return conf.ErrMd5ChecksumInvalid
	}
	f.Items[upItem.Part] = upItem.Data
	f.UploadSize += int64(len(upItem.Data))
	if f.UploadSize >= f.Size && len(f.Items) == f.SliceTotal {
		f.IsSuccess = true
		go f.MergeUp()
	}
	return nil
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
	if f.Md5 != utils.Md5(buffer.Bytes()) {
		fmt.Println("MD5不同")
		return
	}
	//生成缩略图
	if needThumbnail {

	}
	//处理完成删除该信息
}
