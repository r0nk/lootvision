package main

import (
	"bufio"
	"fmt"
	"github.com/TwiN/go-color"
	"math"
	"os"
	"strconv"
	"strings"
)

func readCountsFromFile(filename string) (map[string]int, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	counts := make(map[string]int)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		count, err := strconv.Atoi(fields[0])
		if err != nil {
			panic(err)
		}
		counts[fields[1]] = count
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return counts, nil
}

func writeCountsToFile(filename string, counts map[string]int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for word, count := range counts {
		fmt.Fprintf(writer, "%d\t%s\n", count, word)
	}

	return nil
}

func rarity_color(count int, max int) string {
	ratio := float64(count) / float64(max)
	info := math.Log2(1 / ratio)

	switch {
	case info >= 8:
		return color.Yellow
	case info > 7:
		return color.Purple
	case info > 5:
		return color.Cyan
	case info > 3:
		return color.Green
	case info > 2:
		return color.White
	}
	return color.Gray
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: lootvision <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	counts, err := readCountsFromFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	max := 0
	for _, v := range counts {
		if v > max {
			max = v
		}
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		for _, word := range words {
			counts[word]++
			fmt.Printf("%s ", rarity_color(counts[word], max)+word+color.Reset)
		}
		fmt.Printf("\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
		os.Exit(1)
	}

	err = writeCountsToFile(filename, counts)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}
}
