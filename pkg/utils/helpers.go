package utils

// Sum calculates the sum of a slice of integers.
func Sum(numbers []int) int {
    total := 0
    for _, number := range numbers {
        total += number
    }
    return total
}

// Average calculates the average of a slice of integers.
func Average(numbers []int) float64 {
    if len(numbers) == 0 {
        return 0
    }
    total := Sum(numbers)
    return float64(total) / float64(len(numbers))
}

// Max returns the maximum value in a slice of integers.
func Max(numbers []int) int {
    if len(numbers) == 0 {
        return 0
    }
    max := numbers[0]
    for _, number := range numbers {
        if number > max {
            max = number
        }
    }
    return max
}

// Min returns the minimum value in a slice of integers.
func Min(numbers []int) int {
    if len(numbers) == 0 {
        return 0
    }
    min := numbers[0]
    for _, number := range numbers {
        if number < min {
            min = number
        }
    }
    return min
}