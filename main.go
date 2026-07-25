package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Id          int       `json:"Id"`
	Description []string  `json:"Description"`
	Status      string    `json:"Status"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

func LoadTaskList() ([]Task, error) {
	data, err := os.ReadFile("data.json")
	if os.IsNotExist(err) {
		err = os.WriteFile("data.json", []byte("[]"), 0644)
		if err != nil {
			return nil, err
		}
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	var tasks []Task

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func SaveTaskList(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile("data.json", data, 0644)
}

func AddTask(tasks *[]Task, s []string) {
	var id int

	if len(*tasks) == 0 {
		id = 1
	} else {
		id = (*tasks)[len(*tasks)-1].Id
		id += 1
	}

	task := Task{
		Id:          id,
		Description: s,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	*tasks = append(*tasks, task)
}

func UpdateTask(tasks *[]Task, i int, s []string) error {
	for t := range *tasks {
		task := &(*tasks)[t]
		if task.Id == i {
			task.Description = s
			task.UpdatedAt = time.Now()
			return nil
		}
	}
	err := errors.New("Нет такого id")
	return err
}

func DeleteTask(tasks *[]Task, i int) error {
	for t := range *tasks {
		task := &(*tasks)[t]
		if task.Id == i {
			*tasks = append((*tasks)[:t], (*tasks)[t+1:]...)
			return nil
		}
	}
	err := errors.New("Нет такого id")
	return err
}

func MarkInprogress(tasks *[]Task, i int) error {
	for t := range *tasks {
		task := &(*tasks)[t]
		if task.Id == i {
			(*tasks)[t].Status = "in progress"
			return nil
		}
	}
	err := errors.New("Нет такого id")
	return err
}

func MarkDone(tasks *[]Task, i int) error {
	for t := range *tasks {
		task := &(*tasks)[t]
		if task.Id == i {
			(*tasks)[t].Status = "done"
			return nil
		}
	}
	err := errors.New("Нет такого id")
	return err
}

func MarkTodo(tasks *[]Task, i int) error {
	for t := range *tasks {
		task := &(*tasks)[t]
		if task.Id == i {
			(*tasks)[t].Status = "todo"
			return nil
		}
	}
	err := errors.New("Нет такого id")
	return err
}

func List(tasks *[]Task) {
	fmt.Println(tasks)
}

func ListDone(tasks *[]Task) {
	Done := []Task{}
	for t := range *tasks {
		task := (*tasks)[t]
		if task.Status == "done" {
			Done = append(Done, task)
		}
	}
	fmt.Println(Done)

}

func ListTodo(tasks *[]Task) {
	Todo := []Task{}
	for t := range *tasks {
		task := (*tasks)[t]
		if task.Status == "todo" {
			Todo = append(Todo, task)
		}
	}
	fmt.Println(Todo)

}

func ListInProgress(tasks *[]Task) {
	InProgress := []Task{}
	for t := range *tasks {
		task := (*tasks)[t]
		if task.Status == "in progress" {
			InProgress = append(InProgress, task)
		}
	}
	fmt.Println(InProgress)

}

func main() {

	TaskList, err := LoadTaskList()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Запуск")

loop:
	for {
		var id int
		scanner := bufio.NewScanner(os.Stdin)
		if ok := scanner.Scan(); !ok {
			fmt.Println("ошибка ввода")
			return
		}
		text := scanner.Text()
		fields := strings.Fields(text)

		switch fields[0] {
		case "help":
			fmt.Println("Список команд:\nhelp - список команд\nadd (название) - добавить задачу\nupdate (id) (название) - сменить название задачи\ndelete (id) - удалить задачу\nmark-done/mark-in-progress/mark-todo (id) - изменить состояние задачи\nclear - отчистить список\nlist/list-todo/list-in-progress/list-done - вывести список задач\nend - закончить сессию")
		case "add":
			AddTask(&TaskList, fields[1:])
		case "update":
			id, err = strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println("Это не id")
			}
			if err = UpdateTask(&TaskList, id, fields[2:]); err != nil {
				fmt.Println(err)
			}
		case "delete":
			id, err = strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println("Это не id")
			}
			if err = DeleteTask(&TaskList, id); err != nil {
				fmt.Println(err)
			}
		case "mark-done":
			id, err = strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println("Это не id")
			}
			if err = MarkDone(&TaskList, id); err != nil {
				fmt.Println(err)
			}
		case "mark-in-progress":
			id, err = strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println("Это не id")
			}
			if err = MarkInprogress(&TaskList, id); err != nil {
				fmt.Println(err)
			}
		case "mark-todo":
			id, err = strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println("Это не id")
			}
			if err = MarkTodo(&TaskList, id); err != nil {
				fmt.Println(err)
			}
		case "clear":
			TaskList = []Task{}
		case "end":
			break loop
		case "list":
			List(&TaskList)
		case "list-done":
			ListDone(&TaskList)
		case "list-todo":
			ListTodo(&TaskList)
		case "list-in-progress":
			ListInProgress(&TaskList)
		default:
			fmt.Println("Неизвестно")
		}
	}

	err = SaveTaskList(TaskList)
	if err != nil {
		fmt.Println(err)
		return
	}

}
