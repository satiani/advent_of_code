package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

func part1() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	number_of_twos := 0
	number_of_threes := 0
	for scanner.Scan() {
		text := string(scanner.Text())
		letter_counts := make(map[rune]int)
		for _, char := range text {
			letter_counts[char]++
		}

		two_counted := false
		three_counted := false

		for _, v := range letter_counts {
			if !two_counted && v == 2 {
				number_of_twos++
				two_counted = true
			}

			if !three_counted && v == 3 {
				number_of_threes++
				three_counted = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("checksum is: %d\n", number_of_twos*number_of_threes)
}

type Pair struct {
	index int
	line  string
}

func commonStringFindWorker(stringsList []string, queue <-chan Pair, done <-chan struct{}) {
	for {
		select {
		case pair := <-queue:
			for i := pair.index + 1; i < len(stringsList); i++ {
				runeComparedWith := []rune(stringsList[i])
				differences := 0
				differingPosition := -1
				for pos, char := range pair.line {
					if runeComparedWith[pos] != char {
						differences++
						differingPosition = pos
						if differences > 1 {
							break
						}
					}
				}

				if differences == 1 {
					finalString := make([]rune, len(pair.line)-1)
					copy(finalString, runeComparedWith[:differingPosition])
					copy(finalString[differingPosition:], runeComparedWith[differingPosition+1:])
					fmt.Println(string(finalString))
				}
			}
		case <-done:
			return
		}
	}
}

func part2() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var stringList []string

	for scanner.Scan() {
		stringList = append(stringList, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Bounded parallelism
	numWorkers := 10
	queue := make(chan Pair)
	done := make(chan struct{})
	var wg sync.WaitGroup
	defer func() {
		close(done)
		wg.Wait()
	}()

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			commonStringFindWorker(stringList, queue, done)
			wg.Done()
		}()
	}

	for index, line := range stringList {
		queue <- Pair{index, line}
	}
}

func main() {
	part2()
}
