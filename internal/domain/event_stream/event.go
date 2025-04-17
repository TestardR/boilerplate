package eventstream

type Event struct {
	Type     string
	Payload  []byte
	Metadata map[string]string
}
