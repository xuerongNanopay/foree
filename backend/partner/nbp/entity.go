package nbp

type Response[R any] struct {
	rawRequest  string
	httpStatus  int
	rawResponse int
	data        R
}
