package main

import (
	"flag"
	"log"

	"example.com/maze-solver/pkg/algos"
	"example.com/maze-solver/pkg/maze"
)

// programos entry point (įeigos taškas)
func main() {
	mazeName := flag.String("maze", "small.png", "specify maze name")
	searchAlgo := flag.String("algo", "astar", "specify search algorithm")
	flag.Parse()

	algorithm, err := algos.AlgoFromString(*searchAlgo)
	if err != nil {
		log.Fatal(err)
	}

	maze := maze.New(*mazeName)
	maze.Solve(algorithm)
	maze.SaveToFile("maze_solution.png")
}
