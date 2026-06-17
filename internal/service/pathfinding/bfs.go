package pathfinding

import (
	"errors"

	"01.tomorrow-school.ai/git/msabyrga/lem-in.git/internal/model"
)

// FindRightPaths finds vertex-disjoint paths using node splitting + Edmonds-Karp.
func FindRightPaths(input *model.Input) ([][]string, error) {
	start := input.StartRoom.Name
	end := input.EndRoom.Name

	in := func(n string) string { return n + "_in" }
	out := func(n string) string { return n + "_out" }

	// residual[u][v] = remaining capacity
	residual := make(map[string]map[string]int)

	add := func(u, v string, cap int) {
		if residual[u] == nil {
			residual[u] = make(map[string]int)
		}
		if residual[v] == nil {
			residual[v] = make(map[string]int)
		}
		residual[u][v] += cap
		// reverse edge (capacity 0 initially)
		if _, ok := residual[v][u]; !ok {
			residual[v][u] = 0
		}
	}

	// Collect all room names
	allRooms := make(map[string]bool)
	allRooms[start] = true
	allRooms[end] = true
	for room, neighbors := range input.Links {
		allRooms[room] = true
		for _, n := range neighbors {
			allRooms[n] = true
		}
	}

	// Node splitting: room_in -> room_out with capacity 1
	// (start and end get capacity = large number to allow all ants)
	for room := range allRooms {
		if room == start || room == end {
			add(in(room), out(room), len(input.Links)+1)
		} else {
			add(in(room), out(room), 1)
		}
	}

	// Edges: room_out -> neighbor_in with capacity 1
	for room, neighbors := range input.Links {
		for _, neighbor := range neighbors {
			add(out(room), in(neighbor), 1)
		}
	}

	source := in(start)
	sink := out(end)

	// Edmonds-Karp: BFS to find augmenting paths
	for {
		// BFS
		prev := make(map[string]string)
		visited := map[string]bool{source: true}
		queue := []string{source}
		found := false

		for len(queue) > 0 && !found {
			curr := queue[0]
			queue = queue[1:]

			// sort neighbors for determinism (optional but helps)
			for next, cap := range residual[curr] {
				if cap > 0 && !visited[next] {
					visited[next] = true
					prev[next] = curr
					if next == sink {
						found = true
						break
					}
					queue = append(queue, next)
				}
			}
		}

		if !found {
			break
		}

		// Find bottleneck
		flow := int(^uint(0) >> 1) // max int
		cur := sink
		for cur != source {
			p := prev[cur]
			if residual[p][cur] < flow {
				flow = residual[p][cur]
			}
			cur = p
		}

		// Augment
		cur = sink
		for cur != source {
			p := prev[cur]
			residual[p][cur] -= flow
			residual[cur][p] += flow
			cur = p
		}
	}

	// Extract paths by following used edges (where forward capacity decreased)
	// We detect used edges: out(room) -> in(neighbor) where residual < 1 (i.e., was used)
	usedEdge := make(map[string]map[string]bool)
	for room, neighbors := range input.Links {
		for _, neighbor := range neighbors {
			// edge was used if forward cap = 0 (started at 1)
			if residual[out(room)][in(neighbor)] == 0 {
				if usedEdge[room] == nil {
					usedEdge[room] = make(map[string]bool)
				}
				usedEdge[room][neighbor] = true
			}
		}
	}

	// Walk paths
	var paths [][]string
	for {
		if len(usedEdge[start]) == 0 {
			break
		}

		path := []string{start}
		curr := start
		found := false

		for curr != end {
			moved := false
			for next, used := range usedEdge[curr] {
				if used {
					usedEdge[curr][next] = false
					path = append(path, next)
					curr = next
					moved = true
					break
				}
			}
			if !moved {
				break
			}
			if curr == end {
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
