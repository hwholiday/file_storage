package service

import (
	"filesrv/entity"
)

func (s *service) DelFileByFidAndBucketName(fileInfo *entity.FileInfo) (err error) {
	if err = s.r.FileInfoServer.DelFileInfoByFid(fileInfo.Fid); err != nil {
		return err
	}
	if fileInfo.ExImage.ThumbnailFid != 0 { //删除
		if err = s.r.StorageServer.DelFile(fileInfo.ExImage.ThumbnailFid, fileInfo.BucketName); err != nil {
			return err
		}
	}
	if err = s.r.StorageServer.DelFile(fileInfo.Fid, fileInfo.BucketName); err != nil {
		return err
	}
	return
}

func (s *service) DelFileByFid(fid int64) error {
	fileInfo, err := s.r.FileInfoServer.GetFileInfoByFid(fid)
	if err != nil {
		return err
	}
	if fileInfo != nil {
		return nil
	}
	if err = s.r.FileInfoServer.DelFileInfoByFid(fid); err != nil {
		return err
	}
	if fileInfo.ExImage.ThumbnailFid != 0 { //删除
		if err = s.r.StorageServer.DelFile(fileInfo.ExImage.ThumbnailFid, fileInfo.BucketName); err != nil {
			return err
		}
	}
	if err = s.r.StorageServer.DelFile(fid, fileInfo.BucketName); err != nil {
		return err
	}
	return nil
}
