package conf

import "time"

//文件状态 | mgo 中存在的状态
//1 不存在 |   2 等待上传  3 正在上传  4 存在  5 过期
const (
	FileNotExist = iota + 1
	FileWaitingForUpload
	FileUploading
	FileExists
	FileExpired
)

//有这些文件后缀的需要添加缩略图
var ImageExName = []string{"JPG", "JPEG", "PNG"}

const (
	ThumbnailHeight = 200
	ThumbnailWidth  = 200
)

const (
	MgoContextTimeOut = 5 * time.Second
	FileMaxWaitTime   = 60 * 30 * time.Second
)
