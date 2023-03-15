package web_copy

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileUploader struct {
	// FileField 对应文件在表单中的字段名
	FileField string
	// DstPathFunc 目标路径
	DstPathFunc func(fh *multipart.FileHeader) string
}

// Handle
// 1.	文件，目标文件，err := ctx.Req.FormFile("对应文件在表单中字段名")
// 2.	目标文件, err := os.OpenFile(f.DstPathFunc(目标文件), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
// 3.	_, err = io.CopyBuffer(dst, src, nil)
func (f *FileUploader) Handle() HandlerFunc {
	return func(ctx *Context) {
		src, srcHandle, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = 400
			ctx.RespData = []byte("上传失败，未找到数据")
			log.Println(err)
			return
		}

		dst, err := os.OpenFile(f.DstPathFunc(srcHandle), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败")
			log.Println(err)
			return
		}

		_, err = io.CopyBuffer(dst, src, nil)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败")
			log.Println(err)
			return
		}
		ctx.RespData = []byte("上传成功！")
	}
}

type FileDownloader struct {
	Dir string
}

func (f *FileDownloader) Handle() HandlerFunc {
	return func(ctx *Context) {
		req, _ := ctx.QueryValue("file").String()
		path := filepath.Join(f.Dir, filepath.Clean(req))
		fn := filepath.Base(path)
		// 拼Header
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")

		http.ServeFile(ctx.Resp, ctx.Req, path)
	}
}
