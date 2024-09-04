package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var ctx = context.Background()
var rdb *redis.Client

func InitRedisDB(cfg *Config) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password, // no password set
		DB:       cfg.Redis.DB,       // use default DB
	})
}

func CacheTerminal(terminal *Terminal) error {
	score := terminal.CPUUsage + terminal.MemUsage // 示例：根据CPU和内存使用率评分
	return rdb.ZAdd(ctx, "hot1000", &redis.Z{
		Score:  score,
		Member: terminal.ID,
	}).Err()
}

func UpdateTerminal(terminal *Terminal) error {
	key := "terminal:" + strconv.Itoa(terminal.ID)
	return rdb.HMSet(ctx, key, map[string]interface{}{
		"MAC":            terminal.MAC,
		"CPUUsage":       terminal.CPUUsage,
		"MemUsage":       terminal.MemUsage,
		"LastModifyTime": time.Now(),
	}).Err()
}

func GetHotTerminals() ([]Terminal, error) {
	ids, err := rdb.ZRevRange(ctx, "hot1000", 0, 999).Result()
	if err != nil {
		return nil, err
	}

	terminals := make([]Terminal, len(ids))
	for i, id := range ids {
		key := "terminal:" + id
		data, err := rdb.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		// 转换ID
		terminalID, _ := strconv.Atoi(id)

		// 转换CPUUsage和MemUsage
		cpuUsage, _ := strconv.ParseFloat(data["CPUUsage"], 64)
		memUsage, _ := strconv.ParseFloat(data["MemUsage"], 64)

		terminals[i] = Terminal{
			ID:             terminalID,
			MAC:            data["MAC"],
			CPUUsage:       cpuUsage,
			MemUsage:       memUsage,
			LastModifyTime: time.Now().Truncate(time.Minute),
		}
	}

	return terminals, nil
}
