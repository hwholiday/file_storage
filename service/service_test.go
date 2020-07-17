package service

import (
	"bytes"
	"crypto/md5"
	storage "filesrv/api/pb"
	"filesrv/conf"
	"filesrv/library/log"
	"flag"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
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

func TestService_DownSliceFile(t *testing.T) {
	Convey("GetFileInfo", t, func() {
		fileInfo, err := GetService().GetPbFileInfoByFid(2261944530632704)
		So(err, ShouldBeNil)
		So(fileInfo, ShouldNotBeNil)
		t.Log(fileInfo)
		Convey("DownSliceFile", func() {
			//参数偏移量必须可被1 KB整除。
			//参数限制必须可被1 KB整除。
			//限制不得超过1048576（1 MB）。
			var (
				buf      bytes.Buffer
				limit    = int64(1048576)
				limitEnd = int64(0)
				offset   = int64(1048576)
			)
			num := fileInfo.Size / 1048576
			if fileInfo.Size%1048576 != 0 {
				limitEnd = fileInfo.Size % 1048576
				num++
			}
			for i := int64(0); i < num; i++ {
				var in = &storage.InDownSliceFileItem{
					Fid:    fileInfo.Fid,
					Limit:  limit,
					Offset: offset * i,
				}
				if i+1 == num { //最后一片
					in.Limit = limitEnd
				}
				sData, err := GetService().DownSliceFile(in)
				So(err, ShouldBeNil)
				So(sData, ShouldNotBeNil)
				So(len(sData.Data), ShouldNotEqual, 0)
				buf.Write(sData.Data)
			}
			err = ioutil.WriteFile("info.png", buf.Bytes(), 0777)
			t.Log(buf.Len())
			So(err, ShouldBeNil)
			So(buf.Len(), ShouldNotEqual, 0)
			So(fileInfo.Md5, ShouldEqual, fmt.Sprintf("%x", md5.Sum(buf.Bytes())))

		})
	})
}
