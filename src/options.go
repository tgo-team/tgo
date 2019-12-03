package tgo

import "time"

// Options Options
type Options struct {
	HeartbeatInterval time.Duration // 心跳间隔
}

// NewOptions NewOptions
func NewOptions() *Options {

	return &Options{
		HeartbeatInterval: time.Second * 60,
	}
}

// Option 配置项
type Option func(*Options) error
