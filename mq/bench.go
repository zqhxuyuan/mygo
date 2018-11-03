package main

import (
	"fmt"
	"time"

	"github.com/tylertreat/bench"
	"github.com/tylertreat/bench/requester"
)

func main() {
	r := &requester.RedisPubSubRequesterFactory{
		URL:         ":6379",
		PayloadSize: 500,
		Channel:     "benchmark",
	}

	benchmark := bench.NewBenchmark(r, 10000, 1, 30*time.Second, 0)
	summary, err := benchmark.Run()
	if err != nil {
		panic(err)
	}

	fmt.Println(summary)
	summary.GenerateLatencyDistribution(nil, "redis.txt")
}
