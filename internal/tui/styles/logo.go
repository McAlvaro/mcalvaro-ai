package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// logoLines contains the ASCII art for the mcalvaro-ai logo.
var logoLines = []string{
	`   __  ___         ___     __                             ___     ____ `,
	`  /  |/  /   _____/   |   / /   _   __ ____ _____  ____  /   |   /  _/ `,
	` / /|_/ /   / ___/ /| |  / /   | | / // __ ` + "`" + `/ ___// __ \/ /| |   / /   `,
	`/ /  / /   / /__/ ___ | / /    | |/ // /_/ // /   / /_/ / ___ | _/ /    `,
	`/_/  /_/    \___/_/  |_|/_/     |___/ \__,_//_/    \____/_/  |_|/___/    `,
}

// gradientColors defines the top-to-bottom Capuchino Macchiato gradient:
// Espuma de leche (Foam) -> Crema suave -> Caramelo macchiato -> Canela y toffee -> Moca cálido -> Espresso tostado
var gradientColors = []lipgloss.Color{
	lipgloss.Color("#FFF5EB"), // Espuma de leche cremosa (Steamed Milk Foam)
	lipgloss.Color("#F3DCB7"), // Crema capuchino suave (Cappuccino Crema)
	lipgloss.Color("#E8C088"), // Caramelo macchiato (Golden Caramel)
	lipgloss.Color("#D19B68"), // Canela y toffee (Cinnamon Toffee)
	lipgloss.Color("#B07844"), // Café moca cálido (Warm Mocha)
	lipgloss.Color("#8A5528"), // Espresso tostado (Rich Roasted Espresso)
}

// RenderLogo returns the ASCII logo with a top-to-bottom gradient.
func RenderLogo() string {
	total := len(logoLines)
	if total == 0 {
		return ""
	}

	bands := len(gradientColors)
	var b strings.Builder

	for i, line := range logoLines {
		bandIdx := (i * bands) / total
		if bandIdx >= bands {
			bandIdx = bands - 1
		}
		style := lipgloss.NewStyle().Foreground(gradientColors[bandIdx])
		b.WriteString(style.Render(line))
		if i < total-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}
