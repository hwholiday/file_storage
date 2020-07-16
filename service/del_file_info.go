package service

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
	if fileInfo != nil {
		return nil
	}
	if err = s.r.FileInfoServer.DelFileInfoByFid(fid); err != nil {
		return err
	}
	if err = s.r.StorageServer.DelFile(fid, fileInfo.BucketName); err != nil {
		return err
	}
	return nil
}
