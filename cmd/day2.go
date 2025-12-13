package cmd

import (
	"bufio"
	"log"
	"math"
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
	isFollowUp, _ := cmd.Flags().GetBool("follow-up")
	ranges := readIdRanges(inputFile)
	toWait := len(ranges)

	exitChan := make(chan int)

	for _, r := range ranges {
		if isFollowUp {
			go reportInvalidIdsAnyChainLength(r, exitChan)
			// 41823587595 for actual input is a too-high answer
		} else {
			go reportInvalidIds(r, exitChan)
		}
	}

	result := 0

	for range toWait {
		result += <-exitChan
	}

	log.Printf("The sum of all invalid ids is %d", result)
}

func reportInvalidIds(r IdRange, exitChan chan int) {
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
}

func reportInvalidIdsAnyChainLength(r IdRange, exitChan chan int) {
	upperLen := len(strconv.Itoa(r.Upper))
	subExitChan := make(chan []int)
	toWait := 0
	found := map[int]bool{}

	for i := range upperLen / 2 {
		go reportInvalidIdsForTargetChainLength(r, i+1, subExitChan)
		toWait++
	}
	for range toWait {
		for _, invalidId := range <-subExitChan {
			found[invalidId] = true
		}
	}

	result := 0
	for invalidId := range found {
		result += invalidId
	}

	exitChan <- result
}

func reportInvalidIdsForTargetChainLength(r IdRange, target int, exitChan chan []int) {
	result := []int{}

	lowerAsStr := strconv.Itoa(r.Lower)

	upperAsStr := strconv.Itoa(r.Upper)

	lenBounds := IdRange{
		Lower: len(lowerAsStr),
		Upper: len(upperAsStr),
	}

	if lenBounds.Upper >= target*2 {
		// there exist potential candidates (a target chain fits more than once within the upper bound)
		// now, lets start from the lower bound and grow from there
		targetIdLen := int(math.Max(float64(lenBounds.Lower), 2))
		for targetIdLen <= lenBounds.Upper {
			if targetIdLen%target == 0 {
				// a chain with the target len can bit exactly n times within
				// an id with the target length. Let's search for potential
				// chains to be repeated
				// fmt.Println("Range:", r, "Studying invalid IDs of length", targetIdLen, "due to repeated chains of length", target)

				reps := targetIdLen / target

				// initial chain = 100...0 (target digits)
				targetChain := "1"
				maxChain := "9"
				for range target - 1 {
					targetChain += "0"
					maxChain += "9"
				}
				maxChainInt, _ := strconv.Atoi(maxChain)

				for {
					potentialId := strings.Repeat(targetChain, reps)
					potentialIdInt, _ := strconv.Atoi(potentialId)
					if potentialIdInt > r.Upper {
						// fmt.Println("Range:", r, "Potential id", potentialIdInt, "exceeds upper bound. No more invalid ids of length", targetIdLen, "with target chain length", target)
						break
					}
					if r.contains(potentialIdInt) {
						// fmt.Println("Range:", r, "Invalid id", potentialIdInt, "found, chain length", target)
						result = append(result, potentialIdInt)
					}
					targetChainInt, _ := strconv.Atoi(targetChain)
					targetChain = strconv.Itoa(targetChainInt + 1)
					if targetChainInt >= maxChainInt {
						// fmt.Println("Range:", r, "Exhausted all potential target chains of length", target, "for ids of length", targetIdLen)
						break
					}
				}
			}
			targetIdLen++
		}
	} else {
		// fmt.Println("Range:", r, "No potential candidates for target chain length", target)
	}

	exitChan <- result
}

type IdRange struct {
	Lower int
	Upper int
}

func (r *IdRange) contains(id int) bool {
	return id >= r.Lower && id <= r.Upper
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
