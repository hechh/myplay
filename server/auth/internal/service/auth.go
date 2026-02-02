package service

import (
	"myplay/common/dao/account_data"
	"myplay/common/dao/gen/session_data"
	"myplay/common/pb"
	"myplay/common/table/global_config"
	"myplay/common/token"
	"myplay/server/auth/internal/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hechh/library/crypto"
	"github.com/hechh/library/uerror"
)

// 注册接口
func Prelogin(ctx *gin.Context, req *pb.AuthReq, rsp *pb.AuthRsp) error {
	// 加载配置
	cfg := global_config.SGet(0)
	if cfg == nil {
		return uerror.Err(pb.ErrorCode_TableNotFound, "配置不存在")
	}

	// 查重
	acc, err := account_data.Or(nil, req.Name, req.Phone, req.Email)
	if err != nil {
		return uerror.Err(pb.ErrorCode_MysqlFailed, err.Error())
	}

	// 注册玩家
	if acc == nil {
		// 随机生成密钥 + 加密密码
		secretKey, err := config.GenerateAes128Key()
		if err != nil {
			return err
		}
		passwd, err := crypto.AesEncrypto([]byte(req.Password), secretKey)
		if err != nil {
			return err
		}

		// 生成玩家唯一uid
		uid, err := config.GenerateUid()
		if err != nil {
			return err
		}
		acc = &pb.AccountData{
			Uid:        uid,
			Name:       req.Name,
			Email:      req.Email,
			Phone:      req.Phone,
			Password:   string(passwd),
			SecretKey:  string(secretKey),
			CreateTime: time.Now().Unix(),
			Platform:   req.Platform,
			LoginType:  req.LoginType,
		}
		// 写入数据库
		if err := account_data.Insert(nil, acc); err != nil {
			return err
		}
	}

	// 生成session
	sess := &pb.SessionData{
		Uid:       acc.Uid,
		Name:      req.Name,
		LoginTime: time.Now().Unix(),
		LoginIp:   req.Ip,
		Version:   req.Version,
		DeviceId:  req.DeviceId,
		Platform:  req.Platform,
	}
	if err := session_data.Set(acc.Uid, sess, 0); err != nil {
		return err
	}
	str, err := token.GenToken(sess, config.GetTokenKey())
	if err != nil {
		return err
	}

	// 返回数据
	rsp.Uid = acc.Uid
	rsp.Url = cfg.GateUrl
	rsp.Token = str
	return nil
}
