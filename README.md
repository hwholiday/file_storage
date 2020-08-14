# file_storage
参照 Telegram [Uploading and Downloading Files](https://core.telegram.org/api/files)  
 写的文件服务器（断点续传，分片上传下载,自动生成缩略图等）
```bash
//申请上传文件ID
ApplyFid(info *storage.InApplyFid) (out *storage.OutApplyFid, err error) 
//通过文件MD5获取服务器文件信息
GetPbFileInfoByMd5(md5 string) (fileInfo *storage.FileInfo, err error)
//通过文件ID获取服务器文件信息
GetPbFileInfoByFid(fid int64) (fileInfo *storage.FileInfo, err error)
//上传分片文件
UpSliceFile(in *storage.InUpSliceFileItem) (err error)
//分片下载文件
DownSliceFile(in *storage.InDownSliceFileItem) (out *storage.OutDownSliceFileItem, err error)
//上传不分片文件
UpFile(in *storage.InUpFile) (err error)
//下载不分片文件
GetFile(fid int64) (out *storage.OutDownFile, err error)
//取消上传
CancelByFid(info *storage.InCancelUpload) (err error)
```   
