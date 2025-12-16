package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/vallerion/rscanner"
)

var day6Cmd = &cobra.Command{
	Use:   "day6",
	Short: "cephalopod math",
	Run:   runDay6,
}

func init() {
	rootCmd.AddCommand(day6Cmd)
}

func runDay6(cmd *cobra.Command, args []string) {
	inputFilename, _ := cmd.Flags().GetString("input-file")
	file, err := os.Open(inputFilename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	fs, err := file.Stat()

	scanner := rscanner.NewScanner(file, fs.Size())

	operations := []*Operation{}

	// get operators
	for scanner.Scan() {
		if scanner.Text() != "" {
			break
		}
	}
	log.Trace().Msgf("Scanning operations from last line %q", scanner.Text())
	for _, o := range scanner.Text() {
		if o == ' ' {
			continue
		}
		operator, err := ParseOperator(o)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		operation := &Operation{Operator: operator}
		operation.Init()
		log.Trace().Msgf("Parsed operator: %v", operator)
		operations = append(operations, operation)
	}

	// compute results
	for scanner.Scan() {
		log.Trace().Msgf("Scanning operands from line %v", scanner.Text())
		for i, strNum := range strings.Fields(scanner.Text()) {
			num, _ := strconv.Atoi(strNum)
			operation := operations[i]
			log.Trace().Msgf("Parsed operand: %v for operation %v", num, operation)
			operation.Operate(num)
		}
	}

	if scanner.Err() != nil {
		log.Fatal().Err(err).Send()
	}

	// aggregate operations
	result := 0
	for _, op := range operations {
		result += op.Result
	}

	log.Info().Msgf("The result of the cephalopod math is: %d", result)
}

type Operation struct {
	Operator    Operator
	Result      int
	initialized bool
}

type Operator int

const (
	Unknown Operator = iota
	Multiply
	Sum
)

func (o *Operation) Operate(operand int) {
	o.Result = o.Operator.Apply(o.Result, operand)
	log.Trace().Msgf("New value: %v", o.Result)
}

func (o Operation) String() string {
	return fmt.Sprintf("Operation %v, current accummulation %v", o.Operator, o.Result)
}

func (o *Operation) Init() {
	if o.initialized {
		return
	}
	if o.Operator == Multiply {
		o.Result = 1
	}
	o.initialized = true
}

func ParseOperator(t rune) (Operator, error) {
	switch t {
	case '+':
		return Sum, nil
	case '*':
		return Multiply, nil
	}
	return Unknown, fmt.Errorf("unknown operator: %v", t)
}

func (o Operator) String() string {
	switch o {
	case Multiply:
		return "*"
	case Sum:
		return "+"
	default:
		return "?"
	}
}

func (o Operator) Apply(a, b int) int {
	switch o {
	case Multiply:
		return a * b
	case Sum:
		return a + b
	}
	log.Fatal().Msgf("unknown operator: %v", o)
	panic("")
}
