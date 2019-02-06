type int int

// as of now will place phi in wrong block
func loops(a int) int {
    var b = 10;
    var m int;
    for a > b {
        m = a + b;
        b = m;
    }

    return b;
}
