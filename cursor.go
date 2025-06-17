package nut

type cursor struct {
	value int
}

func (cursor *cursor) next() {
	cursor.value++
}

func (cursor *cursor) previous() {
	cursor.value--
}

func (cursor *cursor) getValue() int {
	return cursor.value
}

func newCursor(value int) *cursor {
	return &cursor{value}
}	