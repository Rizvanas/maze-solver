package algos

import (
	"errors"

	"example.com/maze-solver/pkg/problems"
)

// searchStateSpace - pagrindinė paieškos būsenų erdvėje funkcija.
// Priima:
//   1. Problem interfeisą,
//   2. Algoritmo kurį naudos spendimui enumą
// Gražina: aplankytų būsenų medį, išreikštą hashmap DS
func searchStateSpace(problem problems.Problem, algo Algorithm) (map[string]problems.State, error) {
	// Inicializuojamas paieškos masyvas
	fringe := []problems.State{}
	// Įdedama pradinė būsena
	fringe = append(fringe, problem.GetInitialState())

	// Inicializuojamas aplankytų būsenų hasmapas
	explored := map[string]problems.State{}
	explored[problem.GetInitialState().Describe()] = nil

	solved := false
	// iš paieškos masyvo, nuo paieškos algoritmo priklausančia tvarka, paimama būsena
	for updatedFringe, current, err := takeFromFringe(fringe, algo); !solved; updatedFringe, current, err = takeFromFringe(fringe, algo) {
		fringe = updatedFringe

		if err != nil {
			return nil, err
		}

		// gaunamas visų su gauta būsena galimų veiksmų sąrašas
		actions, err := problem.GetPossibleActions(current)
		if err != nil {
			return nil, err
		}

		// keliaujama per galimų veiksmų sąraša
		for _, action := range actions {
			// gaunama iš dabartinei būsenai pritaikyto veiksmo kylanti būsena
			resultingState, err := problem.GetResultingState(current, action)

			if err != nil {
				return nil, err
			}

			// tikrinama ar gauta būsena yra uždavinio galutinė būsena
			// jeigu taip:
			//   1. galutinė būsena įrašoma į aplankytų būsenų hashmapą
			//   2. nutraukiami abu ciklai - paieška baigta
			if resultingState.Describe() == problem.GetGoalState().Describe() {
				explored[resultingState.Describe()] = current
				solved = true
				break
			}

			// gauta būsena nėra galutinė, todėl tikriname ar ji jau buvo aplankyta
			// jeigu ne:
			//   1. į paieškos masyva pridedama atlikto veiksmo metu gauta būsena
			//   2. aplankytų būsenų hashmape, pažymima iš kur buvo atkeliauta į šią būseną
			if _, visited := explored[resultingState.Describe()]; !visited {

				fringe = append(fringe, resultingState)
				explored[resultingState.Describe()] = current
			}
		}
	}

	// gražinamas aplankytų būsenų hashmapas
	return explored, nil
}

// GraphSearch gauna prie galinės būsenos atvedusių tarpinių būsenų kombinaciją ir
// paverčia jas grafo mazgais, kuriuos gražina
func GraphSearch(problem problems.Problem, algo Algorithm) ([]*problems.Node, error) {
	explored, err := searchStateSpace(problem, algo)
	if err != nil {
		return nil, err
	}

	solution, err := getSolutionNodes(problem, explored)
	if err != nil {
		return nil, err
	}

	return solution, nil
}

// getSolutionNodes yra pagalbinė funkcija, argumentų sąraše priimanti problemą ir ištyrinėtą būsenų kelią ir
// būsenų kelią paverčianti grafo mazgais, kuriuos jinai gražina vartotojui
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

// takeFromFringe - tai pagalbinė funkcija, kuri remdamasi pateiktu paieškos algoritmu simuliuoja tai paieškai
// naudojamos duomenų struktūros elgseną, pvz.:
// algoritmas DFS -> taikomi steko DS loginiai ribojimai
// algoritmas BFS -> taikomi eilės DS loginiai ribojimai
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
