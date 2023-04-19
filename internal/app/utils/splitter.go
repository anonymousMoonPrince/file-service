package utils

func Split(value int64, count int) []int64 {
	partSize := value / int64(count)

	parts := make([]int64, 0, count)
	for i := 0; i < count-1; i++ {
		parts = append(parts, partSize)
	}
	parts = append(parts, value-int64(count-1)*partSize)
	return parts
}
