package fall

type Enviroment string

type Repository interface {
}

type Initializable interface {
	Init() error
}

type Controller interface {
	Configure(r *Router)
}

type UseCase[I any, O any] interface {
	Execute(input I) (O, error)
}

type EnvConfiguration interface {
	Configure(env Enviroment) error
}
