package dijkstra

import (
	"fmt"
	"sort"
)

type Node interface {
	Neighbors() ([]Node, []int)
}

func GetShortestPath(g []Node, start Node, target Node) []Node {
	distances := make(map[Node]int, 0)
	prev := make(map[Node]Node, 0)
	q := make([]Node, 0)
	qmap := make(map[Node]struct{})

	for _, node := range g {
		distances[node] = 10000
		prev[node] = nil
		q = append(q, node)
		qmap[node] = struct{}{}
	}
	distances[start] = 0

	for len(q) > 0 {
		sort.Slice(q, func(i, j int) bool {
			return distances[q[i]] < distances[q[j]]
		})

		u := q[0]
		if u == target {
			path := make([]Node, 0)
			if (prev[u] != nil) || (u == start) {
				for u != nil {
					path = append([]Node{u}, path...)
					u = prev[u]
				}
			}
			if len(path) > 0 {
				fmt.Printf("found path starting at %v\n", path[0])
			} else {
				fmt.Printf("No path\n")
			}
			return path
		}
		q = q[1:]
		delete(qmap, u)
		neighbors, ndistances := u.Neighbors()
		for i, n := range neighbors {
			if _, ok := qmap[n]; ok {
				alt := distances[u] + ndistances[i]
				if alt < distances[n] {
					distances[n] = alt
					prev[n] = u
				}
			}
		}
	}

	return nil
}

func GetShortestDistances(g []Node, start Node) (map[Node]int, map[Node]Node) {
	distances := make(map[Node]int, 0)
	prev := make(map[Node]Node, 0)
	q := make([]Node, 0)
	qmap := make(map[Node]struct{})

	for _, node := range g {
		distances[node] = 10000
		prev[node] = nil
		q = append(q, node)
		qmap[node] = struct{}{}
	}
	distances[start] = 0

	for len(q) > 0 {
		sort.Slice(q, func(i, j int) bool {
			return distances[q[i]] < distances[q[j]]
		})

		u := q[0]
		q = q[1:]
		delete(qmap, u)
		neighbors, ndistances := u.Neighbors()
		for i, n := range neighbors {
			if _, ok := qmap[n]; ok {
				alt := distances[u] + ndistances[i]
				if alt < distances[n] {
					distances[n] = alt
					prev[n] = u
				}
			}
		}
	}

	return distances, prev
}
