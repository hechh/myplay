package middleware

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hechh/framework"
	"github.com/hechh/library/uerror"
)

type Handler[R any, T any] func(*gin.Context, *R, *T) error

func Wrapper[R any, T any](f Handler[R, T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 应答类型
		rsp := new(T)
		irsp, _ := any(rsp).(framework.IResponse)

		// 绑定JSON请求体
		req := new(R)
		if err := bindRequest(c, req); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 调用请求
		if err := f(c, req, rsp); err != nil {
			switch vv := err.(type) {
			case *uerror.UError:
				irsp.SetRspHead(vv.GetCode(), vv.GetMsg())
			default:
				irsp.SetRspHead(-1, vv.Error())
			}
		}
		c.JSON(http.StatusOK, rsp)
	}
}

func bindRequest[R any](c *gin.Context, req *R) error {
	if c.Request.Method == "GET" {
		if err := c.ShouldBindQuery(req); err != nil {
			return err
		}
	} else {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return err
		}
		if len(body) == 0 {
			return nil // 空body允许
		}
		if err := json.Unmarshal(body, req); err != nil {
			return err
		}
	}
	return nil
}
