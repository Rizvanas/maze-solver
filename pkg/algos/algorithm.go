package algos

import "errors"

//--- Paieškos algoritmo enumas -------
type Algorithm byte

const (
	DEPTH_FIRST_SEARCH   Algorithm = 1
	BREADTH_FIRST_SEARCH Algorithm = 2
	ASTAR_SEARCH         Algorithm = 3
)

// String gražina paieškos algoritmo enumo tekstinę eilutę
func (a Algorithm) String() string {
	switch a {
	case DEPTH_FIRST_SEARCH:
		return "dfs"
	case BREADTH_FIRST_SEARCH:
		return "bfs"
	case ASTAR_SEARCH:
		return "astar"
	default:
		return ""
	}
}

// AlgoFromString pagalbinė funkcija, kuri vartotojui pateikus paieškos algoritmą apibūdinančią tekstinę eilutę,
// gražina paieškos algoritmo enumą
func AlgoFromString(str string) (Algorithm, error) {
	switch str {
	case DEPTH_FIRST_SEARCH.String():
		return DEPTH_FIRST_SEARCH, nil
	case BREADTH_FIRST_SEARCH.String():
		return BREADTH_FIRST_SEARCH, nil
	case ASTAR_SEARCH.String():
		return ASTAR_SEARCH, nil
	default:
		return 0, errors.New("algorithm by this name not found")
	}
}

//-------------------------------------
