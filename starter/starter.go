package starter

type Starter interface {
	Start(args ...interface{}) (err error)
}
