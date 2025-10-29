package err

type Msg interface {
	Err() error
	ErrorTitle() string
	ErrorDescription() string
}

type errMsg struct {
	title string
	err   error
}

func (e errMsg) Err() error {
	return e.err
}

func (e errMsg) ErrorTitle() string {
	return e.title
}

func (e errMsg) ErrorDescription() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func NewMsg(title string, err error) Msg {
	return errMsg{title: title, err: err}
}
