package cmd

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	DB_NAME = "db.json"
)

type Expense struct {
	CreatedAt   time.Time   `json:"created_at"`
	Description string      `json:"description"`
	Amount      json.Number `json:"Amount"`
	ID          uint64      `json:"id"`
}

func (e Expense) String() string {
	return fmt.Sprintf(
		"%-3d %-11s %-12s %s",
		e.ID, e.CreatedAt.Format("2006-01-02"),
		e.Description, e.Amount,
	)
}

func toJSONNumber(amount float64) json.Number {
	return json.Number(strconv.FormatFloat(amount, 'f', 2, 64))
}

func createExpense(id uint64, desc string, amount float64) Expense {
	return Expense{
		CreatedAt:   time.Now(),
		Description: desc,
		Amount:      toJSONNumber(amount),
		ID:          id,
	}
}

type Expenses struct {
	Expenses map[uint64]Expense `json:"expenses"`
	Limits   [12]float64        `json:"limits"`
}

var expenses *Expenses

func (e *Expenses) getSortedExpenseIDs() []uint64 {
	ids := make([]uint64, 0, len(e.Expenses))
	for k := range e.Expenses {
		ids = append(ids, k)
	}

	slices.Sort(ids)
	return ids
}

const (
	NOT_OVER      = -1.0
	EXACTLY_EQUAL = 0.0
)

func (e *Expenses) getAmountOverLimit(month time.Month) float64 {
	if e.Limits[month] == 0 {
		return NOT_OVER
	}

	total := e.getSummaryByMonth(month)
	if total < e.Limits[month] {
		return NOT_OVER
	}

	if total == e.Limits[month] {
		return EXACTLY_EQUAL
	}

	return total - e.Limits[month]
}

func (e *Expenses) addExpense(desc string, amount float64) uint64 {
	ids := e.getSortedExpenseIDs()

	id := uint64(1)
	lastIdx := len(ids) - 1

	for cidx, cid := range ids {
		if cid > 1 {
			if _, ok := e.Expenses[cid-1]; !ok {
				id = cid - 1
				break
			}
		}

		if cidx == lastIdx || cid+1 != ids[cidx+1] {
			id = cid + 1
			break
		}
	}

	e.Expenses[id] = createExpense(id, desc, amount)
	fmt.Printf("Expense added successfully (ID: %v)\n", id)

	return id
}

func (e *Expenses) deleteExpense(id uint64) error {
	if _, ok := e.Expenses[id]; !ok {
		return fmt.Errorf("No task with ID %v!\n", id)
	}

	delete(e.Expenses, id)
	return nil
}

func (e *Expenses) getSummary() float64 {
	total := 0.0
	for _, expense := range expenses.Expenses {
		v, _ := expense.Amount.Float64()
		total += v
	}

	return total
}

func (e *Expenses) getSummaryByMonth(month time.Month) float64 {
	total := 0.0
	for _, expense := range expenses.Expenses {
		if expense.CreatedAt.Month() != month {
			continue
		}

		v, _ := expense.Amount.Float64()
		total += v
	}

	return total
}

// Loads a saved JSON database, if it exists
func Load(ctx *cli.Context) error {
	expenses = &Expenses{
		Expenses: make(map[uint64]Expense, 0),
	}

	file, err := os.ReadFile(DB_NAME)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err := json.Unmarshal(file, expenses); err != nil {
		return fmt.Errorf("Error unmarshalling JSON data: %v\n", err)
	}

	return nil
}

// Saves all expenses to a JSON file
func Save(ctx *cli.Context) error {
	if expenses == nil {
		return nil
	}

	data, err := json.Marshal(expenses)
	if err != nil {
		return fmt.Errorf("Error marshalling JSON data: %v\n", err)
	}

	if err := os.WriteFile(DB_NAME, data, 0644); err != nil {
		return err
	}

	return nil
}

func getAmount(ctx *cli.Context) (float64, error) {
	amount := ctx.Float64("amount")
	if amount <= 0.0 {
		return 0.0, errors.New("An expense must have a positive, non-zero amount!")
	}

	return amount, nil
}

func HandleAdd(ctx *cli.Context) error {
	description := ctx.String("description")
	if description == "" {
		description = "an expense"
	}

	if !ctx.IsSet("amount") {
		return errors.New("Must provide the amount of the expense!")
	}

	amount := ctx.Float64("amount")
	if amount <= 0.0 {
		return errors.New("An expense must have a positive, non-zero amount!")
	}

	id := expenses.addExpense(description, amount)

	month := expenses.Expenses[id].CreatedAt.Month()

	over := expenses.getAmountOverLimit(month)
	if over == NOT_OVER {
		return nil
	}

	if over == EXACTLY_EQUAL {
		fmt.Println("Warning! This expense puts you exactly at your monthly spending limit!")
	} else {
		fmt.Printf("Warning! This expense puts you over your your monthly spending limit by $%.2f!\n", over)
	}

	return nil
}

func HandleUpdate(ctx *cli.Context) error {
	if !ctx.IsSet("id") {
		return errors.New("Must provide the ID of the expense to update!")
	}

	id := ctx.Uint64("id")
	if expense, ok := expenses.Expenses[id]; ok {
		if ctx.IsSet("description") {
			expense.Description = ctx.String("description")
		}

		if ctx.IsSet("amount") {
			amount, err := getAmount(ctx)
			if err != nil {
				return err
			}

			expense.Amount = toJSONNumber(amount)
		}

		expenses.Expenses[id] = expense

		over := expenses.getAmountOverLimit(expense.CreatedAt.Month())
		if over == NOT_OVER {
			return nil
		}

		if over == EXACTLY_EQUAL {
			fmt.Println("Warning! This expense puts you exactly at your monthly spending limit!")
		} else {
			fmt.Printf("Warning! This expense puts you over your your monthly spending limit by $%.2f!\n", over)
		}

		return nil
	}

	return fmt.Errorf("No expense with ID %v!\n", id)
}

func HandleList(ctx *cli.Context) error {
	fmt.Println("ID  Date        Description  Amount")

	ids := expenses.getSortedExpenseIDs()
	if len(ids) == 0 {
		fmt.Println("There are no expenses to display!")
		return nil
	}

	for _, id := range ids {
		fmt.Println(expenses.Expenses[id].String())
	}

	return nil
}

func HandleSummary(ctx *cli.Context) error {
	if ctx.IsSet("month") {
		month := ctx.Int("month")
		total := expenses.getSummaryByMonth(time.Month(month))
		fmt.Printf("Total expenses for %s: $%.2f\n", time.Month(month), total)
	} else {
		total := expenses.getSummary()
		fmt.Printf("Total expenses: $%.2f\n", total)
	}

	return nil
}

func HandleDelete(ctx *cli.Context) error {
	if !ctx.IsSet("id") {
		return errors.New("Must provide an ID to delete!")
	}

	id := ctx.Uint64("id")
	if id == 0 {
		return fmt.Errorf("ID must be a non-zero, positive number! (is %d)\n", id)
	}

	if err := expenses.deleteExpense(id); err != nil {
		return err
	}

	fmt.Println("Expense deleted successfully")

	return nil
}

func HandleSetLimit(ctx *cli.Context) error {
	if !ctx.IsSet("amount") {
		return errors.New("Must provide a limit amount!")
	}

	amount, err := getAmount(ctx)
	if err != nil {
		return err
	}

	if ctx.IsSet("month") {
		month := ctx.Int("month")
		if month < 1 || month > 12 {
			return fmt.Errorf("'%d' is not a valid month (must be range from 1-12)\n", month)
		}

		expenses.Limits[month-1] = amount
		fmt.Println("Set monthly spending limit successfully")

		return nil
	}

	for month := range 12 {
		expenses.Limits[month] = amount
	}

	fmt.Println("Set spending limit successfully")
	return nil
}

func HandleSave(ctx *cli.Context) error {
	outputFile := ctx.String("output")
	if outputFile == "" {
		outputFile = "output.csv"
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write([]string{"id", "created_at", "description", "amount"})

	ids := expenses.getSortedExpenseIDs()
	for _, id := range ids {
		expense := expenses.Expenses[id]
		writer.Write([]string{
			strconv.FormatUint(expense.ID, 10),
			expense.CreatedAt.Format("2006-01-02"),
			expense.Description,
			expense.Amount.String(),
		})
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}
