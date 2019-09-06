package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Claim struct {
	id     uint
	left   uint
	top    uint
	width  uint
	height uint
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var claims []Claim
	r := regexp.MustCompile("#([0-9]+) @ ([0-9]+),([0-9]+): ([0-9]+)x([0-9]+)")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		match := r.FindStringSubmatch(scanner.Text())
		if match != nil {
			id, _ := strconv.Atoi(match[1])
			left, _ := strconv.Atoi(match[2])
			top, _ := strconv.Atoi(match[3])
			width, _ := strconv.Atoi(match[4])
			height, _ := strconv.Atoi(match[5])
			claims = append(claims, Claim{
				uint(id), uint(left), uint(top), uint(width), uint(height),
			})
		}
	}
	fmt.Println(claims)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
