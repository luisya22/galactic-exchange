package maputils

func CopyMap[K comparable, V any](originalMap map[K]V) map[K]V {
	copiedMap := make(map[K]V, len(originalMap))
	for key, value := range originalMap {
		copiedMap[key] = value
	}

	return copiedMap
}
