package nut

const (
	block = iota
	void
	meta
)

type node_t struct {
	tag string
	type_t int
	attributes []*attribute
	child *node_t
	sibling *node_t
}