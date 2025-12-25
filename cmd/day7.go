package cmd

import (
	"bufio"
	"errors"
	"os"
	"sync"
	"sync/atomic"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// day7Cmd represents the day7 command
var day7Cmd = &cobra.Command{
	Use:   "day7",
	Short: "Tachyon beams split",
	Run:   runDay7,
}

func init() {
	rootCmd.AddCommand(day7Cmd)
}

func runDay7(cmd *cobra.Command, args []string) {
	inputFilename, _ := cmd.Flags().GetString("input-file")
	isFollowUp, _ := cmd.Flags().GetBool("follow-up")

	startCol, rowCount, splits := parseTachyonInput(inputFilename)
	log.Trace().Msgf("Start column: %d, row count: %d", startCol, rowCount)
	log.Trace().Msgf("All splits %v", &splits)

	if isFollowUp {
		result := countPaths(splits, startCol, rowCount)
		log.Info().Msgf("There are %d possible paths for the particle", result)
	} else {
		result := traceRays(splits, startCol, rowCount)
		log.Info().Msgf("Split %d times", result)
	}
}

func traceRays(splits *sync.Map, startCol, rowCount int) int {
	rays := mapset.NewSet[int]()
	rays.Add(startCol)
	var splitCount atomic.Int32
	for row := range rowCount {
		nextRays := mapset.NewSet[int]()
		for ray := range rays.Iterator().C {
			targetPosition := position{row: row, col: ray}
			_, exists := splits.Load(targetPosition)
			if exists {
				log.Debug().Msgf("Ray split on row %d at col %d", row, ray)
				splitCount.Add(1)
				nextRays.Add(ray - 1)
				nextRays.Add(ray + 1)
			} else {
				log.Debug().Msgf("Ray on row %d at col %d carries on down", row, ray)
				nextRays.Add(ray)
			}
		}
		rays = nextRays
	}
	return int(splitCount.Load())
}

func countPaths(splits *sync.Map, startCol, rowCount int) int {
	paths := mapset.NewSet[*path]()
	paths.Add(&path{ray: startCol, origins: 1})
	for row := range rowCount {
		nextRays := map[int]int{}
		for p := range paths.Iterator().C {
			targetPosition := position{row: row, col: p.ray}
			_, exists := splits.Load(targetPosition)
			if exists {
				addOrIncrease(nextRays, p.ray-1, p.origins)
				addOrIncrease(nextRays, p.ray+1, p.origins)
			} else {
				addOrIncrease(nextRays, p.ray, p.origins)
			}
		}
		paths = mapset.NewSet[*path]()
		for k, v := range nextRays {
			paths.Add(&path{ray: k, origins: v})
		}
	}
	splitCount := 0
	for p := range paths.Iterator().C {
		splitCount += p.origins
	}
	return splitCount
}

func addOrIncrease(m map[int]int, key, delta int) {
	if prev, exists := m[key]; !exists {
		m[key] = delta
	} else {
		m[key] = prev + delta
	}
}

type path struct {
	ray     int
	origins int
}

func parseTachyonInput(inputFilename string) (startCol int, rowCount int, splits *sync.Map) {
	file, err := os.Open(inputFilename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	scanner := bufio.NewScanner(file)
	splits = &sync.Map{}

	for scanner.Scan() {
		if rowCount == 0 {
			startCol, _ = findStart(scanner.Text())
		} else {
			lineSplits, err := parseSplits(scanner.Text())
			if err != nil {
				continue
			}
			for _, pos := range lineSplits {
				log.Debug().Msgf("Split at row %d, col %d", rowCount, pos)
				splits.Store(position{row: rowCount, col: pos}, false)
			}
		}
		rowCount++
	}

	return
}

func findStart(line string) (int, error) {
	for i, r := range line {
		if r == 'S' {
			return i, nil
		}
	}
	return 0, errors.New("starting position not found")
}

func parseSplits(line string) ([]int, error) {
	positions := []int{}
	for i, r := range line {
		if r == '^' {
			positions = append(positions, i)
		}
	}
	if len(positions) == 0 {
		return nil, errors.New("no splits found in line")
	}
	return positions, nil
}

type position struct {
	row int
	col int
}
