type int int

func Add(a int, b int) int {
    return a + b;
}

func identity(a int) int {
    return a;
}

func max(a int, b int) int {
    if a > b {
        return a;
    } else {
        return b;
    }
}

func multipleDefs(a int) int {
    var b = 10;
    var max int;
    if a > b {
        max = a;
    } else {
        max = b;
    }

    return max;
}

// as of now will place phi in wrong block
func loops(a int) int {
    var b = 10;
    for a > b {
        var m = max(a, b);
        b = m;
    }

    return b;
}
