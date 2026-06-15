package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	exec.Command("stty", "-icanon", "-echo").Run()
	defer exec.Command("stty", "icanon", "echo").Run()

	query := ""
	selectedIndex := 0

	for {
		files, _ := os.ReadDir(".")

		filtered := make([]os.DirEntry, 0)
		for _, f := range files {
			if strings.Contains(strings.ToLower(f.Name()), strings.ToLower(query)) {
				filtered = append(filtered, f)
			}
		}

		if selectedIndex >= len(filtered) && len(filtered) > 0 {
			selectedIndex = len(filtered) - 1
		} else if len(filtered) == 0 {
			selectedIndex = 0
		}

		fmt.Print("\033[H\033[2J")

		fmt.Printf("Search: %s\n", query)
		fmt.Println("-------------------")

		for i, f := range filtered {
			if i == selectedIndex {
				fmt.Printf("\033[7m> %s\033[0m\n", f.Name())
			} else {
				fmt.Printf("  %s\n", f.Name())
			}
		}
		fmt.Println("\n[Up/Down] Navigate | [Enter] Select | [Esc] Quit")

		b := make([]byte, 3)
		n, _ := os.Stdin.Read(b)

		if n == 1 {
			switch b[0] {
			case 27:
				return
			case 10, 13:
				fmt.Print("\033[H\033[2J")
				if len(filtered) > 0 {
					fmt.Printf("You selected: %s\n", filtered[selectedIndex].Name())
				}
				return
			case 127, 8:
				if len(query) > 0 {
					query = query[:len(query)-1]
					selectedIndex = 0
				}
			default:
				if b[0] >= 32 && b[0] <= 126 {
					query += string(b[0])
					selectedIndex = 0
				}
			}
		} else if n == 3 && b[0] == 27 && b[1] == 91 {
			if b[2] == 65 {
				if selectedIndex > 0 {
					selectedIndex--
				}
			} else if b[2] == 66 {
				if selectedIndex < len(filtered)-1 {
					selectedIndex++
				}
			}
		}
	}
}