package problems

import (
	"errors"
	"fmt"
)

type Action byte

const (
	UP    Action = 1
	DOWN  Action = 2
	LEFT  Action = 3
	RIGHT Action = 4
)

type State interface {
	Describe() string
}

// Būsenų erdvės paieškos uždavinio interfeisas
type Problem interface {
	GetInitialState() State
	GetGoalState() State
	GetResultingState(state State, action Action) (State, error)
	GetPossibleActions(state State) ([]Action, error)
}

// Labirinto grafo mazgas, kuris įgyvendina State interfeisą
type Node struct {
	X, Y        int
	Connections []*Node
}

func (n Node) Describe() string {
	return fmt.Sprintf("%d;%d", n.X, n.Y)
}

// Problem interfeisą implementinanti struktūra
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

// metodas, kuris konkretizuoja pateiktą būseną ir gražina grafo mazgą (jeigu galima)
func GraphNodeFromState(state State) (Node, error) {
	node, ok := state.(Node)
	if !ok {
		return node, fmt.Errorf("error in GraphProblem.GetpossibleActions takes Node, but received %T", state)
	}
	return node, nil
}
