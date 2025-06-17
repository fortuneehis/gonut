package nut

const (
	eof = iota
	newLine
	node
	lParen
	rParen
	equals
	comma  
	variable
	qString
	indent
)


type token struct {
	kind  int
	value string
}

func newToken(kind int, value string) *token {
	return &token{kind, value}
}