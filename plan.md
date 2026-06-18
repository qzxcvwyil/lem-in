# Implementation Plan

This document created for project planning purposes. Check project description in [Readme](README.md)

## Phase 1

### Input Parsing

#### We have a parts in the input that we need to get

1. Number of ants
2. Start marker (##start)
3. Room Number, Coordinate X, Coordinate Y (1 23 3)
4. Comment marker (#comment)
5. End marker (##end)
6. Links between rooms (0-4)

### Input Validation

#### Notes from Instruction

- Some will have rooms that link to themselves, sending your path-search spinning in circles.
- Some will have too many/too few ants, no ##start or ##end
- Duplicated rooms, links to unknown rooms
- Rooms with invalid coordinates and a variety of other invalid or poorly-formatted input.

_Specific error message's **example: ERROR: invalid data format, invalid number of Ants or ERROR: invalid data format, no start room found**._

## Phase 2

### Check all valid paths with BFS

### Exclude any paths without path to finish

## Phase 3

### Calculate how many ants should go into valid shortest path

### Make an output with result

## Proposed Structure's

### Madina

```
❯ tree .
.
├── README.md
├── cmd
│   └── lem-in
│       └── main.go             // entrypoint 
├── go.mod
├── internal
│   ├── app
│   │   └── app.go              // App orchestration
│   ├── model                   // Domain models
│   │   └── input.go 
│   ├── parser                  // Input parsing and turning to domain 
│   │   └── parser.go
│   ├── service                 // services for solving
│   │   ├── pathfinding
│   │   │   └── bfs.go
│   │   └── solver
│   │       └── solver.go
│   └── validation              // Input validation, links logic check 
│       └── validation.go
└── plan.md

11 directories, 10 files
```

### Adam

```
lem-in/
├── cmd/
│   └── lem-in/
│       └── main.go          // entrypoint: args -> app.Run -> print error/result
├── internal/
│   ├── app/
│   │   └── run.go           // orchestration: parse -> validate -> solve -> simulate
│   │
│   ├── domain/
│   │   ├── farm.go          // Room, Farm
│   │   ├── path.go          // Path, Plan
│   │   └── move.go          // Move
│   │
│   ├── parser/
│   │   ├── parser.go        // Parse(lines) -> Farm
│   │   └── validate.go      // semantic/input validation
│   │
│   ├── graph/
│   │   ├── bfs.go           // shortestDistance, distancesToEnd
│   │   └── paths.go         // candidate path search
│   │
│   ├── solver/
│   │   ├── selector.go      // choose compatible path set
│   │   └── distribute.go    // ants distribution across paths
│   │
│   ├── simulator/
│   │   └── simulate.go      // turn-by-turn movement generation
│   │
│   └── input/
│       └── printer.go       // raw input + moves formatting/output
│
├── pkg/
│   └── errs/
│       └── invalid.go       // ERROR: invalid data format...
│
└── README.md
```
