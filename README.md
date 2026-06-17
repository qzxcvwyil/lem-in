# Lem-in

Upon successfully finding the quickest path, lem-in will display the content of the file passed as argument and each move the ants make from room to room.

Standard packages used only.

Make sure your Go version is **1.24.2** or above

## Program start:

1. Clone project: `git clone https://01.tomorrow-school.ai/git/msabyrga/lem-in.git`
2. Open directory: `cd lem-in`
3. Start the app: `go run cmd/main.go examples/example00.txt`

If something doesn't work contact **Telegram**: **
@wqzxyl**

## Project structure:

1. **cmd/** - main.go file
2. **internal/**
   - **app** - orchestration layer
   - **model** - domain types
   - **parser** - input file parsing
   - **validation** - input validation
   - **service/pathfinding** - BFS-based path finding
   - **service/solver** - ant distribution and simulation
3. **examples/** - example input files

---

Made by Bakhtiyar.S