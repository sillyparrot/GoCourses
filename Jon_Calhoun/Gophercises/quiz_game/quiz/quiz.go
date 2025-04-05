// Package quiz creates a quiz based on the given csv file (default is problems.csv).
// The user will have to answer the questions one by one.
// The user can choose to shuffle the questions and to have a time limit.
// At the end of the quiz, the time elapsed and the number of questions correct will be shown.
package quiz

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

// RetrieveProblems is used to read a csv file and retrieve two columns from it.
// question and answer.
func RetrieveProblems(filename string) ([][]string, error) {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)
	problems, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return problems, nil
}

// CreateQuiz function.
func CreateQuiz(in *os.File, problems [][]string, shuffle bool, timeout int) {
	var num_correct int
	if in == nil {
		in = os.Stdin
	}

	if shuffle {
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	ch_finished := make(chan bool)
	startTime := time.Now()

	var timer *time.Timer
	if timeout != 0 {
		timer = time.NewTimer(time.Duration(timeout) * time.Second)
	}

	// Process each problem
	go func() {
		for i, problem := range problems {
			var answer string
			fmt.Printf("Problem #%d: %s=", i, problem[0])
			_, err := fmt.Fscan(in, &answer)
			if err != nil {
				log.Fatal(err)
			}
			if strings.TrimSpace(answer) == strings.TrimSpace(problem[1]) {
				num_correct++
			}
		}
		ch_finished <- true
	}()

	select {
	case <-ch_finished:
		elapsed := time.Since(startTime)
		fmt.Printf("\nComplete! Elapsed time: %.2f seconds\n", elapsed.Seconds())
	case <-func() <-chan time.Time {
		if timer != nil {
			return timer.C // use timer if set
		}
		return make(chan time.Time) // Dummy channel to prevent blocking
	}():
	}

	fmt.Printf("You scored %d out of %d\n", num_correct, len(problems))
}

func Quiz() {
	filename := flag.String("file", "problems.csv", "csv file containing the problems")
	timeout := flag.Int("timeout", 0, "the time limit in seconds to complete all problems")
	shuffle := flag.Bool("shuffle", false, "shuffle the questions")

	flag.Parse()

	problems, err := RetrieveProblems(*filename)
	if err != nil {
		log.Fatal(err)
	}

	CreateQuiz(nil, problems, *shuffle, *timeout)
}
