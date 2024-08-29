package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sort"
	"time"
)

type Stats struct {
	uniqueURLs     map[string]bool
	domainCount    map[string]int
	ipCount        map[string]int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run lab1.go log1.txt log2.txt ...")
		return
	}

	start := time.Now()
	stats := Stats{
		uniqueURLs:  make(map[string]bool),
		domainCount: make(map[string]int),
		ipCount:     make(map[string]int),
	}

	for _, file := range os.Args[1:] {
		fmt.Printf("Reading %s...\n", file)
		readLog(file, &stats)
	}

	printSummary(&stats, start)
}

func readLog(filename string, stats *Stats) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", filename)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if  err != nil {
			fmt.Printf("Error reading file: %s, because: %s\n", filename, err)
			break
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 4 {
			continue
		}

		ip := fields[2]
		rawURL := fields[3]

		// Track unique URLs
		stats.uniqueURLs[rawURL] = true

		// Parse and track domains
		domain := strings.Split(rawURL, "/")[2]
		
		stats.domainCount[domain]++
		
		// Track IPs (crawlers)
		stats.ipCount[ip]++
	}
}

func printSummary(stats *Stats, start time.Time) {
	fmt.Printf("\n* Unique URLs: %d\n", len(stats.uniqueURLs))
	fmt.Printf("* Unique Domains: %d\n", len(stats.domainCount))

	fmt.Println("* Top 10 Websites:")
	printTop(stats.domainCount, 10)

	fmt.Println("* Top 5 crawlers:")
	printTop(stats.ipCount, 5)

	duration := time.Since(start).Seconds()
	fmt.Printf("\nCompleted in %.1fs.\n", duration)
}

func printTop(counter map[string]int, top int) {
	type kv struct {
		Key   string
		Value int
	}
	// Convert the map to a slice of kv pairs
	var kvSlice []kv
	for k, v := range counter {
		kvSlice = append(kvSlice, kv{Key: k, Value: v})
	}

	// Sort the slice by Value in descending order
	sort.Slice(kvSlice, func(i, j int) bool {
		return kvSlice[i].Value > kvSlice[j].Value
	})

	// Print the top N elements
	for i := 0; i < top && i < len(kvSlice); i++ {
		fmt.Printf("    - %s\n", kvSlice[i].Key)
	}
}
