package main

import (
	"flag"
	"net"
	"runtime"
	"time"

	"github.com/rpcxio/rpcx-benchmark/grpc/pb"
	"github.com/smallnest/rpcx/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"runtime/pprof"
	"os/signal"
	"syscall"
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
	f, err := os.Create("cpu.prof")
   	if err != nil {
         log.Fatal("could not create CPU profile: ", err)
        }

    	// Start CPU profiling
    	if err := pprof.StartCPUProfile(f); err != nil {
        	log.Fatal("could not start CPU profile: ", err)
    	}
    	defer pprof.StopCPUProfile()
	lis, err := net.Listen("tcp", *host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterHelloServer(s, &Hello{})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//s.Serve(lis)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-signalChan
	// Gracefully stop the server

	// Gracefully stop the server
	s.GracefulStop()
	pprof.StopCPUProfile()
 	memFile, err := os.Create("mem.prof")
    	if err != nil {
        	log.Fatal("could not create memory profile: ", err)
    	}
    	defer memFile.Close()

    	runtime.GC() // Run garbage collection to get up-to-date memory stats
    	if err := pprof.WriteHeapProfile(memFile); err != nil {
        	log.Fatal("could not write memory profile: ", err)
    	}

    	//log.Println("Profiling data saved. Exiting...")
}
