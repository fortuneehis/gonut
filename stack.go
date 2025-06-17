package nut

type stack[T any] struct {
	items []T
}

func (stack *stack[T]) push(item T) {
	stack.items = append(stack.items, item)
}

func (stack *stack[T]) get(index int) T {
	if index > stack.size() {
		panic("Stack Index Out Of Bounds")
	}
	return stack.items[index]
}
  
func (stack *stack[T]) top() T {
	return stack.get(stack.size() - 1)

}

func (stack *stack[T]) empty() bool {
	return len(stack.items) == 0 
}
 

func (stack *stack[T]) pop() {
	stack.items = stack.items[:len(stack.items) - 1]
}

func (stack *stack[T]) size() int {
	return len(stack.items)
}

func newStack[T any]() *stack[T] {
	return &stack[T]{items: []T{}}
}