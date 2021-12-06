package actions

type Action interface {
	Run() error
}
