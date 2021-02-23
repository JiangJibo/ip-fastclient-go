package domain

import "time"

type FastMetric struct {
	Memory      int32
	DecryptCost int64
	LoadCost    int64
}

// 计算启动耗时
func (metric *FastMetric) TimeCost(startTime int64) int64 {
	endTime := time.Millisecond.Milliseconds()
	timeCost := endTime - startTime
	return timeCost
}