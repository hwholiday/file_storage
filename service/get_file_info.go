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

func (s *service) GetFileInfoByMd5(md5 string) (fileInfo *entity.FileInfo, err error) {
	if fileInfo, err = s.r.FileInfoServer.GetFileInfoByMd5(md5); err == nil {
		if fileInfo.ExpiredTime > 0 && utils.GetTimeUnix() > fileInfo.ExpiredTime {
			//文件已经过期,删除文件
			fileInfo.Status = conf.FileExpired
			go s.DelFileByFidAndBucketName(fileInfo.Fid, fileInfo.BucketName)
		}
	}
	return
}
