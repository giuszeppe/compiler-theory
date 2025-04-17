package main

import (
        "errors"
)

// Generic stack for type any (Go 1.18+)
type Stack[T any] struct {
    items []T
}

// Push adds an item to the top of the stack
func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

// Pop removes and returns the top item
func (s *Stack[T]) Pop() (T, error) {
    if len(s.items) == 0 {
        var zero T
            return zero, errors.New("stack is empty")
    }
    last := s.items[len(s.items)-1]
        s.items = s.items[:len(s.items)-1]
        return last, nil
}

// Peek returns the top item without removing it
func (s *Stack[T]) Peek() (T, error) {
    if len(s.items) == 0 {
        var zero T
            return zero, errors.New("stack is empty")
    }
    return s.items[len(s.items)-1], nil
}

// IsEmpty checks if the stack is empty
func (s *Stack[T]) IsEmpty() bool {
    return len(s.items) == 0
}

// Size returns the number of items in the stack
func (s *Stack[T]) Size() int {
    return len(s.items)
}

