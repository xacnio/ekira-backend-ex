package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
	"time"
)

var ctx = context.Background()

type RedisCon struct {
	rdb *redis.Client
}

type RedisDatabase string

const (
	RedisDatabasePhoneVerification RedisDatabase = "phone_verification"
	RedisDatabaseOthers            RedisDatabase = "others"
)

var DatabaseLists = map[RedisDatabase]int{
	"phone-verification": 14,
	"others":             15,
}

func NewRConnectionDB(database RedisDatabase) RedisCon {
	db, ok := DatabaseLists[database]
	if !ok {
		db = DatabaseLists["others"]
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASS"), // no password set
		DB:       db,
	})
	return RedisCon{rdb: rdb}
}

func NewRConnection() RedisCon {
	db := 0
	if os.Getenv("REDIS_DB") != "" {
		db, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASS"), // no password set
		DB:       db,
	})
	return RedisCon{rdb: rdb}
}

func (rdb *RedisCon) RClose() {
	rdb.rdb.Close()
}

func (rdb *RedisCon) RPing() error {
	_, errP := rdb.rdb.Ping(ctx).Result()
	if errP != nil {
		return errP
	}
	return nil
}

func (rdb *RedisCon) RExists(key string) (bool, error) {
	res, err := rdb.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if res == 0 {
		return false, nil
	}
	return true, nil
}

func (rdb *RedisCon) RGet(key string) (string, error) {
	val2, err := rdb.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	} else {
		return val2, nil
	}
}

func (rdb *RedisCon) RGetTTL(key string) (time.Duration, error) {
	res, err := rdb.rdb.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (rdb *RedisCon) RSet(key, value string) error {
	errs := rdb.rdb.Set(ctx, key, value, 0).Err()
	return errs
}

func (rdb *RedisCon) RDel(key string) error {
	errs := rdb.rdb.Del(ctx, key).Err()
	return errs
}

func (rdb *RedisCon) RSetTTL(key, value string, ttl_second time.Duration) error {
	errs := rdb.rdb.Set(ctx, key, value, ttl_second*time.Second).Err()
	return errs
}

type RedisKey struct {
	Key  string
	Data string
}

func (rdb *RedisCon) RGetAllKeys() ([]RedisKey, error) {
	res, err := rdb.rdb.Do(ctx, "KEYS", "*").StringSlice()
	if err != nil {
		return []RedisKey{}, err
	}
	if len(res) == 0 {
		return []RedisKey{}, errors.New("n/a")
	}
	results := make([]RedisKey, len(res))
	for i := 0; i < len(res); i++ {
		val2, err := rdb.rdb.Get(ctx, res[i]).Result()
		if err != nil {
			continue
		} else {
			results = append(results, RedisKey{Key: res[i], Data: val2})
		}
	}
	return results, nil
}
