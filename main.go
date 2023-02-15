package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type room struct {
	name       string
	neighbours []room
	visited    bool
	x          string
	y          string
}

var rooms []room
var paths []string
var startRoom room
var endRoom room

func isRoom(line string) bool {
	temp := strings.Split(line, " ")
	if len(temp) != 3 {
		return false
	}
	_, err1 := strconv.Atoi(temp[1])
	_, err2 := strconv.Atoi(temp[2])
	if err1 != nil || err2 != nil {
		return false
	}
	return true
}

func isLink(line string) bool {
	temp := strings.Split(line, "-")
	if len(temp) != 2 {
		return false
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("USAGE: go run lem-in.go <filename.txt>")
	} else {
		reader, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
		}
		farm := string(reader)
		lines := strings.Split(farm, "\n")
		start := ""
		end := ""
		n := 0

		// Parsing the text file into a graph (network of nodes)
		for i, line := range lines {
			switch {
			case line == "##start":
				if i > 0 { // nr of ants always above ##start
					n, err = strconv.Atoi(lines[i-1])
					if err != nil {
						fmt.Println("Error fetching number of ants.")
						return // call some generic "invalid input function"?
					}
				}
				// Find next valid room
				for j := i + 1; j < len(lines); j++ {
					if isRoom(lines[j]) {
						start = strings.Split(lines[j], " ")[0]
					}
				}
				if start == "" {
					fmt.Println("No start.")
					return
				}
			case line == "##end":
				// Find next valid room
				if i < len(lines)-3 { // end can not be at last line, must be X amout of links after end
					for j := i + 1; j < len(lines); j++ {
						if isRoom(lines[j]) {
							end = strings.Split(lines[j], " ")[0]
						}
					}
					if end == "" {
						fmt.Println("No end.")
						return
					}
				} else {
					fmt.Println("Invalid map format") // call error function?
					return
				}
			case isRoom(line):
				// add room to graph
				// also validate temp for correct room format? no invalid characters etc
				temp := strings.Split(line, " ")
				var r room
				r.name = temp[0]
				r.x = temp[1]
				r.y = temp[2]
				r.visited = false
				rooms = append(rooms, r)
			case isLink(line):
				// add connection between nodes
				temp := strings.Split(line, "-")
				for _, room := range rooms {
					if room.name == temp[0] {
						for _, sroom := range rooms {
							if sroom.name == temp[1] {
								room.neighbours = append(room.neighbours, sroom)
							}
						}
					} else if room.name == temp[1] {
						for _, sroom := range rooms {
							if sroom.name == temp[0] {
								room.neighbours = append(room.neighbours, sroom)
							}
						}
					}
				}
			}
		}
		findStart(start)
		findEnd(end)
		DFS(startRoom, endRoom)
		m := intoMap(simplePaths)
		combinationsGraph(m)
		findNIPC(allCombinations[0])
		// Fastest path combination
		FPC := findFastest(NIPC)
		distributeAnts(FPC, n)
	}
}

// just a comment

type fastPath struct {
	itenerary List
	capacity  int
	ants      []ant
}

type fastPathRoom struct {
	name string
	ant  ant
	next *fastPathRoom
}

type ant struct {
	num int
}

type List struct {
	head *fastPathRoom
}

func (l *List) Insert(name string) {
	list := &fastPathRoom{name: name, next: nil}
	if l.head == nil {
		l.head = list
	} else {
		p := l.head
		for p.next != nil {
			p = p.next
		}
		p.next = list
	}
}

func Show(l *List) {
	p := l.head
	for p != nil {
		fmt.Printf("-> %v ", p.name)
		p = p.next
	}
}

var fastPaths []fastPath

func distributeAnts(FPC [][]string, n int) {
	var ants []ant
	for i := 1; i <= n; i++ {
		var a ant
		a.num = i
		ants = append(ants, a)
	}
	for _, p := range FPC {
		var temp fastPath
		temp.itenerary = List{}
		for _, r := range p {
			temp.itenerary.Insert(r)
		}
		temp.capacity = len(p)
		fastPaths = append(fastPaths, temp)
	}
	for i := range ants {
		var fastest fastPath
		fastest.capacity = 0
		for _, p := range fastPaths {
			if fastest.capacity == 0 {
				fastest = p
			} else {
				if p.capacity+len(p.ants) < len(fastest.ants)+fastest.capacity {
					fastest = p
				}
			}
		}
		for _, p := range fastPaths {
			if p.itenerary == fastest.itenerary {
				p.ants = append(p.ants, ants[i])
			}
		}
	}
}

var moves string

func steps(fastPaths []fastPath, ants int) {
	for ants > 0 {
		for _, p := range fastPaths {
			pl := p.itenerary.head
			if len(p.ants) > 0 {
				addAnt(p, pl, p.ants[0], &ants)
				if len(p.ants) > 1 {
					p.ants = p.ants[1:]
				} else {
					p.ants = nil
				}
			} else {
				pushAnts(p, pl, &ants)
			}
		}

	}
}

func pushAnts(p fastPath, pl *fastPathRoom, n *int) {
	if pl.ant.num == 0 && pl.next != nil {
		pl = pl.next
		pushAnts(p, pl, n)
	}
	if pl.ant.num == 0 && pl.next == nil {
		return
	}
	if pl.next != nil && pl.next.ant.num == 0 {
		pl.next.ant = pl.ant
		pl.ant.num = 0
		pl = pl.next
		moves += "L" + strconv.Itoa(pl.ant.num) + "-" + pl.next.name + " "

		pushAnts(p, pl, n)
	}
	if pl.ant.num != 0 && pl.next == nil {
		pl.ant.num = 0
		*n--
		return
	}
	if pl.ant.num != 0 && pl.next != nil {
		pl = pl.next
		pushAnts(p, pl, n)
	}
}

func addAnt(p fastPath, pl *fastPathRoom, ant ant, n *int) {
	if pl.ant.num == 0 {
		pl.ant = ant
		moves += "L" + strconv.Itoa(ant.num) + "-" + pl.name + " "
		return
	}
	if pl.next == nil {
		pl.ant.num = 0
		*n--
		addAnt(p, p.itenerary.head, ant, n)
	}
	for pl != nil {
		addAnt(p, pl, ant, n)
	}
	return
}

func findStart(start string) {
	for _, r := range rooms {
		if r.name == start {
			startRoom = r
		}
	}
}

func findEnd(end string) {
	for _, r := range rooms {
		if r.name == end {
			endRoom = r
		}
	}
}

var currentPath []string
var simplePaths [][]string

func DFS(current room, end room) {
	if current.visited == true {
		return
	}
	current.visited = true
	currentPath = append(currentPath, current.name)
	if current.name == end.name && current.x == end.x && current.y == end.y {
		simplePaths = append(simplePaths, currentPath)
		current.visited = false
		currentPath = currentPath[:(len(currentPath) - 1)]
		return
	}
	for _, adj := range current.neighbours {
		DFS(adj, end)
	}
	currentPath = currentPath[:(len(currentPath) - 1)]
	current.visited = false
	return
}

func intoMap(simplePaths [][]string) map[string][][]string {
	// Remove starting point
	for _, path := range simplePaths {
		path = path[1:]
	}
	// Make map of each starting point next to start
	m := make(map[string][][]string)
	for _, p := range simplePaths {
		if _, found := m[p[0]]; found == true {
			m[p[0]] = append(m[p[0]], p)
		} else {
			pr := [][]string{p}
			m[p[0]] = pr
		}
	}
	return m
}

type path struct {
	source        string
	itenerary     []string
	compared      bool
	friends       []path
	friendsSource string
}

var allCombinations []path

func combinationsGraph(m map[string][][]string) {

	for k, v := range m {
		for _, e := range v {
			var p path
			p.source = k
			p.itenerary = e
			p.compared = false
			allCombinations = append(allCombinations, p)
		}
	}
	for _, p := range allCombinations {
		for k := range m {
			if p.source == k {
				p.friendsSource = nextKey(k, m)
			}
		}
	}
	for _, p := range allCombinations {
		for _, op := range allCombinations {
			if p.friendsSource == op.source {
				p.friends = append(p.friends, op)
			}
		}
	}
}

func nextKey(currentKey string, m map[string][][]string) string {
	found := false
	for k := range m {
		if found == true {
			return k
		}
		if k == currentKey {
			found = true
		}
	}
	return ""
}

// Non intersecting path combinations
var NIPC [][][]string
var combo [][]string

func findNIPC(current path) {
	if current.compared == true {
		return
	}
	current.compared = true
	combo = append(combo, current.itenerary)
	if current.friendsSource == "" {
		NIPC = append(NIPC, combo)
		current.compared = false
		combo = combo[:(len(combo) - 1)]
		return
	}
	for _, adj := range current.friends {
		findNIPC(adj)
	}
	combo = combo[:(len(combo) - 1)]
	current.compared = false
	return
}

// Equal tells whether a and b contain the same elements.
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func areIntersected(paths [][]string, other_path []string) bool {
	for _, path := range paths {
		for i := 0; i < len(path); i++ {
			for j := 0; j < len(other_path); j++ {
				for k := 0; k < len(path[i]); k++ {
					for l := 0; l < len(other_path[j]); l++ {
						if path[i][k] == other_path[j][l] {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

var sortedNIPC map[float64][][]string

func findFastest(NIPC [][][]string) [][]string {
	var average float64
	for _, combo := range NIPC {
		average = lenAll(combo) / float64(len(NIPC))
		sortedNIPC[average] = combo
	}
	averages := make([]float64, 0, len(sortedNIPC))
	for k := range sortedNIPC {
		averages = append(averages, k)
	}
	sort.Float64s(averages)
	fastest := sortedNIPC[averages[0]]
	return fastest
}

func lenAll(combo [][]string) float64 {
	var l float64
	l = 0
	for _, path := range combo {
		l += float64(len(path))
	}
	return l
}
