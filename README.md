# file_storage
参照 Telegram [Uploading and Downloading Files](https://core.telegram.org/api/files)  
 写的文件服务器（断点续传，分片上传下载,自动生成缩略图等）

申请上传文件
---
请求
   - 文件名
   - 文件大小 （字节）
   - 文件扩展名 
   - 是否需要缩略图
   - 文件MD5

返回
   - 文件id
   - 上传过期时间，到时间未上传完，自动清理


取消上传文件
---
请求
   - 文件ID
    
返回
   - 成功 or 失败
   

获取文件信息
---
请求
   - 文件ID
    
返回
   - 文件信息
   
上传文件
---
请求

   - 文件ID
   - 文件分片ID file_part
   - 文件MD5
   - 文件内容 part_size
   - 是否是最后一片
###### 所有分片必须具有相同的大小（part_size），并且必须满足以下条件：part_size % 1024 = 0 （可被1KB整除524288 % part_size = 0（512KB必须可以被part_size整除）如果最后一部分的大小小于part_size，则不必满足这些条件。每个部分都应具有序列号file_part，其值的范围为1到3000。

服务器处理 
        
- 每个分片接收完毕的时候都检查下服务器接受文件MD5是否相等
- 检查该分片是否已经上传，是则不做任何处理,不是则把文件存入内存，再将已上传文件大小累加
- 判断已上传文件大小是否等于文件总大小 （标记文件已经完成）
- 按照文件分片顺序拼接文件
- 检查文件MD5是否相等，上传到文件服务器,文件上传结束

返回
   - 成功 or 失败


下载文件
---
请求
   - 文件ID
   - limit  可被1024整除
   - offset 可被1024整除
   - offset 不能超过 1048576（1 MB）

服务器处理
- minio GetObject 可配置 
```
var opt minio.GetObjectOptions
err := opt.SetRange(start, end)
GetObject(bucketName, fileName, opt)
```


返回
   - 文件信息     
 


## 对外方法 
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
