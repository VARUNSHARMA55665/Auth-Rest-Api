package helpers

import (
	"auth-rest-api/resources"
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func NewRedisClient(url string) (*redis.Client, error) {
	opts := &redis.Options{
		Addr: url,
	}

	client := redis.NewClient(opts)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func SetRedis(key string, val string, timeInMins int) error {
	err := client.Set(context.Background(), key, val, time.Duration(int64(timeInMins)*int64(time.Minute))).Err()
	return err
}

func GetRedis(key string) *redis.StringCmd {
	return client.Get(context.Background(), key)
}

func DelRedis(uemail string) *redis.IntCmd {
	return client.Del(context.Background(), uemail)
}

func Exists(key string) *redis.IntCmd {
	return client.Exists(context.Background(), key)
}

func SetRedisInSeconds(key string, val string, timeInSecs int) error {
	err := client.Set(context.Background(), key, val, time.Duration(int64(timeInSecs)*int64(time.Second))).Err()
	return err
}

func GetStatus() error {
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	if _, err := client.Ping(ctx).Result(); err != nil {
		return err
	}
	return nil
}

func InitRedis() error {
	env := os.Getenv("GO_ENV")

	redisUrl := resources.GetConfig().GetString("config." + env + ".redisUrl")
	var err error
	client, err = NewRedisClient(redisUrl)
	if err != nil {
		return err
	}

	log.Println("Redis Connect Success")

	return nil
}
