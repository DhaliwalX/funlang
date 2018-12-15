// +build ignore

package test

func Add(a int, b int) int {
	return a + b
}

func Pointer(a *int) *int {
	b := &a
	return *b
}


