package utils

type PageIndex struct {
	Limit  int
	Offset int
}

func (p *PageIndex) Getimit() int {
	return p.Limit
}

func (p *PageIndex) GeOffset() int {
	return p.Offset
}
