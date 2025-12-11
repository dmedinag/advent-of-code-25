package cmd

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// day1Cmd represents the day1 command
var day1Cmd = &cobra.Command{
	Use:   "day1",
	Short: "Get password to enter safe",
	Run:   day1run,
}

func init() {
	rootCmd.AddCommand(day1Cmd)
}

func day1run(cmd *cobra.Command, args []string) {
	followUp, _ := cmd.Flags().GetBool("follow-up")
	inputFile, _ := cmd.Flags().GetString("input-file")
	instructions := readRotations(inputFile)
	// for _, instr := range instructions {
	// 	fmt.Println(instr.String())
	// }

	currentPosition := 50
	password := 0

	for _, instruction := range instructions {
		previousPosition := currentPosition
		if instruction.rotation == Left {
			currentPosition -= instruction.distance
		} else {
			currentPosition += instruction.distance
		}
		// fmt.Printf("Position after applying %v: %d\n", instruction, currentPosition)
		if !followUp {
			if currentPosition%100 == 0 {
				password++
			}
		} else {
			// case 1: from one sign to the opposite
			if previousPosition > 0 && currentPosition <= 0 {
				// fmt.Println("crossed from positive to negative")
				password++
			}
			if previousPosition < 0 && currentPosition >= 0 {
				// fmt.Println("crossed from negative to positive")
				password++
			}
			// case 2: overflowing
			excess := currentPosition / 100
			// if excess != 0 {
			// 	fmt.Printf("overflowed by %d\n", excess)
			// }
			password += int(math.Abs(float64(excess)))
		}

		currentPosition %= 100

	}
	fmt.Printf("The password is: %d\n", password)
}

type Rotation int

const (
	Left Rotation = iota
	Right
)

func (r Rotation) String() string {
	switch r {
	case Left:
		return "Left"
	default:
		return "Right"
	}
}

func RotationFromRune(r rune) Rotation {
	switch r {
	case 'L':
		return Left
	case 'R':
		return Right
	default:
		panic("invalid rotation rune")
	}
}

type Instruction struct {
	rotation Rotation
	distance int
}

func (i *Instruction) String() string {
	return fmt.Sprintf("%v %d", i.rotation, i.distance)
}

func parseInstruction(s string) Instruction {
	rotation := RotationFromRune([]rune(s)[0])
	distance, _ := strconv.Atoi(s[1:])
	return Instruction{
		rotation: rotation,
		distance: distance,
	}
}

func readRotations(filename string) []Instruction {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	instructions := []Instruction{}

	for scanner.Scan() {
		instructions = append(instructions, parseInstruction(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return instructions
}
