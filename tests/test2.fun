type int int

func Factorial(a int) int {
    if a < 0 || a == 0 {
        return 1;
    }
    return a * Factorial(a - 1);
}