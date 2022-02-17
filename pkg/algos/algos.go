package algos

import (
	"errors"

	"example.com/maze-solver/pkg/problems"
)

func GraphDFS(problem problems.Problem) ([]*problems.Node, error) {
	fringe, explored, err := DFS(problem)
	if err != nil {
		return nil, err
	}

	solution, err := getSolutionNodes(explored, fringe.GetArray())
	if err != nil {
		return nil, err
	}

	return solution, nil
}

func DFS(problem problems.Problem) (stack, map[string]int, error) {
	fringe := make(stack, 0)
	fringe = fringe.Push(problem.GetInitialState())
	explored := map[string]int{}

	goalReached := false
	for currentState, ok := fringe.Peek(); !goalReached; currentState, ok = fringe.Peek() {
		if !ok {
			return nil, nil, errors.New("no solution could be found")
		}

		actions, err := problem.GetPossibleActions(currentState)
		if err != nil {
			return nil, nil, err
		}

		explored[currentState.Describe()] = 1

		visitedCount := 0
		for _, action := range actions {
			resultingState, err := problem.GetResultingState(currentState, action)
			if err != nil {
				return nil, nil, err
			}

			if resultingState.Describe() == problem.GetGoalState().Describe() {
				explored[resultingState.Describe()] = 1
				fringe = fringe.Push(resultingState)
				goalReached = true
				break
			}

			if _, visited := explored[resultingState.Describe()]; !visited {
				fringe = fringe.Push(resultingState)
			} else {
				visitedCount++
			}

			if visitedCount == len(actions) {
				fringe, _, ok = fringe.Pop()
				if !ok {
					break
				}
			}
		}
	}

	return fringe, explored, nil
}

func GraphBFS(problem problems.Problem) {
	panic("GraphBFS not implemented")
}

func Graph_DJIKSTRA(problem problems.Problem) {
	panic("Graph_DJIKSTRA not implemented")
}

type stack []problems.State

func (s stack) Push(value problems.State) []problems.State {
	return append(s, value)
}

func (s stack) Pop() (stack, problems.State, bool) {
	l := len(s)
	if l == 0 {
		return nil, nil, false
	}
	return s[:l-1], s[l-1], true
}

func (s stack) Peek() (problems.State, bool) {
	l := len(s)
	if l == 0 {
		return nil, false
	}
	return s[l-1], true
}

func (s stack) GetArray() []problems.State {
	return s
}

func getSolutionNodes(explored map[string]int, statePath []problems.State) ([]*problems.Node, error) {
	solution := make([]*problems.Node, 0)
	for _, state := range statePath {
		if _, ok := explored[state.Describe()]; !ok {
			continue
		}
		node, err := problems.GraphNodeFromState(state)
		if err != nil {
			return nil, err
		}
		solution = append(solution, &node)
	}

	return solution, nil
}
