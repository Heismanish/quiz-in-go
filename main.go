package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func problemPuller(filename string) ([]problem, error) {
	// 1 .open the file
	if fObj, err := os.Open(filename); err == nil {
		// 2. create a new reader
		csvR := csv.NewReader(fObj)

		// 3. it will need to read file
		if cLines, err := csvR.ReadAll(); err == nil {
			// 4. call the parseProblem
			return parseProblem(cLines), nil
		} else {
			return nil, fmt.Errorf("error in reading data in csv from %s file %s", filename, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening %s file %s", filename, err.Error())
	}

}

func main() {
	fmt.Print("Welcome to the quiz ")

	// 1. Input the name of the file
	fName := flag.String("f", "quiz.csv", "path of csv file")

	// 2. Set the duration of the timer
	timer := flag.Int("t", 30, "timer of the quiz")
	flag.Parse()

	// 3. Pull the problems from the file (calling our problem puller func)
	problems, err := problemPuller(*fName)

	// 4. Handle the errors
	if err != nil {
		exit(fmt.Sprintf("something went wrong: %s", err.Error()))
	}

	// 5. Create a variable to count our correct answers
	correctAns := 0
	// 6. Using the duration of the timer, we want to intailise the timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)

	// 7.  Loop thorugh the problems, print the questions, we will accept the answers
problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d : %s =", i+1, p.q)

		go func() {
			fmt.Scanf("%s", &answer)
			ansC <- answer
		}()

		select {
		case <-tObj.C:
			fmt.Println()
			break problemLoop
		case iAns := <-ansC:
			if iAns == p.a {
				correctAns++
			}
			if i == len(problems) {
				close(ansC)
			}
		}
	}

	// 8. We will calculate and print out the result
	fmt.Printf("Your result is %d out of %d \n", correctAns, len(problems))

	// Close the channel
	close(ansC)

}

func parseProblem(lines [][]string) []problem {
	// go over all the lines and parse them, with problem structs
	r := make([]problem, len(lines))

	for i := 0; i < len(lines); i++ {
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}
	return r
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
