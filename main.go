package main

import (
	"context"
	"flag"
	"log"

	"github.com/gadumitrachioaiei/slotserver/slot"

	ourdynamodb "github.com/gadumitrachioaiei/slotserver/dynamodb"
	ourredis "github.com/gadumitrachioaiei/slotserver/redis"
	"github.com/gadumitrachioaiei/slotserver/server"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-redis/redis/v9"
)

const (
	defaultAddr      = "localhost:8080"
	defaultRedisAddr = "redis:6379"
	defaultChips     = 1000
)

var httpAddr = flag.String("http", defaultAddr, "HTTP service address")
var redisAddr = flag.String("redis", defaultRedisAddr, "Redis service address")
var useDynamoDB = flag.Bool("dynamodb", false, "Use dynamodb for users")

func main() {
	flag.Parse()
	var (
		userService slot.UserService
		err         error
	)
	if *useDynamoDB {
		dynamoDBClient, err := dynamoDBClient()
		if err != nil {
			log.Fatalf("unable to instantiate dynamodb client, %v", err)
		}
		userService, err = ourdynamodb.NewUsers(dynamoDBClient, defaultChips)
	} else {
		rdb := redis.NewClient(&redis.Options{
			Addr: *redisAddr,
		})
		defer rdb.Close()
		userService, err = ourredis.NewUsers(rdb, defaultChips)
	}
	if err != nil {
		log.Fatalf("cannot use users service: %v", err)
	}
	s := server.New(*httpAddr, userService)
	if err := s.Start(); err != nil {
		log.Fatalf("cannot start webserver: %v", err)
	}
}

func dynamoDBClient() (*dynamodb.Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	svc := dynamodb.NewFromConfig(sdkConfig)
	return svc, nil
}
