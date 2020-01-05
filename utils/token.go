package utils

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

var tokenDB *redis.Client

type token struct{}

func init() {
	tokenDB = redis.NewClient(&redis.Options{
		Addr:     Config.RedisAddr,
		Password: Config.RedisPWD,
		DB:       Config.RedisTokenId,
	})
	go checkTokenDBAlive()
}

func checkTokenDBAlive() {
	for {
		_, err := tokenDB.Ping().Result()
		if err != nil {
			panic(err)
		} else {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "Token Database Alive")
		}
		time.Sleep(time.Minute * 30)
	}
}

func (*token) SetHandShake(id string) string {
	key := "HandShake_" + id
	token := GetRandomString(64)
	err := tokenDB.Set(key, token, time.Minute*2).Err()
	if err != nil {
		return ""
	}
	err = tokenDB.Set(key+"_expire", TimeStamp(2), time.Minute*2).Err()
	if err != nil {
		return ""
	}
	return token
}

func (*token) GetHandShake(id string) (token string, expire int64) {
	key := "HandShake_" + id
	var script = redis.NewScript(`
		local res = {}
		res[1] = redis.call("GET", KEYS[1])
		res[2] = redis.call("GET", KEYS[2])
		redis.call("DEL", KEYS[1])
		redis.call("DEL", KEYS[2])
		return res
    `)
	res, err := script.Run(tokenDB, []string{key, key + "_expire"}).Result()
	if err != nil || res.([]interface{})[0] == nil || res.([]interface{})[1] == nil {
		return "", 0
	}
	token = res.([]interface{})[0].(string)
	expire, _ = strconv.ParseInt(res.([]interface{})[1].(string), 10, 64)
	return
}

func (*token) Set(id string) string {
	key := id
	token := GetRandomString(64)
	err := tokenDB.Set(key, token, time.Minute*10).Err()
	if err != nil {
		return ""
	} else {
		return token
	}
}

func (*token) SetOverride(id, task string, expire int64) bool {
	err := tokenDB.Set(id+"_"+task, "OVERRIDE", time.Minute*time.Duration(expire)).Err()
	return err == nil
}

func (*token) Check(id, token string) bool {
	val, err := tokenDB.Get(id).Result()
	if err != nil {
		return false
	} else {
		return val == token
	}
}

func (*token) Delete(id string) bool {
	_, err := tokenDB.Del(id).Result()
	return err == nil
}

var Token = new(token)
