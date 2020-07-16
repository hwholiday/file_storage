package fileInfo

import (
	"filesrv/conf"
	"filesrv/entity"
	"filesrv/library/database/mongo"
	"filesrv/library/log"
	"filesrv/library/utils"
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"testing"
)

var fi *fileInfo

func TestMain(m *testing.M) {
	_ = flag.Set("conf", "./../../cmd/filesrv.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.NewLogger(conf.Conf.Log)
	fi = NewFileInfo(mongo.NewMongo(conf.Conf.Mongo), conf.Conf)
	os.Exit(m.Run())
}

func TestFileInfo_InsertFileInfo(t *testing.T) {
	Convey("TestFileInfo_InsertFileInfo", t, func() {
		f := &entity.FileInfo{
			Fid:         1,
			Name:        "123.jpg",
			BucketName:  "s_1",
			Size:        1234,
			ContentType: "123/jpg",
			Md5:         "123wdeqwe12313qd",
			ExName:      "jpg",
			IsImage:     true,
			ExImage: entity.ImageEx{
				High:           200,
				Width:          300,
				ThumbnailFid:   0,
				ThumbnailHigh:  0,
				ThumbnailWidth: 0,
			},
			SliceTotal:  20,
			ExpiredTime: 0,
			Status:      conf.FileExists,
			CreateTime:  utils.GetTimeUnix(),
			UpdateTime:  utils.GetTimeUnix(),
		}
		err := fi.InsertFileInfo(f)
		So(err, ShouldBeNil)
	})
}

func TestFileInfo_UpdateFileInfoStatusByFid(t *testing.T) {
	Convey("TestFileInfo_UpdateFileInfoStatusByFid", t, func() {
		err := fi.UpdateFileInfoStatusByFid(1, conf.FileUploading)
		So(err, ShouldBeNil)
	})
}

func TestFileInfo_UpdateFileInfoByFid(t *testing.T) {
	Convey("TestFileInfo_UpdateFileInfoByFid", t, func() {
		err := fi.UpdateFileInfoByFid(1, bson.D{{"$set", bson.D{
			{"ex_image.thumbnail_fid", 2},
			{"ex_image.thumbnail_high", 100},
			{"ex_image.thumbnail_width", 100},
		}}})
		So(err, ShouldBeNil)
	})
}

func TestFileInfo_GetFileInfoByFid(t *testing.T) {
	Convey("TestFileInfo_GetFileInfoByFid", t, func() {
		info, err := fi.GetFileInfoByFid(1)
		So(err, ShouldBeNil)
		t.Log(info)
	})
}

func TestFileInfo_GetFileInfoByMd5(t *testing.T) {
	Convey("TestFileInfo_GetFileInfoByMd5", t, func() {
		info, err := fi.GetFileInfoByMd5("123wdeqwe12313qd")
		So(err, ShouldBeNil)
		So(info, ShouldNotBeNil)
		t.Log(info)
	})
}
