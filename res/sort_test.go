package res

import "testing"

func TestSort(t *testing.T) {
    a := []int{0,1,2,3,4}
    b := make([]int, len(a))
    copy(b, a)
    sort(a, b)

    t.Log(a)
    t.Log(b)
}