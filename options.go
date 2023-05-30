package main

import "github.com/pion/webrtc/v3"

type Options struct {
	Signaler   string             `json:"signaler"`
	ICEServers []webrtc.ICEServer `json:"ices"`
	LoggerName string             `json:"logger"`
	Debug      bool               `json:"debug"`
}
