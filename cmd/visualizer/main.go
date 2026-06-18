package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
	roomOrder              []string // insertion order, used for stable legend ordering
	links                  [][2]string
	antPos                 map[int]string // ant id -> current room name
	moves                  []map[int]string
	moveLines              []string // raw "Lx-y Lz-w ..." text per turn, for the side log
	start                  string
	end                    string
	totalAnts              int
	minX, maxX, minY, maxY int
}

type Canvas struct {
	w, h  int
	chars [][]rune
	color [][]string
}

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	cyan   = "\033[36m"
	yellow = "\033[33m"
	green  = "\033[32m"
	red    = "\033[31m"
	gray   = "\033[90m"
	white  = "\033[97m"
)

func main() {
	state := parse(os.Stdin)
	if len(state.rooms) == 0 {
		fmt.Println("No colony data received. Pipe lem-in's output into the visualizer:")
		fmt.Println("  go run . examples/example00.txt | go run ./cmd/visualizer")
		return
	}
	animate(state)
}

func newCanvas(w, h int) *Canvas {
	c := &Canvas{w: w, h: h}
	c.chars = make([][]rune, h)
	c.color = make([][]string, h)
	for y := 0; y < h; y++ {
		c.chars[y] = make([]rune, w)
		c.color[y] = make([]string, w)
		for x := 0; x < w; x++ {
			c.chars[y][x] = ' '
			c.color[y][x] = ""
		}
	}
	return c
}

func (c *Canvas) set(x, y int, ch rune, col string) {
	if x < 0 || x >= c.w || y < 0 || y >= c.h {
		return
	}
	if c.chars[y][x] != ' ' && (ch == '-' || ch == '|' || ch == '\\' || ch == '/') {
		return
	}
	c.chars[y][x] = ch
	c.color[y][x] = col
}

func (c *Canvas) drawLine(x0, y0, x1, y1 int, ch rune, col string) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x1 < x0 {
		sx = -1
	}
	if y1 < y0 {
		sy = -1
	}
	x, y := x0, y0
	if dx > dy {
		errAcc := dx / 2
		for i := 0; i <= dx; i++ {
			c.set(x, y, ch, col)
			errAcc -= dy
			if errAcc < 0 {
				y += sy
				errAcc += dx
			}
			x += sx
		}
	} else {
		errAcc := dy / 2
		for i := 0; i <= dy; i++ {
			c.set(x, y, ch, col)
			errAcc -= dx
			if errAcc < 0 {
				x += sx
				errAcc += dy
			}
			y += sy
		}
	}
}

func (c *Canvas) text(x, y int, s string, col string) {
	for i, ch := range s {
		c.set(x+i, y, ch, col)
	}
}

func (c *Canvas) reserve(x, y int) {
	if x < 0 || x >= c.w || y < 0 || y >= c.h {
		return
	}
	c.chars[y][x] = ' '
	c.color[y][x] = ""
}

func canvasHeight(s *State) int {
	const minHeight = 28
	const maxHeight = 50

	uniqueY := make(map[int]struct{})
	for _, r := range s.rooms {
		uniqueY[r.y] = struct{}{}
	}

	needed := len(uniqueY)*2 + 4
	if needed < minHeight {
		return minHeight
	}
	if needed > maxHeight {
		return maxHeight
	}
	return needed
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func parse(r *os.File) *State {
	s := &State{
		rooms:  make(map[string]Room),
		antPos: make(map[int]string),
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

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

		if lineNum == 0 {
			ants, _ = strconv.Atoi(line)
			lineNum++
			continue
		}

		if strings.HasPrefix(line, "L") && looksLikeMoveLine(line) {
			movesSection = true
		}

		if movesSection {
			turn := make(map[int]string)
			parts := strings.Fields(line)
			for _, p := range parts {
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
			s.moveLines = append(s.moveLines, line)
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

		if strings.Contains(line, "-") && !strings.Contains(line, " ") {
			parts := strings.SplitN(line, "-", 2)
			if len(parts) == 2 {
				s.links = append(s.links, [2]string{parts[0], parts[1]})
			}
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 3 {
			x, errX := strconv.Atoi(parts[1])
			y, errY := strconv.Atoi(parts[2])
			if errX != nil || errY != nil {
				continue
			}
			r := Room{name: parts[0], x: x, y: y}
			if _, exists := s.rooms[parts[0]]; !exists {
				s.roomOrder = append(s.roomOrder, parts[0])
			}
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

	s.totalAnts = ants
	for i := 1; i <= ants; i++ {
		s.antPos[i] = s.start
	}

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

func looksLikeMoveLine(line string) bool {
	fields := strings.Fields(line)
	for _, f := range fields {
		if !strings.HasPrefix(f, "L") {
			return false
		}
		rest := strings.TrimPrefix(f, "L")
		dash := strings.Index(rest, "-")
		if dash <= 0 {
			return false
		}
		if _, err := strconv.Atoi(rest[:dash]); err != nil {
			return false
		}
	}
	return len(fields) > 0
}

func animate(s *State) {
	mapWidth := 100
	mapHeight := canvasHeight(s)
	labelPad := 1

	scaleX := func(x int) int {
		if s.maxX == s.minX {
			return mapWidth / 2
		}
		return (x-s.minX)*(mapWidth-2*labelPad-12)/(s.maxX-s.minX) + labelPad + 2
	}
	scaleY := func(y int) int {
		if s.maxY == s.minY {
			return mapHeight / 2
		}
		return (y-s.minY)*(mapHeight-3)/(s.maxY-s.minY) + 1
	}

	pos := make(map[string][2]int)
	for name, r := range s.rooms {
		pos[name] = [2]int{scaleX(r.x), scaleY(r.y)}
	}

	render := func(turnIdx int) {
		canvas := newCanvas(mapWidth, mapHeight)

		for _, link := range s.links {
			p1, ok1 := pos[link[0]]
			p2, ok2 := pos[link[1]]
			if !ok1 || !ok2 {
				continue
			}
			canvas.drawLine(p1[0], p1[1], p2[0], p2[1], '.', gray)
		}

		sortedNames := append([]string{}, s.roomOrder...)
		sort.Strings(sortedNames) // stable, deterministic draw order
		for _, name := range sortedNames {
			p := pos[name]
			switch name {
			case s.start:
				canvas.set(p[0], p[1], 'S', bold+red)
			case s.end:
				canvas.set(p[0], p[1], 'E', bold+green)
			default:
				canvas.set(p[0], p[1], 'o', white)
			}
			canvas.text(p[0]+1, p[1], name, dim+cyan)
		}
		for _, name := range sortedNames {
			p := pos[name]
			canvas.reserve(p[0]+1+len([]rune(name)), p[1])
		}

		roomAnts := make(map[string][]int)
		for id, room := range s.antPos {
			roomAnts[room] = append(roomAnts[room], id)
		}
		for name, antIDs := range roomAnts {
			p, ok := pos[name]
			if !ok {
				continue
			}
			label := fmt.Sprintf("(%d)", len(antIDs))
			col := bold + yellow
			if name == s.end {
				col = bold + green
			}
			canvas.text(p[0], p[1]+1, label, col)
		}

		clearScreen()

		fmt.Printf("%s%s=== Lem-in Visualizer ===  Turn %d / %d ===%s\n",
			bold, cyan, turnIdx, len(s.moves), reset)
		fmt.Printf("%sS%s start   %sE%s end   %so%s room   %s(n)%s ants waiting there\n\n",
			bold+red, reset, bold+green, reset, white, reset, yellow, reset)

		printCanvas(canvas)

		atEnd := len(roomAnts[s.end])
		fmt.Printf("\nAnts at end: %s%d%s / %d\n", green, atEnd, reset, s.totalAnts)

		if turnIdx > 0 && turnIdx <= len(s.moveLines) {
			fmt.Printf("%sMoves this turn:%s %s\n", dim, reset, s.moveLines[turnIdx-1])
		}
	}

	render(0)
	time.Sleep(900 * time.Millisecond)

	for i, turn := range s.moves {
		for id, room := range turn {
			s.antPos[id] = room
		}
		render(i + 1)
		time.Sleep(600 * time.Millisecond)
	}

	fmt.Printf("\n%s%sDone! All %d ants reached %s.%s\n", bold, green, s.totalAnts, s.end, reset)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J\033[3J")
	os.Stdout.Sync()
}

func printCanvas(c *Canvas) {
	var sb strings.Builder
	for y := 0; y < c.h; y++ {
		lastColor := ""
		for x := 0; x < c.w; x++ {
			ch := c.chars[y][x]
			col := c.color[y][x]
			if col != lastColor {
				if lastColor != "" {
					sb.WriteString(reset)
				}
				if col != "" {
					sb.WriteString(col)
				}
				lastColor = col
			}
			sb.WriteRune(ch)
		}
		if lastColor != "" {
			sb.WriteString(reset)
		}
		sb.WriteByte('\n')
	}
	fmt.Print(sb.String())
}
