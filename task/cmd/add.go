package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/anwyho/gophercises/task/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to your task list",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := db.CreateTask(task)
		switch {
		case err != nil:
			fmt.Printf("Couldn't add task: %s", err.Error())
			os.Exit(1)
		default:
			fmt.Printf("Added \"%s\" to your task list.", task)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
