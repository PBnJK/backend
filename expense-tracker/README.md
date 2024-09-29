# Expense Tracker
CLI expenses-tracking tool for the
[roadmap.sh expense tracker](https://roadmap.sh/projects/expense-tracker)
project.

## Running
### 1. Clone
```bash
git clone https://github.com/pbnjk/backend.git
```
### 2. Build
```bash
cd backend/expense-tracker
go build
```
### 3. Run!
```bash
# Adds an expense
./expense-tracker add --description "Lunch" --amount 20
# OUTPUT: Expense added successfully (ID: 1)

# Updates an expense
./expense-tracker update --id 1 --amount 30

# Delete expenses by ID
./expense-tracker delete --id 1

# List expenses
./expense-tracker list

# Print a summary of the expenses
./expense-tracker summary

# Set a monthly spending limit!
./expense-tracker set-limit --amount 1000 --month 8

# ...or yearly (omit month)
./expense-tracker set-limit --amount 1000
```

## DB Format
Expenses are stored in a .json file called "db.json". The format of the JSON
structure is as follows:
```json
{
	"expenses": {
		"1": {
			"created_at": "2024-09-28T20:01:33.427304798-03:00",
			"description": "Lunch",
			"amount":20.00,
			"id": 1,
		},
		"2": {
			"createdAt": "2024-09-28T20:03:15.509254764-03:00",
			"description": "Dinner",
			"amount":30.00,
			"id": 2,
		}
	}
	"limits":[
		1500,
		10000,
		7900,
		...
	]
}
```
