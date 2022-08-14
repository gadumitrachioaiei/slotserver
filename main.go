package main

import (
	"flag"
	"log"

	ourredis "github.com/gadumitrachioaiei/slotserver/redis"
	"github.com/gadumitrachioaiei/slotserver/server"

	"github.com/go-redis/redis/v9"
)

const (
	defaultAddr      = "localhost:8080"
	defaultRedisAddr = "redis:6379"
	defaultChips     = 1000
)

var httpAddr = flag.String("http", defaultAddr, "HTTP service address")
var redisAddr = flag.String("redis", defaultRedisAddr, "Redis service address")

func main() {
	flag.Parse()
	rdb := redis.NewClient(&redis.Options{
		Addr: *redisAddr,
	})
	defer rdb.Close()
	userService, err := ourredis.NewUsers(rdb, defaultChips)
	if err != nil {
		log.Fatalf("cannot use users service: %v", err)
	}
	s := server.New(*httpAddr, userService)
	if err := s.Start(); err != nil {
		log.Fatalf("cannot start webserver: %v", err)
	}
}
