package findpath

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/afero"
)

var (
	fs = afero.NewOsFs()
)

type Level interface {
	Cell(x, y int) int
	Goal() Pos
	Start() Pos
	Print(w io.Writer, playerPos *Pos)
}

type level struct {
	cells [][]int
	start Pos
	goal  Pos
}

func (l *level) Print(w io.Writer, playerPos *Pos) {
	for i := range (*l).cells {
		for j := range l.cells[i] {
			if playerPos != nil && playerPos.Y == i && playerPos.X == j {
				io.WriteString(w, "P")
				continue
			}
			io.WriteString(w, fmt.Sprintf("%s", cellTypeChar(l.cells[i][j])))
		}
		if len(l.cells[i]) > 1 {
			io.WriteString(w, "\n")
		}
	}
}

func (l *level) Cell(x, y int) int {
	if x < 0 || y < 0 {
		return -1
	}
	if len((*l).cells)-1 < y || len((*l).cells[y])-1 < x {
		return -1
	}

	return (*l).cells[y][x]
}

func (l *level) Goal() Pos {
	return l.goal
}

func (l *level) Start() Pos {
	return l.start
}

func OpenLevel(filepath string) (Level, error) {
	content, err := afero.ReadFile(fs, filepath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) < 1 {
		return nil, errors.New("empty level file")
	}

	size := len(lines[0])
	cells := make([][]int, size)
	start := new(Pos)
	goal := new(Pos)

	for i := range lines {
		if len(lines[i]) != size {
			return nil, errors.New("error level file format: all lines should be of same length")
		}
		cells[i] = make([]int, size)
		for x := range lines[i] {
			switch lines[i][x] {
			case '*':
				cells[i][x] = CellTypeWall
			case 'G':
				*goal = NewPos(x, i)
				cells[i][x] = CellTypeGoal
			case 'S':
				*start = NewPos(x, i)
				cells[i][x] = CellTypeStart
			case 'X':
				cells[i][x] = CellTypeWall
			case ' ':
				cells[i][x] = CellTypeEmpty
			default:
				return nil, fmt.Errorf("error level file format: unknown character '%c'", lines[i][x])
			}
		}
	}

	if goal == nil {
		return nil, errors.New("error level content: missing goal")
	} else if start == nil {
		return nil, errors.New("error level content: missing start")
	}

	return &level{
		goal:  *goal,
		start: *start,
		cells: cells,
	}, nil
}

type Pos struct {
	X int
	Y int
}

func (p Pos) String() string {
	return fmt.Sprintf("[%d;%d]", p.X, p.Y)
}

func NewPos(x, y int) Pos {
	return Pos{x, y}
}

const (
	CellTypeGoal  = 0
	CellTypeStart = 1
	CellTypeWall  = 2
	CellTypeEmpty = 3
)

func cellTypeChar(cellType int) string {
	switch cellType {
	case CellTypeGoal:
		return "G"
	case CellTypeStart:
		return "S"
	case CellTypeWall:
		return "X"
	case CellTypeEmpty:
		return " "
	default:
		return "?"
	}
}
