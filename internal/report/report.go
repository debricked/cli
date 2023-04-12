package report

type IReporter interface {
	Order(args IOrderArgs) error
}

type IOrderArgs interface{}
