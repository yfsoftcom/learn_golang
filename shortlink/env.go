package main

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	S Storage
}

func getEnv() *Env {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	pass := os.Getenv("REDIS_PASS")
	if pass == "" {
		pass = ""
	}
	dbS := os.Getenv("REDIS_DB")
	if dbS == "" {
		dbS = "0"
	}
	db, err := strconv.Atoi(dbS)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connect to redis %s, %s, %d", addr, pass, db)

	return &Env{
		S: NewRedisCli(addr, pass, db),
	}
}
