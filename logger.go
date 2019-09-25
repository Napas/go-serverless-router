package routing

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}

type NilLogger struct{}

func (*NilLogger) Print(...interface{})          {}
func (*NilLogger) Printf(string, ...interface{}) {}
func (*NilLogger) Println(...interface{})        {}
