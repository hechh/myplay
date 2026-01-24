package login_lock

import (
	"fmt"
	"myplay/common/constant"
	"myplay/common/pb"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hechh/library/myredis"
	"github.com/hechh/library/uerror"
)

var (
	luascript = redis.NewScript(`
    if redis.call("GET", KEYS[1]) == ARGV[1] then
        return redis.call("DEL", KEYS[1])
    else
        return 0
    end
    `)
)

func Lock(uid uint64, nodeId uint32, expire time.Duration) error {
	client := myredis.Get(constant.REDIS_DB_ACCOUNT)
	if client == nil {
		return uerror.Err(pb.ErrorCode_RedisFailed, "redis数据库不存在")
	}
	key := fmt.Sprintf("lock:%d", uid)
	ok, err := client.SetNX(key, uid, expire)
	if err != nil {
		return err
	}
	if !ok {
		return uerror.Err(pb.ErrorCode_LockFailed, "分布式锁抢占失败")
	}
	return nil
}

func Unlock(uid uint64, nodeId uint32) error {
	client := myredis.Get(constant.REDIS_DB_ACCOUNT)
	if client == nil {
		return uerror.Err(pb.ErrorCode_RedisFailed, "redis数据库不存在")
	}
	key := fmt.Sprintf("lock:%d", uid)
	_, err := client.Run(luascript, key, nodeId)
	return err
}
