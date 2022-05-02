package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"path"
	"project/cms/internal/proto"
	"project/pkg/logger"
	"project/pkg/util/files"
)

func (h *Handler) UploadFile(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		c.JSON(RespWithMsg(InvalidParam, ""))
		return
	}
	if f.Size > 2<<20 {
		c.JSON(RespWithMsg(OverSize, "单个文件最大不能超过2M"))
		return
	}
	file, _ := f.Open()
	defer file.Close()
	b, _ := io.ReadAll(file)
	remotePath := "file/" + files.GenFilePath(b) + path.Ext(f.Filename)
	err = h.cos.PutObject(c, remotePath, bytes.NewReader(b))
	if err != nil {
		logger.FromContext(c).Error("cos.PutObject error", nil, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, &proto.UploadResp{
		Host: h.cdn,
		Path: remotePath,
	})
}

func (h *Handler) UploadImage(c *gin.Context) {
	f, err := c.FormFile("image")
	if err != nil {
		c.JSON(RespWithMsg(InvalidParam, ""))
		return
	}
	if f.Size > 500<<10 {
		c.JSON(RespWithMsg(OverSize, "单张图片限制500KB以内"))
		return
	}
	file, _ := f.Open()
	defer file.Close()
	b, _ := io.ReadAll(file)
	ext, ok := files.CheckImage(b)
	if !ok {
		c.JSON(RespWithMsg(UnsupportedType, "无效的图片类型，仅支持jpg/png/gif格式"))
		return
	}
	remotePath := "img/" + files.GenFilePath(b) + "." + ext
	err = h.cos.PutObject(c, remotePath, bytes.NewReader(b))
	if err != nil {
		logger.FromContext(c).Error("cos.PutObject error", nil, err)
		c.JSON(RespWithErr(err))
		return
	}
	c.JSON(OK, &proto.UploadResp{
		Host: h.cdn,
		Path: remotePath,
	})
}
