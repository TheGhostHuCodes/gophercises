package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type problem struct {
	question string
	answer   string
}

func parseCsvProblems(filename string) []problem {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Could not open CSV file: '%v'", err)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Could not parse csv file: '%v'", err)
	}

	return parseCsvLines(lines)
}

func parseCsvLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i, line := range lines {
		problems[i] = problem{
			question: strings.TrimSpace(line[0]),
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return problems
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "A CSV file in the format 'question,answer'.")
	timeLimit := flag.Int("limit", 30, "A time limit for the quiz, in seconds.")
	flag.Parse()

	problems := parseCsvProblems(*csvFilename)
	fmt.Printf("%v\n", problems)

	reader := bufio.NewReader(os.Stdin)

	var answered int
	var correct int

	fmt.Println("Press enter to start the quiz.")
	_, _ = reader.ReadString('\n')

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i, problem.question)
		answered++

		answerCh := make(chan string)
		go func() {
			guess, _ := reader.ReadString('\n')
			answerCh <- guess
		}()

		select {
		case <-sigs:
			fmt.Printf("\nExiting early.")
			fmt.Printf("\n%v correct out of %v total.\n", correct, answered)
		case <-timer.C:
			fmt.Printf("\nOut of time.")
			fmt.Printf("\n%v correct out of %v total.\n", correct, answered)
			return
		case guess := <-answerCh:
			if strings.TrimSpace(guess) == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("%v correct out of %v total.\n", correct, answered)
}
