package service

import (
	"bytes"
	storage "filesrv/api/pb"
	"filesrv/conf"
	"filesrv/entity"
	"filesrv/library/log"
	"filesrv/library/utils"
	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"image"
	"image/jpeg"
)

func (s *service) UpFile(in *storage.InUpFile) (err error) {
	if in.Md5 != utils.Md5(in.Data) {
		return conf.ErrMd5ChecksumInvalid
	}
	var f *entity.FileInfo
	f, err = s.GetFileInfoByFid(in.Fid)
	if err != nil {
		return
	}
	//上传文件
	if err = s.r.StorageServer.UpFile(f.Fid, f.BucketName, in.Data); err != nil {
		log.GetLogger().Info("[UpFile] UpFileNotSlice", zap.Any(f.BucketName, f.Fid))
		_ = s.r.FileInfoServer.DelFileInfoByFid(f.Fid)
		return
	}
	if err = s.r.FileInfoServer.UpdateFileInfoStatusByFid(f.Fid, conf.FileExists); err != nil {
		log.GetLogger().Info("[UpFile]  UpdateFileInfoStatusByFid", zap.Any(f.BucketName, f.Fid))
		_ = s.r.StorageServer.DelFile(f.Fid, f.BucketName)
		return
	}
	//生成缩略图
	if f.IsImage {
		var img image.Image
		img, _, err = image.Decode(bytes.NewReader(in.Data))
		if err != nil {
			log.GetLogger().Error("[NewFileItem] UpThumbnail Decode", zap.Any(f.BucketName, f.Fid), zap.Error(err))
			return
		}
		// height 为 0 保持宽高比
		reImg := imaging.Thumbnail(img, conf.ThumbnailWidth, conf.ThumbnailHeight, imaging.NearestNeighbor)
		var buf bytes.Buffer
		if err = jpeg.Encode(&buf, reImg, nil); err != nil {
			log.GetLogger().Error("[NewFileItem] UpThumbnail Encode", zap.Any(f.BucketName, f.Fid), zap.Error(err))
			return
		}
		var thumbnailFid = utils.GetSnowFlake().GetId()
		if err = s.r.StorageServer.UpFile(thumbnailFid, f.BucketName, buf.Bytes()); err != nil {
			log.GetLogger().Info("[NewFileItem] MergeUp UpFileNotSlice", zap.Any(f.BucketName, thumbnailFid))
			return
		}
		if err = s.r.FileInfoServer.UpdateFileInfoByFid(f.Fid, bson.D{{"$set", bson.D{
			{"ex_image.thumbnail_fid", thumbnailFid},
			{"ex_image.thumbnail_height", conf.ThumbnailWidth},
			{"ex_image.thumbnail_width", conf.ThumbnailHeight},
		}}}); err != nil {
			log.GetLogger().Info("[UpFile]  UpdateFileInfoByFid", zap.Any(f.BucketName, thumbnailFid))
			_ = s.r.StorageServer.DelFile(thumbnailFid, f.BucketName)
			return
		}
	}
	log.GetLogger().Debug("[UpFile]  success", zap.Any(f.BucketName, f.Fid))
	return
}
