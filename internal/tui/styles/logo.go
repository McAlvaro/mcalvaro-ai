package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// logoLines contains the ASCII art for the mcalvaro-ai logo.
var logoLines = []string{
	"MMMMMMMM               MMMMMMMM                                   AAA               lllllll                                                                                                 AAA               IIIIIIIIII",
	"M:::::::M             M:::::::M                                  A:::A              l:::::l                                                                                                A:::A              I::::::::I",
	"M::::::::M           M::::::::M                                 A:::::A             l:::::l                                                                                               A:::::A             I::::::::I",
	"M:::::::::M         M:::::::::M                                A:::::::A            l:::::l                                                                                              A:::::::A            II::::::II",
	"M::::::::::M       M::::::::::M    cccccccccccccccc           A:::::::::A            l::::lvvvvvvv           vvvvvvvaaaaaaaaaaaaa  rrrrr   rrrrrrrrr      ooooooooooo                   A:::::::::A             I::::I  ",
	"M:::::::::::M     M:::::::::::M  cc:::::::::::::::c          A:::::A:::::A           l::::l v:::::v         v:::::v a::::::::::::a r::::rrr:::::::::r   oo:::::::::::oo                A:::::A:::::A            I::::I  ",
	"M:::::::M::::M   M::::M:::::::M c:::::::::::::::::c         A:::::A A:::::A          l::::l  v:::::v       v:::::v  aaaaaaaaa:::::ar:::::::::::::::::r o:::::::::::::::o              A:::::A A:::::A           I::::I  ",
	"M::::::M M::::M M::::M M::::::Mc:::::::cccccc:::::c        A:::::A   A:::::A         l::::l   v:::::v     v:::::v            a::::arr::::::rrrrr::::::ro:::::ooooo:::::o             A:::::A   A:::::A          I::::I  ",
	"M::::::M  M::::M::::M  M::::::Mc::::::c     ccccccc       A:::::A     A:::::A        l::::l    v:::::v   v:::::v      aaaaaaa:::::a r:::::r     r:::::ro::::o     o::::o            A:::::A     A:::::A         I::::I  ",
	"M::::::M   M:::::::M   M::::::Mc:::::c                   A:::::AAAAAAAAA:::::A       l::::l     v:::::v v:::::v     aa::::::::::::a r:::::r     rrrrrrro::::o     o::::o           A:::::AAAAAAAAA:::::A        I::::I  ",
	"M::::::M    M:::::M    M::::::Mc:::::c                  A:::::::::::::::::::::A      l::::l      v:::::v:::::v     a::::aaaa::::::a r:::::r            o::::o     o::::o          A:::::::::::::::::::::A       I::::I  ",
	"M::::::M     MMMMM     M::::::Mc::::::c     ccccccc    A:::::AAAAAAAAAAAAA:::::A     l::::l       v:::::::::v     a::::a    a:::::a r:::::r            o::::o     o::::o         A:::::AAAAAAAAAAAAA:::::A      I::::I  ",
	"M::::::M               M::::::Mc:::::::cccccc:::::c   A:::::A             A:::::A   l::::::l       v:::::::v      a::::a    a:::::a r:::::r            o:::::ooooo:::::o        A:::::A             A:::::A   II::::::II",
	"M::::::M               M::::::M c:::::::::::::::::c  A:::::A               A:::::A  l::::::l        v:::::v       a:::::aaaa::::::a r:::::r            o:::::::::::::::o       A:::::A               A:::::A  I::::::::I",
	"M::::::M               M::::::M  cc:::::::::::::::c A:::::A                 A:::::A l::::::l         v:::v         a::::::::::aa:::ar:::::r             oo:::::::::::oo       A:::::A                 A:::::A I::::::::I",
	"MMMMMMMM               MMMMMMMM    ccccccccccccccccAAAAAAA                   AAAAAAAllllllll          vvv           aaaaaaaaaa  aaaarrrrrrr               ooooooooooo        AAAAAAA                   AAAAAAAIIIIIIIIII",
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
