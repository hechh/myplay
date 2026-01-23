package router_data

import (
	"fmt"

	"github.com/hechh/framework"
	"github.com/hechh/library/myredis"
	"github.com/spf13/cast"
)

func GetRouterKey(uid string) string {
	return fmt.Sprintf("router:%s", uid)
}

func GetRouter(dbname string, uid uint64) (string, error) {
	client := myredis.Get(dbname)
	str := cast.ToString(uid)
	return client.Get(GetRouterKey(str))
}

func SaveRouter(dbname string, data map[string]framework.IRouter) error {
	if len(data) <= 0 {
		return nil
	}
	args := []any{}
	for key, item := range data {
		buf, err := item.Marshal()
		if err != nil {
			return err
		}
		args = append(args, GetRouterKey(key), string(buf))
	}
	if err := myredis.Get(dbname).MSet(args...); err != nil {
		return err
	}
	for _, item := range data {
		item.SetStatus(false)
	}
	return nil
}
