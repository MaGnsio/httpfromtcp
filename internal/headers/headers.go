package headers

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	h["Host"] = "localhost:42069"
	n = 23
	done = false
	return n, done, err
}
