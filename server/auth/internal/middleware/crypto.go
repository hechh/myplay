package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hechh/library/crypto"
	"github.com/hechh/library/mlog"
	"github.com/hechh/library/pool"
)

// 加密中间件
func Crypto(getSecretKey func() string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取包体
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Request.Body.Close()

		// 解密请求
		deBody, err := crypto.AesDecrypto(body, []byte(getSecretKey()))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(deBody))
		c.Request.ContentLength = int64(len(deBody))

		// 响应自动加密
		crw := &cryptoWriter{
			ResponseWriter: c.Writer,
			body:           pool.GetBytes(),
		}
		c.Writer = crw

		// 继续调用其他中间件和handler
		c.Next()

		// 加密响应（只在成功状态码时加密）
		crw.response(c, getSecretKey())
		c.Writer = crw.ResponseWriter
	}
}

type cryptoWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *cryptoWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *cryptoWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func (w *cryptoWriter) response(c *gin.Context, secret string) {
	defer pool.PutBytes(w.body) // 回收对象

	if w.Status() != http.StatusOK || w.body.Len() <= 0 {
		return
	}
	// 加密应答
	enBody, err := crypto.AesEncrypto(w.body.Bytes(), []byte(secret))
	if err != nil {
		mlog.Errorf("解密错误：%v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if _, err := w.ResponseWriter.Write(enBody); err != nil {
		mlog.Errorf("解密错误：%v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	header := w.ResponseWriter.Header()
	header.Set("Content-Type", "application/octet-stream") // 二进制
	header.Set("Content-Encoding", "aes-gcm")              // 设置加密方式
	header.Del("Content-Length")                           // 清除原来的头，避免重复
}
