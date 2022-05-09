package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"ULID-with-sequence-GoLang/crockford"
	"strings"

	"os"
	"strconv"
	"time"
)

func main() {
	fmt.Print("Please enter test duration in seconds: ")
	myScanner := bufio.NewScanner(os.Stdin)
	myScanner.Scan()
	line := myScanner.Text()
	duration, _ := strconv.ParseFloat(line, 64)
	startTime := time.Now().Unix()
	endTime := float64(startTime) + duration

	channel := make(chan string)
	counter := 0
	go generator(channel)

	for float64(time.Now().Unix()) <= endTime {
		fmt.Println(<-channel)
		counter += 1
	}
	fmt.Println("Seconds per ULID:", duration/float64(counter))
	fmt.Println("Number of ULIDs per second:", float64(counter)/duration)
}

func generator(channel chan string) {
	newSequence := 0
	oldTimestamp := 0

	for newSequence <= 32767 { // serial search of new_sequence
		newTimestamp := int(time.Now().Unix() * 1000)   // timestamp calculation
		newRandomness, _ := rand.Prime(rand.Reader, 65) // randomness calculation
		if oldTimestamp == newTimestamp {               // within a millisecond
			newSequence += 1 // sequence increment
		} else { // new millisecond
			newSequence = 0
		}
		binTimestamp := fmt.Sprintf("%b", newTimestamp)
		binTimestamp = "0000" + binTimestamp
		crockTimestamp := crockford.Encode(binTimestamp)
		binSequence := fmt.Sprintf("%b", newSequence)
		binSequence = strings.Repeat("0", 15-len(binSequence)) + binSequence
		crockSequence := crockford.Encode(binSequence)
		binRandomness := fmt.Sprintf("%b", newRandomness)
		crockRandomness := crockford.Encode(binRandomness)
		newUlid := strings.Repeat("0", 10-len(crockTimestamp)) + crockTimestamp +
			strings.Repeat("0", 3-len(crockSequence)) + crockSequence +
			strings.Repeat("0", 13-len(crockRandomness)) + crockRandomness

		channel <- newUlid
		oldTimestamp = newTimestamp
	}
	c := 0
	for oldTimestamp == int(time.Now().Unix()*1000) {
		c += 1
	}
	go generator(channel)
}
