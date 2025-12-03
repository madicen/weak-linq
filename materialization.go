package weaklinq

import "reflect"

//----------------------------------------------------------------------------//
// Materialization                                                            //
//----------------------------------------------------------------------------//

////////////////////////////////////////////////////////////////////////////////

// AndAssignToMap assigns the result of iterating over the MapIterable to a map.
// The type of the result map must be the same as the type of the key and value
// selector results in the MapIterable. If there is a key or value type mismatch,
// or the result is not a pointer to a map, this function will panic. If the
// MapIterable is set to overwrite, this function will attempt to build a slice
// of the value type, or append to the value if it already exists.
func (iterable MapIterable[T]) AndAssignToMap(result any) {

	res := reflect.ValueOf(result)
	if res.Kind() != reflect.Pointer || res.IsNil() {
		panic("'result' must be a pointer to a map")
	}

	mapValue := reflect.Indirect(res)
	if mapValue.Kind() != reflect.Map {
		panic("'result' must be a pointer to a map")
	}

	m := reflect.Indirect(res)

	for item := range iterable.itemIterable.Seq {

		key := iterable.keySelector(item)
		val := iterable.itemSelector(item)
		reflectKey := reflect.ValueOf(key)
		reflectVal := reflect.ValueOf(val)

		if iterable.overwrite {
			m.SetMapIndex(reflectKey, reflectVal)

		} else {
			existingValue := m.MapIndex(reflectKey)

			if !existingValue.IsValid() {
				sliceType := reflect.SliceOf(reflect.TypeOf(val))
				newSlice := reflect.MakeSlice(sliceType, 0, 1)
				existingValue = newSlice
			}

			appendedSlice := reflect.Append(existingValue, reflect.ValueOf(val))
			m.SetMapIndex(reflectKey, appendedSlice)
		}
	}

	res.Elem().Set(m)

	/*
		// Value
		result := make(map[K]V)
		linq.From([]T{...}).
			GroupThis("ItemField").
			By("KeyField").
			AndAssignToMap(&result)

		// List
		result := make(map[K][]V)
		linq.From([]T{...}).
			GroupListsOf("ItemField").
			By("KeyField").
			AndAssignToMap(&result)
	*/
}

////////////////////////////////////////////////////////////////////////////////

// AndAssignToSlice assigns the result of iterating over the Iterable to a
// slice. The type of the result slice must be the same as the type of the
// item in the Iterable. If there is a type mismatch, or the result is not
// a pointer to a slice, this function will panic.
func (iterable Iterable[T]) AndAssignToSlice(result any) {

	res := reflect.ValueOf(result)
	if res.Kind() != reflect.Pointer || res.IsNil() {
		panic("'result' must be a pointer to a slice")
	}

	sliceValue := reflect.Indirect(res)
	if sliceValue.Kind() != reflect.Slice {
		panic("'result' must be a pointer to a slice")
	}

	s := reflect.Indirect(res)

	for item := range iterable.Seq {
		s.Set(reflect.Append(s, reflect.ValueOf(item)))
	}

	res.Elem().Set(s)

	/*
		result := make([]V, 0)
		linq.From([]T{...}).
			Get("ItemField").
			AndAssignToSlice(&result)
	*/
}
