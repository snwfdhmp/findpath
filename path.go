package findpath

type Path interface {
	Pathway() []Pos
	Add(Pos)
	AddList(pos []Pos)
	IsValid() bool
	String() string
}

type path struct {
	pathway []Pos
}

func NewPath() Path {
	return &path{
		pathway: make([]Pos, 0),
	}
}

func NewPathFromPathway(pathway []Pos) Path {
	return &path{
		pathway: pathway,
	}
}

func (p *path) Add(pos Pos) {
	(*p).pathway = append(p.pathway, pos)
}

func (p *path) AddList(pos []Pos) {
	(*p).pathway = append(p.pathway, pos...)
}

func (p *path) IsValid() bool {
	for curPos := 0; curPos < len(p.pathway)-1; curPos++ {
		if abs(p.pathway[curPos].X-p.pathway[curPos+1].X)+abs(p.pathway[curPos].Y-p.pathway[curPos+1].Y) != 1 {
			return false
		}
	}
	return true
}

func (p *path) Pathway() []Pos {
	return (*p).pathway
}

func (p *path) String() string {
	output := ""
	for i := range p.pathway {
		output += p.pathway[i].String()
		if i < len(p.pathway)-1 {
			output += "=>"
		}
	}

	return output
}

func abs(n int) int {
	if n <= 0 {
		return n * -1
	}
	return n
}
