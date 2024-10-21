package email

import (
	"bytes"
	"io"
	"mime"
)

type mimeEncoder struct {
	mime.WordEncoder
}

var (
	bEncoding = mimeEncoder{mime.BEncoding}
	qEncoding = mimeEncoder{mime.QEncoding}
)

type Encoding string
type header map[string][]string

const (
	QuotedPrintable Encoding = "quoted-printable"
	Base64          Encoding = "base64"
	Unencoded       Encoding = "8bit"
)

type file struct {
	Name     string
	Header   map[string][]string
	CopyFunc func(w io.Writer) error
}

type part struct {
	contentType string
	copier      func(io.Writer) error
	encoding    Encoding
}

type Message struct {
	header      header
	parts       []*part
	attachments []*file
	embedded    []*file
	charset     string
	encoding    Encoding
	hEncoder    mimeEncoder
	buf         bytes.Buffer
}
