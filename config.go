package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"proxy/services"
	"proxy/utils"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app     *kingpin.Application
	service *services.ServiceItem
	cmd     *exec.Cmd
)

func initConfig() (err error) {
	//keygen
	if len(os.Args) > 1 {
		if os.Args[1] == "keygen" {
			utils.Keygen()
			os.Exit(0)
		}
	}

	//define  args
	// tcpArgs := services.TCPArgs{}
	httpArgs := services.HTTPArgs{}
	// tunnelServerArgs := services.TunnelServerArgs{}
	// tunnelClientArgs := services.TunnelClientArgs{}
	// tunnelBridgeArgs := services.TunnelBridgeArgs{}
	// udpArgs := services.UDPArgs{}
	socksArgs := services.SocksArgs{}
	//build srvice args
	app = kingpin.New("proxy", "happy with proxy")
	app.Author("snail").Version(APP_VERSION)
	debug := app.Flag("debug", "debug log output").Default("false").Bool()
	daemon := app.Flag("daemon", "run proxy in background").Default("false").Bool()
	forever := app.Flag("forever", "run proxy in forever,fail and retry").Default("false").Bool()
	logfile := app.Flag("log", "log file path").Default("").String()

	//########http#########
	http := app.Command("http", "proxy on http mode")
	httpArgs.Parent = http.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	httpArgs.CertFile = http.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	httpArgs.KeyFile = http.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	httpArgs.LocalType = http.Flag("local-type", "local protocol type <tls|tcp|kcp>").Default("tcp").Short('t').Enum("tls", "tcp", "kcp")
	httpArgs.ParentType = http.Flag("parent-type", "parent protocol type <tls|tcp|ssh|kcp>").Short('T').Enum("tls", "tcp", "ssh", "kcp")
	httpArgs.Always = http.Flag("always", "always use parent proxy").Default("false").Bool()
	httpArgs.Timeout = http.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Default("2000").Int()
	httpArgs.HTTPTimeout = http.Flag("http-timeout", "check domain if blocked , http request timeout milliseconds when connect to host").Default("3000").Int()
	httpArgs.Interval = http.Flag("interval", "check domain if blocked every interval seconds").Default("10").Int()
	httpArgs.Blocked = http.Flag("blocked", "blocked domain file , one domain each line").Default("blocked").Short('b').String()
	httpArgs.Direct = http.Flag("direct", "direct domain file , one domain each line").Default("direct").Short('d').String()
	httpArgs.AuthFile = http.Flag("auth-file", "http basic auth file,\"username:password\" each line in file").Short('F').String()
	httpArgs.Auth = http.Flag("auth", "http basic auth username and password, mutiple user repeat -a ,such as: -a user1:pass1 -a user2:pass2").Short('a').Strings()
	httpArgs.PoolSize = http.Flag("pool-size", "conn pool size , which connect to parent proxy, zero: means turn off pool").Short('L').Default("0").Int()
	httpArgs.CheckParentInterval = http.Flag("check-parent-interval", "check if proxy is okay every interval seconds,zero: means no check").Short('I').Default("3").Int()
	httpArgs.Local = http.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()
	httpArgs.SSHUser = http.Flag("ssh-user", "user for ssh").Short('u').Default("").String()
	httpArgs.SSHKeyFile = http.Flag("ssh-key", "private key file for ssh").Short('S').Default("").String()
	httpArgs.SSHKeyFileSalt = http.Flag("ssh-keysalt", "salt of ssh private key").Short('s').Default("").String()
	httpArgs.SSHPassword = http.Flag("ssh-password", "password for ssh").Short('A').Default("").String()
	httpArgs.KCPKey = http.Flag("kcp-key", "key for kcp encrypt/decrypt data").Short('B').Default("encrypt").String()
	httpArgs.KCPMethod = http.Flag("kcp-method", "kcp encrypt/decrypt method").Short('M').Default("3des").String()
	httpArgs.LocalIPS = http.Flag("local bind ips", "if your host behind a nat,set your public ip here avoid dead loop").Short('g').Strings()
	httpArgs.AuthURL = http.Flag("auth-url", "http basic auth username and password will send to this url,response http code equal to 'auth-code' means ok,others means fail.").Default("").String()
	httpArgs.AuthURLTimeout = http.Flag("auth-timeout", "access 'auth-url' timeout milliseconds").Default("3000").Int()
	httpArgs.AuthURLOkCode = http.Flag("auth-code", "access 'auth-url' success http code").Default("204").Int()
	httpArgs.AuthURLRetry = http.Flag("auth-retry", "access 'auth-url' fail and retry count").Default("1").Int()

	//########tcp#########
	// tcp := app.Command("tcp", "proxy on tcp mode")
	// tcpArgs.Parent = tcp.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	// tcpArgs.CertFile = tcp.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	// tcpArgs.KeyFile = tcp.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	// tcpArgs.Timeout = tcp.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Short('e').Default("2000").Int()
	// tcpArgs.ParentType = tcp.Flag("parent-type", "parent protocol type <tls|tcp|kcp|udp>").Short('T').Enum("tls", "tcp", "udp", "kcp")
	// tcpArgs.LocalType = tcp.Flag("local-type", "local protocol type <tls|tcp|kcp>").Default("tcp").Short('t').Enum("tls", "tcp", "kcp")
	// tcpArgs.PoolSize = tcp.Flag("pool-size", "conn pool size , which connect to parent proxy, zero: means turn off pool").Short('L').Default("0").Int()
	// tcpArgs.CheckParentInterval = tcp.Flag("check-parent-interval", "check if proxy is okay every interval seconds,zero: means no check").Short('I').Default("3").Int()
	// tcpArgs.Local = tcp.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()
	// tcpArgs.KCPKey = tcp.Flag("kcp-key", "key for kcp encrypt/decrypt data").Short('B').Default("encrypt").String()
	// tcpArgs.KCPMethod = tcp.Flag("kcp-method", "kcp encrypt/decrypt method").Short('M').Default("3des").String()

	//########udp#########
	// udp := app.Command("udp", "proxy on udp mode")
	// udpArgs.Parent = udp.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	// udpArgs.CertFile = udp.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	// udpArgs.KeyFile = udp.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	// udpArgs.Timeout = udp.Flag("timeout", "tcp timeout milliseconds when connect to parent proxy").Short('t').Default("2000").Int()
	// udpArgs.ParentType = udp.Flag("parent-type", "parent protocol type <tls|tcp|udp>").Short('T').Enum("tls", "tcp", "udp")
	// udpArgs.PoolSize = udp.Flag("pool-size", "conn pool size , which connect to parent proxy, zero: means turn off pool").Short('L').Default("0").Int()
	// udpArgs.CheckParentInterval = udp.Flag("check-parent-interval", "check if proxy is okay every interval seconds,zero: means no check").Short('I').Default("3").Int()
	// udpArgs.Local = udp.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()

	//########tunnel-server#########
	// tunnelServer := app.Command("tserver", "proxy on tunnel server mode")
	// tunnelServerArgs.Parent = tunnelServer.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	// tunnelServerArgs.CertFile = tunnelServer.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	// tunnelServerArgs.KeyFile = tunnelServer.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	// tunnelServerArgs.Timeout = tunnelServer.Flag("timeout", "tcp timeout with milliseconds").Short('t').Default("2000").Int()
	// tunnelServerArgs.IsUDP = tunnelServer.Flag("udp", "proxy on udp tunnel server mode").Default("false").Bool()
	// tunnelServerArgs.Key = tunnelServer.Flag("k", "client key").Default("default").String()
	// tunnelServerArgs.Route = tunnelServer.Flag("route", "local route to client's network, such as :PROTOCOL://LOCAL_IP:LOCAL_PORT@[CLIENT_KEY]CLIENT_LOCAL_HOST:CLIENT_LOCAL_PORT").Short('r').Default("").Strings()
	// tunnelServerArgs.Mux = tunnelServer.Flag("mux", "use multiplexing mode").Default("false").Bool()

	//########tunnel-client#########
	// tunnelClient := app.Command("tclient", "proxy on tunnel client mode")
	// tunnelClientArgs.Parent = tunnelClient.Flag("parent", "parent address, such as: \"23.32.32.19:28008\"").Default("").Short('P').String()
	// tunnelClientArgs.CertFile = tunnelClient.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	// tunnelClientArgs.KeyFile = tunnelClient.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	// tunnelClientArgs.Timeout = tunnelClient.Flag("timeout", "tcp timeout with milliseconds").Short('t').Default("2000").Int()
	// tunnelClientArgs.Key = tunnelClient.Flag("k", "key same with server").Default("default").String()
	// tunnelClientArgs.Mux = tunnelClient.Flag("mux", "use multiplexing mode").Default("false").Bool()

	//########tunnel-bridge#########
	// tunnelBridge := app.Command("tbridge", "proxy on tunnel bridge mode")
	// tunnelBridgeArgs.CertFile = tunnelBridge.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	// tunnelBridgeArgs.KeyFile = tunnelBridge.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	// tunnelBridgeArgs.Timeout = tunnelBridge.Flag("timeout", "tcp timeout with milliseconds").Short('t').Default("2000").Int()
	// tunnelBridgeArgs.Local = tunnelBridge.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()
	// tunnelBridgeArgs.Mux = tunnelBridge.Flag("mux", "use multiplexing mode").Default("false").Bool()

	//########ssh#########
	socks := app.Command("socks", "proxy on ssh mode")
	socksArgs.Parent = socks.Flag("parent", "parent ssh address, such as: \"23.32.32.19:22\"").Default("").Short('P').String()
	socksArgs.ParentType = socks.Flag("parent-type", "parent protocol type <tls|tcp|kcp|ssh>").Default("tcp").Short('T').Enum("tls", "tcp", "kcp", "ssh")
	socksArgs.LocalType = socks.Flag("local-type", "local protocol type <tls|tcp|kcp>").Default("tcp").Short('t').Enum("tls", "tcp", "kcp")
	socksArgs.Local = socks.Flag("local", "local ip:port to listen").Short('p').Default(":33080").String()
	socksArgs.UDPParent = socks.Flag("udp-parent", "udp parent address, such as: \"23.32.32.19:33090\"").Default("").Short('X').String()
	socksArgs.UDPLocal = socks.Flag("udp-local", "udp local ip:port to listen").Short('x').Default(":33090").String()
	socksArgs.CertFile = socks.Flag("cert", "cert file for tls").Short('C').Default("proxy.crt").String()
	socksArgs.KeyFile = socks.Flag("key", "key file for tls").Short('K').Default("proxy.key").String()
	socksArgs.SSHUser = socks.Flag("ssh-user", "user for ssh").Short('u').Default("").String()
	socksArgs.SSHKeyFile = socks.Flag("ssh-key", "private key file for ssh").Short('S').Default("").String()
	socksArgs.SSHKeyFileSalt = socks.Flag("ssh-keysalt", "salt of ssh private key").Short('s').Default("").String()
	socksArgs.SSHPassword = socks.Flag("ssh-password", "password for ssh").Short('A').Default("").String()
	socksArgs.Always = socks.Flag("always", "always use parent proxy").Default("false").Bool()
	socksArgs.Timeout = socks.Flag("timeout", "tcp timeout milliseconds when connect to real server or parent proxy").Default("5000").Int()
	socksArgs.Interval = socks.Flag("interval", "check domain if blocked every interval seconds").Default("10").Int()
	socksArgs.Blocked = socks.Flag("blocked", "blocked domain file , one domain each line").Default("blocked").Short('b').String()
	socksArgs.Direct = socks.Flag("direct", "direct domain file , one domain each line").Default("direct").Short('d').String()
	socksArgs.AuthFile = socks.Flag("auth-file", "http basic auth file,\"username:password\" each line in file").Short('F').String()
	socksArgs.Auth = socks.Flag("auth", "socks auth username and password, mutiple user repeat -a ,such as: -a user1:pass1 -a user2:pass2").Short('a').Strings()
	socksArgs.KCPKey = socks.Flag("kcp-key", "key for kcp encrypt/decrypt data").Short('B').Default("encrypt").String()
	socksArgs.KCPMethod = socks.Flag("kcp-method", "kcp encrypt/decrypt method").Short('M').Default("3des").String()
	socksArgs.LocalIPS = socks.Flag("local bind ips", "if your host behind a nat,set your public ip here avoid dead loop").Short('g').Strings()
	socksArgs.AuthURL = socks.Flag("auth-url", "auth username and password will send to this url,response http code equal to 'auth-code' means ok,others means fail.").Default("").String()
	socksArgs.AuthURLTimeout = socks.Flag("auth-timeout", "access 'auth-url' timeout milliseconds").Default("3000").Int()
	socksArgs.AuthURLOkCode = socks.Flag("auth-code", "access 'auth-url' success http code").Default("204").Int()
	socksArgs.AuthURLRetry = socks.Flag("auth-retry", "access 'auth-url' fail and retry count").Default("0").Int()

	//parse args
	serviceName := kingpin.MustParse(app.Parse(os.Args[1:]))
	flags := log.Ldate
	if *debug {
		flags |= log.Lshortfile | log.Lmicroseconds
	} else {
		flags |= log.Ltime
	}
	log.SetFlags(flags)

	if *logfile != "" {
		f, e := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if e != nil {
			log.Fatal(e)
		}
		log.SetOutput(f)
	}
	if *daemon {
		args := []string{}
		for _, arg := range os.Args[1:] {
			if arg != "--daemon" {
				args = append(args, arg)
			}
		}
		cmd = exec.Command(os.Args[0], args...)
		cmd.Start()
		f := ""
		if *forever {
			f = "forever "
		}
		log.Printf("%s%s [PID] %d running...\n", f, os.Args[0], cmd.Process.Pid)
		os.Exit(0)
	}
	if *forever {
		args := []string{}
		for _, arg := range os.Args[1:] {
			if arg != "--forever" {
				args = append(args, arg)
			}
		}
		go func() {
			for {
				if cmd != nil {
					cmd.Process.Kill()
				}
				cmd = exec.Command(os.Args[0], args...)
				cmdReaderStderr, err := cmd.StderrPipe()
				if err != nil {
					log.Printf("ERR:%s,restarting...\n", err)
					continue
				}
				cmdReader, err := cmd.StdoutPipe()
				if err != nil {
					log.Printf("ERR:%s,restarting...\n", err)
					continue
				}
				scanner := bufio.NewScanner(cmdReader)
				scannerStdErr := bufio.NewScanner(cmdReaderStderr)
				go func() {
					for scanner.Scan() {
						fmt.Println(scanner.Text())
					}
				}()
				go func() {
					for scannerStdErr.Scan() {
						fmt.Println(scannerStdErr.Text())
					}
				}()
				if err := cmd.Start(); err != nil {
					log.Printf("ERR:%s,restarting...\n", err)
					continue
				}
				pid := cmd.Process.Pid
				log.Printf("worker %s [PID] %d running...\n", os.Args[0], pid)
				if err := cmd.Wait(); err != nil {
					log.Printf("ERR:%s,restarting...", err)
					continue
				}
				log.Printf("%s [PID] %d unexpected exited, restarting...\n", os.Args[0], pid)
				time.Sleep(time.Second * 5)
			}
		}()
		return
	}
	if *logfile == "" {
		poster()
	}
	//regist services and run service
	services.Regist("http", services.NewHTTP(), httpArgs)
	// services.Regist("tcp", services.NewTCP(), tcpArgs)
	// services.Regist("udp", services.NewUDP(), udpArgs)
	// services.Regist("tserver", services.NewTunnelServerManager(), tunnelServerArgs)
	// services.Regist("tclient", services.NewTunnelClient(), tunnelClientArgs)
	// services.Regist("tbridge", services.NewTunnelBridge(), tunnelBridgeArgs)
	services.Regist("socks", services.NewSocks(), socksArgs)
	service, err = services.Run(serviceName)
	if err != nil {
		log.Fatalf("run service [%s] fail, ERR:%s", serviceName, err)
	}
	return
}

func poster() {
	fmt.Printf(`
		########  ########   #######  ##     ## ##    ## 
		##     ## ##     ## ##     ##  ##   ##   ##  ##  
		##     ## ##     ## ##     ##   ## ##     ####   
		########  ########  ##     ##    ###       ##    
		##        ##   ##   ##     ##   ## ##      ##    
		##        ##    ##  ##     ##  ##   ##     ##    
		##        ##     ##  #######  ##     ##    ##    
		
		v%s`+" by snail , blog : http://www.host900.com/\n\n", APP_VERSION)
}
