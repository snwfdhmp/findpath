package main

import (
	"fmt"

	"github.com/snwfdhmp/findpath"
)

type challenger struct {
	crossed map[int]map[int]bool
}

func (c *challenger) setCrossed(pos findpath.Pos) {
	if _, ok := (*c).crossed[pos.Y]; !ok {
		(*c).crossed[pos.Y] = make(map[int]bool)
	}

	(*c).crossed[pos.Y][pos.X] = true
}

func (c *challenger) isCrossed(pos findpath.Pos) bool {
	if _, ok := (*c).crossed[pos.Y]; !ok {
		return false
	}

	return (*c).crossed[pos.Y][pos.X]
}

type move struct {
	targetPos    findpath.Pos
	previousMove *move
}

func (m *move) isFirstMove() bool {
	return m.previousMove == nil
}

func NewMove(previous *move, target findpath.Pos) *move {
	return &move{
		targetPos:    target,
		previousMove: previous,
	}
}

func (c *challenger) FindPath(lvl findpath.Level) findpath.Path {
	moves := []*move{NewMove(nil, lvl.Goal())}
	c.setCrossed(lvl.Goal())

	fmt.Printf("Searching...\n")
Walk:
	for cur := 0; cur < len(moves); cur++ {
		fmt.Printf("- %s\n", moves[cur].targetPos.String())
		newMoves := c.PossibleMoves(lvl, moves[cur])
		for _, move := range newMoves {
			if c.isCrossed(move.targetPos) {
				continue
			}

			c.setCrossed(move.targetPos)
			moves = append(moves, move)

			if move.targetPos == lvl.Start() {
				break Walk
			}
		}
	}

	path := findpath.NewPath()
	fmt.Printf("Rewinding...\n")
	for move := moves[len(moves)-1]; move != nil; move = move.previousMove {
		fmt.Printf("pos: %s\n", move.targetPos.String())
		path.Add(move.targetPos)
	}

	return path
}

func (c *challenger) PossibleMoves(lvl findpath.Level, curMove *move) []*move {
	moves := make([]*move, 0)
	//try up
	directions := [][]int{
		{curMove.targetPos.X, curMove.targetPos.Y + 1}, // up
		{curMove.targetPos.X + 1, curMove.targetPos.Y}, // right
		{curMove.targetPos.X, curMove.targetPos.Y - 1}, // down
		{curMove.targetPos.X - 1, curMove.targetPos.Y}, // left
	}

	for _, d := range directions { //for each direction
		if cellUp := lvl.Cell(d[0], d[1]); findpath.IsViableCellType(cellUp) { //if target cell is viable
			moves = append(moves, NewMove(curMove, findpath.NewPos(d[0], d[1])))
		}
	}

	return moves
}

func main() {
	lvl, err := findpath.OpenLevel("testdata/map_1.txt")
	if err != nil {
		fmt.Printf("fatal: %s\n", err)
		return
	}

	challenger := newChallenger()

	if err := findpath.Challenge(challenger, lvl); err != nil {
		fmt.Printf("failed: %v\n", err)
		return
	}
}

func newChallenger() *challenger {
	return &challenger{
		crossed: make(map[int]map[int]bool),
	}
}
