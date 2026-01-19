package notes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Writer handles writing markdown notes to disk
type Writer struct {
	notesDir      string
	workplaceName string
}

// NewWriter creates a new note writer
func NewWriter(notesDir, workplaceName string) *Writer {
	return &Writer{
		notesDir:      notesDir,
		workplaceName: workplaceName,
	}
}

// WriteNote writes a note to disk
func (w *Writer) WriteNote(note *Note) error {
	if note.FilePath == "" {
		note.FilePath = filepath.Join(w.notesDir, GenerateFilename(note.Date, w.workplaceName))
	}

	content := w.generateMarkdown(note)
	return os.WriteFile(note.FilePath, []byte(content), 0644)
}

// CreateTodayNote creates a new note for today
func (w *Writer) CreateTodayNote(date time.Time) *Note {
	note := NewNote(date, w.workplaceName)
	note.FilePath = filepath.Join(w.notesDir, GenerateFilename(date, w.workplaceName))
	return note
}

// generateMarkdown generates the markdown content for a note
func (w *Writer) generateMarkdown(note *Note) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("id: %s\n", note.ID))
	sb.WriteString("aliases: []\n")
	sb.WriteString("tags:\n")
	for _, tag := range note.Tags {
		sb.WriteString(fmt.Sprintf("  - %s\n", tag))
	}
	sb.WriteString(fmt.Sprintf("date: %s\n", note.Date.Format("2006-01-02")))
	sb.WriteString("---\n\n")

	// Title
	sb.WriteString(fmt.Sprintf("# %s\n\n", note.Title))

	// Summary fields
	sb.WriteString(fmt.Sprintf("summary::%s\n\n", formatInlineSummary(note.Summary)))
	sb.WriteString(fmt.Sprintf("yesterday's summary::%s\n\n", formatInlineSummary(note.YesterdaySummary)))

	// Pending Work section
	sb.WriteString("## Pending Work\n\n")
	for _, item := range note.PendingWork {
		sb.WriteString(fmt.Sprintf("- [ ] %s\n", item.Text))
	}
	sb.WriteString("\n")

	// Work Completed section
	sb.WriteString("## Work Completed\n\n")
	for _, item := range note.CompletedWork {
		sb.WriteString(fmt.Sprintf("- [x] %s\n", item.Text))
	}
	sb.WriteString("\n")

	return sb.String()
}

// formatInlineSummary formats the summary for inline display
func formatInlineSummary(summary string) string {
	if summary == "" {
		return ""
	}
	return " " + summary
}

// UpdateSummary updates the summary field in an existing note
func (w *Writer) UpdateSummary(note *Note, summary string) error {
	note.Summary = summary
	return w.WriteNote(note)
}

// UpdateYesterdaySummary updates the yesterday's summary field
func (w *Writer) UpdateYesterdaySummary(note *Note, summary string) error {
	note.YesterdaySummary = summary
	return w.WriteNote(note)
}

// MovePendingToCompleted moves all pending items to completed for an item
func (w *Writer) MovePendingToCompleted(note *Note, index int) error {
	note.MarkItemCompleted(index)
	return w.WriteNote(note)
}

// AddPendingItem adds a pending item to a note and saves
func (w *Writer) AddPendingItem(note *Note, text string) error {
	note.AddPendingItem(text)
	return w.WriteNote(note)
}

// AddCompletedItem adds a completed item to a note and saves
func (w *Writer) AddCompletedItem(note *Note, text string) error {
	note.AddCompletedItem(text)
	return w.WriteNote(note)
}
