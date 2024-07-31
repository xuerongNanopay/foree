package idm

type IDMClient interface {
	Transfer(req IDMRequest) error
}
