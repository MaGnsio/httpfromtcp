package request

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
)

var crlf = []byte("\r\n")

const bufferSize = 1024

var (
	ErrorMalformedRequestLine = fmt.Errorf("malformed request-line")
	ErrorInvalidMethod        = fmt.Errorf("invalid method")
	ErrorInvalidHTTPVersion   = fmt.Errorf("invalid HTTP version")
	ErrorEmptyRequest         = fmt.Errorf("empty request")
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

func (rl *RequestLine) ValidMethod() bool {
	return rl.Method == strings.ToUpper(rl.Method)
}

type Request struct {
	RequestLine
	state parserState
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.state = StateDone

		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) Done() bool {
	return r.state == StateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		state: StateInit,
	}

	// NOTE: buffer could exceed available memory for very large requests.
	buf, bufLen := make([]byte, bufferSize), 0
	for !request.Done() {
		// reading into the buffer
		n, err := reader.Read(buf[bufLen:])
		log.Printf("Read %d bytes, error: %v\n", n, err)
		if n == 0 {
			request.state = StateDone
			break
		}
		if err != nil {
			return nil, err
		}
		bufLen += n

		// parse from the buffer
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		// shift the buffer to remove parsed data
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, crlf)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(crlf)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrorMalformedRequestLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrorInvalidHTTPVersion
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HTTPVersion:   string(httpParts[1]),
	}

	if !rl.ValidMethod() {
		return nil, read, ErrorInvalidMethod
	}

	return rl, read, nil
}
