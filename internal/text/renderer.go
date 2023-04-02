package text

type Renderer interface {
	Render(s string) string
}

type Passthrough struct{}

func (p *Passthrough) Render(s string) string {
	return s
}
