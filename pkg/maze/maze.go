package maze

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"example.com/maze-solver/pkg/algos"
	"example.com/maze-solver/pkg/problems"
)

// Fasadinė struktūra, abstrahuojanti labirinto paveiksliuko apdorojimą,
// pavertimą grafu, paieškos funkcijų kvietimus ir sprendimo saugojimą faile
type Maze struct {
	Img      image.Image
	Solution []*problems.Node
	graph    *problems.GraphProblem
}

func New(filename string) *Maze {
	fmt.Println("Opening image...")
	file, err := os.Open(fmt.Sprintf("../img/input/%s", filename))
	if err != nil {
		log.Fatal("Error: File could not be opened")
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Converting image to graph...")
	pixels, err := getPixels(img)

	if err != nil {
		log.Fatal("Error: image could not be decoded")
	}

	graphProblem := pixelsToGraph(pixels)

	maze := Maze{Img: img, Solution: nil, graph: graphProblem}
	return &maze
}

func (maze *Maze) Solve(algo algos.Algorithm) {
	fmt.Println("Finding Exit...")
	var solution []*problems.Node = nil
	var err error
	switch algo {
	case algos.ASTAR_SEARCH:
		solution, err = algos.Astar(maze.graph)
	default:
		solution, err = algos.GraphSearch(maze.graph, algo)
	}
	if err != nil {
		log.Fatal(err)
	}
	maze.Solution = solution
	fmt.Println("Exit found!!!")
}

func (maze Maze) SaveToFile(filename string) {
	if maze.Solution == nil {
		log.Fatal("no maze solution found")
	}
	pathColor := color.RGBA{46, 225, 87, 255}
	outputFile, err := os.Create(fmt.Sprintf("../img/output/%s", filename))
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	for i := 1; i < len(maze.Solution); i++ {
		xStart, xEnd := sort(maze.Solution[i].X, maze.Solution[i-1].X)
		yStart, yEnd := sort(maze.Solution[i].Y, maze.Solution[i-1].Y)

		if xStart == xEnd {
			for y := yStart; y <= yEnd; y++ {
				maze.Img.(draw.Image).Set(xStart, y, pathColor)
			}
		} else {
			for x := xStart; x <= xEnd; x++ {
				maze.Img.(draw.Image).Set(x, yStart, pathColor)
			}
		}
	}

	png.Encode(outputFile, maze.Img)
	fmt.Println("Writing solution to file...")
	writeSolutionToFile(maze.Img, maze.Solution)
	fmt.Println("Done!!!")
}

// pagalbinis metodas, įrašantis labirinto sprendimo kelią į png paveiksliuką
func writeSolutionToFile(img image.Image, solution []*problems.Node) {
	pathColor := color.RGBA{46, 225, 87, 255}
	outputFile, err := os.Create("../img/output/maze_solution.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	for i := 1; i < len(solution); i++ {
		xStart, xEnd := sort(solution[i].X, solution[i-1].X)
		yStart, yEnd := sort(solution[i].Y, solution[i-1].Y)

		if xStart == xEnd {
			for y := yStart; y <= yEnd; y++ {
				img.(draw.Image).Set(xStart, y, pathColor)
			}
		} else {
			for x := xStart; x <= xEnd; x++ {
				img.(draw.Image).Set(x, yStart, pathColor)
			}
		}
	}

	png.Encode(outputFile, img)
}

// pagalbinis metodas, paverčiantis pateiktą png failą, pikselių matrica.
func getPixels(img image.Image) ([][]Pixel, error) {
	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	var pixels [][]Pixel

	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// metodas konvertuoja spalvų gamą į 8 bitų formatą ir gražina pikselio objektą
func rgbaToPixel(r, g, b, a uint32) Pixel {
	return Pixel{
		R: int(r / 257),
		G: int(g / 257),
		B: int(b / 257),
		A: int(a / 257),
	}
}

// metodas paverčia pateiktą pikselių matricą, grafu ir
// gražina grafo paieškos uždavinio objektą
func pixelsToGraph(pixels [][]Pixel) *problems.GraphProblem {
	var graph = &problems.GraphProblem{}
	var verticals = map[int]*problems.Node{}
	var lastHorizontal = &problems.Node{}

	initalIsSet := false

	for row := 0; row < len(pixels); row++ {
		for col := 0; col < len(pixels[row]); col++ {
			pixel := pixels[row][col]
			if pixel.IsBlack() {
				lastHorizontal = nil
				delete(verticals, col)
				continue
			}

			switch {
			case row == 0 || col == 0 || row == len(pixels)-1 || col == len(pixels[row])-1:
				var newNode *problems.Node
				if initalIsSet {
					graph.Goal = problems.Node{X: col, Y: row}
					newNode = &graph.Goal
				} else {
					graph.Initial = problems.Node{X: col, Y: row}
					newNode = &graph.Initial
				}
				if node, ok := verticals[col]; ok {
					node.Connections = append(node.Connections, newNode)
					newNode.Connections = append(newNode.Connections, node)
					delete(verticals, col)
				}
				lastHorizontal = newNode
				if row < len(pixels)-1 && pixels[row+1][col].IsWhite() {
					verticals[col] = newNode
				}
				initalIsSet = true
			default:
				left, right := pixels[row][col-1], pixels[row][col+1]
				top, bottom := pixels[row-1][col], pixels[row+1][col]

				if isDeadEndOrIntersection(left, right, top, bottom) {
					node := problems.Node{
						X: col,
						Y: row,
					}

					if lastHorizontal != nil {
						lastHorizontal.Connections = append(lastHorizontal.Connections, &node)
						node.Connections = append(node.Connections, lastHorizontal)
					}
					lastHorizontal = &node

					if nodeAbove, ok := verticals[col]; ok {
						node.Connections = append(node.Connections, nodeAbove)
						nodeAbove.Connections = append(nodeAbove.Connections, &node)
						delete(verticals, col)
					}

					if pixels[row+1][col].IsWhite() {
						verticals[col] = &node
					}
				}
			}
		}
	}

	return graph
}

// pikselį atspindinti struktūra
type Pixel struct {
	R int
	G int
	B int
	A int
}

func (p Pixel) IsWhite() bool {
	return p.R == 255 && p.G == 255 && p.B == 255
}

func (p Pixel) IsBlack() bool {
	return p.R == 0 && p.G == 0 && p.B == 0
}

// -----------------------------

// shouldPlaceNode - pagalbinė funkciją, kuri tikrina
// ar kelią atspindintis pikselis yra kampe abra sankryžoje
func isDeadEndOrIntersection(left, right, top, bottom Pixel) bool {
	wallCount := 0
	hasHorizontalPath, hasVerticalPath := false, false

	if left.IsBlack() {
		wallCount++
	} else {
		hasHorizontalPath = true
	}

	if right.IsBlack() {
		wallCount++
	} else {
		hasHorizontalPath = true
	}

	if top.IsBlack() {
		wallCount++
	} else {
		hasVerticalPath = true
	}

	if bottom.IsBlack() {
		wallCount++
	} else {
		hasVerticalPath = true
	}

	return wallCount == 3 || (hasHorizontalPath && hasVerticalPath)
}

// pagalbinė funkcija rūšiuojanti du kintamuosius
func sort(n1, n2 int) (int, int) {
	if n1 < n2 {
		return n1, n2
	}
	return n2, n1
}

func init() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}
