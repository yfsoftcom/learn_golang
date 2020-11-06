## Notes For Interview

### Golang
- GPM
  - M(Machine): 一个内核线程，由Go的调度器进行管理，G运行在M上
  - P(Processer): 一个逻辑处理器，调度器将其绑定了一个本地的routine队列，其主要指责就是读取可执行的routine，并执行它，其占用了一个CPU线程的资源
  - G(Goroutine): 一个运行中的routine，每个g都会有自己的内存和状态，通常一个P同一时间只会允许1个G运行

  #### CPU线程和内核线程的区别
  - 内核线程可以成千上万，但是CPU线程只和CPU的数量有关

  #### 其他相关概念
  - 全局routine队列
    - 在本地队列无法容纳的g时存放，由P进行有条件的读取
      - 当P处于空闲状态，即本地队列被清空的情况下，会从全局队列中抓取一部分到本地执行
      - 非空闲状态下的P会有一定的概率（1/61）从全局中获取到g进行执行
  - 本地routine队列
    - 程序运行之初，调度器会平均分配routine到各个P的Local队列中，该队列有一个上限（默认是255），超出的话，会放入到全局的队列中
    - 空闲中的P会从其他不空闲中的P中偷取一部分g执行
  - 如果G中的任务发生了阻塞了，GPM会发生什么变化
    - 用户态的阻塞：
      - 比如wait/chan等
      - [TODO] 该g会被标记为阻塞状态，P会跳过该g，执行后续的任务，直到该g完成之后被唤醒。
    - 系统级的阻塞：
      - P会剥离出去，原先的MG会等待g执行完成再放到P的队列中
      - go调度器会创建一个新的MG来绑定这个P
      - g执行完了之后，原有的M会暂存，供之后调用
  - 如果g执行了一个网络I/O阻塞，会发生什么？
    - g被分离之后，会放到带有网络轮询器的routine中进行执行（可能就是epull）

  #### Reference
  - https://mp.weixin.qq.com/s/_ujmGibYT3s61dBkIIeayw
- GC
  - 实现模型？
  - 何时触发？
  - 能否手动执行？
  - SWT
  - write barrier 的作用

  #### Reference
  - https://liujiacai.net/blog/2018/08/04/incremental-gc/
- reflect

- channel
  #### 发送者
  - 检查有没有##等待中##(意思就是有空闲的接受者，因此可以判定该通道已经空了)的接受者等待接收
    - 如果有，直接发送给接受者，然后继续执行当前的g
    - 如果没有，则尝试将数据存储到chan的临时区
      - 如果chan已经满了，则将当前的g挂起，放到runq队列中，等待被读取时唤醒
      - 如果chan没有满，将数据copy到chan的临时区进行存储，继续当前的g

  #### 接受者
  - 如发送者休眠，则唤起发送者，然后将数据读取过来
  - 如果没有休眠者，则从chan的开头部分取出数据（以保证按照顺序进行）
  - 如果两者都没有，则会标记当前的g挂起，等待发送者唤醒

  #### reference
  - https://speakerdeck.com/kavya719/understanding-channels?slide=87
- defer
  - 符合FILO的规则，最先压入的会最后执行
  - defer 后的语句虽然不会立刻执行，但会被立即计算

- slice/array

- make/new

- sync.Map

- init

### Computer

- Process/Thread/Routine

- 用户态/内核态

### Network

- TCP/UDP

- grpc/http/http2.0

- keep alive

### Data Struct


### Algorithm(算法)


### Docker

- docker 原理

- docker网络，如何分配ip

- docker如何优化一个Dockerfile使其减少构建时间

- docker 和传统的 VM 的区别

### K8S

- 如何实现选主的

- 如何实现容器之间跨机器相互访问的

- Ingress 实现了什么，有什么作用

- 如何灰度发布

- 如何创建存储

- 如何感知服务的状态

- 如何使用k8s中的插件