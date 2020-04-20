// the output
/*
 # use multi proc cores
 50 * 10000 ~ 9.64s
 500 * 10000 ~ 111.46s

 # use single proc core
 50 * 10000 ~ 13.27
 500 * 10000 ~ 184.05
*/

package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "grpc-foo/foo"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	wg sync.WaitGroup
)

const (
	networkType = "tcp"
	address     = "localhost:5009"
	parallel    = 50    //连接并行度
	times       = 10000 //每连接请求次数
)

type Data struct{}

func main() {

	// only one logic proc core, but actully it will switch automatically of the pythically.
	log.Printf("runtime.NumCPU(): %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	currTime := time.Now()

	//并行请求
	for i := 0; i < int(parallel); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exe()
		}()
	}
	wg.Wait()

	log.Printf("time taken: %.2f ", time.Now().Sub(currTime).Seconds())
}

func exe() {
	//建立连接
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	for i := 0; i < int(times); i++ {
		SayHello(client)
	}
}

// implement the interface
// normally feedback with the hi!
func SayHello(client pb.GreeterClient) {
	var request pb.HelloRequest
	r := rand.Intn(parallel)
	request.Name = fmt.Sprintf("%d%s", int32(r), ":Name")

	response, _ := client.SayHello(context.Background(), &request) //调用远程方法

	//判断返回结果是否正确
	if id, _ := strconv.Atoi(strings.Split(response.Message, ":")[0]); id != r {
		log.Printf("response error  %#v", response)
	}

}
