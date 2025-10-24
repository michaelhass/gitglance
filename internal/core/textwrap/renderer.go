package textwrap

type Renderer interface {
	Render(s ...string) string
}

type Passthrough struct{}

func (p *Passthrough) Render(s ...string) string {
	var output string
	for _, value := range s {
		output += value
	}
	return output
}
