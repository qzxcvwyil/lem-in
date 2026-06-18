package solver

import (
	"fmt"
	"strings"
)

type ant struct {
	id       int
	path     []string
	position int
	delay    int
}

func Solve(input [][]string, ants int) (string, error) {
	distributedAnts := distribute(input, ants)

	allAnts := make([]ant, 0)
	antID := 1

	for i, count := range distributedAnts {
		for j := range count {
			allAnts = append(allAnts, ant{
				id:       antID,
				path:     input[i],
				position: 1,
				delay:    j,
			})
			antID++
		}
	}

	var sb strings.Builder
	done := 0
	turn := 0

	for {
		if done == len(allAnts) {
			break
		}

		first := true
		for i := range allAnts {
			if turn < allAnts[i].delay {
				continue
			}

			if allAnts[i].position >= len(allAnts[i].path) {
				continue
			}

			if !first {
				sb.WriteRune(' ')
			}

			fmt.Fprintf(&sb, "L%d-%s", allAnts[i].id, allAnts[i].path[allAnts[i].position])
			first = false
			allAnts[i].position++

			if allAnts[i].position >= len(allAnts[i].path) {
				done++
			}
		}
		turn++

		sb.WriteRune('\n')
	}

	return sb.String(), nil
}

func distribute(paths [][]string, ants int) []int {
	assigned := make([]int, len(paths))

	for range ants {
		bestPath := 0
		for i := range paths {
			if len(paths[i])+assigned[i] < len(paths[bestPath])+assigned[bestPath] {
				bestPath = i
			}
		}
		assigned[bestPath]++
	}

	return assigned
}
