package main

import (
	"flag"
	"fmt"
	"encoding/csv"
	"strings"
	"os"
	"time"
)

func main() {
	csvFileNameName := flag.String("csv", "problems.csv", "a csv file in the format of 'questions, answers'")
	timeLimit := flag.Int("time", 10, "the time limit for the quiz")
	
	flag.Parse()
	file, err := os.Open(*csvFileNameName)
	if err != nil{
		exit(fmt.Sprintf("Failed to open %s\n", *csvFileNameName))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil{
		exit("Failed to read csv file")
	}
	problems:= parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)	

	count := 0
problemLoop:
	for i, p := range(problems){
		fmt.Printf("Problem #%d: %s = \n", i + 1, p.q)
		answerCh := make(chan string)
    
		go func(){
			var answer string
		  fmt.Scanln(&answer)
			answerCh <- answer
		}()
		select{
		  case <-timer.C:
			  fmt.Println()
			  break problemLoop
		case answer := <-answerCh:
			  
		    if answer == p.a{
			    count++
		    }
		}
	}
	fmt.Printf("You got %d correct out of %d problems\n", count, len(problems))
}

func parseLines(lines [][]string) []problem{
	ret := make([]problem, len(lines))
	for i, line := range lines{
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct{
	q string
	a string
}

func exit(msg string){
	fmt.Println(msg)
	os.Exit(1)
}