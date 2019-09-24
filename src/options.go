package tgo

import "time"

type Options struct {
	HeartbeatInterval time.Duration // 心跳间隔
}

func NewOptions() *Options {

	return &Options{
		HeartbeatInterval: time.Second*60,
	}
}

type Option func(*Options) error