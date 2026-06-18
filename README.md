# Lem-in

Upon successfully finding the quickest path, lem-in will display the content of the file passed as argument and each move the ants make from room to room.

Standard packages used only.

Make sure your Go version is **1.24.2** or above

## Program start:

1. Clone project: `git clone https://github.com/qzxcvwyil/lem-in.git`
2. Open directory: `cd lem-in`
3. Start the app: `go run cmd/main.go examples/example00.txt`

If something doesn't work contact **Telegram**: **
@wqzxyl**

## Project structure:

1. **cmd/** - main.go file
   - **visualizer** - bonus ant farm visualizer (reads lem-in's stdout and animates it)
2. **internal/**
   - **app** - orchestration layer
   - **model** - domain types
   - **parser** - input file parsing
   - **validation** - input validation
   - **service/pathfinding** - BFS-based path finding
   - **service/solver** - ant distribution and simulation
3. **examples/** - example input files

## Bonus: visualizer

```
go run cmd/main.go examples/example05.txt | go run cmd/visualizer/main.go
```

Renders an animated map of the colony directly in the terminal: rooms are drawn
at their real coordinates from the input file, connected by actual lines
(not just a single dash), each room is labeled with its name, and the number
of ants currently waiting in a room is shown right under it. The current
turn's moves (`Lx-y ...`) are printed below the map so you can follow along
with the raw output too.

---

Made by Madina Sabyrgali