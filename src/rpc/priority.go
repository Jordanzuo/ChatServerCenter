package rpc

// 优先级
type Priority int8

const (
	// 高优先级
	Con_HighPriority Priority = 1 + iota

	// 中优先级
	Con_MediumPriority

	// 低优先级
	Con_LowPriority
)
