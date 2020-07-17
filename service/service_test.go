package service

import (
	"crypto/md5"
	storage "filesrv/api/pb"
	"filesrv/conf"
	"filesrv/library/log"
	"flag"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"sync"
	"testing"
	"time"
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

var waitSync sync.WaitGroup

func TestService_ApplyFid(t *testing.T) {
	Convey("ApplyFid", t, func() {
		fileData, err := os.Open("./../res/665419.png")
		So(err, ShouldBeNil)
		defer fileData.Close()
		fi, err := fileData.Stat()
		So(err, ShouldBeNil)
		buf := make([]byte, fi.Size())
		_, err = fileData.Read(buf)
		So(err, ShouldBeNil)
		var sliceSize uint32
		num := fi.Size() / (512 * 1024)
		sliceSize = 512 * 1024
		if fi.Size()%(512*1024) != 0 {
			num++
		}
		out, err := GetService().ApplyFid(&storage.InApplyFid{
			Name:        fi.Name(),
			Size:        fi.Size(),
			ExName:      "png",
			Md5:         fmt.Sprintf("%x", md5.Sum(buf)),
			SliceTotal:  int32(num),
			Height:      800,
			Width:       800,
			ExpiredTime: 0,
		})
		So(err, ShouldBeNil)
		So(out, ShouldNotBeNil)
		t.Log(out)
		if out.Status == conf.FileExists {
			t.Log("file exists")
			return
		}
		Convey("UpSliceFile", func() {
			var startTime = time.Now().UnixNano() / 1e6
			for i := 1; i <= int(num); i++ {
				waitSync.Add(1)
				go func(index int) {
					data := (index - 1) * int(sliceSize)
					if int(num) == index {
						t.Log("send ", "fid ", out.Fid, "index ", index)
						err = GetService().UpSliceFile(&storage.InUpSliceFileItem{
							Fid:  out.Fid,
							Part: int32(index),
							Data: buf[data:],
							Md5:  fmt.Sprintf("%x", md5.Sum(buf[data:])),
						})
						if err != nil {
							t.Error(err)
						}
					} else {
						t.Log("send ", "fid ", out.Fid, "index ", index)
						err = GetService().UpSliceFile(&storage.InUpSliceFileItem{
							Fid:  out.Fid,
							Part: int32(index),
							Data: buf[data : data+int(sliceSize)],
							Md5:  fmt.Sprintf("%x", md5.Sum(buf[data:data+int(sliceSize)])),
						})
						if err != nil {
							t.Error(err)
						}
					}
					waitSync.Done()
				}(i)
			}
			waitSync.Wait()
			t.Log("time consuming (millisecond)", time.Now().UnixNano()/1e6-startTime)
		})
	})
	select {}
}

func TestService_GetPbFileInfoByFid(t *testing.T) {
	Convey("TestService_GetPbFileInfoByFid", t, func() {
		fileInfo, err := GetService().GetPbFileInfoByFid(2261944530632704)
		So(err, ShouldBeNil)
		So(fileInfo, ShouldNotBeNil)
		t.Log(fileInfo)
	})
}
