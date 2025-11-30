package misc

import (
	"iter"

	"golang.org/x/exp/constraints"
)

func Range[T constraints.Integer](start, end T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := start; i < end; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func Collect[T any](seq iter.Seq[T]) []T {
	var out []T
	for v := range seq {
		out = append(out, v)
	}
	return out
}

func Zip[A, B any](ita iter.Seq[A], itb iter.Seq[B]) iter.Seq2[A, B] {
	return func(yield func(A, B) bool) {
		next_a, stop_a := iter.Pull(ita)
		next_b, stop_b := iter.Pull(itb)
		for {
			value_a, ok_a := next_a()
			value_b, ok_b := next_b()
			if !ok_a || !ok_b || !yield(value_a, value_b) {
				break
			}
		}
		stop_a()
		stop_b()
	}
}

func IterateMapValues[K comparable, V any](mymap map[K]V, iterator iter.Seq[K]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for key := range iterator {
			value, exists := mymap[key]
			if !exists || !yield(value) {
				return
			}
		}
	}
}

func IterateSlice[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range slice {
			if !yield(v) {
				return
			}
		}
	}
}

func IterateValues[T any](values ...T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
}
