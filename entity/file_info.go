package entity

type FileInfo struct {
	Fid         int64   `bson:"fid"`          //文件ID
	Name        string  `bson:"name"`         //文件名
	BucketName  string  `bson:"bucket_name"`  //桶名
	Size        int64   `bson:"size"`         //文件总大小
	ContentType string  `bson:"content_type"` //文件信息
	Md5         string  `bson:"md5"`          //文件MD5
	ExName      string  `bson:"ex_name"`      //文件扩展名
	ExImage     ImageEx `bson:"ex_image"`     //图片文件扩展信息
	SliceTotal  int     `bson:"slice_total"`  // 1 为不分片文件  (1~3000)
	ExpiredTime int64   `bson:"expired_time"` //过期时间 设置为0 文件永久不过期
	Status      int32   `bson:"status"`       //conf 文件状态 查看解析
	CreateTime  int64   `bson:"create_time"`
	UpdateTime  int64   `bson:"update_time"` //如果最后一次操作时间+FileMaxWaitTime大于当前时间,并且文件状态不是已经完成,自动删除
}

func (f FileInfo) TableName() string {
	return "file_info"
}

type ImageEx struct {
	High           int   `bson:"high"`
	Width          int   `bson:"width"`
	ThumbnailFid   int64 `bson:"thumbnail_fid"`
	ThumbnailHigh  int   `bson:"thumbnail_high"`
	ThumbnailWidth int   `bson:"thumbnail_width"`
}
