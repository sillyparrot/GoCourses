package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	filename := flag.String("file", "problems.csv", "csv file containing the problems")
	timeout := flag.Int("timeout", 0, "the time limit in seconds to complete all problems")
	shuffle := flag.Bool("shuffle", false, "shuffle the questions")

	flag.Parse()

	// Open the CSV file
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)
	problems, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var num_correct int
	num_problems := len(problems)
	problem_num := 1

	if *shuffle {
		rand.Shuffle(num_problems, func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	ch_finished := make(chan bool)
	startTime := time.Now()

	var timer *time.Timer
	if *timeout != 0 {
		timer = time.NewTimer(time.Duration(*timeout) * time.Second)
	}

	// Process each problem
	go func() {
		for _, problem := range problems {
			var answer string
			fmt.Printf("Problem #%d: %s=", problem_num, problem[0])
			fmt.Scan(&answer)
			if strings.TrimSpace(answer) == strings.TrimSpace(problem[1]) {
				num_correct++
			}
			problem_num++
		}
		ch_finished <- true
	}()

	select {
	case <-ch_finished:
		elapsed := time.Since(startTime)
		fmt.Printf("Complete! Elapsed time: %.2f seconds", elapsed.Seconds())
	case <-func() <-chan time.Time {
		if timer != nil {
			return timer.C // use timer if set
		}
		return make(chan time.Time) // Dummy channel to prevent blocking
	}():
	}

	fmt.Printf("\nYou scored %d out of %d\n", num_correct, num_problems)

}
