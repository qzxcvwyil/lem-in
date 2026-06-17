package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Room struct {
	name string
	x, y int
}

type State struct {
	rooms                  map[string]Room
	links                  []([2]string)
	antPos                 map[int]string // ant id -> room name
	moves                  []map[int]string
	start                  string
	end                    string
	minX, maxX, minY, maxY int
}

func main() {
	state := parse()
	animate(state)
}

func parse() *State {
	s := &State{
		rooms:  make(map[string]Room),
		antPos: make(map[int]string),
	}

	scanner := bufio.NewScanner(os.Stdin)
	var ants int
	lineNum := 0
	nextIsStart := false
	nextIsEnd := false
	movesSection := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		// First line: number of ants
		if lineNum == 0 {
			ants, _ = strconv.Atoi(line)
			lineNum++
			continue
		}

		// Detect moves section (starts with L)
		if strings.HasPrefix(line, "L") {
			movesSection = true
		}

		if movesSection {
			turn := make(map[int]string)
			parts := strings.Fields(line)
			for _, p := range parts {
				// Lx-room
				p = strings.TrimPrefix(p, "L")
				dash := strings.Index(p, "-")
				if dash < 0 {
					continue
				}
				id, err := strconv.Atoi(p[:dash])
				if err != nil {
					continue
				}
				room := p[dash+1:]
				turn[id] = room
			}
			s.moves = append(s.moves, turn)
			continue
		}

		if line == "##start" {
			nextIsStart = true
			continue
		}
		if line == "##end" {
			nextIsEnd = true
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Link
		if strings.Contains(line, "-") && !strings.Contains(line, " ") {
			parts := strings.SplitN(line, "-", 2)
			if len(parts) == 2 {
				s.links = append(s.links, [2]string{parts[0], parts[1]})
			}
			continue
		}

		// Room: name x y
		parts := strings.Fields(line)
		if len(parts) == 3 {
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			r := Room{name: parts[0], x: x, y: y}
			s.rooms[parts[0]] = r
			if nextIsStart {
				s.start = parts[0]
				nextIsStart = false
			} else if nextIsEnd {
				s.end = parts[0]
				nextIsEnd = false
			}
		}
	}

	// Initial positions: all ants at start
	for i := 1; i <= ants; i++ {
		s.antPos[i] = s.start
	}

	// Compute bounds
	first := true
	for _, r := range s.rooms {
		if first {
			s.minX, s.maxX = r.x, r.x
			s.minY, s.maxY = r.y, r.y
			first = false
		}
		if r.x < s.minX {
			s.minX = r.x
		}
		if r.x > s.maxX {
			s.maxX = r.x
		}
		if r.y < s.minY {
			s.minY = r.y
		}
		if r.y > s.maxY {
			s.maxY = r.y
		}
	}

	return s
}

func animate(s *State) {
	width := 60
	height := 20

	scaleX := func(x int) int {
		if s.maxX == s.minX {
			return width / 2
		}
		return (x-s.minX)*(width-4)/(s.maxX-s.minX) + 2
	}
	scaleY := func(y int) int {
		if s.maxY == s.minY {
			return height / 2
		}
		return (y-s.minY)*(height-3)/(s.maxY-s.minY) + 1
	}

	// ANSI colors
	const (
		reset  = "\033[0m"
		bold   = "\033[1m"
		cyan   = "\033[36m"
		yellow = "\033[33m"
		green  = "\033[32m"
		red    = "\033[31m"
		gray   = "\033[90m"
	)

	render := func(turn int) {
		// Build grid
		grid := make([][]rune, height)
		for i := range grid {
			grid[i] = make([]rune, width)
			for j := range grid[i] {
				grid[i][j] = ' '
			}
		}

		// Place room dots
		for _, r := range s.rooms {
			cx := scaleX(r.x)
			cy := scaleY(r.y)
			if cy >= 0 && cy < height && cx >= 0 && cx < width {
				grid[cy][cx] = '·'
			}
		}

		// Place ants (just show number or 'A' if >9)
		roomAnts := make(map[string][]int)
		for id, room := range s.antPos {
			roomAnts[room] = append(roomAnts[room], id)
		}

		// Clear screen and move cursor to top
		fmt.Print("\033[H\033[2J")

		// Header
		fmt.Printf("%s%s=== Lem-in Visualizer === Turn: %d/%d ===%s\n\n",
			bold, cyan, turn, len(s.moves), reset)

		// Draw links as dashes (simple)
		for _, link := range s.links {
			r1, ok1 := s.rooms[link[0]]
			r2, ok2 := s.rooms[link[1]]
			if !ok1 || !ok2 {
				continue
			}
			x1, y1 := scaleX(r1.x), scaleY(r1.y)
			x2, y2 := scaleX(r2.x), scaleY(r2.y)
			// Draw midpoint as '-'
			mx, my := (x1+x2)/2, (y1+y2)/2
			if my >= 0 && my < height && mx >= 0 && mx < width {
				if grid[my][mx] == ' ' {
					grid[my][mx] = '-'
				}
			}
		}

		// Print grid with colors
		for row := 0; row < height; row++ {
			for col := 0; col < width; col++ {
				ch := grid[row][col]
				switch ch {
				case '·':
					// Check if any ant here
					for name, r := range s.rooms {
						if scaleX(r.x) == col && scaleY(r.y) == row {
							ants := roomAnts[name]
							if len(ants) > 0 {
								if name == s.end {
									fmt.Printf("%s%s%d%s", bold, green, len(ants), reset)
								} else {
									fmt.Printf("%s%s%d%s", bold, yellow, len(ants), reset)
								}
							} else if name == s.start {
								fmt.Printf("%s%sS%s", bold, red, reset)
							} else if name == s.end {
								fmt.Printf("%s%sE%s", bold, green, reset)
							} else {
								fmt.Printf("%s·%s", gray, reset)
							}
							goto nextCol
						}
					}
				case '-':
					fmt.Printf("%s-%s", gray, reset)
					goto nextCol
				default:
					fmt.Print(string(ch))
				}
			nextCol:
			}
			fmt.Println()
		}

		// Legend
		fmt.Printf("\n%sS%s = start  %sE%s = end  %s#%s = ants (count)%s\n",
			red, reset, green, reset, yellow, reset, reset)
		fmt.Printf("Ants at end: %s%d%s / %d\n",
			green, len(roomAnts[s.end]), reset, len(s.antPos))
	}

	// Initial state (turn 0)
	render(0)
	time.Sleep(800 * time.Millisecond)

	// Animate each turn
	for i, turn := range s.moves {
		// Apply moves
		for id, room := range turn {
			s.antPos[id] = room
		}
		render(i + 1)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n\033[1m\033[32mDone! All ants reached the end.\033[0m")
}
