package findpath

func IsViableCellType(cellType int) bool {
	if cellType == CellTypeGoal || cellType == CellTypeEmpty || cellType == CellTypeStart {
		return true
	}

	return false
}
