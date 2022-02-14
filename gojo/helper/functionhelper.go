package helper

// WrapUnaryAsync wraps the Action in a function accepting an interface{} parameter since type information is lost
//after being sent to the port
func WrapUnaryAsync[T any](function func(T)) func(interface{}) {
	return func(a interface{}) {
		function(a.(T))
	}
}

// WrapUnarySync wraps the Action in a function accepting an interface{} parameter and returns an interface{} parameter
//since type information is lost after being sent to the port
func WrapUnarySync[T any, R any](function func(a T) R) func(interface{}) interface{} {
	return func(a interface{}) interface{} {
		return function(a.(T))
	}
}

// WrapBinarySync wraps the Action in a function accepting interface{} parameters and returns an interface{} parameter
//since type information is lost after being sent to the port
func WrapBinarySync[T any, S any, R any](function func(T, S) R) func(interface{}, interface{}) interface{} {
	return func(a interface{}, b interface{}) interface{} {
		return function(a.(T), b.(S))
	}
}

// WrapBinaryAsync wraps the Action in a function accepting interface{} parameters
//since type information is lost after being sent to the port
func WrapBinaryAsync[T any, R any](function func(T, R)) func(interface{}, interface{}) {
	return func(a interface{}, b interface{}) {
		function(a.(T), b.(R))
	}
}

// WrapTernarySync wraps the Action in a function accepting interface{} parameters and returns an interface{} parameter
//since type information is lost after being sent to the port
func WrapTernarySync[T any, S any, R any, U any](function func(T, S, R) U) func(interface{}, interface{}, interface{}) interface{} {
	return func(a interface{}, b interface{}, c interface{}) interface{} {
		return function(a.(T), b.(S), c.(R))
	}
}

// WrapTernaryAsync wraps the Action in a function accepting interface{} parameters
//since type information is lost after being sent to the port
func WrapTernaryAsync[T any, S any, R any](function func(T, S, R)) func(interface{}, interface{}, interface{}) {
	return func(a interface{}, b interface{}, c interface{}) {
		function(a.(T), b.(S), c.(R))
	}
}
