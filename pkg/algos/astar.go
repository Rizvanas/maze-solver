package algos

import (
	"container/heap"

	"example.com/maze-solver/pkg/problems"
)

func Astar(problem problems.Problem) ([]*problems.Node, error) {
	initial := problem.GetInitialState()
	current := StateCost{State: initial}
	actionCost := 0.0
	costToGoal, err := problem.CalculateCostToGoal(initial)
	if err != nil {
		return nil, err
	}
	explored := map[string]problems.State{}
	current.TotalCost = costToGoal
	current.PathCost = 0.0
	statesWithCost := []*StateCost{&current}
	fringe := make(PriorityQueue, len(statesWithCost))

	for i, killMe := range statesWithCost {
		fringe[i] = killMe
		fringe[i].Index = i
	}

	heap.Init(&fringe)

	explored[current.State.Describe()] = nil

	for solved := false; !solved; {
		next := heap.Pop(&fringe).(*StateCost)
		current = *next

		actions, err := problem.GetPossibleActions(current.State)
		if err != nil {
			return nil, err
		}

		for _, action := range actions {
			resultingState, err := problem.GetResultingState(current.State, action)
			if err != nil {
				return nil, err
			}

			actionCost, err = problem.GetActionCost(current.State, action)
			if err != nil {
				return nil, err

			}
			costToGoal, err = problem.CalculateCostToGoal(resultingState)
			if err != nil {
				return nil, err
			}

			if resultingState.Describe() == problem.GetGoalState().Describe() {
				explored[resultingState.Describe()] = current.State
				solved = true
				break
			}

			if _, visited := explored[resultingState.Describe()]; !visited {
				explored[resultingState.Describe()] = current.State

				pathCost := current.PathCost + actionCost
				state := StateCost{
					State:     resultingState,
					TotalCost: costToGoal + pathCost,
					PathCost:  pathCost,
				}
				heap.Push(&fringe, &state)
			} else if index := IndexOfState(&statesWithCost, resultingState); index != -1 {
				// update pq with new cost
				pathCost := current.PathCost + actionCost
				statesWithCost[index].PathCost = pathCost
				statesWithCost[index].TotalCost = costToGoal + pathCost

				for i, killMe := range statesWithCost {
					fringe[i] = killMe
					fringe[i].Index = i
				}

				heap.Init(&fringe)
			}
		}
	}

	solution, err := getSolutionNodes(problem, explored)
	if err != nil {
		return nil, err
	}

	return solution, nil
}

// patikrina ar prioritetine eile turi elementa
func IndexOfState(statesWithCost *[]*StateCost, state problems.State) int {
	for i, sc := range *statesWithCost {
		if sc.State.Describe() == state.Describe() {
			return i
		}
	}
	return -1
}
