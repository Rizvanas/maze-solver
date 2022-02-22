package problems

// veiksmo enum
type Action byte

const (
	UP    Action = 1
	DOWN  Action = 2
	LEFT  Action = 3
	RIGHT Action = 4
)

// Būsenos interfeisas
type State interface {
	Describe() string
}

// Būsenų erdvės paieškos uždavinio interfeisas
type Problem interface {
	GetInitialState() State
	GetGoalState() State
	GetResultingState(state State, action Action) (State, error)
	GetActionCost(state State, action Action) (float64, error)
	GetPossibleActions(state State) ([]Action, error)
	CalculateCostToGoal(state State) (float64, error)
}
