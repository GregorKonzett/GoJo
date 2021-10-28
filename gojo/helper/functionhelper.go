package helper

func WrapUnarySend[T any](function func(T)) func(interface{}) {
	return func(a interface{}) {
		function(a.(T))
	}
}

func WrapBinarySend[T any, R any](function func(T, R)) func(interface{}, interface{}) {
	return func(a interface{}, b interface{}) {
		function(a.(T), b.(R))
	}
}

func WrapBinaryRecv[T any, R any](function func(T) R) func(interface{}) interface{} {
	return func(a interface{}) interface{} {
		return function(a.(T))
	}
}
