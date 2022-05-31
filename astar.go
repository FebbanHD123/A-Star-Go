package main

import (
	"fmt"
	"gfx"
	"math"
	"math/rand"
	"time"
)

type Path struct {
	nodes []Node
}

type Node struct {
	gScore                     float64
	fScore                     float64
	x, y                       int
	predecessorX, predecessorY int
}

var obstacles []Node
var openList []Node
var closedList []Node

func main() {

	rand.Seed(time.Now().UnixMilli())
	for i := 0; i < 30000; i++ {
		node := Node{
			x: rand.Intn(800),
			y: rand.Intn(600),
		}
		obstacles = append(obstacles, node)
	}

	gfx.Fenster(800, 600)
	path, found := findPath(1, 1, rand.Intn(800), rand.Intn(600))
	fmt.Println("found")
	if found {
		gfx.Stiftfarbe(0, 0, 0)
		gfx.Cls()
		for _, node := range path.nodes {
			gfx.Stiftfarbe(0, 255, 0)
			gfx.Punkt(uint16(node.x), uint16(node.y))
		}
		for true {

		}
	} else {
		fmt.Println("No path found")
	}
}

func findPath(xStart, yStart, xGoal, yGoal int) (Path, bool) {
	goalNode := GetNode(xGoal, yGoal)
	startNode := GetNode(xStart, yStart)
	openList = append(openList, startNode)

	if IsNodeObstacle(goalNode) {
		fmt.Println("Goal is obstacle")
		return Path{}, false
	}

	for len(openList) > 0 {
		currentNode := RemoveMinFromOpenList()
		if currentNode.equals(goalNode) {
			return recustructPath(currentNode), true
		}
		closedList = append(closedList, currentNode)
		expandNode(&goalNode, &currentNode)

		gfx.UpdateAus()
		gfx.Stiftfarbe(0, 0, 0)
		gfx.Cls()
		gfx.Stiftfarbe(0, 255, 0)
		for _, node := range recustructPath(currentNode).nodes {
			gfx.Punkt(uint16(node.x), uint16(node.y))
		}
		for _, obstacle := range obstacles {
			gfx.Stiftfarbe(150, 0, 0)
			gfx.Punkt(uint16(obstacle.x), uint16(obstacle.y))
		}
		gfx.UpdateAn()
	}
	return Path{}, false

}

func expandNode(goalNode, node *Node) {
	for _, successor := range node.getNeighbors() {
		if IsNodeClosed(successor) {
			continue
		}
		tentativG := node.gScore + 1
		if tentativG > successor.gScore && successor.gScore >= 0 {
			continue
		}
		successor.SetCameFrom(*node)
		successor.gScore = tentativG
		successor.fScore = tentativG + successor.distanceTo(*goalNode)
		if !IsNodeOpen(successor) {
			openList = append(openList, successor)
		}
	}
}

func recustructPath(currentNode Node) Path {
	var nodes []Node
	var current = currentNode
	nodes = append(nodes, current)
	for {
		current = current.GetCameFrom()

		nodes = append(nodes[:1], nodes[0:]...)
		nodes[0] = current
		if current.predecessorX < 0 {
			break
		}
	}
	return Path{nodes: nodes}
}

func (n Node) equals(node Node) bool {
	return n.x == node.x && n.y == node.y
}

func (n Node) distanceTo(node Node) float64 {
	return math.Sqrt(math.Pow(float64(n.x-node.x), 2) + math.Pow(float64(n.y-node.y), 2))
}

func (n Node) getNeighbors() []Node {

	var neighbors []Node
	for dx := 0; dx < 3; dx++ {
		for dy := 0; dy < 3; dy++ {
			x := n.x + 1 - dx
			y := n.y + 1 - dy
			if x >= 800 || y >= 600 || x < 0 || y < 0 {
				continue
			}
			node := GetNode(x, y)
			if !IsNodeObstacle(node) && !IsNodeClosed(node) {
				neighbors = append(neighbors, node)
			}
		}
	}

	return neighbors
}

func GetNode(x, y int) Node {
	node := Node{
		gScore:       -1,
		fScore:       -1,
		x:            x,
		y:            y,
		predecessorX: -1,
	}
	for _, n := range openList {
		if n.equals(node) {
			return n
		}
	}
	for _, n := range closedList {
		if n.equals(node) {
			return n
		}
	}
	return node
}

func IsNodeClosed(node Node) bool {
	for _, n := range closedList {
		if n.equals(node) {
			return true
		}
	}
	return false
}

func IsNodeOpen(node Node) bool {
	for _, n := range openList {
		if n.equals(node) {
			return true
		}
	}
	return false
}

func IsNodeObstacle(node Node) bool {
	for _, n := range obstacles {
		if n.equals(node) {
			return true
		}
	}
	return false
}

func RemoveMinFromOpenList() Node {
	min := openList[0]
	var index int
	for i, node := range openList {
		if node.fScore >= 0 && node.fScore < min.fScore {
			min = node
			index = i
		}
	}
	openList = append(openList[:index], openList[index+1:]...)
	return min
}

func (n *Node) SetCameFrom(node Node) {
	n.predecessorX = node.x
	n.predecessorY = node.y
}

func (n Node) GetCameFrom() Node {
	return GetNode(n.predecessorX, n.predecessorY)
}
