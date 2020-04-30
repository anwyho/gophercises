/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/anwyho/gophercises/task/db"
	"github.com/spf13/cobra"
)

type Set map[int]struct{}
type Empty struct{}

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task as done",
	Run: func(cmd *cobra.Command, args []string) {
		ids := make(Set)
		tasks, err := db.AllTasks()
		if err != nil || len(tasks) == 0 {
			fmt.Println("Couldn't find any tasks.")
			os.Exit(1)
		}
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			switch {
			case err != nil:
				fmt.Printf("\"%s\" is not a task number.\n", arg)
			case len(tasks) < id || id < 1:
				fmt.Printf("Task %d does not exist.\n", id)
			default:
				ids[id] = Empty{}
			}
		}
		if len(ids) == 0 {
			os.Exit(1)
		}
		isSuccessful := true
		for id := range ids {
			task := tasks[id-1] // bounds are handled above
			err := db.DeleteTask(task.Key)
			switch {
			case err != nil:
				fmt.Printf("Failed to mark task as done: %s\n", err.Error())
				isSuccessful = false
			default:
				fmt.Printf("Completed \"%s\" Key=%d\n", task.Value, task.Key)
			}
		}
		if !isSuccessful {
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
