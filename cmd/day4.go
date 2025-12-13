package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// day4Cmd represents the day4 command
var day4Cmd = &cobra.Command{
	Use:   "day4",
	Short: "Accessible rolls",
	Run:   runDay4,
}

func init() {
	rootCmd.AddCommand(day4Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// day4Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// day4Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runDay4(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input-file")
	// isFollowUp, _ := cmd.Flags().GetBool("follow-up")

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	rows := make([][]rune, 3)

	filledRows := 0
	middleRowIndex := 1

	c := make(chan int)
	toWait := 0

	for scanner.Scan() {
		toWait++
		rows[(middleRowIndex-1)%3] = []rune(scanner.Text())
		filledRows++
		middleRowIndex++
		if filledRows < 2 {
			// until we have at least 2 rows we can't assess much
			continue
		}
		if filledRows == 2 {
			// special case: first row (no rolls above)
			go countAccessibleRolls(nil, rows[0], rows[1], c, toWait-2)
			continue
		}
		go countAccessibleRolls(rows[(middleRowIndex-1)%3], rows[middleRowIndex%3], rows[(middleRowIndex+1)%3], c, toWait-2)
	}

	// special case: last row (no rolls below)
	go countAccessibleRolls(rows[(middleRowIndex)%3], rows[(middleRowIndex+1)%3], nil, c, toWait-1)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	accessibleRolls := 0
	for range toWait {
		accessibleRolls += <-c
	}

	log.Printf("There are %d accessible rolls in the map", accessibleRolls)
}

func countAccessibleRolls(prev, cur, next []rune, c chan int, rowIdx int) {
	accessibleRolls := 0
	if prev == nil {
		prev = make([]rune, len(cur))
		for i := range len(cur) {
			prev[i] = '.'
		}
	}
	if next == nil {
		next = make([]rune, len(cur))
		for i := range len(cur) {
			next[i] = '.'
		}
	}
	for i, r := range cur {
		switch r {
		case '.': // empty space, nothing to check
			continue
		case '@': // roll, check if it's accessible
			rollPrev, rollCur, rollNext := []rune{}, []rune{}, []rune{}

			if i == 0 {
				rollPrev = append(rollPrev, '.')
				rollCur = append(rollCur, '.')
				rollNext = append(rollNext, '.')
			} else {
				rollPrev = append(rollPrev, prev[i-1])
				rollCur = append(rollCur, cur[i-1])
				rollNext = append(rollNext, next[i-1])
			}
			rollPrev = append(rollPrev, prev[i])
			rollCur = append(rollCur, cur[i])
			rollNext = append(rollNext, next[i])
			if i == len(cur)-1 {
				rollPrev = append(rollPrev, '.')
				rollCur = append(rollCur, '.')
				rollNext = append(rollNext, '.')
			} else {
				rollPrev = append(rollPrev, prev[i+1])
				rollCur = append(rollCur, cur[i+1])
				rollNext = append(rollNext, next[i+1])
			}
			if isRollAccessible(rollPrev, rollCur, rollNext) {
				accessibleRolls++
			}
		default:
			panic(fmt.Sprintf("what is this? I don't know what %c is", r))
		}
	}
	// log.Printf("Found %v accessible rolls in row %v:\n%c\n%c\n%c\n", accessibleRolls, rowIdx, prev, cur, next)
	c <- accessibleRolls
}

func isRollAccessible(prev, cur, next []rune) bool {
	adjacentRolls := 0
	for i := range 3 {
		if prev[i] == '@' {
			adjacentRolls++
		}
		if next[i] == '@' {
			adjacentRolls++
		}
		if i != 1 && cur[i] == '@' {
			adjacentRolls++
		}
	}
	return adjacentRolls < 4
}
