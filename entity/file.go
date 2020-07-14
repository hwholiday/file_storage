package entity

type FileApplyFid struct {
	Fid    int64 `json:"fid"`
	Status int   `json:"status"` //1 不存在   2 等待上传  3 正在上传  4 存在  5 过期
}
