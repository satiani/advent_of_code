package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func part1() {
	file, err := os.Open("./input.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var endResult int32 = 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		integer, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		endResult = endResult + int32(integer)
	}

	fmt.Println(endResult)
}

func part2() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	history := make(map[int32]int32)
	var total int32 = 0

	scanner := bufio.NewScanner(file)
	for {
		for scanner.Scan() {
			integer, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}

			total = total + int32(integer)
			if _, ok := history[total]; ok {
				fmt.Printf("The first repeating total is %d", total)
				return
			}
			history[total] = 0
		}
		file.Seek(0, 0)
		scanner = bufio.NewScanner(file)
	}
}

func main() {
	part1()
	part2()
}
