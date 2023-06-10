package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"time"

	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	"github.com/shynome/wgortc"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"remoon.net/xhe/pkg/config"
	"remoon.net/xhe/signaler"
)

func main() {
	id := flag.String("id", "", "")
	flag.Parse()
	if *id == "" {
		panic("id is required")
	}
	endpoint := fmt.Sprintf("https://test@signaler.slive.fun?t=%s", *id)
	config := config.Config{
		Device: config.Device{
			PrivateKey: "4G6w8LSJazn887NIkRKtcgnTSXGc1RenvbzV3YUUb3M=",
			Address:    "192.168.4.2/24",
		},
		Peers: []config.Peer{
			{
				PublicKey:           "0RguZc+bKUliW8KzTudHhYxxqj+Fnb3vuKoXeNY1IHE=",
				AllowedIPs:          []string{"192.168.4.1/32"},
				Endpoint:            endpoint,
				PersistentKeepalive: "5",
			},
		},
	}
	tdev, tnet, err := netstack.CreateNetTUN(
		[]netip.Addr{netip.MustParseAddr("192.168.4.2")},
		[]netip.Addr{netip.MustParseAddr("1.1.1.1")},
		device.DefaultMTU)
	try.To(err)

	logger := device.NewLogger(device.LogLevelError, "call ")
	signaler := signaler.New("")
	bind := wgortc.NewBind(signaler)

	dev := device.NewDevice(tdev, bind, logger)
	defer dev.Close()

	try.To(dev.IpcSet(config.String()))
	try.To(dev.Up())

	client := http.Client{
		Transport: &http.Transport{DialContext: tnet.DialContext},
		Timeout:   10 * time.Second,
	}

	resp := try.To1(client.Get("http://192.168.4.1"))
	body := try.To1(io.ReadAll(resp.Body))
	fmt.Println(string(body), "from golang")
	// eval ok
	resp2 := try.To1(client.Post("http://192.168.4.1/xhe-eval", "text/plain", bytes.NewBufferString(`resolve(9999)`)))
	assert.Equal(resp2.StatusCode, 200)
	body2 := try.To1(io.ReadAll(resp2.Body))
	assert.Equal(string(body2), "9999")
	// eval wrong
	resp3 := try.To1(client.Post("http://192.168.4.1/xhe-eval", "text/plain", bytes.NewBufferString(`resolve(`)))
	assert.Equal(resp3.StatusCode, 500)
	body3 := try.To1(io.ReadAll(resp3.Body))
	fmt.Println("wrong eval result:", string(body3))
}
