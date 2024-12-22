package utils

// ConvertListOfObjectsToListOfPointers converts a list of objects to a list of pointers
func ConvertListOfObjectsToListOfPointers[T any](list *[]T) []*T {
	if list == nil {
		return nil
	}
	pointerList := make([]*T, len(*list))
	for i := range *list {
		pointerList[i] = &(*list)[i]
	}
	return pointerList
}
