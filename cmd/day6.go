package cmd

import (
	"fmt"
	"os"
	"strconv"
	"unicode"

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
	isFollowUp, _ := cmd.Flags().GetBool("follow-up")
	file, err := os.Open(inputFilename)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	fs, err := file.Stat()

	scanner := rscanner.NewScanner(file, fs.Size())

	// get operators
	for scanner.Scan() {
		if scanner.Text() != "" {
			break
		}
	}

	operations := parseOperations(scanner.Text())

	for scanner.Scan() {
		line := scanner.Text()
		parseOperands(operations, line, isFollowUp)
	}

	if scanner.Err() != nil {
		log.Fatal().Err(err).Send()
	}

	// aggregate operations
	result := 0
	for _, op := range operations {
		partialResult := op.Result
		if isFollowUp {
			partialResult = op.GetVerticalResult()
		}
		result += partialResult
	}

	log.Info().Msgf("The result of the cephalopod math is: %d", result)
}

func parseOperations(line string) []*Operation {
	operations := []*Operation{}
	operandSize := 0
	operator := Unknown

	log.Trace().Msgf("Scanning operations from last line %q", line)
	for _, o := range line {
		if o == ' ' {
			operandSize++
			continue
		}
		newOperator, err := ParseOperator(o)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		if operandSize == 0 {
			// special case for first operator
			operandSize++
			operator = newOperator
			continue
		}
		operation := &Operation{Operator: operator, operandSize: operandSize - 1}
		operation.Init()
		log.Trace().Msgf("Created operation: %v", operation)
		log.Trace().Msgf("Parsed operator: %v, computing operand size", newOperator)
		operations = append(operations, operation)
		operandSize = 1
		operator = newOperator
	}

	operation := &Operation{Operator: operator, operandSize: operandSize}
	operation.Init()
	log.Trace().Msgf("Created final operation: %v", operation)
	operations = append(operations, operation)

	return operations
}

func parseOperands(operations []*Operation, line string, isFollowUp bool) {
	log.Trace().Msgf("Scanning operands from line %v", line)
	cursor := 0
	for _, op := range operations {
		numStr := line[cursor : cursor+op.operandSize]
		cursor += op.operandSize + 1
		log.Trace().Msgf("Parsed operand: %q for operation %v", numStr, op)
		if isFollowUp {
			op.OperateVertical(numStr)
		} else {
			num, _ := strconv.Atoi(numStr)
			op.Operate(num)
		}
	}
}

type Operation struct {
	Operator    Operator
	Result      int
	operandSize int
	cache       [][]rune
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

func (o *Operation) OperateVertical(operand string) {
	parsedOperand := []rune{}
	for _, digit := range operand {
		parsedOperand = append(parsedOperand, digit)
	}

	o.cache = append(o.cache, parsedOperand)
}

func (o *Operation) GetVerticalResult() int {
	operands := []int{}
	// operandsFromCache
	for i := range o.operandSize {
		operand := ""
		for _, row := range o.cache {
			if unicode.IsNumber(row[i]) {
				operand = string(row[i]) + operand
			}
		}
		operandInt, err := strconv.Atoi(operand)
		if err != nil {
			log.Panic().Err(err).Msgf("failed to parse operand %q", operand)
		}
		log.Debug().Str("operation", o.String()).Msgf("Found operand %d", operandInt)
		operands = append(operands, operandInt)
	}
	for _, operand := range operands {
		o.Result = o.Operator.Apply(o.Result, operand)
	}
	return o.Result
}

func (o Operation) String() string {
	return fmt.Sprintf("Operation {%v, size %d}", o.Operator, o.operandSize)
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
