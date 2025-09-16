package main

import (
	"fmt"
	"site-monitor/internal/checker"
	"sync"
	"time"
)

func main() {
	targets := []string{
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/500",
		"https://httpbin.org/delay/2",
	}

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		var wg sync.WaitGroup
		wg.Add(len(targets))

		for _, target := range targets {
			go func(url string) {
				defer wg.Done()
				result := checker.CheckSite(url)
				fmt.Println(result)
			}(target)
		}

		wg.Wait()
	}
}
