package err

type Msg interface {
	Err() error
	ErrorTitle() string
	ErrorDescription() string
}
