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
	flag.Parse()

	problems := parseCsvProblems(*csvFilename)
	fmt.Printf("%v\n", problems)

	reader := bufio.NewReader(os.Stdin)

	var answered int
	var correct int

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Ending early!")
		fmt.Printf("%v correct out of %v total.\n", correct, answered)
		os.Exit(0)
	}()

	for i, problem := range problems {
		answered++
		fmt.Printf("Problem #%d: %s = ", i, problem.question)
		guess, _ := reader.ReadString('\n')
		if strings.TrimSpace(guess) == problem.answer {
			correct++
		}
	}

	fmt.Printf("%v correct out of %v total.\n", correct, answered)
}
