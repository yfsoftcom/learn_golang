/*
a pool for sharing resources.
*/
package main

import (
	"errors"
	"io"
	"log"
	"sync"
)

type Pool struct {
	m sync.Mutex
	// close flag
	closed    bool
	resources chan io.Closer
	factory   func() (io.Closer, error)
}

type Connection struct {
	ID int
}

func (conn *Connection) Close() error {
	log.Printf("Connection: %d closed!\n", conn.ID)
	return nil
}

func (pool *Pool) Close() {
	pool.m.Lock()
	defer pool.m.Unlock()

	if pool.closed {
		return
	}
	pool.closed = true
	close(pool.resources)
	for r := range pool.resources {
		r.Close()
	}
	log.Printf("Pool: closed!\n")
}

func (pool *Pool) Release(conn *Connection) {

	if pool.closed {
		return
	}

	pool.m.Lock()
	defer pool.m.Unlock()
	select {
	case pool.resources <- conn:
		log.Printf("Pool: Release, Connection id: %d!\n", conn.ID)
	default:
		log.Printf("Pool overflow! close the release: %d\n", conn.ID)
		conn.Close()
	}

}

func (pool *Pool) Get() (*Connection, error) {

	if pool.closed {
		return nil, errors.New("pool closed!")
	}
	// 这里不需要 lock
	// pool.m.Lock()
	// defer pool.m.Unlock()
	// 这里如果没有 select ，第一次获取会发生阻塞，被 lock 之后造成 deadlock
	select {
	case conn, ok := <-pool.resources:
		if !ok {
			log.Println("pool closed")
		}

		return conn.(*Connection), nil

	default:
		closer, err := pool.factory()
		if err != nil {
			return nil, err
		}
		return closer.(*Connection), nil
	}
	return nil, errors.New("unexcepted ")
}

func New(factory func() (io.Closer, error), size uint8) (*Pool, error) {
	return &Pool{
		factory:   factory,
		resources: make(chan io.Closer, size),
	}, nil
}

func main() {
	log.Println("start")
	i := 0
	pool, _ := New(func() (io.Closer, error) {

		i += 1
		return &Connection{
			ID: i,
		}, nil
	}, 10)

	conns := []*Connection{}
	for j := 0; j < 30; j++ {
		conn, err := pool.Get()
		if err != nil {
			log.Fatal(err)
		}
		conns = append(conns, conn)
		log.Printf("Connection do %d", conn.ID)
	}
	for _, conn := range conns {
		pool.Release(conn)
	}

	pool.Close()
}
