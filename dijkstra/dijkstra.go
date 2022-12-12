package dijkstra

import (
	"fmt"
	"sort"
)

type Node interface {
	Neighbors() ([]Node, []int)
}

/*
function Dijkstra(Graph, source):
 2
 3      for each vertex v in Graph.Vertices:
 4          dist[v] ← INFINITY
 5          prev[v] ← UNDEFINED
 6          add v to Q
 7      dist[source] ← 0
 8
 9      while Q is not empty:
10          u ← vertex in Q with min dist[u]
11          remove u from Q
12
13          for each neighbor v of u still in Q:
14              alt ← dist[u] + Graph.Edges(u, v)
15              if alt < dist[v]:
16                  dist[v] ← alt
17                  prev[v] ← u
18
19      return dist[], prev[]
*/

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

/*
1  S ← empty sequence
2  u ← target
3  if prev[u] is defined or u = source:          // Do something only if the vertex is reachable
4      while u is defined:                       // Construct the shortest path with a stack S
5          insert u at the beginning of S        // Push the vertex onto the stack
6          u ← prev[u]                           // Traverse from target to source
*/
