package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

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
	isFollowUp, _ := cmd.Flags().GetBool("follow-up")

	if isFollowUp {
		rollMap := readMap(inputFile)
		registerNeighbors(&(rollMap.Rolls))
		accessibleRolls := findAccessibleRolls(&rollMap)
		log.Printf("There are %d accessible rolls in the map", len(accessibleRolls))
	} else {
		base(inputFile)
	}
}

func base(inputFile string) {
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

func readMap(inputFile string) RollMap {
	result := map[Coordinates]*Roll{}

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	c := make(chan []Roll)
	rowCount := 0
	var colCount int

	for scanner.Scan() {
		colCount = len(scanner.Text())
		go parseRollsFromRow(rowCount, scanner.Text(), c)
		rowCount++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// log.Println("Found rolls on the following coordinates:")
	for range rowCount {
		for _, r := range <-c {
			result[r.Position] = &r
			// log.Printf("(%d, %d)\n", r.Position.Row, r.Position.Col)
		}
	}

	return RollMap{
		Rows:  rowCount,
		Cols:  colCount,
		Rolls: result,
	}
}

func parseRollsFromRow(row int, rawElements string, c chan []Roll) {
	rolls := []Roll{}
	for col, r := range []rune(rawElements) {
		if r == '@' {
			rolls = append(rolls, Roll{
				Position: Coordinates{Row: row, Col: col},
			})
		}
	}
	c <- rolls
}

func registerNeighbors(rolls *map[Coordinates]*Roll) {
	for pos, r := range *rolls {
		for _, target := range pos.Adjacent() {
			if other, found := (*rolls)[target]; found {
				// log.Printf("%v and %v found as neighbors\n", pos, target)
				other.registerNeighbor(r)
				r.registerNeighbor(other)
			}
		}
	}
}

func findAccessibleRolls(rollMap *RollMap) []*Roll {
	accessibleRollsSet := map[Coordinates]*Roll{}
	candidates := []*Roll{}

	for _, r := range (*rollMap).Rolls {
		candidates = append(candidates, r)
	}

	log.Printf("Got %v candidates to assess, map:\n%v\n", len(candidates), *rollMap)

	for len(candidates) > 0 {
		nextCandidates := map[Coordinates]*Roll{}
		for _, roll := range candidates {
			// log.Printf("Assessing roll at %v, with neighbors:\n", roll.Position)
			// for _, n := range roll.neighbors {
			// 	if !n.IsRemoved() {
			// 		log.Printf(" - %v\n", n.Position)
			// 	}
			// }
			if roll.tryRemove() {
				// log.Printf("Roll at %v is accessible, map:\n%v\n", roll.Position, *rollMap)
				accessibleRollsSet[roll.Position] = roll
				for c, r := range roll.neighbors {
					if !r.IsRemoved() {
						// log.Printf("Registering neighbor at %v as candidate for next round", c)
						nextCandidates[c] = r
						// } else {
						// log.Printf("Neighbor at %v is already removed, skipping", c)
					}
				}
				// } else {
				// log.Printf("Roll at %v is not accessible", roll.Position)
			}
		}
		candidates = []*Roll{}
		for _, r := range nextCandidates {
			candidates = append(candidates, r)
		}
	}
	accessibleRolls := []*Roll{}
	for _, r := range accessibleRollsSet {
		accessibleRolls = append(accessibleRolls, r)
	}
	return accessibleRolls
}

type Roll struct {
	Position  Coordinates
	neighbors map[Coordinates]*Roll
	removed   bool
}

func (r Roll) String() string {
	if r.neighbors == nil {
		r.neighbors = make(map[Coordinates]*Roll)
	}
	return fmt.Sprintf("%v roll@%v with %d neighbors", r.removed, r.Position, len(r.neighbors))
}

func (r Roll) IsRemoved() bool {
	return r.removed
}

func (r *Roll) tryRemove() bool {
	if r.removed {
		// fail fast if this roll was already detected as accessible
		return true
	}
	if len(r.neighbors) < 4 {
		r.removed = true
		return true
	}
	removedNeighbors := 0
	for _, n := range r.neighbors {
		if n.IsRemoved() {
			removedNeighbors++
			if (len(r.neighbors) - removedNeighbors) < 4 {
				r.removed = true
				return true
			}
		}
	}
	return false
}

func (r *Roll) registerNeighbor(other *Roll) {
	if r.neighbors == nil {
		r.neighbors = make(map[Coordinates]*Roll, 8)
		return
	}
	r.neighbors[other.Position] = other
}

type Coordinates struct {
	Row int
	Col int
}

func (c Coordinates) String() string {
	return fmt.Sprintf("(%d, %d)", c.Row, c.Col)
}

func (c Coordinates) Adjacent() []Coordinates {
	result := []Coordinates{
		{Row: c.Row + 1, Col: c.Col + 1},
		{Row: c.Row, Col: c.Col + 1},
		{Row: c.Row + 1, Col: c.Col},
	}
	if c.Row > 0 {
		result = append(result, Coordinates{Row: c.Row - 1, Col: c.Col})
		result = append(result, Coordinates{Row: c.Row - 1, Col: c.Col + 1})
	}
	if c.Col > 0 {
		result = append(result, Coordinates{Row: c.Row, Col: c.Col - 1})
		result = append(result, Coordinates{Row: c.Row + 1, Col: c.Col - 1})
	}
	if c.Row > 0 && c.Col > 0 {
		result = append(result, Coordinates{Row: c.Row - 1, Col: c.Col - 1})
	}
	return result
}

type RollMap struct {
	Rows  int
	Cols  int
	Rolls map[Coordinates]*Roll
}

func (m RollMap) String() string {
	s := "  "
	for y := range m.Cols {
		s += strconv.Itoa(y)
	}
	s += "\n"
	for x := range m.Rows {
		s += strconv.Itoa(x) + " "
		for y := range m.Cols {
			if r, found := m.Rolls[Coordinates{Row: x, Col: y}]; found {
				if r.removed {
					s += "x"
				} else {
					s += "@"
				}
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
}
