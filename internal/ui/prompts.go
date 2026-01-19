package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/sandepten/work-obsidian-noter/internal/notes"
)

// Prompter handles interactive CLI prompts
type Prompter struct{}

// NewPrompter creates a new prompter
func NewPrompter() *Prompter {
	return &Prompter{}
}

// ConfirmCompletion asks if a work item was completed
func (p *Prompter) ConfirmCompletion(item notes.WorkItem) (bool, error) {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Did you complete: \"%s\"", item.Text),
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// SelectPendingItems allows selecting multiple pending items to mark as done
func (p *Prompter) SelectPendingItems(items []notes.WorkItem) ([]int, error) {
	if len(items) == 0 {
		return nil, nil
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "> {{ .Text | cyan }}",
		Inactive: "  {{ .Text }}",
		Selected: "{{ .Text | green }}",
	}

	var selectedIndices []int

	fmt.Println("\nReview pending items (press 'q' when done):")

	for i, item := range items {
		completed, err := p.ConfirmCompletion(item)
		if err != nil {
			return selectedIndices, err
		}
		if completed {
			selectedIndices = append(selectedIndices, i)
		}
	}

	// Suppress unused variable warning
	_ = templates

	return selectedIndices, nil
}

// PromptForNewItem asks for a new work item
func (p *Prompter) PromptForNewItem() (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter new work item (leave empty to skip)",
	}

	result, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			return "", nil
		}
		return "", err
	}

	return result, nil
}

// ConfirmAction asks for a yes/no confirmation
func (p *Prompter) ConfirmAction(message string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// SelectFromList allows selecting an item from a list
func (p *Prompter) SelectFromList(label string, items []string) (int, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return -1, err
	}

	return index, nil
}

// DisplayWorkItems shows a formatted list of work items
func (p *Prompter) DisplayWorkItems(pending, completed []notes.WorkItem) {
	fmt.Println("\n--- Pending Work ---")
	if len(pending) == 0 {
		fmt.Println("  No pending items")
	} else {
		for i, item := range pending {
			fmt.Printf("  %d. [ ] %s\n", i+1, item.Text)
		}
	}

	fmt.Println("\n--- Completed Work ---")
	if len(completed) == 0 {
		fmt.Println("  No completed items")
	} else {
		for i, item := range completed {
			fmt.Printf("  %d. [x] %s\n", i+1, item.Text)
		}
	}
	fmt.Println()
}

// DisplayMessage shows a message to the user
func (p *Prompter) DisplayMessage(message string) {
	fmt.Println(message)
}

// DisplayError shows an error message
func (p *Prompter) DisplayError(message string) {
	fmt.Printf("Error: %s\n", message)
}

// DisplaySuccess shows a success message
func (p *Prompter) DisplaySuccess(message string) {
	fmt.Printf("Success: %s\n", message)
}
