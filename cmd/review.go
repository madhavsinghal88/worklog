package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/sandepten/work-obsidian-noter/internal/notes"
	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review pending items from previous notes",
	Long: `Manually review and process pending items from previous notes
without creating a new note or generating summaries.`,
	RunE: runReview,
}

func init() {
	rootCmd.AddCommand(reviewCmd)
}

func runReview(cmd *cobra.Command, args []string) error {
	today := time.Now().Truncate(24 * time.Hour)

	// Find the most recent previous note
	previousNote, err := parser.FindMostRecentNote(today)
	if err != nil {
		return fmt.Errorf("error finding previous note: %w", err)
	}

	if previousNote == nil {
		prompter.DisplayMessage("No previous notes found.")
		return nil
	}

	fmt.Printf("Reviewing note: %s (Date: %s)\n\n", filepath.Base(previousNote.FilePath), previousNote.Date.Format("2006-01-02"))

	if !previousNote.HasPendingWork() {
		prompter.DisplayMessage("No pending items to review.")
		prompter.DisplayWorkItems(previousNote.PendingWork, previousNote.CompletedWork)
		return nil
	}

	fmt.Println("Review pending items:\n")

	completedIndices, err := prompter.SelectPendingItems(previousNote.PendingWork)
	if err != nil {
		return fmt.Errorf("error reviewing items: %w", err)
	}

	if len(completedIndices) == 0 {
		prompter.DisplayMessage("No items marked as completed.")
		return nil
	}

	// Sort indices in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(completedIndices)))

	// Mark items as completed
	for _, idx := range completedIndices {
		item := previousNote.PendingWork[idx]
		previousNote.CompletedWork = append(previousNote.CompletedWork, notes.WorkItem{
			Text:      item.Text,
			Completed: true,
		})
		// Remove from pending
		previousNote.PendingWork = append(previousNote.PendingWork[:idx], previousNote.PendingWork[idx+1:]...)
	}

	// Save the note
	if err := writer.WriteNote(previousNote); err != nil {
		return fmt.Errorf("error saving note: %w", err)
	}

	prompter.DisplaySuccess(fmt.Sprintf("Marked %d item(s) as completed.", len(completedIndices)))

	// Show updated state
	prompter.DisplayWorkItems(previousNote.PendingWork, previousNote.CompletedWork)

	return nil
}
