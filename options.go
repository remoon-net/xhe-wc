package main

import "github.com/pion/webrtc/v3"

type Options struct {
	// Deprecated: 转而使用 config.Link
	Signaler   string             `json:"signaler"`
	ICEServers []webrtc.ICEServer `json:"ices"`
	LoggerName string             `json:"logger"`
	Debug      bool               `json:"debug"`
}
