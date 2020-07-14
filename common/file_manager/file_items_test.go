package file_manager

import (
	"crypto/md5"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

var waitSync sync.WaitGroup

func TestNewFileItem(t *testing.T) {

	fileData, err := os.Open("./825220.png")
	CheckErr(err)
	defer fileData.Close()
	fi, err := fileData.Stat()
	CheckErr(err)
	fmt.Println(fi.Size())
	buf := make([]byte, fi.Size())
	_, err = fileData.Read(buf)
	CheckErr(err)

	var sliceSize uint32
	num := fi.Size() / (512 * 1024)
	sliceSize = 512 * 1024
	if fi.Size()%(512*1024) != 0 {
		num++
	}
	fileItem := NewFileItem(&FileItem{
		Fid:        123,
		Size:       fi.Size(),
		SliceTotal: int(num),
		Md5:        fmt.Sprintf("%x", md5.Sum(buf)),
	})
	fmt.Println("sliceSize", sliceSize)
	fmt.Println("计划分片个数", num)
	var startTime = time.Now().UnixNano() / 1e6
	for i := 1; i <= int(num); i++ {
		waitSync.Add(1)
		go func(index int) {
			data := (index - 1) * int(sliceSize)
			fmt.Println(index, "success >>>> ", data, "---------", data+int(sliceSize))
			if int(num) == index {
				fmt.Println(fileItem.AddItem(&FileUploadItem{
					Fid:  123,
					Part: index,
					Data: buf[data:],
					Md5:  fmt.Sprintf("%x", md5.Sum(buf[data:])),
				}))
			} else {
				fmt.Println(fileItem.AddItem(&FileUploadItem{
					Fid:  123,
					Part: index,
					Data: buf[data : data+int(sliceSize)],
					Md5:  fmt.Sprintf("%x", md5.Sum(buf[data:data+int(sliceSize)])),
				}))
			}
			waitSync.Done()
		}(i)
	}
	waitSync.Wait()
	fmt.Println("消耗时间  ", time.Now().UnixNano()/1e6-startTime, "》》毫秒")
	select {}
}
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
