package idm

type IDMClient interface {
	TransferOut(req IDMRequest) error
}
