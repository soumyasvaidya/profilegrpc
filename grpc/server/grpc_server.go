package main

import (
	"flag"
	"net"
	"runtime"
	"time"

	rpcx "log"
	"net/http"

	_ "net/http/pprof"

	"github.com/rpcxio/rpcx-benchmark/grpc/pb"
	"github.com/smallnest/rpcx/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	host  = flag.String("s", "127.0.0.1:8972", "listened ip and port")
	delay = flag.Duration("delay", 0, "delay to mock business processing")
)

type Hello struct{}

func (t *Hello) Say(ctx context.Context, args *pb.BenchmarkMessage) (reply *pb.BenchmarkMessage, err error) {
	s := "OK"
	var i int32 = 100
	args.Field1 = s
	args.Field2 = i
	if *delay > 0 {
		time.Sleep(*delay)
	} else {
		runtime.Gosched()
	}
	return args, nil
}

func main() {
	flag.Parse()
	go func() {
		rpcx.Println(http.ListenAndServe("localhost:6060", nil)) // Enable pprof on http://localhost:6060
	}()

	lis, err := net.Listen("tcp", *host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterHelloServer(s, &Hello{})
	s.Serve(lis)
	//rpcx.Fatal(http.ListenAndServe(":8080", nil)) // Start your server
}
