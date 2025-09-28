package redis_utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Xadd(client * redis.Client,data []byte)error{
	ctx := context.Background()
	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: "websites",
		Values: map[string]any{
			"site": string(data),
		},
	}).Result()
	if err != nil{
		return err
	}
	return nil
}