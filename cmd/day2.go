package cmd

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// day2Cmd represents the day2 command
var day2Cmd = &cobra.Command{
	Use:   "day2",
	Short: "Invalid ids",
	Run:   runDay2,
}

func init() {
	rootCmd.AddCommand(day2Cmd)
}

func runDay2(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input-file")
	ranges := readIdRanges(inputFile)
	toWait := len(ranges)

	exitChan := make(chan int)

	for _, r := range ranges {
		go func(r IdRange, exitChan chan int) {
			sum := 0
			lower := r.Lower

			for {
				// 1. find first potential invalid id,
				// i.e. first number in the range with an even number of digits
				candidate := strconv.Itoa(lower)

				if len(candidate)%2 != 0 {
					candidate = "0" + candidate
				}

				// 2. determine target invalid id (for ABCXYZ, it'd be ABCABC)
				halfCandidate := candidate[:len(candidate)/2]
				candidate = halfCandidate + halfCandidate

				// 3. determine whether the target invalid id is within range (lower <= target <= upper)

				candidateInt, _ := strconv.Atoi(candidate)
				if candidateInt > r.Upper {
					break
				}
				if candidateInt >= r.Lower {
					// log.Println("found invalid id candidate", candidateInt, "in range", r)
					sum += candidateInt
				}
				// 4. find next candidate (AB[C+1]AB[C+1]), see if it's within range, abort when it isn't
				nextHalfCandidateInt, _ := strconv.Atoi(halfCandidate)
				nextHalfCandidateInt++
				nextHalfCandidate := strconv.Itoa(nextHalfCandidateInt)
				nextCandidateInt, _ := strconv.Atoi(nextHalfCandidate + nextHalfCandidate)
				if nextCandidateInt > r.Upper {
					break
				}
				lower = nextCandidateInt
			}

			// if sum == 0 {
			// 	log.PrintLn("no invalid ids found in range %v", r)
			// }
			exitChan <- sum
		}(r, exitChan)
	}

	result := 0

	for range toWait {
		result += <-exitChan
	}

	log.Printf("There are %d invalid ids", result)
}

type IdRange struct {
	Lower int
	Upper int
}

func IdRangeFromString(s string) IdRange {
	idPair := strings.Split(s, "-")
	if idPair == nil || len(idPair) != 2 {
		log.Fatal("there's no id range on %q", s)
	}
	lower, _ := strconv.Atoi(idPair[0])
	upper, _ := strconv.Atoi(idPair[1])
	return IdRange{
		Lower: lower,
		Upper: upper,
	}
}

func readIdRanges(filename string) []IdRange {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	ranges := []IdRange{}

	for scanner.Scan() {
		rangeStrings := strings.Split(scanner.Text(), ",")
		for _, r := range rangeStrings {
			ranges = append(ranges, IdRangeFromString(r))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ranges
}
