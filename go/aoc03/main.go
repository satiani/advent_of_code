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

// IterateCoordinates iterates over each of the discrete x,y coordinates that
// a claim has and passes it on to a function in the form of a Coordinate instance.
// The function can return false at any time to break out of the loop, otherwise it
// should return true.
func (c *Claim) IterateCoordinates(f func(c Coordinate) bool) {
	for x := c.left + 1; x <= c.left+c.width; x++ {
		for y := c.top + 1; y <= c.top+c.height; y++ {
			result := f(Coordinate{x, y})
			if !result {
				return
			}
		}
	}
}

type Coordinate struct {
	x uint
	y uint
}

// ParseClaims takes the raw claims from the input text file and return a slice of
// Claim instances.
func ParseClaims() []Claim {
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

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return claims
}

// AggregateClaims returns a map keyed with XxY coordinates of the square with the number of claims
// on them. If a square does not exist assume no claims on it exist. It also returns the width and
// height of the square discovered through iterating over all the claims
func AggregateClaims(claims []Claim) (aggregates map[Coordinate]uint, width uint, height uint) {
	aggregates = make(map[Coordinate]uint, 10)
	for index := range claims {
		claim := &claims[index]
		impliedWidth := claim.left + claim.width
		if impliedWidth > width {
			width = impliedWidth
		}
		impliedHeight := claim.top + claim.height
		if impliedHeight > height {
			height = impliedHeight
		}

		claim.IterateCoordinates(func(c Coordinate) bool {
			aggregates[c]++
			return true
		})
	}

	return
}

func FindNonOverlappingClaims(
	claims []Claim,
	aggregates map[Coordinate]uint,
) (nonOverlapping []*Claim) {
	for i := range claims {
		claim := &claims[i]
		hasContendingClaims := false
		claim.IterateCoordinates(func(c Coordinate) bool {
			if aggregates[c] > 1 {
				hasContendingClaims = true
				return false
			}
			return true
		})

		if !hasContendingClaims {
			nonOverlapping = append(nonOverlapping, claim)
		}
	}
	return
}

func main() {
	claims := ParseClaims()
	aggregates, height, width := AggregateClaims(claims)
	countClaimedTwiceOrMore := 0
	for index := range aggregates {
		if aggregates[index] > 1 {
			countClaimedTwiceOrMore++
		}
	}

	fmt.Printf("The number of squares claimed twice or more is: %d\n", countClaimedTwiceOrMore)
	fmt.Printf("Height: %d\n", height)
	fmt.Printf("Width: %d\n", width)

	nonOverlapping := FindNonOverlappingClaims(claims, aggregates)
	fmt.Println("The non-overlapping claim IDs are:")
	for _, claim := range nonOverlapping {
		fmt.Println(claim.id)
	}
}
