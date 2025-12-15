package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// day5Cmd represents the day5 command
var day5Cmd = &cobra.Command{
	Use:   "day5",
	Short: "fresh products in catalog",
	Run:   runDay5,
}

func init() {
	rootCmd.AddCommand(day5Cmd)
}

func runDay5(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input-file")
	isFollowUp, _ := cmd.Flags().GetBool("follow-up")

	if isFollowUp {
		log.Fatal().Msg("Follow-up not implemented yet for day 5")
	}

	intervals, products := parseInput(inputFile)

	staleProducts := findStaleProducts(intervals, products)

	log.Info().Msgf("Found %d stale products", len(staleProducts))
	log.Info().Msgf("So there are %d fresh products", len(products)-len(staleProducts))
}

func parseInput(inputFilename string) ([]*Interval, []int) {
	file, err := os.Open(inputFilename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer file.Close()

	intervals := []*Interval{}
	products := []int{}

	scanner := bufio.NewScanner(file)
	readingProducts := false
	merged := 0

out:
	for scanner.Scan() {
		if readingProducts {
			product, _ := strconv.Atoi(scanner.Text())
			products = append(products, product)
			log.Trace().Msgf("Registered product %d", product)
			continue
		}
		if scanner.Text() == "" {
			readingProducts = true
			continue
		}
		newInterval := ParseInterval(scanner.Text())
		for _, interval := range intervals {
			if interval.Merge(newInterval) {
				merged++
				log.Trace().Msgf("Merged interval %v to form %v", newInterval, interval)
				continue out
			}
		}
		log.Trace().Msgf("Registered new interval %v", newInterval)
		intervals = append(intervals, &newInterval)
	}

	log.Debug().Msgf("Registered %d products, and %d intervals of fresh products (merged into %v)", len(products), len(intervals)+merged, len(intervals))

	return intervals, products
}

func findStaleProducts(intervals []*Interval, products []int) []int {
	staleProducts := []int{}
out:
	for _, p := range products {
		for _, i := range intervals {
			if i.Contains(p) {
				log.Trace().Msgf("Product %d\tis fresh (in interval %v)", p, i)
				continue out
			}
		}
		staleProducts = append(staleProducts, p)
	}
	return staleProducts
}

type Interval struct {
	Lower int
	Upper int
}

func ParseInterval(input string) Interval {
	bounds := strings.Split(input, "-")
	if len(bounds) != 2 {
		log.Fatal().Msgf("Invalid interval format: %s", input)
	}
	lower, _ := strconv.Atoi(bounds[0])
	upper, _ := strconv.Atoi(bounds[1])
	return Interval{
		Lower: lower,
		Upper: upper,
	}
}

func (i Interval) Contains(target int) bool {
	return target >= i.Lower && target <= i.Upper
}

func (i Interval) String() string {
	return fmt.Sprintf("(%d, %d)", i.Lower, i.Upper)
}

func (i *Interval) Merge(other Interval) bool {
	extended := false
	if other.Upper > i.Upper && other.Lower <= i.Upper {
		// other Interval _extends_ the current to the right
		// 3-5 merge 4-6 = 3-6
		// 3-5 merge 5-6 = 3-6
		// 3-5 merge 1-6 = 3-6 <- this branch will cover only the upper part
		log.Trace().Msgf("Extending %v to the right:\n\tinput %v\n\tnew interval: (%d, %d)", *i, other, i.Lower, other.Upper)
		(*i).Upper = other.Upper
		extended = true
	}
	if other.Lower < i.Lower && other.Upper >= i.Lower {
		// same as the upper branch, but for the lower bound
		log.Trace().Msgf("Extending %v to the left:\n\tgiven %v\n\tnew interval: (%d, %d)", *i, other, other.Lower, i.Upper)
		(*i).Lower = other.Lower
		extended = true
	}
	return extended
}
