package redis_lock

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"time"
)

func LockRequest(ctx context.Context, redisCli *redis.Client, requestKey string, reqStruct interface{}) (func(), error) {
	v := strconv.Itoa(rand.Int())
	reqStr, err := json.Marshal(reqStruct)
	if err != nil {
		return nil, err
	}
	reqSum := fmt.Sprintf("%x", md5.Sum(reqStr))
	key := "lock_request:" + requestKey + ":" + reqSum
	//max life-time 3 sec
	count, err := redisCli.SetNX(ctx, key, v, 3*time.Second).Result()
	if err != nil {
		return nil, err
	}
	if !count {
		return nil, nil
	}
	// run return fn in defer
	return func() {
		rv, err := redisCli.Get(ctx, key).Result()
		if errors.Is(err, redis.Nil) {
			return
		}
		if err != nil {
			return
		}
		if rv != v {
			return
		}
		redisCli.Del(ctx, key)
	}, nil
}
