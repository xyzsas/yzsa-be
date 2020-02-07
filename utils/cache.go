package utils

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var cacheDB *redis.Client

type cache struct{}

func init() {
	cacheDB = redis.NewClient(&redis.Options{
		Addr:     Config.RedisAddr,
		Password: Config.RedisPWD,
		DB:       Config.RedisCacheId,
	})
	go checkCacheDBAlive()
}

func checkCacheDBAlive() {
	for {
		_, err := cacheDB.Ping().Result()
		if err != nil {
			panic(err)
		} else {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Cache Database Alive")
		}
		time.Sleep(time.Minute * 30)
	}
}

// 运行一个Lua脚本
func (*cache) Run(script string, input []string) interface{} {
	var s = redis.NewScript(script)
	res, err := s.Run(cacheDB, input).Result()
	if err == nil {
		return res
	} else {
		return nil
	}
}

// 判断Key存在
func (c *cache) Exist(key string) bool {
	res, _ := cacheDB.Exists(key).Result()
	return res > 0
}

// 删除Key
func (*cache) Delete(key string) bool {
	_, err := cacheDB.Del(key).Result()
	return err == nil
}

// 向Hash中插入一组元素
func (c *cache) HSetMany(key string, value map[string]interface{}) bool {
	for k, v := range value {
		err := cacheDB.HSet(key, k, v).Err()
		if err != nil {
			fmt.Println(err)
			c.Delete(key)
			return false
		}
	}
	return true
}

// 读出Hash中所有元素
func (c *cache) HGetAll(key string) map[string]string {
	res, _ := cacheDB.HGetAll(key).Result()
	return res
}

// 取Set中Key的数目
func (c *cache) SCard(key string) int64 {
	res, _ := cacheDB.SCard(key).Result()
	return res
}

// 向Set中添加一个元素
func (c *cache) SAdd(key, value string) bool {
	err := cacheDB.SAdd(key, value).Err()
	return err == nil
}

// 检查某元素是否在Set中
func (c *cache) SExist(key, value string) bool {
	res, err := cacheDB.SIsMember(key, value).Result()
	if err != nil {
		return false
	} else {
		return res
	}
}

var Cache = new(cache)
