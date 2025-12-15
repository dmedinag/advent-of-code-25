package cmd

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
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

	log.Info().Msgf("The sum of all invalid ids is %d", result)
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
			log.Debug().Str("range", r.String()).Msgf("found invalid id candidate %d", candidateInt)
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

	if sum == 0 {
		log.Debug().Msgf("no invalid ids found in range %v", r)
	}
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

	logger := log.With().Int("target_chain_length", target).Str("range", r.String()).Logger()

	if lenBounds.Upper >= target*2 {
		// there exist potential candidates (a target chain fits more than once within the upper bound)
		// now, lets start from the lower bound and grow from there
		targetIdLen := int(math.Max(float64(lenBounds.Lower), 2))
		for targetIdLen <= lenBounds.Upper {
			if targetIdLen%target == 0 {
				// a chain with the target len can bit exactly n times within
				// an id with the target length. Let's search for potential
				// chains to be repeated
				logger.Debug().Msgf("Studying invalid IDs of length %d", targetIdLen)

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
						logger.Debug().Msgf("Potential id %d exceeds upper bound. No more invalid ids of length %d with target chain length", potentialIdInt, targetIdLen)
						break
					}
					if r.contains(potentialIdInt) {
						logger.Debug().Msgf("Invalid id %v found", potentialIdInt)
						result = append(result, potentialIdInt)
					}
					targetChainInt, _ := strconv.Atoi(targetChain)
					targetChain = strconv.Itoa(targetChainInt + 1)
					if targetChainInt >= maxChainInt {
						logger.Debug().Msgf("Exhausted all potential target chains of length for ids of length %d", targetIdLen)
						break
					}
				}
			}
			targetIdLen++
		}
	} else {
		logger.Debug().Msgf("No potential candidates for target chain length")
	}

	exitChan <- result
}

type IdRange struct {
	Lower int
	Upper int
}

func (r IdRange) String() string {
	return "[" + strconv.Itoa(r.Lower) + "-" + strconv.Itoa(r.Upper) + "]"
}

func (r *IdRange) contains(id int) bool {
	return id >= r.Lower && id <= r.Upper
}

func IdRangeFromString(s string) IdRange {
	idPair := strings.Split(s, "-")
	if idPair == nil || len(idPair) != 2 {
		log.Fatal().Msgf("there's no id range on %q", s)
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
		log.Fatal().Err(err).Send()
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
		log.Fatal().Err(err).Send()
	}
	return ranges
}
