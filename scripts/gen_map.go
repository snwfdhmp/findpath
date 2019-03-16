package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/murlokswarm/log"
	"github.com/snwfdhmp/findpath"
	"github.com/spf13/afero"
)

var (
	randomCellTypeLibrary = []int{findpath.CellTypeEmpty, findpath.CellTypeEmpty,
		findpath.CellTypeWall, findpath.CellTypeWall, findpath.CellTypeWall}
)

func randomCellType() int {
	return randomCellTypeLibrary[rand.Intn(len(randomCellTypeLibrary))]
}

var (
	fs = afero.NewOsFs()
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	generateMap(80, 20, 4)
}

func generateMap(sizeX, sizeY int, complexity int) {
	// allocate 2d slice for board
	board := make([][]int, sizeY)
	for i := range board {
		board[i] = make([]int, sizeX)
	}

	startPos := randPos(sizeX-1, sizeY-1)
	path := findpath.NewPath()
	path.Add(startPos)
	previousPos := &startPos
	for curPathIndex := 0; curPathIndex <= complexity; curPathIndex++ {
		// moveTo := randPos(sizeX, sizeY)
		moveTo := randPosAwayFrom(sizeX, sizeY, sizeX/3, 3*sizeX/4, sizeY/4, 3*sizeY/4, previousPos.X, previousPos.Y)
		pathway := makePathBlind(*previousPos, moveTo)
		previousPos = &moveTo
		path.AddList(pathway.Pathway()[1:])
	}
	goalPos := *previousPos

	fmt.Printf("start: %s => goal: %s\nsafe path: %s\n", startPos, goalPos, path)

	avoid := make(map[findpath.Pos]bool)
	for _, pos := range path.Pathway() {
		avoid[pos] = true
	}

	for y := range board {
		for x := range board[y] {
			pos := findpath.NewPos(x, y)
			_, isProtected := avoid[pos]
			switch {
			case pos == startPos:
				board[y][x] = findpath.CellTypeStart
				continue
			case pos == goalPos:
				board[y][x] = findpath.CellTypeGoal
				continue
			case isProtected:
				board[y][x] = findpath.CellTypeEmpty
			default:
				switch {
				case y == 0 || x == 0 || y == sizeY-1 || x == sizeX-1:
					board[y][x] = findpath.CellTypeWall
				default:
					board[y][x] = randomCellType()
				}
			}
		}
	}

	lvl, err := findpath.NewLevel(board)
	if err != nil {
		log.Errorf("error creating level: %s", err)
	}

	lvl.Print(os.Stdout, nil)
}

func randPos(mapSizeX, mapSizeY int) findpath.Pos {
	return findpath.NewPos(1+rand.Intn(mapSizeX-1), 1+rand.Intn(mapSizeY-1)) //generate random X and Y in given range, and return new position
}

func randPosAwayFrom(mapSizeX, mapSizeY, awayMinX, awayMaxX, awayMinY, awayMaxY, fromX, fromY int) findpath.Pos {
	direction := rand.Intn(2) * -1
	awayDist := awayMinX
	if awayMinX != awayMaxX {
		awayDist = rand.Intn(awayMaxX - awayMinX)
	}
	randX := max(min(fromX+(direction*awayDist), mapSizeX-2), 1) //compute safe value

	direction = rand.Intn(2) * -1
	awayDist = awayMinY
	if awayMinY != awayMaxY {
		awayDist = rand.Intn(awayMaxY - awayMinY)
	}
	randY := max(min(fromY+(direction*awayDist), mapSizeY-2), 1) //compute safe value

	return findpath.NewPos(randX, randY)
}

//makePathBlind return the path to go from 'from' to 'to' without taking care of the cells' †ypes
func makePathBlind(from, to findpath.Pos) findpath.Path {
	distX := to.X - from.X
	distY := to.Y - from.Y

	//create path's slice
	moves := make([]findpath.Pos, abs(distX)+abs(distY)+3)
	moves[0] = from
	deltaX := 0 //this represents the distance walked by X
	deltaY := 0 //this represents the distance walked by Y
	for i := 1; i < len(moves); i++ {
		distX = to.X - moves[i-1].X
		distY = to.Y - moves[i-1].Y
		directions := make(map[int]int) //possible direction to take

		// log.Infof("makePathBlind: distanceX=%d distanceY=%d", distX, distY)
		switch {
		case distX < 0:
			directions[Left] = Left
		case distX >= 0:
			directions[Right] = Right
		}
		switch {
		case distY < 0:
			directions[Down] = Down
		case distY > 0:
			directions[Up] = Up
		}
		// log.Infof("makePathBlind: directions: %s", strDirectionMap(directions))
		if len(directions) < 1 {
			log.Warn("makePathBlind:  no more direction to go, exiting for loop")
			break
		}

		randDir := rand.Intn(len(directions))
		move := findpath.NewPos(moves[i-1].X, moves[i-1].Y)
		for dir := range directions {
			if randDir != 0 {
				randDir--
				continue
			}

			switch dir {
			case Up:
				deltaY++
				move.Y++
			case Down:
				deltaY--
				move.Y--
			case Right:
				deltaX++
				move.X++
			case Left:
				deltaX--
				move.X--
			}
			// log.Infof("going %s", strDirection(dir))
			break
		}

		moves[i] = move
	}

	path := findpath.NewPathFromPathway(moves)
	// fmt.Printf("path output: %s\n", path)
	if !path.IsValid() {
		log.Error("path is not valid !")
	}

	return path
}

const (
	//directions
	Up    = 0
	Right = 1
	Down  = 2
	Left  = 3
)

func strDirection(direction int) string {
	switch direction {
	case Up:
		return "up"
	case Down:
		return "down"
	case Right:
		return "right"
	case Left:
		return "left"
	default:
		return "?unknown_direction?"
	}
}

func strDirectionMap(directions map[int]int) string {
	output := ""
	for dir := range directions {
		switch dir {
		case Up:
			output += "up "
		case Down:
			output += "down "
		case Right:
			output += "righ "
		case Left:
			output += "left "
		default:
			output += "?unknown_direction? "
		}
	}

	return output
}

func abs(n int) int {
	if n < 0 {
		return n * -1
	}
	return n
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
