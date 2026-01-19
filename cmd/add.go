package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new pending work item",
	Long:  `Add a new pending work item to today's note.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	today := time.Now().Truncate(24 * time.Hour)
	taskText := strings.Join(args, " ")

	// Get or create today's note
	todayNote, err := parser.FindTodayNote(today)
	if err != nil {
		return fmt.Errorf("error finding today's note: %w", err)
	}

	if todayNote == nil {
		todayNote = writer.CreateTodayNote(today)
		fmt.Println("Creating today's note...")
	}

	// Add the new item
	todayNote.AddPendingItem(taskText)

	// Save the note
	if err := writer.WriteNote(todayNote); err != nil {
		return fmt.Errorf("error saving note: %w", err)
	}

	prompter.DisplaySuccess(fmt.Sprintf("Added pending item: %s", taskText))
	return nil
}
