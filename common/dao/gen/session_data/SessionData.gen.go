/*
* 本代码由pbtool工具生成，请勿手动修改
 */

package session_data

import (
	"fmt"
	pb "myplay/common/pb"
	"time"

	"github.com/hechh/library/myredis"
	"github.com/hechh/library/uerror"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
)

func GetKey(Uid uint64) string {
	return fmt.Sprintf("session:%d", Uid)
}

func Get(Uid uint64) (*pb.SessionData, error) {
	// 获取redis连接
	client := myredis.Get("myplay")
	if client == nil {
		return nil, uerror.New(-1, "myplay数据库不存在")
	}

	// 加载数据
	str, err := client.Get(GetKey(Uid))
	if err != nil {
		return nil, err
	}

	// 解析数据
	if len(str) > 0 {
		data := &pb.SessionData{}
		err = proto.Unmarshal([]byte(str), data)
		return data, err
	}
	return nil, nil
}

func Set(Uid uint64, val *pb.SessionData, expiration time.Duration) error {
	// 获取redis连接
	client := myredis.Get("myplay")
	if client == nil {
		return uerror.New(-1, "myplay数据库不存在")
	}

	// 编码数据
	buf, err := proto.Marshal(val)
	if err != nil {
		return err
	}

	// 存储数据
	return client.Set(GetKey(Uid), buf, expiration)
}

func Del(Uid uint64) error {
	// 获取redis连接
	client := myredis.Get("myplay")
	if client == nil {
		return uerror.New(-1, "myplay数据库不存在")
	}

	// 删除数据
	_, err := client.Del(GetKey(Uid))
	return err
}

func MSet(vals map[string]*pb.SessionData) error {
	// 获取redis连接
	client := myredis.Get("myplay")
	if client == nil {
		return uerror.New(-1, "myplay数据库不存在")
	}

	// 解析数据
	args := []any{}
	for key, val := range vals {
		args = append(args, key)
		if buf, err := proto.Marshal(val); err != nil {
			return err
		} else {
			args = append(args, buf)
		}
	}

	// 批量储存数据
	return client.MSet(args...)
}

func MGet(keys ...string) (map[string]*pb.SessionData, error) {
	// 获取redis连接
	client := myredis.Get("myplay")
	if client == nil {
		return nil, uerror.New(-1, "myplay数据库不存在")
	}

	// 批量加载数据
	values, err := client.MGet(keys...)
	if err != nil {
		return nil, err
	}

	// 解析数据
	rets := map[string]*pb.SessionData{}
	for i, key := range keys {
		value := values[i]
		if value == nil {
			continue
		}
		item := &pb.SessionData{}
		err := proto.Unmarshal([]byte(cast.ToString(value)), item)
		if err == nil {
			rets[key] = item
		}
	}
	return rets, nil
}
