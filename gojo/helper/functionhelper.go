package helper

func WrapUnaryAsync[T any](function func(T)) func(interface{}) {
	return func(a interface{}) {
		function(a.(T))
	}
}

func WrapUnarySync[T any, R any](function func(a T) R) func(interface{}) interface{} {
	return func(a interface{}) interface{} {
		return function(a.(T))
	}
}

func WrapBinarySync[T any, S any, R any](function func(T, S) R) func(interface{}, interface{}) interface{} {
	return func(a interface{}, b interface{}) interface{} {
		return function(a.(T), b.(S))
	}
}

func WrapBinaryAsync[T any, R any](function func(T, R)) func(interface{}, interface{}) {
	return func(a interface{}, b interface{}) {
		function(a.(T), b.(R))
	}
}

func WrapTernarySync[T any, S any, R any, U any](function func(T, S, R) U) func(interface{}, interface{}, interface{}) interface{} {
	return func(a interface{}, b interface{}, c interface{}) interface{} {
		return function(a.(T), b.(S), c.(R))
	}
}

func WrapTernaryAsync[T any, S any, R any](function func(T, S, R)) func(interface{}, interface{}, interface{}) {
	return func(a interface{}, b interface{}, c interface{}) {
		function(a.(T), b.(S), c.(R))
	}
}
