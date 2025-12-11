package cmd

import (
	"bufio"
	"fmt"
	"log"
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
	day1Cmd.Flags().BoolP("follow-up", "f", false, "runs the follow up")
	day1Cmd.Flags().StringP("input-file", "i", "inputs/01", "select file to parse")
}

func day1run(cmd *cobra.Command, args []string) {
	followUp, _ := cmd.Flags().GetBool("follow-up")
	if followUp {
		day1FollowUp(cmd, args)
	} else {
		day1Base(cmd, args)
	}
}

func day1Base(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input-file")
	instructions := readInput(inputFile)
	// for _, instr := range instructions {
	// 	fmt.Println(instr.String())
	// }

	currentPosition := 50
	timesAtZero := 0

	for _, instruction := range instructions {
		if instruction.rotation == Left {
			currentPosition -= instruction.distance
		} else {
			currentPosition += instruction.distance
		}
		// fmt.Printf("Position after applying %v: %d\n", instruction, currentPosition)
		currentPosition = currentPosition % 100
		if currentPosition == 0 {
			timesAtZero++
		}
	}
	fmt.Printf("The password is: %d\n", timesAtZero)
}

func day1FollowUp(cmd *cobra.Command, args []string) {
	fmt.Println("day1 follow up, not implemented")
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

func readInput(filename string) []Instruction {
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
