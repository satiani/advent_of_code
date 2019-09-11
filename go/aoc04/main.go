package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

const REF_TIME = "2006-01-02 15:04"

type EventType int

const (
	BeginsShift EventType = iota
	Sleeps
	WakesUp
)

type GuardEvent struct {
	timestamp *time.Time
	id        int
	eventType EventType
}

func (s GuardEvent) String() string {
	return fmt.Sprintf("%v %v %v", *(s.timestamp), s.id, s.eventType)
}

type GuardEvents []GuardEvent

func NewGuardEvents() GuardEvents {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	timeStatusSplitter := regexp.MustCompile(`\[([^]]+)\] (.*)`)
	guardIdParser := regexp.MustCompile(`Guard #([0-9]+) begins.*`)
	var guardEvents GuardEvents

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		match := timeStatusSplitter.FindAllStringSubmatch(scanner.Text(), -1)

		timeMatch := match[0][1]
		parsedTime, _ := time.Parse(REF_TIME, timeMatch)

		event := match[0][2]
		var parsedEventType EventType
		var id int

		switch event {
		case "wakes up":
			parsedEventType = WakesUp
		case "falls asleep":
			parsedEventType = Sleeps
		default:
			parsedEventType = BeginsShift
			match = guardIdParser.FindAllStringSubmatch(event, -1)
			id, err = strconv.Atoi(match[0][1])
		}

		guardEvents = append(guardEvents, GuardEvent{
			&parsedTime,
			id,
			parsedEventType,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Sort(guardEvents)

	return guardEvents
}

func (evs GuardEvents) Len() int {
	return len(evs)
}

func (evs GuardEvents) Swap(i, j int) {
	evs[i], evs[j] = evs[j], evs[i]
}

func (evs GuardEvents) Less(i, j int) bool {
	return evs[i].timestamp.Before(*evs[j].timestamp)
}

type MinutesSleptByGuard map[int]map[int]int

func calculateMinutesSleptByGuard(guardEvents GuardEvents) MinutesSleptByGuard {
	minutesSleptByGuard := make(MinutesSleptByGuard)

	var currentGuardId int
	var sleepStartMinute int
	hasJustWokenUp := false
	for index := range guardEvents {
		event := &guardEvents[index]
		if event.id != 0 {
			currentGuardId = event.id
		} else {
			event.id = currentGuardId
		}

		switch event.eventType {
		case WakesUp:
			hasJustWokenUp = true
		case Sleeps:
			sleepStartMinute = event.timestamp.Minute()
		default:
		}

		if hasJustWokenUp {
			minutesSlept := minutesSleptByGuard[event.id]
			if minutesSlept == nil {
				minutesSlept = make(map[int]int)
				minutesSleptByGuard[event.id] = minutesSlept
			}
			eventMinute := event.timestamp.Minute()
			for i := 0; i < eventMinute-sleepStartMinute; i++ {
				minutesSlept[sleepStartMinute+i]++
			}
			hasJustWokenUp = false
		}
	}

	return minutesSleptByGuard
}

func part1() {
	guardEvents := NewGuardEvents()
	minutesSleptByGuard := calculateMinutesSleptByGuard(guardEvents)

	var mostMinutesAsleep int
	var guardMostAsleep int
	var minuteMostSlept int

	for guard, minutesSlept := range minutesSleptByGuard {
		totalSlept := 0
		mostSleepsPerMinute := 0
		localMinuteMostSlept := 0
		for minute, total := range minutesSlept {
			totalSlept = totalSlept + total
			if total > mostSleepsPerMinute {
				mostSleepsPerMinute = total
				localMinuteMostSlept = minute
			}
		}
		if totalSlept > mostMinutesAsleep {
			mostMinutesAsleep = totalSlept
			guardMostAsleep = guard
			minuteMostSlept = localMinuteMostSlept
		}
	}

	fmt.Println(minuteMostSlept * guardMostAsleep)
}

func part2() {
	guardEvents := NewGuardEvents()
	minutesSleptByGuard := calculateMinutesSleptByGuard(guardEvents)

	var mostSleepsPerMinute int
	var minuteMostSlept int
	var guardWithMinuteMostSlept int

	for guard, minutesSlept := range minutesSleptByGuard {
		for minute, total := range minutesSlept {
			if total > mostSleepsPerMinute {
				mostSleepsPerMinute = total
				minuteMostSlept = minute
				guardWithMinuteMostSlept = guard
			}
		}
	}

	fmt.Println(minuteMostSlept * guardWithMinuteMostSlept)
}

func main() {
	part1()
	part2()
}
