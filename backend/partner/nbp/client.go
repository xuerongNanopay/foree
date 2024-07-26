package nbp

type NBPClient interface {
	Hello() (*HelloResponse, error)
}
