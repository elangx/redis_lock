package redis_lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestRequestLock(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "1234",
		DB:       0,
	})
	//first lock
	fn1, err1 := LockRequest(context.Background(), client, "test", map[string]interface{}{
		"key1": "abc",
		"key2": "def",
	})
	if err1 != nil {
		t.Error("err:", err1)
		return
	}
	//nil means lock fail
	if fn1 == nil {
		t.Error("lock failed")
	}
	// second lock when the first do not release,it should fail
	fn2, err2 := LockRequest(context.Background(), client, "test", map[string]interface{}{
		"key1": "abc",
		"key2": "def",
	})
	if err2 != nil {
		t.Error("err:", err2)
		return
	}
	//fail
	if fn2 != nil {
		t.Error("double lock")
		return
	}
	//unlock the first lock action
	fn1()

	//auto unlock after 3 seconds
	fn3, err3 := LockRequest(context.Background(), client, "test", map[string]interface{}{
		"key1": "abc",
		"key2": "def",
	})
	if err3 != nil {
		t.Error("err:", err3)
		return
	}
	if fn3 == nil {
		t.Error("re-lock failed")
		return
	}
	fn4, _ := LockRequest(context.Background(), client, "test", map[string]interface{}{
		"key1": "abc",
		"key2": "def",
	})
	//fail in 3 sec
	if fn4 != nil {
		t.Error("auto unlock double lock failed")
	}
	time.Sleep(4 * time.Second)

	//re-lock should success
	fn5, err5 := LockRequest(context.Background(), client, "test", map[string]interface{}{
		"key1": "abc",
		"key2": "def",
	})
	if err5 != nil {
		t.Error("err:", err5)
		return
	}
	if fn5 == nil {
		t.Error("re-lock failed")
		return
	}
	//that's the right way to unlock
	defer fn5()
}
