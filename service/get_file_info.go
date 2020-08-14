package service

import (
	storage "filesrv/api/pb"
	"filesrv/conf"
	"filesrv/entity"
	"filesrv/library/utils"
)

func (s *service) GetFileInfoByFid(fid int64) (fileInfo *entity.FileInfo, err error) {
	if fileInfo, err = s.r.FileInfoServer.GetFileInfoByFid(fid); err == nil {
		if fileInfo != nil {
			if fileInfo.ExpiredTime > 0 && utils.GetTimeUnix() > fileInfo.ExpiredTime {
				//文件已经过期,删除文件
				fileInfo.Status = conf.FileExpired
				go s.DelFileByFidAndBucketName(fileInfo)
			}
		}
	}
	return
}

func (s *service) GetFileInfoByMd5(md5 string) (fileInfo *entity.FileInfo, err error) {
	if fileInfo, err = s.r.FileInfoServer.GetFileInfoByMd5(md5); err == nil {
		if fileInfo != nil {
			if fileInfo.ExpiredTime > 0 && utils.GetTimeUnix() > fileInfo.ExpiredTime {
				//文件已经过期,删除文件
				fileInfo.Status = conf.FileExpired
				go s.DelFileByFidAndBucketName(fileInfo)
			}
		}
	}
	return
}

func (s *service) GetFileInfoByMd5NotAutoClear(md5 string) (fileInfo *entity.FileInfo, err error) {
	fileInfo, err = s.r.FileInfoServer.GetFileInfoByMd5(md5)
	return
}

func (s *service) GetPbFileInfoByFid(fid int64) (fileInfo *storage.FileInfo, err error) {
	var info *entity.FileInfo
	if info, err = s.r.FileInfoServer.GetFileInfoByFid(fid); err != nil {
		return
	}
	fileInfo = s.convertFileInfoDataToPbFileInfo(info)
	return
}

func (s *service) GetPbFileInfoByMd5(md5 string) (fileInfo *storage.FileInfo, err error) {
	var info *entity.FileInfo
	if info, err = s.r.FileInfoServer.GetFileInfoByMd5(md5); err != nil {
		return
	}
	fileInfo = s.convertFileInfoDataToPbFileInfo(info)
	return
}

func (s *service) convertFileInfoDataToPbFileInfo(f *entity.FileInfo) *storage.FileInfo {
	st := &storage.FileInfo{
		Fid:         f.Fid,
		Name:        f.Name,
		BucketName:  f.BucketName,
		Size:        f.Size,
		ContentType: f.ContentType,
		Md5:         f.Md5,
		ExName:      f.ExName,
		IsImage:     f.IsImage,
		SliceTotal:  f.SliceTotal,
		ExpiredTime: f.ExpiredTime,
		Status:      f.Status,
		CreateTime:  f.CreateTime,
		UpdateTime:  f.UpdateTime,
	}
	if st.IsImage {
		st.ExImage = &storage.ImageEx{
			Height:          f.ExImage.Height,
			Width:           f.ExImage.Width,
			ThumbnailFid:    f.ExImage.ThumbnailFid,
			ThumbnailHeight: f.ExImage.ThumbnailHeight,
			ThumbnailWidth:  f.ExImage.ThumbnailWidth,
		}
	}
	return st
}
