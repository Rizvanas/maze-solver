package problems

import (
	"errors"
	"fmt"
	"math"
)

type GraphProblem struct {
	Initial Node
	Goal    Node
}

// metodas gražinantis pradinę problemos būseną
func (p GraphProblem) GetInitialState() State {
	return p.Initial
}

// metodas gražinantis siekiamą problemos būseną
func (p GraphProblem) GetGoalState() State {
	return p.Goal
}

// metodas, gražinantis būseną, kuri gaunama pateiktai būsenai pritaikius pateiktą veiksmą
func (p GraphProblem) GetResultingState(state State, action Action) (State, error) {
	node, err := GraphNodeFromState(state)
	if err != nil {
		return nil, err
	}
	child, err := actionToConnection(node, action)
	return child, err
}

// metodas, gražinantis su pateikta būsena galimus atlikti veiksmus
func (p GraphProblem) GetPossibleActions(state State) ([]Action, error) {
	node, err := GraphNodeFromState(state)
	if err != nil {
		return nil, err
	}

	var actions []Action

	for _, connection := range node.Connections {
		action, err := connectionToAction(node, *connection)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}

	return actions, nil
}

func (p GraphProblem) GetActionCost(state State, action Action) (float64, error) {
	result, err := p.GetResultingState(state, action)
	if err != nil {
		return 0, err
	}

	node1, err := GraphNodeFromState(state)
	if err != nil {
		return 0, err
	}

	node2, err := GraphNodeFromState(result)
	if err != nil {
		return 0, err
	}

	return distanceBetweenNodes(node1, node2), nil
}

func (p GraphProblem) CalculateCostToGoal(state State) (float64, error) {
	current, err := GraphNodeFromState(state)
	if err != nil {
		return 0.0, nil
	}

	goalX, goalY := p.Goal.X, p.Goal.Y
	currentX, currentY := current.X, current.Y
	deltaX, deltaY := goalX-currentX, goalY-currentY

	return math.Sqrt(float64(deltaX*deltaX + deltaY*deltaY)), nil
}

// metodas, kuris konkretizuoja pateiktą būseną ir gražina grafo mazgą (jeigu galima)
func GraphNodeFromState(state State) (Node, error) {
	node, ok := state.(Node)
	if !ok {
		return node, fmt.Errorf("error in GraphProblem.GetpossibleActions takes Node, but received %T", state)
	}
	return node, nil
}

// connectionToAction - pagalbinė funkcija, gražinanti veiksmą,
// kurį naudojant galima gauti gretimą mazgą
func connectionToAction(node Node, connection Node) (Action, error) {
	switch {
	case connection.X == node.X && connection.Y < node.Y:
		return UP, nil
	case connection.X == node.X && connection.Y > node.Y:
		return DOWN, nil
	case connection.X < node.X && connection.Y == node.Y:
		return LEFT, nil
	case connection.X > node.X && connection.Y == node.Y:
		return RIGHT, nil
	default:
		return 0, errors.New("childToAction: child is in the wrong place")
	}
}

// actionToConnection - pagalbinė funkcija, gražinanti gretimą mazgą,
// kuris gaunamas einamajam mazgui pritaikius veiksmą
func actionToConnection(node Node, action Action) (Node, error) {
	for _, connection := range node.Connections {
		connectionPos, err := connectionToAction(node, *connection)
		if err != nil {
			return Node{}, err
		}
		if connectionPos == action {
			return *connection, nil
		}
	}
	return Node{}, errors.New("chould not find state that corresponds to your given action")
}

func distanceBetweenNodes(node1 Node, node2 Node) float64 {
	return math.Abs(float64(node1.X-node2.X)) + math.Abs(float64(node1.Y-node2.Y))
}
