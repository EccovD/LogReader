package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
)

type Data struct {
	Exception struct {
		Message string `json:"Message"`
	} `json:"exception"`
	CustomerId string `json:"CustomerId"`
}

type Count struct {
	Key   string
	Count int
}

func main() {
	fileContent, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	var jsonData []Data
	err = json.Unmarshal(fileContent, &jsonData)
	if err != nil {
		log.Fatalf("Error unmarshalliing %v", err)
	}

	formRegex := regexp.MustCompile(`Form \d+`)

	formCounts := make(map[string]int)
	customerCounts := make(map[string]int)

	for _, data := range jsonData {
		matches := formRegex.FindAllString(data.Exception.Message, -1)
		for _, match := range matches {
			formCounts[match]++
		}
		customerCounts[data.CustomerId]++
	}

	topN := func(countMap map[string]int, n int) []Count {
		var counts []Count
		for key, count := range countMap {
			counts = append(counts, Count{Key: key, Count: count})
		}
		sort.Slice(counts, func(i, j int) bool {
			return counts[i].Count > counts[j].Count
		})
		if len(counts) > n {
			counts = counts[:n]
		}
		return counts
	}

	topForms := topN(formCounts, 5)
	topCustomers := topN(customerCounts, 5)

	fmt.Println("Top 5 Forms:")
	for _, form := range topForms {
		fmt.Printf("%s: %d\n", form.Key, form.Count)
	}
	fmt.Println("\nTop 5 CustomerIDs:")
	for _, customer := range topCustomers {
		fmt.Printf("%s: %d\n", customer.Key, customer.Count)
	}
}
