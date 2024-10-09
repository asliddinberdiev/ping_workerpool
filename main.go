package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asliddinberdiev/ping_workerpool/workerpool"
)

const (
	INTERVAL        = time.Second * 30
	REQUEST_TIMEOUT = time.Second * 3
	WORKERS_COUNT   = 3
)

var urls = []string{
	"https://agsu.uz",
	"https://qalqonacademy.uz",
	"https://google.com",
	"https://golang.org",
}

func main() {
	results := make(chan workerpool.Result)
	workerPool := workerpool.New(WORKERS_COUNT, REQUEST_TIMEOUT, results)

	workerPool.Init()

	go generateJobs(workerPool)
	go proccessResults(results)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	workerPool.Stop()
}

func proccessResults(results chan workerpool.Result) {
	go func() {
		for result := range results {
			fmt.Println(result.Info())
		}
	}()
}

func generateJobs(wp *workerpool.Pool) {
	for {
		for _, url := range urls {
			wp.Push(workerpool.Job{URL: url})
		}

		time.Sleep(INTERVAL)
	}
}
