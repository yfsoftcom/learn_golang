package main

import (
	"log"
	"sync"
	"time"
	"math/rand"
	"strconv"
	"strings"
	"fmt"

	pb "grpc-foo/foo"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	wg sync.WaitGroup
)

const (
	networkType = "tcp"
	server      = "127.0.0.1"
	port        = "5009"
	parallel    = 50        //连接并行度
	times       = 100000    //每连接请求次数
)

type Data struct{}

func main() {

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
	conn, _ := grpc.Dial(server + ":" + port)
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	fmt.Println("%T", client)
	for i := 0; i < int(times); i++ {
		SayHello(client)
	}
}

func SayHello(client pb.GreeterClient) {
	var request pb.HelloRequest
	r := rand.Intn(parallel)
	request.Name = fmt.Sprintf("%d%s",int32(r), ":Name")

	response, _ := client.SayHello(context.Background(), &request) //调用远程方法

	//判断返回结果是否正确
	if id, _ := strconv.Atoi(strings.Split(response.Message, ":")[0]); id != r {
			log.Printf("response error  %#v", response)
	}

}
