package utils

func GeneratePointerArrayFrom(array []string) []*string {
	var pointerArray []*string

	for _, element := range array {
		pointerArray = append(pointerArray, &element)
	}

	return pointerArray
}
