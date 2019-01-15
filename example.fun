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
