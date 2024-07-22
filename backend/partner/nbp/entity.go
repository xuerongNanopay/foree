package nbp

type Request[R any] struct {
	httpStatus  int
	rawResponse int
	data        R
}
