package service

import (
	"filesrv/conf"
	"filesrv/entity"
	"filesrv/library/utils"
)

func (s *service) GetFileInfoByFid(fid int64) (fileInfo *entity.FileInfo, err error) {
	if fileInfo, err = s.r.FileInfoServer.GetFileInfoByFid(fid); err == nil {
		if fileInfo.ExpiredTime > 0 && utils.GetTimeUnix() > fileInfo.ExpiredTime {
			//文件已经过期,删除文件
			fileInfo.Status = conf.FileExpired
			go s.DelFileByFidAndBucketName(fid, fileInfo.BucketName)
		}
	}
	return
}

func (s *service) DelFileByFidAndBucketName(fid int64, bucketName string) (err error) {
	if err = s.r.FileInfoServer.DelFileInfoByFid(fid); err != nil {
		return err
	}
	if err = s.r.StorageServer.DelFile(fid, bucketName); err != nil {
		return err
	}
	return
}

func (s *service) DelFileByFid(fid int64) error {
	fileInfo, err := s.r.FileInfoServer.GetFileInfoByFid(fid)
	if err != nil {
		return err
	}
	if err = s.r.FileInfoServer.DelFileInfoByFid(fid); err != nil {
		return err
	}
	if err = s.r.StorageServer.DelFile(fid, fileInfo.BucketName); err != nil {
		return err
	}
	return nil
}
