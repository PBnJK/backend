package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
)

type Scores struct {
	Easy   int `json:"easy"`
	Medium int `json:"medium"`
	Hard   int `json:"hard"`
}

var (
	userInput string
	scores    Scores
)

func getCustomScore() int {
loop:
	for {
		fmt.Print("Type your custom number of chances: ")
		fmt.Scanln(&userInput)

		guess, err := strconv.Atoi(userInput)
		if err != nil {
			fmt.Printf("'%s' isn't even a number! Try again:\n\n", userInput)
			goto loop
		}

		if guess < 1 || guess > 100 {
			fmt.Print("You can't have less than 1 chance! Try again:\n\n")
			goto loop
		}

		return guess
	}
}

const (
	EASY   = 0
	MEDIUM = 1
	HARD   = 2
	CUSTOM = 3

	DB_NAME = "db.json"
)

func playGame() {
	answer := rand.IntN(100) + 1
	diff := EASY

	numberOfChances := 0
	chances := 1

	fmt.Print("\nOK! I'm thinking of a number between 1 and 100...\n\n")

loop:
	for {
		fmt.Println("Please select a difficulty level:")
		fmt.Println("1. Easy (10 chances);")
		fmt.Println("2. Medium (5 chances);")
		fmt.Println("3. Hard (3 chances);")
		fmt.Println("4. Custom.")
		fmt.Print("> ")

		fmt.Scanln(&userInput)
		fmt.Println()
		switch userInput {
		case "1", "E", "e":
			fmt.Println("Alright! You have selected the easy difficulty level.")
			diff = EASY
			numberOfChances = 10
		case "2", "M", "m":
			fmt.Println("Alright! You have selected the medium difficulty level.")
			diff = EASY
			numberOfChances = 5
		case "3", "H", "h":
			fmt.Println("Alright! You have selected the hard difficulty level.")
			diff = HARD
			numberOfChances = 3
		case "4", "C", "c":
			numberOfChances = getCustomScore()
			diff = CUSTOM
			fmt.Printf("\nAlright! You have %d chances then.\n", numberOfChances)
		default:
			fmt.Printf("No such option %s. Try again:\n\n", userInput)
			continue loop
		}

		fmt.Print("Let's start the game!\n\n")

	gameLoop:
		if chances > numberOfChances {
			fmt.Printf("Aww... seems like you couldn't guess the number :(\nIt was %d!\n\n", answer)
			break loop
		}

		fmt.Print("Enter your guess: ")
		fmt.Scanln(&userInput)

		guess, err := strconv.Atoi(userInput)
		if err != nil {
			fmt.Printf("'%s' isn't even a number! Try again:\n\n", userInput)
			goto gameLoop
		}

		if guess < 1 || guess > 100 {
			fmt.Print("You should guess a number from 1 to 100! Try again:\n\n")
			goto gameLoop
		}

		switch {
		case guess == answer:
			fmt.Printf("Congratulations! You guessed the correct number in %d attempts\n\n", chances)
			break loop
		case guess < answer:
			fmt.Printf("Incorrect! The number is greater than %d\n\n", guess)
		case guess > answer:
			fmt.Printf("Incorrect! The number is less than %d\n\n", guess)
		}

		chances++
		goto gameLoop
	}

	switch diff {
	case EASY:
		if scores.Easy > chances {
			fmt.Printf("New high score on easy! %d guesses -> %d!\n\n", scores.Easy, chances)
			scores.Easy = chances
		}
	case MEDIUM:
		if scores.Medium > chances {
			fmt.Printf("New high score on medium! %d guesses -> %d!\n\n", scores.Medium, chances)
			scores.Medium = chances
		}
	case HARD:
		if scores.Hard > chances {
			fmt.Printf("New high score on hard! %d guesses -> %d!\n\n", scores.Hard, chances)
			scores.Hard = chances
		}
	}
}

func seeHighScore() {
	fmt.Println("\nHigh-Scores:")
	fmt.Printf("Easy.... %02d\n", scores.Easy)
	fmt.Printf("Medium.. %02d\n", scores.Medium)
	fmt.Printf("Hard.... %02d\n\n", scores.Hard)
}

func loadHighScores() {
	file, err := os.ReadFile(DB_NAME)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		scores.Easy = 10
		scores.Medium = 5
		scores.Hard = 3
	}

	if err := json.Unmarshal(file, &scores); err != nil {
		return
	}
}

func saveHighScores() {
	data, err := json.Marshal(&scores)
	if err != nil {
		return
	}

	if err := os.WriteFile(DB_NAME, data, 0644); err != nil {
		return
	}
}

func main() {
	loadHighScores()

	fmt.Print("Welcome to the great Number Guessing Game!\n\n")

loop:
	for {
		fmt.Println("Choose from the options below:")
		fmt.Println("1. Play the Game;")
		fmt.Println("2. See High-Scores;")
		fmt.Println("3. Exit.")
		fmt.Print("> ")

		fmt.Scanln(&userInput)
		switch userInput {
		case "1", "P", "p":
			playGame()
		case "2", "S", "s":
			seeHighScore()
		case "3", "E", "e":
			break loop
		default:
			fmt.Printf("\nNo such option %s. Try again:\n\n", userInput)
			continue loop
		}
	}

	saveHighScores()
}
