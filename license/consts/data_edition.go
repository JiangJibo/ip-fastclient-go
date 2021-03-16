package consts

type code int

const (

	// 基本版本, 默认值
	common code = iota + 1

	// 流量调度版本
	route

	// 位置增强版本
	trace
)
