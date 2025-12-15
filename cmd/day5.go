package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
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

	intervals, products := parseInput(inputFile)

	if isFollowUp {
		freshCount := 0
		for _, i := range intervals.ToSlice() {
			moreProducts := i.Upper - i.Lower + 1
			freshCount += moreProducts
			log.Debug().Msgf("Fresh products due to %v: %d", i, moreProducts)
		}
		log.Info().Msgf("There are %d different fresh products", freshCount)
	} else {
		staleProducts := findStaleProducts(intervals, products)

		log.Info().Msgf("There are %d fresh products", len(products)-len(staleProducts))
	}
}

func parseInput(inputFilename string) (mapset.Set[*Interval], []int) {
	file, err := os.Open(inputFilename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer file.Close()

	intervals := mapset.NewSet[*Interval]()
	products := []int{}

	scanner := bufio.NewScanner(file)
	readingProducts := false
	merged := 0

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
		log.Trace().Msgf("Registered new interval %v", newInterval)
		intervals.Add(&newInterval)
	}

	intervals, newMerged := compactIntervals(intervals)
	merged += newMerged

	log.Debug().Msgf("Registered %d products, and %d intervals of fresh products (merged into %v)", len(products), intervals.Cardinality()+merged, intervals.Cardinality())

	return intervals, products
}

func compactIntervals(intervals mapset.Set[*Interval]) (mapset.Set[*Interval], int) {
	items := intervals.ToSlice()

	if len(items) <= 1 {
		return intervals, 0
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Lower < items[j].Lower
	})

	mergedCount := 0
	result := mapset.NewSet[*Interval]()

	current := items[0]
	result.Add(current)

	for i := 1; i < len(items); i++ {
		next := items[i]

		if current.Merge(*next) {
			mergedCount++
		} else {
			current = next
			result.Add(current)
		}
	}

	return result, mergedCount
}

func findStaleProducts(intervals mapset.Set[*Interval], products []int) []int {
	staleProducts := []int{}
out:
	for _, p := range products {
		for _, i := range intervals.ToSlice() {
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
	if other.Upper > i.Upper && other.Lower <= i.Upper+1 {
		// other Interval _extends_ the current to the right
		// 3-5 merge 4-6 = 3-6
		// 3-5 merge 5-6 = 3-6
		// 3-5 merge 1-6 = 3-6 <- this branch will cover only the upper part
		log.Trace().Msgf("Extending %v to the right:\n\tinput %v\n\tnew interval: (%d, %d)", *i, other, i.Lower, other.Upper)
		(*i).Upper = other.Upper
		extended = true
	}
	if other.Lower < i.Lower && other.Upper >= i.Lower-1 {
		// same as the upper branch, but for the lower bound
		log.Trace().Msgf("Extending %v to the left:\n\tgiven %v\n\tnew interval: (%d, %d)", *i, other, other.Lower, i.Upper)
		(*i).Lower = other.Lower
		extended = true
	}
	if other.Lower >= i.Lower && other.Upper <= i.Upper {
		// already included
		return true
	}
	return extended
}
