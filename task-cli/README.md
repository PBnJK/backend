# Task CLI
CLI task-tracking tool for the
[roadmap.sh task tracker](https://roadmap.sh/projects/task-tracker) project.

## Running
### 1. Clone
```bash
git clone https://github.com/pbnjk/backend.git
```
### 2. Build
```bash
cd backend/task-cli
go build
```
### 3. Run!
```bash
# Adds a task
./task-cli add "A very hard task"
# OUTPUT: Task added successfully! ID: 1

# Updates a task
./task-cli update 1 "Actually, not very hard at all"

# Marks a task as...
./task-cli mark-in-progress 1 # In progress
./task-cli mark-done 1 # Done

# Delete tasks by ID
./task-cli delete 1

# ...or by status
./task-cli delete todo
./task-cli delete in-progress
./task-cli delete done

# List tasks
./task-cli list

# Listing by status
./task-cli list todo
./task-cli list in-progress
./task-cli list done

# You can also do verbose printing (adds date of creation/updating)
./task-cli list --verbose
```

## DB Format
Tasks are stored in a .json file called "db.json". The format of the JSON
structure is as follows:
```json
{
	"tasks": {
		"1": {
			"createdAt": "2024-09-28T20:01:33.427304798-03:00",
			"updatedAt": "2024-09-28T20:01:33.427304876-03:00",
			"desc": "wow cool task",
			"id": 1,
			"status": 0
		},
		"2": {
			"createdAt": "2024-09-28T20:03:15.509254764-03:00",
			"updatedAt": "2024-09-28T20:03:30.999780767-03:00",
			"desc": "another task wow",
			"id": 2,
			"status": 1
		}
	}
}
```

## Some miscellaneous notes:
- Tasks are indexed by ID. The tool tries to use free IDs when they go out of
use, so indexes should be in a sequential order;
- Statuses are represented as a numeric ID from 0 to 2 ("todo", "in-progress"
and "done" respectively);
- Descriptions *can* be arbitrarily long, but the list command pads them to 48
characters by default, so they'll look ugly if larger than that.
