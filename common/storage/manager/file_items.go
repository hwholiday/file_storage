package manager

import (
	"bytes"
	"filesrv/conf"
	"filesrv/library/log"
	"filesrv/library/utils"
	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"image"
	"image/jpeg"
	"sort"
	"sync"
	"time"
)

type FileItem struct {
	mu            *sync.Mutex
	Fid           int64  //文件ID
	BucketName    string //桶名
	Size          int64  //文件总大小
	UploadSize    int64  //已经上传大小
	Md5           string //文件MD5
	IsImage       bool
	SliceTotal    int   // 1   为不分片文件  (1~3000)
	SliceSize     int   //上传除开最后一片的大小,用来判断最后一片外的每片大小是否相等
	IsSuccess     bool  //上传完成
	AutoClearTime int64 //到这个点没上传完成,自动删除
	Items         map[int][]byte
	autoTime      *time.Timer
}

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
		m.send(s.Fid)
		//未上传完成被自动清理
		log.GetLogger().Debug("[NewFileItem] AutoClear", zap.Any("fid", s.Fid))
		_ = m.r.FileInfoServer.DelFileInfoByFid(s.Fid)
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
		sortItems []int
		data      = make([]byte, 0, f.Size)
		buffer    = bytes.NewBuffer(data)
	)

	for k, _ := range f.Items {
		sortItems = append(sortItems, k)
	}
	sort.Ints(sortItems)
	for _, v := range sortItems {
		buffer.Write(f.Items[v])
	}
	f.autoTime.Stop()
	defer func() {
		//不管是否上传完成结束都要删除内存中的数据
		m.send(f.Fid)
	}()
	if f.Md5 != utils.Md5(buffer.Bytes()) {
		return
	}
	//上传文件
	if err := m.r.StorageServer.UpFileNotSlice(f.Fid, f.BucketName, buffer.Bytes()); err != nil {
		log.GetLogger().Info("[NewFileItem] MergeUp UpFileNotSlice", zap.Any(f.BucketName, f.Fid))
		_ = m.r.FileInfoServer.DelFileInfoByFid(f.Fid)
		return
	}
	if err := m.r.FileInfoServer.UpdateFileInfoStatusByFid(f.Fid, conf.FileExists); err != nil {
		log.GetLogger().Info("[NewFileItem] MergeUp UpdateFileInfoStatusByFid", zap.Any(f.BucketName, f.Fid))
		_ = m.r.StorageServer.DelFile(f.Fid, f.BucketName)
		return
	}
	//生成缩略图
	if f.IsImage {
		f.UpThumbnail(buffer.Bytes())
	}
}

func (f *FileItem) UpThumbnail(data []byte) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.GetLogger().Error("[NewFileItem] UpThumbnail Decode", zap.Any(f.BucketName, f.Fid), zap.Error(err))
		return
	}
	// height 为 0 保持宽高比
	reImg := imaging.Thumbnail(img, conf.ThumbnailWidth, 0, imaging.NearestNeighbor)
	var buf bytes.Buffer
	if err = jpeg.Encode(&buf, reImg, nil); err != nil {
		log.GetLogger().Error("[NewFileItem] UpThumbnail Encode", zap.Any(f.BucketName, f.Fid), zap.Error(err))
		return
	}
	var thumbnailFid = utils.GetSnowFlake().GetId()
	if err := m.r.StorageServer.UpFileNotSlice(thumbnailFid, f.BucketName, buf.Bytes()); err != nil {
		log.GetLogger().Info("[NewFileItem] MergeUp UpFileNotSlice", zap.Any(f.BucketName, thumbnailFid))
		return
	}
	if err := m.r.FileInfoServer.UpdateFileInfoByFid(f.Fid, bson.D{{"$set", bson.D{
		{"ex_image.thumbnail_fid", thumbnailFid},
		{"ex_image.thumbnail_high", reImg.Bounds().Dx()},
		{"ex_image.thumbnail_width", reImg.Bounds().Dy()},
	}}}); err != nil {
		log.GetLogger().Info("[NewFileItem] MergeUp UpdateFileInfoByFid", zap.Any(f.BucketName, thumbnailFid))
		_ = m.r.StorageServer.DelFile(thumbnailFid, f.BucketName)
		return
	}
}
