package gredis

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/youngking/gin-blog/pkg/setting"
	"time"
)

var RedisCoon *redis.Pool

func SetUp() error {
	RedisCoon = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeoout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				c.Close()
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return err
			}
			return nil
		},
	}
	return nil
}

func Set(key string, data interface{}, time int) error {
	coon := RedisCoon.Get()
	defer coon.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = coon.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = coon.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}
	return nil

}
func Exist(key string) bool {
	coon := RedisCoon.Get()
	defer coon.Close()

	exist, err := redis.Bool(coon.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exist
}

func Get(key string) ([]byte, error) {
	coon := RedisCoon.Get()
	defer coon.Close()
	value, err := redis.Bytes(coon.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return value, nil
}

func Delete(key string) (bool, error) {
	coon := RedisCoon.Get()
	value, err := redis.Bool(coon.Do("DEL", key))
	if err != nil {
		return false, err
	}
	return value, nil
}

func LikeDelete(key string) error {
	coon := RedisCoon.Get()
	defer coon.Close()

	keys, err := redis.Strings(coon.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}
	for _, key := range keys {
		_, err := Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
