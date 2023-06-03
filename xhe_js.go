package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/netip"
	"net/url"
	"os"
	"syscall/js"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	promise "github.com/nlepage/go-js-promise"
	"github.com/remoon-net/xhe/pkg/config"
	"github.com/remoon-net/xhe/signaler"
	"github.com/shynome/wgortc"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
)

var dev *device.Device

func main() {
	js.Global().Set("XheConnect", js.FuncOf(connect))
	<-make(chan any)
}

func connect(this js.Value, args []js.Value) (p any) {
	p, resolve, reject := promise.New()
	go func() {
		defer err2.Catch(func(err error) {
			reject(err.Error())
		})

		if len(args) == 0 {
			reject("config is required")
			return
		}

		var config = try.To1(getConfig[config.Config](args[0]))

		if config.Address == "" {
			reject("address is required")
			return
		}
		addr := try.To1(parseIP(config.Address))
		tdev, tnet, err := netstack.CreateNetTUN(
			[]netip.Addr{addr},
			[]netip.Addr{netip.MustParseAddr("1.1.1.1")},
			device.DefaultMTU)
		try.To(err)

		var options Options
		if len(args) == 2 {
			options = try.To1(getConfig[Options](args[1]))
		}
		if options.Debug {
			err2.SetErrorTracer(os.Stderr)
		}
		if options.LoggerName == "" {
			options.LoggerName = "xhe-connect"
		}
		options.LoggerName += " "

		s := signaler.New(options.Signaler)
		bind := wgortc.NewBind(s)
		bind.ICEServers = options.ICEServers

		logger := device.NewLogger(device.LogLevelError, options.LoggerName)
		dev := device.NewDevice(tdev, bind, logger)

		wgConf := config.String()
		try.To(dev.IpcSet(wgConf))
		try.To(dev.Up())

		xwg := NewXheWireguard(dev, tnet)
		resolve(xwg.ToJS())
	}()
	return
}

type XheWireguard struct {
	dev *device.Device
	net *netstack.Net
}

func NewXheWireguard(dev *device.Device, net *netstack.Net) *XheWireguard {
	return &XheWireguard{
		dev: dev,
		net: net,
	}
}

func (xwg *XheWireguard) ToJS() (root js.Value) {
	root = js.Global().Get("Object").New()
	root.Set("ListenTCP", js.FuncOf(xwg.ListenTCP))
	root.Set("IpcGet", js.FuncOf(xwg.IpcGet))
	return root
}

func (n *XheWireguard) IpcGet(this js.Value, args []js.Value) (p any) {
	p, resolve, reject := promise.New()
	go func() {
		defer err2.Catch(func(err error) {
			reject(err.Error())
		})
		config := try.To1(n.dev.IpcGet())
		resolve(config)
	}()
	return
}

func (n *XheWireguard) ListenTCP(this js.Value, args []js.Value) (p any) {
	p, resolve, reject := promise.New()
	var port int = 80
	if len(args) >= 1 {
		port = args[0].Int()
	}
	go func() {
		defer err2.Catch(func(err error) {
			reject(err.Error())
		})
		l := try.To1(n.net.ListenTCP(&net.TCPAddr{Port: port}))
		s := NewTCPServer(l)
		s.net = n.net
		resolve(s.ToJS())
	}()
	return
}

type TCPServer struct {
	listener *gonet.TCPListener
	net      *netstack.Net
	mux      *http.ServeMux
}

func NewTCPServer(l *gonet.TCPListener) *TCPServer {
	return &TCPServer{
		listener: l,
	}
}

func (l *TCPServer) ToJS() (root js.Value) {
	root = js.Global().Get("Object").New()
	root.Set("Serve", js.FuncOf(l.Serve))
	root.Set("Close", js.FuncOf(l.Close))
	root.Set("ServeReady", js.FuncOf(l.ServeReady))
	root.Set("ReverseProxy", js.FuncOf(l.ReverseProxy))
	root.Set("HandleEval", js.FuncOf(l.HandleEval))
	return
}

func (l *TCPServer) Serve(this js.Value, args []js.Value) (p any) {
	p, resolve, reject := promise.New()
	go func() {
		defer err2.Catch(func(err error) {
			reject(err.Error())
		})
		l.mux = http.NewServeMux()
		try.To(http.Serve(l.listener, l.mux))
		resolve("exited")
	}()
	return
}

func (l *TCPServer) ServeReady(this js.Value, args []js.Value) any {
	return l.mux != nil
}

func (l *TCPServer) Close(this js.Value, args []js.Value) (p any) {
	p, resolve, reject := promise.New()
	go func() {
		if err := l.listener.Close(); err != nil {
			reject(err.Error())
			return
		}
		resolve("closed")
	}()
	return
}

func (l *TCPServer) ReverseProxy(this js.Value, args []js.Value) (p any) {
	p, resolve, reject := promise.New()
	go func() {
		defer err2.Catch(func(err error) {
			reject(err.Error())
		})
		if len(args) < 2 {
			reject("path and host is required")
			return
		}
		path := args[0].String()
		remote := try.To1(url.Parse(args[1].String()))
		var proxy http.Handler = httputil.NewSingleHostReverseProxy(remote)
		if path != "/" {
			proxy = http.StripPrefix(path, proxy)
		}
		l.mux.Handle(path, proxy)
		resolve(path)
	}()
	return
}

func (l *TCPServer) HandleEval(this js.Value, args []js.Value) (p any) {
	path := args[0].String()
	l.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		defer err2.Catch(func(err error) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
		})
		content := try.To1(io.ReadAll(r.Body))
		j := try.To1(Eval(string(content)))
		fmt.Fprint(w, j)
	})
	return
}

func Eval(content string) (s string, err error) {
	defer err2.Handle(&err)
	f := js.Global().Get("Function").New("resolve", "reject", fmt.Sprintf(`"use strict";%s;resolve();`, content))
	p := js.Global().Get("Promise").New(f)
	v := try.To1(promise.Await(p))
	s = js.Global().Get("JSON").Call("stringify", v).String()
	return
}
