package render

import (
	"strings"

	"charm.land/lipgloss/v2"
)

const (
	minimumBandLength = 2
	logoGlowWidth     = 3
)

func Band(width, length int, double bool) string {
	bar := Center(Paint(strings.Repeat("=", min(width, max(minimumBandLength, length))), Muted), width)

	if double {
		return bar + "\n" + bar
	}

	return bar
}

func LogoHighlight(progress, frames int) string {
	column := progress * lipgloss.Width(Logo) / frames
	lines := strings.Split(Logo, "\n")

	for row, line := range lines {
		var text strings.Builder

		for index, char := range []rune(line) {
			color := Accent

			if index >= column && index < column+logoGlowWidth {
				color = Foreground
			}

			text.WriteString(Paint(string(char), color))
		}

		lines[row] = text.String()
	}

	return strings.Join(lines, "\n")
}

func Thinking(frame int, animated bool) string {
	const text = "Thinking…"
	if !animated {
		return text
	}

	frames := []rune("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏")

	return string(frames[frame%len(frames)]) + " " + text
}
