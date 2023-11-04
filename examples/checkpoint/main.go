package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	. "github.com/twoojoo/ttask"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()
	item := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	T(T(T(T(T(
		FromItem("cp-example", item).
			WithStorage(NewRedisStorage(client)),
		MapArray(func(x int) int {
			return x * 3
		})),
		Print[[]int]("before >")),
		Checkpoint("reduce", ReduceArray("", func(acc *string, x int) string {
			return *acc + strconv.Itoa(x)
		}))),
		Delay[string](5*time.Second)),
		Print[string]("after >"),
	).Catch(func(i *Inner, e error) {
		log.Fatal(e)
	}).Run(ctx)
}
