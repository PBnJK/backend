package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

type TaskStatus int

const (
	STATUS_TODO        TaskStatus = 0
	STATUS_IN_PROGRESS TaskStatus = 1
	STATUS_DONE        TaskStatus = 2
)

func (s TaskStatus) String() string {
	switch s {
	case STATUS_TODO:
		return "To-do"
	case STATUS_IN_PROGRESS:
		return "In Progress"
	case STATUS_DONE:
		return "Done"
	default:
		return "???"
	}
}

const (
	DB_NAME = "db.json"
)

type Task struct {
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Description string     `json:"desc"`
	Id          uint64     `json:"id"`
	Status      TaskStatus `json:"status"`
}

func createTask(id uint64, desc string) Task {
	return Task{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Description: desc,
		Id:          id,
		Status:      STATUS_TODO,
	}
}

func (t Task) String(verbose bool) string {
	if verbose {
		return fmt.Sprintf(
			"%-4d %-48s %-12s %-20s %-20s", t.Id, t.Description, t.Status.String(),
			t.CreatedAt.Format("2006-01-02 15:04:05"), t.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	} else {
		return fmt.Sprintf("%-4d %-48s %s", t.Id, t.Description, t.Status.String())
	}
}

type Tasks struct {
	Tasks map[uint64]Task `json:"tasks"`
}

var tasks *Tasks = nil

func getSortedTasksIDs() []uint64 {
	ids := make([]uint64, 0, len(tasks.Tasks))
	for k := range tasks.Tasks {
		ids = append(ids, k)
	}

	slices.Sort(ids)
	return ids
}

func (t *Tasks) addTask(desc string) {
	ids := getSortedTasksIDs()

	id := uint64(1)
	lastIdx := len(ids) - 1

	for cidx, cid := range ids {
		if cid > 1 {
			if _, ok := t.Tasks[cid-1]; !ok {
				id = cid - 1
				break
			}
		}

		if cidx == lastIdx || cid+1 != ids[cidx+1] {
			id = cid + 1
			break
		}
	}

	tasks.Tasks[id] = createTask(id, desc)
	fmt.Printf("Task added successfully! ID: %v\n", id)
}

func (t *Tasks) updateTask(id uint64, newDesc string) error {
	if task, ok := t.Tasks[id]; ok {
		task.Description = newDesc
		task.UpdatedAt = time.Now()

		t.Tasks[id] = task

		return nil
	}

	return fmt.Errorf("No task with ID %v!\n", id)
}

func (t *Tasks) deleteTasksByID(id uint64) error {
	if _, ok := t.Tasks[id]; !ok {
		return fmt.Errorf("No task with ID %v!\n", id)
	}

	delete(t.Tasks, id)
	return nil
}

func (t *Tasks) deleteTasksByStatus(status TaskStatus) {
	for k, v := range t.Tasks {
		if v.Status == status {
			delete(t.Tasks, k)
		}
	}
}

func (t *Tasks) markTaskAs(id uint64, status TaskStatus) error {
	if task, ok := t.Tasks[id]; ok {
		task.Status = status
		task.UpdatedAt = time.Now()

		t.Tasks[id] = task

		return nil
	}

	return fmt.Errorf("No task with ID %v!\n", id)
}

func (t *Tasks) listByStatus(verbose bool, status TaskStatus) {
	printListHeader(verbose)

	hasTask := false

	ids := getSortedTasksIDs()
	for _, id := range ids {
		task := t.Tasks[id]

		if task.Status == status {
			fmt.Println(task.String(verbose))
			hasTask = true
		}
	}

	if !hasTask {
		fmt.Println("There are no tasks to display!")
	}
}

// Loads a saved JSON database, if it exists
func Load(ctx *cli.Context) error {
	tasks = &Tasks{
		Tasks: make(map[uint64]Task, 0),
	}

	file, err := os.ReadFile(DB_NAME)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err := json.Unmarshal(file, tasks); err != nil {
		return fmt.Errorf("Error unmarshalling JSON data: %v\n", err)
	}

	return nil
}

// Saves all tasks to a JSON file
func Save(ctx *cli.Context) error {
	if tasks == nil {
		return nil
	}

	data, err := json.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("Error marshalling JSON data: %v\n", err)
	}

	if err := os.WriteFile(DB_NAME, data, 0644); err != nil {
		return err
	}

	return nil
}

func getIdFromString(id string) (uint64, error) {
	if id == "" {
		return 0, errors.New("Must provide task ID")
	}

	idNum, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("Couldn't convert ID '%v' to number!", id)
	}

	return uint64(idNum), nil
}

func HandleAdd(ctx *cli.Context) error {
	if ctx.Args().Get(0) == "" {
		return errors.New("Must provide a task description")
	}

	tasks.addTask(ctx.Args().Get(0))
	return nil
}

func HandleUpdate(ctx *cli.Context) error {
	desc := ctx.Args().Get(1)
	if desc == "" {
		return errors.New("Must provide task ID and updated description")
	}

	id, err := getIdFromString(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	if err := tasks.updateTask(id, desc); err != nil {
		return err
	}

	return nil
}

func HandleDelete(ctx *cli.Context) error {
	id, err := getIdFromString(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	if err := tasks.deleteTasksByID(id); err != nil {
		return err
	}

	return nil
}

func HandleDeleteDone(ctx *cli.Context) error {
	tasks.deleteTasksByStatus(STATUS_DONE)
	return nil
}

func HandleDeleteTodo(ctx *cli.Context) error {
	tasks.deleteTasksByStatus(STATUS_TODO)
	return nil
}

func HandleDeleteInProgress(ctx *cli.Context) error {
	tasks.deleteTasksByStatus(STATUS_IN_PROGRESS)
	return nil
}

func HandleMarkInProgress(ctx *cli.Context) error {
	id, err := getIdFromString(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	if err := tasks.markTaskAs(id, STATUS_IN_PROGRESS); err != nil {
		return err
	}

	return nil
}

func HandleMarkDone(ctx *cli.Context) error {
	id, err := getIdFromString(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	if err := tasks.markTaskAs(id, STATUS_DONE); err != nil {
		return err
	}

	return nil
}

func printListHeader(verbose bool) {
	if verbose {
		fmt.Printf("%-4s %-48s %-12s %-20s %-20s\n", "ID", "DESCRIPTION", "STATUS", "CREATED AT", "UPDATED AT")
	} else {
		fmt.Printf("%-4s %-48s %s\n", "ID", "DESCRIPTION", "STATUS")
	}
}

func HandleList(ctx *cli.Context) error {
	verbose := ctx.Bool("verbose")

	printListHeader(verbose)

	ids := getSortedTasksIDs()
	if len(ids) == 0 {
		fmt.Println("There are no tasks to display!")
		return nil
	}

	for _, id := range ids {
		fmt.Println(tasks.Tasks[id].String(verbose))
	}

	return nil
}

func HandleListDone(ctx *cli.Context) error {
	tasks.listByStatus(ctx.Bool("verbose"), STATUS_DONE)
	return nil
}

func HandleListTodo(ctx *cli.Context) error {
	tasks.listByStatus(ctx.Bool("verbose"), STATUS_TODO)
	return nil
}

func HandleListInProgress(ctx *cli.Context) error {
	tasks.listByStatus(ctx.Bool("verbose"), STATUS_IN_PROGRESS)
	return nil
}
