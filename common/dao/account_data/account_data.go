package account_data

import (
	"myplay/common/constant"
	"myplay/common/pb"

	"github.com/go-xorm/xorm"
	"github.com/hechh/library/database"
	"github.com/hechh/library/uerror"
)

func init() {
	database.Register(constant.MYSQL_DB_ACCOUNT, &pb.AccountData{})
}

func Query(sess *xorm.Session, uid uint64) (*pb.AccountData, error) {
	if sess == nil {
		client := database.Get(constant.MYSQL_DB_ACCOUNT)
		if client == nil {
			return nil, uerror.Err(-1, "mysql数据库未建立连接")
		}
		sess = client.NewSession()
		defer sess.Close()
	}
	item := &pb.AccountData{}
	has, err := sess.Where("uid = ?", uid).Get(item)
	if err != nil {
		return nil, uerror.Err(-1, err.Error())
	}
	if !has {
		return nil, nil
	}
	return item, nil
}

func Insert(sess *xorm.Session, acc *pb.AccountData) error {
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
		return uerror.Err(-1, "创建账号失败%v", acc)
	}
	return nil
}

// 更新指定字段值
func Update(sess *xorm.Session, data *pb.AccountData, cols ...string) error {
	if sess == nil {
		client := database.Get(constant.MYSQL_DB_ACCOUNT)
		if client == nil {
			return uerror.Err(-1, "mysql数据库未建立连接")
		}
		sess = client.NewSession()
		defer sess.Close()
	}
	_, err := sess.Table(new(pb.AccountData)).Where("uid=?", data.Uid).MustCols(cols...).Update(data)
	return err
}

func Or(sess *xorm.Session, name string, phone, email string) (*pb.AccountData, error) {
	if sess == nil {
		client := database.Get(constant.MYSQL_DB_ACCOUNT)
		if client == nil {
			return nil, uerror.Err(-1, "mysql数据库未建立连接")
		}
		sess = client.NewSession()
		defer sess.Close()
	}
	if len(name) > 0 {
		sess = sess.Or("name=?", name)
	}
	if len(phone) > 0 {
		sess = sess.Or("phone = ?", phone)
	}
	if len(email) > 0 {
		sess = sess.Or("email = ?", email)
	}
	item := &pb.AccountData{}
	has, err := sess.Get(item)
	if err != nil {
		return nil, uerror.Err(-1, err.Error())
	}
	if !has {
		return nil, nil
	}
	return item, nil
}
