package cmd

import (
	"bufio"
	"log"
	"math"
	"os"

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
	isFollowUp, _ := cmd.Flags().GetBool("follow-up")

	banks := readBatteryBanks(inputFile)

	// for i, b := range banks {
	// 	fmt.Printf("Bank %d: %v\n", i, b)
	// }
	commsChan := make(chan int)

	var nBatteries int
	if isFollowUp {
		nBatteries = 12
	} else {
		nBatteries = 2
	}

	for _, bank := range banks {
		go reportMaxBankJoltageNBatteries(bank, nBatteries, commsChan)
	}

	joltage := 0
	for range len(banks) * nBatteries {
		joltage += <-commsChan
	}
	log.Printf("Total joltage using max %d batteries per bank: %d\n", nBatteries, joltage)
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

func reportMaxBankJoltageNBatteries(bank []int, nBatteries int, c chan int) {
	// for a bank with batteries b1, b2, ..., bN, find the max across b0... b(N-nBatteries)
	// For example for a bank consisting of batteries b0,b1,b2,b3 (N=4); and nBatteries=2,
	// the first battery to be activated _must be_ within {b0,b1,b2} so that there exists a second battery to be activated
	// Once the first battery has been found, we should repeat the process to find the rest of batteries.
	// If the first activated battery is bX, we'll call this same function recursively on b(X+1)..bN with nBatteries-1
	// until there are no more batteries left to activate

	relevantBatteries := bank[:len(bank)-nBatteries+1]
	maxIndex := -1
	maxValue := -1
	for i, b := range relevantBatteries {
		if b == 9 {
			maxIndex = i
			maxValue = b
			break
		}
		if b > maxValue {
			maxIndex = i
			maxValue = b
		}
	}
	// fmt.Printf("Selected #%d battery %d @ %d from bank %v\n", 3-nBatteries, maxValue, maxIndex, bank)
	if nBatteries > 1 {
		go reportMaxBankJoltageNBatteries(bank[maxIndex+1:], nBatteries-1, c)
	}
	joltage := maxValue * int(math.Pow10(nBatteries-1))
	c <- joltage
}
