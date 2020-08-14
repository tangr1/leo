package command

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jung-kurt/gofpdf"
	"github.com/spf13/cobra"
)

func CalculationCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cal [页数] [难度]",
		Short: "生成一篇新的计算小超市",
		Long:  "生成一篇新的计算小超市",
		Args:  cobra.MinimumNArgs(0),
		Run:   calculation,
	}
}

func calculation(cmd *cobra.Command, args []string) {
	rand.Seed(time.Now().UnixNano())
	pages := 1
	if len(args) >= 1 {
		var err error
		if pages, err = strconv.Atoi(args[0]); err != nil {
			color.Red("参数错误，页数必须要是正整数")
			return
		}
	}
	functions := make([]func(int) (int, string, int), 0)
	hard := 1
	if len(args) > 1 {
		var err error
		if hard, err = strconv.Atoi(args[1]); err != nil {
			color.Red("参数错误，难度只能是1,2,3,4三种选择")
			return
		}
	}
	functions = append(functions, oneOnePlusNoCarry)
	functions = append(functions, oneOnePlusCarry)
	functions = append(functions, TwoOnePlusNoCarry)
	functions = append(functions, TwoOnePlusCarry)
	functions = append(functions, TwoTwoPlusNoCarry)
	functions = append(functions, TwoTwoPlusCarry)
	functions = append(functions, oneOneSubstractNoCarry)
	functions = append(functions, TwoOneSubstractNoCarry)
	functions = append(functions, TwoOneSubstractCarry)
	functions = append(functions, TwoOneSubstractCarry)
	functions = append(functions, TwoTwoSubstractNoCarry)
	functions = append(functions, TwoTwoSubstractCarry)
	max := hard * 10
	pdf := gofpdf.New("P", "mm", "A4", "")
	for i := 0; i < pages; i++ {
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 24)
		pdf.Cell(180, 15, "           Calculation Supermarket by Dad")
		pdf.Ln(-1)
		pdf.SetFont("Arial", "", 20)
		pdf.Cell(50, 15, "Name:")
		pdf.Cell(50, 15, "Date:")
		pdf.Cell(50, 15, "Time:")
		pdf.Cell(50, 15, "Score:")
		pdf.Ln(-1)
		pdf.SetFont("Arial", "B", 18)
		pdf.Cell(65, 10, "Part I")
		pdf.Ln(-1)
		generateProblem(pdf, max, functions)
		pdf.Ln(-1)
		pdf.SetFont("Arial", "B", 18)
		pdf.Cell(65, 10, "Part II")
		pdf.Ln(-1)
		generateProblem(pdf, max, functions)
	}
	pdf.OutputFileAndClose("calculation.pdf")
	color.Green("成功生成计算小超市，请在工具同一目录寻找calculation.pdf")
}

func generateProblem(pdf *gofpdf.Fpdf, max int, functions []func(int) (int, string, int)) {
	pdf.SetFont("Courier", "B", 14)
	for i := 0; i < 30; i++ {
		problem(i+1, max, functions[rand.Intn(len(functions))], pdf)
	}
}

func fillColumn(count int, text string, start int) string {
	column := make([]string, 0)
	for i := 0; i < count; i++ {
		column = append(column, " ")
	}
	for i := 0; i < count; i++ {
		if i >= len(text)+start {
			break
		}
		if i >= start {
			column[i] = string(text[i-start])
		}
	}
	return strings.Join(column[:], "")
}

func problem(index int, max int, calculate func(int) (int, string, int), pdf *gofpdf.Fpdf) {
	first, symbol, second := calculate(max)
	firstStr := fillColumn(4, fmt.Sprintf("%d", first), 1)
	indexStr := fillColumn(4, fmt.Sprintf("(%d)", index), 0)
	secondStr := fillColumn(4, fmt.Sprintf("%d", second), 1)
	pdf.Cell(65, 10, fmt.Sprintf("%s%s%s%s=", indexStr, firstStr, symbol, secondStr))
	if index%3 == 0 {
		pdf.Ln(-1)
	}
}

func oneOnePlusNoCarry(max int) (int, string, int) {
	for {
		first := rand.Intn(9) + 1
		second := rand.Intn(9) + 1
		if first+second < 10 {
			return first, "+", second
		}
	}
}

func oneOnePlusCarry(max int) (int, string, int) {
	for {
		first := rand.Intn(9) + 1
		second := rand.Intn(9) + 1
		if first+second >= 10 {
			return first, "+", second
		}
	}
}

func TwoOnePlusNoCarry(max int) (int, string, int) {
	for {
		first := rand.Intn(max-1) + 11
		second := rand.Intn(9) + 1
		if first%10+second < 10 {
			return first, "+", second
		}
	}
}

func TwoOnePlusCarry(max int) (int, string, int) {
	for {
		first := rand.Intn(max-1) + 11
		second := rand.Intn(9) + 1
		if first%10+second >= 10 {
			return first, "+", second
		}
	}
}

func TwoTwoPlusNoCarry(max int) (int, string, int) {
	for {
		first := rand.Intn(max-1) + 11
		second := rand.Intn(max-1) + 11
		if first%10+second%10 < 10 {
			return first, "+", second
		}
	}
}

func TwoTwoPlusCarry(max int) (int, string, int) {
	for {
		first := rand.Intn(max-1) + 11
		second := rand.Intn(max-1) + 11
		if first%10+second >= 10 {
			return first, "+", second
		}
	}
}

func oneOneSubstractNoCarry(max int) (int, string, int) {
	for {
		first := rand.Intn(9) + 1
		second := rand.Intn(9) + 1
		if first-second >= 0 {
			return first, "-", second
		}
	}
}

func TwoOneSubstractNoCarry(max int) (int, string, int) {
	for {
		first := rand.Intn(max-1) + 11
		second := rand.Intn(9) + 1
		if first%10-second >= 0 {
			return first, "-", second
		}
	}
}

func TwoOneSubstractCarry(max int) (int, string, int) {
	for {
		first := rand.Intn(max-1) + 11
		second := rand.Intn(9) + 1
		if first%10-second < 0 {
			return first, "-", second
		}
	}
}

func TwoTwoSubstractNoCarry(max int) (int, string, int) {
	for {
		first := rand.Intn(max-1) + 11
		second := rand.Intn(max-1) + 11
		if first > second && first%10 >= second%10 {
			return first, "-", second
		}
	}
}

func TwoTwoSubstractCarry(max int) (int, string, int) {
	for {
		first := rand.Intn(max-1) + 11
		second := rand.Intn(max-1) + 11
		if max <= 10 || (first > second && first%10 < second%10) {
			return first, "+", second
		}
	}
}
