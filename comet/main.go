package main

import (
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
	"im/libs/perf"
	"runtime"
)

var (
	DefaultServer *Server
	// Debug bool
)

func main() {
	flag.Parse()

	if err := InitConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 设置cpu 核数
	runtime.GOMAXPROCS(Conf.Base.MaxProc)
	// 使用logrus包

	// log.Info("111 noteworthy happened!")
	// 加入性能监控
	perf.Init(Conf.Base.PprofBind)

	if err := InitLogicRpc(Conf.RpcLogicAddrs); err != nil {

		log.Panicf("InitLogicRpc Fatal error: %s \n", err)
	}

	// new server
	Buckets := make([]*Bucket, Conf.Bucket.Num)

	for i := 0; i < Conf.Bucket.Num; i++ {
		Buckets[i] = NewBucket(BucketOptions{
			ChannelSize: Conf.Bucket.Channel,
			RoomSize:    Conf.Bucket.Room,
			RoutineAmount: Conf.Bucket.RoutineAmount,
			RoutineSize: Conf.Bucket.RoutineSize,
		})
	}
	operator := new(DefaultOperator)
	DefaultServer = NewServer(Buckets, operator, ServerOptions{
		WriteWait:       Conf.Base.WriteWait,
		PongWait:        Conf.Base.PongWait,
		PingPeriod:      Conf.Base.PingPeriod,
		MaxMessageSize:  Conf.Base.MaxMessageSize,
		ReadBufferSize:  Conf.Base.ReadBufferSize,
		WriteBufferSize: Conf.Base.WriteBufferSize,
		BroadcastSize:   Conf.Base.BroadcastSize,
	})

	// log.Infof("server %v", DefaultServer)
	// log.Panicf("buckets :%v", buckets)
	log.Info("start InitPushRpc")
	if err := InitPushRpc(Conf.RpcPushAdds); err != nil {
		log.Fatal(err)
	}
	if err := InitWebsocket(Conf.Websocket.Bind); err != nil {
		log.Fatal(err)
	}






}
