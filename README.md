# 简介

暴露浏览器中的接口到 xhe wireguard vpn 网络中

# 如何使用

```js
let xhe = await XheConnect(
  {
    PrivateKey: "CFp1j2epz2sUJ8ovPeUgiStto8EOAAnmvGWp+SSECmI=",
    Address: "192.168.4.1/24",
    Peers: [
      {
        PublicKey: "TUpLBfCKwL0joxi+nOsE3+wgxjhlIZdtlcftup/lRik=",
        AllowedIPs: ["192.168.4.2/32"],
      },
    ],
  },
  {
    signaler: `https://test:test@signaler.slive.fun?t=${device_id}`,
  }
);
let server = await xhe.ListenTCP(80);
server.Serve().catch(() => {
  // donothing
});
if (!server.ServeReady()) {
  throw new Error("server is not ready");
}
// 反向代理浏览器接口到 192.168.4.1:80 里
await server.ReverseProxy("/", `http://127.0.0.1:${port}/`);
```

之后你可以在 golang 里调用浏览器里的接口, 或者通过 [xhe vpn](https://github.com/remoon-net/xhe) tun 组网后
直接使用 `curl http://192.168.4.1:80` 访问

```golang
...
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
}
```
