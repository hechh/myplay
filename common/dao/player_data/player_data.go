package player_data

import (
	"myplay/common/constant"
	"myplay/common/pb"

	"github.com/go-xorm/xorm"
	"github.com/hechh/library/database"
	"github.com/hechh/library/uerror"
)

func init() {
	database.Register(constant.MYSQL_DB_ACCOUNT, &pb.PlayerData{})
}

func Query(sess *xorm.Session, uid uint64) (*pb.PlayerData, error) {
	if sess == nil {
		client := database.Get(constant.MYSQL_DB_ACCOUNT)
		if client == nil {
			return nil, uerror.Err(-1, "mysql数据库未建立连接")
		}
		sess = client.NewSession()
		defer sess.Close()
	}
	item := &pb.PlayerData{}
	has, err := sess.Where("uid = ?", uid).Get(item)
	if err != nil {
		return nil, uerror.Err(-1, err.Error())
	}
	if !has {
		return nil, nil
	}
	return item, nil
}

func Insert(sess *xorm.Session, acc *pb.PlayerData) error {
	if sess == nil {
		client := database.Get(constant.MYSQL_DB_ACCOUNT)
		if client == nil {
			return uerror.Err(-1, "mysql数据库未建立连接")
		}
		sess = client.NewSession()
		defer sess.Close()
	}
	// 创建账号
	affected, err := sess.Insert(acc)
	if err != nil {
		return err
	}
	if affected == 0 {
		return uerror.Err(-1, "数据插入失败%v", acc)
	}
	return nil
}

// 更新指定字段值
func Update(sess *xorm.Session, data *pb.PlayerData, cols ...string) error {
	if sess == nil {
		client := database.Get(constant.MYSQL_DB_ACCOUNT)
		if client == nil {
			return uerror.Err(-1, "mysql数据库未建立连接")
		}
		sess = client.NewSession()
		defer sess.Close()
	}
	_, err := sess.Table(&pb.PlayerData{}).Where("uid=?", data.Uid).MustCols(cols...).Update(data)
	return err
}
