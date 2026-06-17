package pathfinding

import (
	"errors"

	"01.tomorrow-school.ai/git/msabyrga/lem-in.git/internal/model"
)

// FindRightPaths finds all non-overlapping paths from start to end
// using a flow-based BFS (Edmonds-Karp style).
func FindRightPaths(input *model.Input) ([][]string, error) {
	start := input.StartRoom.Name
	end := input.EndRoom.Name

	// Build residual graph
	residual := make(map[string]map[string]bool)
	for room, neighbors := range input.Links {
		if residual[room] == nil {
			residual[room] = make(map[string]bool)
		}
		for _, neighbor := range neighbors {
			residual[room][neighbor] = true
			if residual[neighbor] == nil {
				residual[neighbor] = make(map[string]bool)
			}
		}
	}

	// Repeatedly find augmenting paths and flip edges
	for {
		cameFrom := make(map[string]string)
		visited := make(map[string]bool)
		queue := []string{start}
		visited[start] = true
		found := false

		for len(queue) > 0 {
			room := queue[0]
			queue = queue[1:]

			if room == end {
				found = true
				break
			}

			for neighbor, available := range residual[room] {
				if available && !visited[neighbor] {
					visited[neighbor] = true
					cameFrom[neighbor] = room
					queue = append(queue, neighbor)
				}
			}
		}

		if !found {
			break
		}

		// Augment along the found path
		current := end
		for current != start {
			prev := cameFrom[current]
			residual[prev][current] = false
			residual[current][prev] = true
			current = prev
		}
	}

	// Extract used edges: forward edge was used if it's now false and reverse is true
	used := make(map[string]map[string]bool)
	for room, neighbors := range input.Links {
		for _, neighbor := range neighbors {
			if !residual[room][neighbor] && residual[neighbor][room] {
				if used[room] == nil {
					used[room] = make(map[string]bool)
				}
				used[room][neighbor] = true
			}
		}
	}

	// Walk paths from start to end using used edges
	var paths [][]string
	for {
		path := []string{start}
		current := start
		found := false

		for current != end {
			moved := false
			for neighbor, isUsed := range used[current] {
				if isUsed {
					path = append(path, neighbor)
					used[current][neighbor] = false
					current = neighbor
					moved = true
					break
				}
			}
			if !moved {
				break
			}
			if current == end {
				found = true
				break
			}
		}

		if !found {
			break
		}
		paths = append(paths, path)
	}

	if len(paths) == 0 {
		return nil, errors.New("no valid paths found")
	}

	return paths, nil
}
