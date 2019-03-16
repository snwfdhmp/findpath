package findpath

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// Challenger represents a challenger interface
type Challenger interface {
	FindPath(lvl Level) Path
}

func Challenge(c Challenger, lvl Level) error {
	path := c.FindPath(lvl)
	if path == nil {
		return errors.New("path is nil")
	}

	pathString := make([]string, 0)
	for i := range path.Pathway() {
		moveString := fmt.Sprintf("[%d;%d]", path.Pathway()[i].X, path.Pathway()[i].Y)
		fmt.Printf("Step %d: %s\n", i, moveString)
		lvl.Print(os.Stdout, &path.Pathway()[i])
		time.Sleep(1 * time.Second)
		pathString = append(pathString, moveString)
	}

	fmt.Printf("path: %s\n", strings.Join(pathString, "=>"))

	if err := ValidatePath(lvl, path); err != nil {
		return fmt.Errorf("path validation error: %v", err)
	}

	fmt.Println("path validated !")

	return nil
}

func ValidatePath(lvl Level, path Path) error {
	if !path.IsValid() {
		return errors.New("invalid path")
	}

	lenPathway := len(path.Pathway())
	if lenPathway < 2 {
		return fmt.Errorf("pathway too short (length=%d, should be >2)", lenPathway)
	}

	for _, pos := range path.Pathway() {
		if !IsViableCellType(lvl.Cell(pos.X, pos.Y)) {
			return errors.New("path goes through unviable cells")
		}
	}

	pathStart := path.Pathway()[0]
	pathEnd := path.Pathway()[lenPathway-1]
	if !(lvl.Start() == pathStart && lvl.Goal() == pathEnd) && !(lvl.Start() == pathEnd && lvl.Goal() == pathStart) {
		return fmt.Errorf(`path should start with lvl.Start() (should: [%d;%d], cur: [%d;%d]) and end with lvl.Goal() (should: [%d;%d], cur: [%d;%d])`, lvl.Start().X, lvl.Start().Y, pathStart.X, path.Pathway()[0].Y, lvl.Goal().X, lvl.Goal().Y, pathEnd.X, pathEnd.Y)
	}

	return nil
}
