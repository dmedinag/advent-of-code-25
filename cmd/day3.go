package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// day3Cmd represents the day3 command
var day3Cmd = &cobra.Command{
	Use:   "day3",
	Short: "Broken elevators, maximum joltage",
	Run:   day3Run,
}

func init() {
	rootCmd.AddCommand(day3Cmd)
}

func day3Run(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input-file")
	// isFollowUp, _ := cmd.Flags().GetBool("follow-up")

	banks := readBatteryBanks(inputFile)

	// for i, b := range banks {
	// 	fmt.Printf("Bank %d: %v\n", i, b)
	// }
	commsChan := make(chan int)

	for _, bank := range banks {
		go reportMaxBankJoltageTwoBatteries(bank, commsChan)
	}

	joltage := 0
	for range len(banks) {
		joltage += <-commsChan
	}
	log.Printf("Total joltage: %d\n", joltage)
}

func readBatteryBanks(filename string) [][]int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	batteryBanks := [][]int{}

	for scanner.Scan() {
		batteryBanks = append(batteryBanks, parseBatteryBank(scanner.Bytes()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return batteryBanks
}

func parseBatteryBank(text []byte) []int {
	bank := make([]int, len(text))
	for i, b := range text {
		// ascii numbers start at 48
		bank[i] = int(b) - 48
	}
	return bank
}

func reportMaxBankJoltageTwoBatteries(bank []int, c chan int) {
	maxTens := -1
	maxUnits := -1
	for i, b := range bank {
		if b > maxTens && i < len(bank)-1 {
			maxUnits = -1
			maxTens = b
		} else {
			if b > maxUnits {
				maxUnits = b
			}
		}
	}
	joltage, _ := strconv.Atoi(fmt.Sprintf("%d%d", maxTens, maxUnits))
	// fmt.Printf("Max joltage for bank %v: %v\n", bank, joltage)
	c <- joltage
}
