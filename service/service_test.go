package service

import (
	storage "filesrv/api/pb"
	"filesrv/conf"
	"filesrv/library/log"
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = flag.Set("conf", "./../cmd/filesrv.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.NewLogger(conf.Conf.Log)
	NewService(conf.Conf)
	os.Exit(m.Run())
}

func TestService_ApplyFid(t *testing.T) {
	Convey("TestService_ApplyFid", t, func() {
		out, err := GetService().ApplyFid(&storage.InApplyFid{
			Name:        "qaz.jpg",
			Size:        123451123,
			ExName:      "jpg",
			Md5:         "qwertyu",
			SliceTotal:  10,
			ExpiredTime: 0,
		})
		So(err, ShouldBeNil)
		So(out, ShouldNotBeNil)
		t.Log(out)
	})
}
