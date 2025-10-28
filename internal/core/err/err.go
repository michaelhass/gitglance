package err

type Msg interface {
	Error() error
	ErrorTitle() string
	ErrorDescription() string
}
