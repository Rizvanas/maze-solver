package algos

import (
	"errors"

	"example.com/maze-solver/pkg/problems"
)

func SearchStateSpace(problem problems.Problem, algo Algorithm) (map[string]problems.State, error) {
	fringe := []problems.State{}
	fringe = append(fringe, problem.GetInitialState())
	explored := map[string]problems.State{}
	explored[problem.GetInitialState().Describe()] = nil

	solved := false
	for updatedFringe, current, err := takeFromFringe(fringe, algo); !solved; updatedFringe, current, err = takeFromFringe(fringe, algo) {
		fringe = updatedFringe

		if err != nil {
			return nil, err
		}

		actions, err := problem.GetPossibleActions(current)
		if err != nil {
			return nil, err
		}

		for _, action := range actions {
			resultingState, err := problem.GetResultingState(current, action)

			if err != nil {
				return nil, err
			}

			if resultingState.Describe() == problem.GetGoalState().Describe() {
				explored[resultingState.Describe()] = current
				solved = true
				break
			}

			// patikrini ar gauta busena jau buvo aplankyta
			if _, visited := explored[resultingState.Describe()]; !visited {
				fringe = append(fringe, resultingState)
				explored[resultingState.Describe()] = current
			}
		}
	}
	return explored, nil
}

func GraphDFS(problem problems.Problem) ([]*problems.Node, error) {
	explored, err := SearchStateSpace(problem, DEPTH_FIRST_SEARCH)
	if err != nil {
		return nil, err
	}

	solution, err := getSolutionNodes(problem, explored)
	if err != nil {
		return nil, err
	}

	return solution, nil
}

func GraphBFS(problem problems.Problem) ([]*problems.Node, error) {
	explored, err := SearchStateSpace(problem, BREADTH_FIRST_SEARCH)
	if err != nil {
		return nil, err
	}

	solution, err := getSolutionNodes(problem, explored)
	if err != nil {
		return nil, err
	}

	return solution, nil
}

func Graph_DJIKSTRA(problem problems.Problem) {
	panic("Graph_DJIKSTRA not implemented")
}

func getSolutionNodes(problem problems.Problem, explored map[string]problems.State) ([]*problems.Node, error) {
	solution := make([]*problems.Node, 0)
	current := problem.GetGoalState()
	ok := false

	for current.Describe() != problem.GetInitialState().Describe() {
		node, err := problems.GraphNodeFromState(current)
		if err != nil {
			return nil, err
		}
		solution = append([]*problems.Node{&node}, solution...)
		current, ok = explored[current.Describe()]
		if !ok {
			return nil, errors.New("getSolutionNodes: something went wrong")
		}
	}

	node, err := problems.GraphNodeFromState(current)
	if err != nil {
		return nil, err
	}
	solution = append([]*problems.Node{&node}, solution...)

	return solution, nil
}

type Algorithm byte

const (
	DEPTH_FIRST_SEARCH   Algorithm = 1
	BREADTH_FIRST_SEARCH Algorithm = 2
	ASTAR_SEARCH         Algorithm = 3
)

func addToFringe(fringe *[]problems.State, state problems.State) []problems.State {
	return append(*fringe, state)
}

func takeFromFringe(fringe []problems.State, algo Algorithm) ([]problems.State, problems.State, error) {
	l := len(fringe)
	if l == 0 {
		return nil, nil, errors.New("fringe is empty")
	}

	switch algo {
	case DEPTH_FIRST_SEARCH:
		return fringe[:l-1], fringe[l-1], nil
	case BREADTH_FIRST_SEARCH:
		return fringe[1:], fringe[0], nil
	default:
		return nil, nil, errors.New("ASTAR_SEARCH fringe unimplemented")
	}
}

func peekFringe(fringe []problems.State, algo Algorithm) (problems.State, error) {
	l := len(fringe)
	if l == 0 {
		return nil, errors.New("fringe is empty")
	}

	switch algo {
	case DEPTH_FIRST_SEARCH:
		return fringe[l-1], nil
	case BREADTH_FIRST_SEARCH:
		return fringe[0], nil
	default:
		return nil, errors.New("ASTAR_SEARCH fringe not yet implemented")
	}
}
