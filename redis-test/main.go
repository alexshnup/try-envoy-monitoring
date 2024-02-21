package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type JsonData struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func main() {

	//get os env REDISHOST
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "127.0.0.1"
	}
	//get os env REDISPORT
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisReadOnly := false
	envReadOnly := os.Getenv("REDIS_READONLY")
	if envReadOnly != "" {
		redisReadOnly = true
	}
	redisKey := os.Getenv("REDIS_KEY")
	if redisKey == "" {
		redisKey = "key"
	}
	// Initialize a Redis client.
	rdb := redis.NewClient(&redis.Options{
		// Addr:     "localhost:6379", // Redis server address
		Addr:     redisHost + ":" + redisPort, // Redis server address
		Password: "",                          // No password set
		DB:       0,                           // Default DB
	})

	// // for {
	// // Reading the value back from Redis.
	// val, err := rdb.Get(ctx, "key").Result()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("key", val)

	//parse json from  val to JsonData

	if !redisReadOnly {
		jsonData := JsonData{}
		// //parse
		// err = json.Unmarshal([]byte(val), &jsonData)
		// if err != nil {
		// 	// panic(err)
		// 	fmt.Println("err", err)
		// }
		jsonData.Name = "test"
		jsonData.Count = time.Now().Second()*1000 + time.Now().Nanosecond()/1000000

		fmt.Printf("JsonData %+v\n", jsonData)

		//marshal
		jsonDataBytes, err := json.Marshal(jsonData)
		if err != nil {
			// panic(err)
			fmt.Println("err", err)
		}

		// Writing a value to Redis.
		err = rdb.Set(ctx, redisKey, jsonDataBytes, 0).Err()
		if err != nil {
			panic(err)
		}
	}

	// Reading the value back from Redis.
	val, err := rdb.Get(ctx, redisKey).Result()
	if err != nil {
		// panic(err)
		fmt.Println("err", err)
	}

	fmt.Printf("%s=%s\n", redisKey, val)

	// 	time.Sleep(1 * time.Second)
	// }
}
