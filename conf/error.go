package conf

import "errors"

//file manager error info
var (
	ErrFilePartsInvalid         = errors.New("FILE_PARTS_INVALID")           //无效的零件数。该值不在1..3000
	ErrFilePartTooBig           = errors.New("FILE_PART_TOO_BIG")            //已超出文件部分内容的大小限制（512 KB）
	ErrFilePartEmpty            = errors.New("FILE_PART_EMPTY")              //发送的文件部分为空（512 KB）
	ErrFilePartSizeInvalid512KB = errors.New("FILE_PART_SIZE_INVALID-512KB") //不能按part_size平均分配
	ErrFilePartSizeInvalid1KB   = errors.New("FILE_PART_SIZE_INVALID-1KB")   //part_size不能被1KB整除
	ErrFilePartSizeChanged      = errors.New("FILE_PART_SIZE_CHANGED")       //分片大小与同一文件中先前零件之一的大小不同
	ErrMd5ChecksumInvalid       = errors.New("MD5_CHECKSUM_INVALID")         //文件的校验和与md5_checksum参数不匹配
	ErrFileUploadCompleted      = errors.New("FILE_UPLOAD_COMPLETED")        //文件已经上传完成
	ErrFileUploading            = errors.New("FILE_UPLOADING")               //文件正在上传
	ErrFilePartUploadCompleted  = errors.New("FILE_PART_UPLOAD_COMPLETED")   //文件分片已经上传完成
	ErrFileIdInvalid            = errors.New("FILE_ID_INVALID")              //文件地址无效
	ErrOffsetInvalid            = errors.New("OFFSET_INVALID")               //偏移值无效
	ErrLimitInvalid             = errors.New("LIMIT_INVALID")                //限制值无效
)

// minio err info
var (
	ErrFileSizeInvalid = errors.New("FILE_SIZE_INVALID") //文件服务器存的文件大小,不等于文件大小

)
