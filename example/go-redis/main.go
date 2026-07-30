package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go-redis/redis"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run . <command>")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  set")
		fmt.Println("  setex")
		fmt.Println("  get")
		fmt.Println("  ttl")
		fmt.Println("  hset")
		fmt.Println("  hget")
		fmt.Println("  hgetall")
		fmt.Println("  publish")
		fmt.Println("  subscribe")
		fmt.Println("  leaderboard")
		os.Exit(1)
	}

	ctx := context.Background()

	setting := redis.Setting{
		Addr:     "localhost:6380",
		Username: "cache-service",
		Password: "cache123",
		PathCert: "../../docker/certs/ca.crt",
	}

	repo, err := redis.NewRepository(setting)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	switch os.Args[1] {

	// Strings
	case "set":
		err := repo.SET(
			ctx,
			"cache:strings:title",
			"Exploring Redis Features",
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("OK")

	case "setex":
		err := repo.SETEX(
			ctx,
			"cache:strings:temporary",
			"Hello Redis",
			10*time.Second,
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("OK")

	case "get":
		value, err := repo.GET(
			ctx,
			"cache:strings:title",
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(value)

	case "ttl":
		ttl, err := repo.TTL(
			ctx,
			"cache:strings:temporary",
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(ttl)

	// Hashes
	case "hset":
		err := repo.HSET(
			ctx,
			"cache:user:1001",
			map[string]interface{}{
				"name":  "Andrian",
				"email": "andrian@example.com",
				"role":  "admin",
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("OK")

	case "hget":
		value, err := repo.HGET(
			ctx,
			"cache:user:1001",
			"name",
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(value)

	case "hgetall":
		profile, err := repo.HGETALL(
			ctx,
			"cache:user:1001",
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(profile)

	// Pub/Sub
	case "publish":
		err := repo.PUBLISH(
			ctx,
			"cache:notifications",
			"Hello Redis!",
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Message published")

	case "subscribe":
		log.Fatal(
			repo.SUBSCRIBE(
				ctx,
				"cache:notifications",
			),
		)

	// Sorted Sets
	case "leaderboard":
		err := repo.ZADD(
			ctx,
			"cache:leaderboard",
			1500,
			"Alice",
		)
		if err != nil {
			log.Fatal(err)
		}

		err = repo.ZADD(
			ctx,
			"cache:leaderboard",
			3200,
			"Bob",
		)
		if err != nil {
			log.Fatal(err)
		}

		err = repo.ZADD(
			ctx,
			"cache:leaderboard",
			2100,
			"Charlie",
		)
		if err != nil {
			log.Fatal(err)
		}

		err = repo.ZADD(
			ctx,
			"cache:leaderboard",
			1800,
			"David",
		)
		if err != nil {
			log.Fatal(err)
		}

		ranking, err := repo.ZREVRANGE(
			ctx,
			"cache:leaderboard",
			0,
			2,
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(ranking)

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
