package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Color palette - Modern dark theme with vibrant accents
var (
	// Primary colors
	Purple    = lipgloss.Color("#9D4EDD")
	Cyan      = lipgloss.Color("#00D9FF")
	Green     = lipgloss.Color("#00FF9F")
	Yellow    = lipgloss.Color("#FFE66D")
	Red       = lipgloss.Color("#FF6B6B")
	Orange    = lipgloss.Color("#FF9F43")
	Pink      = lipgloss.Color("#FF6B9D")
	Blue      = lipgloss.Color("#4ECDC4")

	// Neutral colors
	White     = lipgloss.Color("#FFFFFF")
	Gray      = lipgloss.Color("#6C757D")
	DarkGray  = lipgloss.Color("#495057")
	Subtle    = lipgloss.Color("#383838")
)

// Icons for different states
const (
	IconPending   = "○"
	IconCompleted = "✓"
	IconAdd       = "+"
	IconWarning   = "⚠"
	IconInfo      = "ℹ"
	IconSuccess   = "✓"
	IconError     = "✗"
	IconArrow     = "→"
	IconBullet    = "•"
)

// Base styles
var (
	// Title styles
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Purple)

	SubtitleStyle = lipgloss.NewStyle().
		Foreground(Gray).
		Italic(true)

	// Header for sections
	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Cyan)

	// Card style for containing content
	CardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Subtle).
		Padding(0, 1)

	// Pending work card
	PendingCardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Yellow).
		Padding(0, 1)

	// Completed work card
	CompletedCardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Green).
		Padding(0, 1)

	// Task item styles
	PendingItemStyle = lipgloss.NewStyle().
		Foreground(Yellow)

	CompletedItemStyle = lipgloss.NewStyle().
		Foreground(Green)

	// Status message styles
	SuccessStyle = lipgloss.NewStyle().
		Foreground(Green).
		Bold(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(Red).
		Bold(true)

	WarningStyle = lipgloss.NewStyle().
		Foreground(Yellow).
		Bold(true)

	InfoStyle = lipgloss.NewStyle().
		Foreground(Cyan)

	// Badge styles
	CountBadgeStyle = lipgloss.NewStyle().
		Foreground(White).
		Background(Purple).
		Padding(0, 1).
		Bold(true)

	PendingBadgeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000")).
		Background(Yellow).
		Padding(0, 1).
		Bold(true)

	CompletedBadgeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000")).
		Background(Green).
		Padding(0, 1).
		Bold(true)

	// Muted text
	MutedStyle = lipgloss.NewStyle().
		Foreground(Gray)

	// Summary box - compact inline style
	SummaryStyle = lipgloss.NewStyle().
		Foreground(Gray).
		Italic(true)

	// Application header
	AppHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(Purple).
		Background(lipgloss.Color("#1a1a2e")).
		Padding(0, 2).
		MarginBottom(1)

	// Divider
	DividerStyle = lipgloss.NewStyle().
		Foreground(Subtle)

	// Empty state
	EmptyStateStyle = lipgloss.NewStyle().
		Foreground(Gray).
		Italic(true)

	// Prompt style
	PromptStyle = lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true)
)

// Helper functions for rendering

// RenderTitle renders a styled title
func RenderTitle(text string) string {
	return TitleStyle.Render(text)
}

// RenderHeader renders a section header
func RenderHeader(text string) string {
	return HeaderStyle.Render(text)
}

// RenderSuccess renders a success message with icon
func RenderSuccess(text string) string {
	return SuccessStyle.Render(IconSuccess + " " + text)
}

// RenderError renders an error message with icon
func RenderError(text string) string {
	return ErrorStyle.Render(IconError + " " + text)
}

// RenderWarning renders a warning message with icon
func RenderWarning(text string) string {
	return WarningStyle.Render(IconWarning + " " + text)
}

// RenderInfo renders an info message with icon
func RenderInfo(text string) string {
	return InfoStyle.Render(IconInfo + " " + text)
}

// RenderPendingItem renders a pending task item
func RenderPendingItem(index int, text string) string {
	icon := PendingItemStyle.Render(IconPending)
	num := MutedStyle.Render(fmt.Sprintf("%2d.", index))
	return fmt.Sprintf("  %s %s %s", num, icon, text)
}

// RenderCompletedItem renders a completed task item
func RenderCompletedItem(index int, text string) string {
	icon := CompletedItemStyle.Render(IconCompleted)
	num := MutedStyle.Render(fmt.Sprintf("%2d.", index))
	return fmt.Sprintf("  %s %s %s", num, icon, CompletedItemStyle.Render(text))
}

// RenderEmptyState renders an empty state message
func RenderEmptyState(text string) string {
	return EmptyStateStyle.Render(text)
}

// RenderDivider renders a horizontal divider
func RenderDivider(width int) string {
	divider := ""
	for i := 0; i < width; i++ {
		divider += "─"
	}
	return DividerStyle.Render(divider)
}

// RenderSummary renders a compact inline summary
func RenderSummary(title, content string) string {
	label := InfoStyle.Bold(true).Render(title + ":")
	return label + " " + SummaryStyle.Render(content)
}

// RenderBadge renders a count badge
func RenderBadge(count int, style lipgloss.Style) string {
	return style.Render(fmt.Sprintf(" %d ", count))
}
