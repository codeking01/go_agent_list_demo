package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func generateMAC() string {
	mac := fmt.Sprintf("00:%02x:%02x:%02x:%02x:%02x",
		rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
	return mac
}

func simulateTerminal(id int) *Terminal {
	return &Terminal{
		ID:       id,
		MAC:      generateMAC(),
		CPUUsage: rand.Float64() * 100,
		MemUsage: rand.Float64() * 100,
	}
}

func main() {
	// 加载配置
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 初始化MySQL和Redis
	InitPostGreSQLDB(config)
	InitRedisDB(config)

	// 模拟5万+数据
	for i := 1; i <= 50000; i++ {
		terminal := simulateTerminal(i)

		// 更新缓存
		err := UpdateTerminal(terminal)
		if err != nil {
			return
		}

		// 根据评分缓存前1000个高频终端
		err = CacheTerminal(terminal)
		if err != nil {
			return
		}

		// 模拟关键字段变化记录 (e.g., MAC 地址变化)
		if rand.Intn(1000) < 5 { // 例如每1000个终端中5个MAC发生变化
			oldMAC := terminal.MAC
			terminal.MAC = generateMAC()

			// 更新缓存和数据库
			err := UpdateTerminal(terminal)
			if err != nil {
				return
			}
			logChange := &ChangeLog{
				TerminalID: terminal.ID,
				OldMAC:     oldMAC,
				NewMAC:     terminal.MAC,
				Timestamp:  time.Now().Unix(),
			}
			err = AgentLogChange(logChange)
			if err != nil {
				return
			}
		}
	}

	// 定期持久化 (例如每隔5分钟持久化一次)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		terminals, _ := GetHotTerminals()

		for _, terminal := range terminals {
			err := AgentSaveTerminal(&terminal)
			if err != nil {
				return
			}
		}
	}
}
