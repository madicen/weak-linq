package weaklinq

import (
	"reflect"
)

//----------------------------------------------------------------------------//
// Flatten                                                                    //
//----------------------------------------------------------------------------//

////////////////////////////////////////////////////////////////////////////////

// FlattenThese returns a new Iterable where the items are flattened by the
// items returned by the selector.
func (iterable Iterable[T]) FlattenThese(selector func(T) Iterable[any]) Iterable[any] {

	return Iterable[any]{
		Seq: func(yield func(any) bool) {
			iterable.Seq(func(item T) bool {
				selector(item).Seq(func(item any) bool {
					return yield(item)
				})
				return true
			})
		},
	}

	/*
		linq.From([]T{...}).
			FlattenThese(
				func(item T) Iterable[any] {
					return linq.From([]any{...})
				},
			)
	*/
}

////////////////////////////////////////////////////////////////////////////////

// Flatten returns a new Iterable where the items are flattened by the
// items returned by the field name.
// If T is not a struct, or fieldName is not found, this function will panic.
func (iterable Iterable[T]) Flatten(fieldName string) Iterable[any] {

	return iterable.FlattenThese(
		func(item T) Iterable[any] {
			result := getFieldNameFunc[T](fieldName)(item)
			ref := reflect.ValueOf(result)
			if ref.Kind() != reflect.Slice {
				return From([]any{result})
			}
			return Iterable[any]{
				Seq: func(yield func(any) bool) {
					for i := 0; i < ref.Len(); i++ {
						if !yield(ref.Index(i).Interface()) {
							return
						}
					}
				},
			}
		},
	)

	/*
		linq.From([]T{...}).
			Flatten("ItemField")
	*/
}
