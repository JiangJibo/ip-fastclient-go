package domain

import (
	"fmt"
	"log"
	"time"
)

type FastMetric struct {
	Memory      int32
	DecryptCost int64
	LoadCost    int64
}

// 计算启动耗时
func (metric *FastMetric) TimeCost(startTime int64) int64 {
	endTime := time.Now().UnixNano() / 1e6
	timeCost := endTime - startTime
	return timeCost
}

func (metric *FastMetric) doubleInfo(msg string) {
	log.Print(msg)
	fmt.Println(msg)
}

func (metric *FastMetric) Echo() {
	decryptCostInfo := fmt.Sprintf("[fast-ip2geo-client] | metric | load1 cost: %d %s", metric.DecryptCost, "milliseconds")
	loadCostInfo := fmt.Sprintf("[fast-ip2geo-client] | metric | load2 cost: %d %s", metric.LoadCost, "milliseconds")
	metric.doubleInfo(decryptCostInfo)
	metric.doubleInfo(loadCostInfo)
}
