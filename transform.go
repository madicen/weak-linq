package weaklinq

//----------------------------------------------------------------------------//
// Transform                                                                  //
//----------------------------------------------------------------------------//

////////////////////////////////////////////////////////////////////////////////

// GetThese returns a new Iterable where the items are transformed by selector.
func (iterable Iterable[T]) GetThese(selector func(T) any) Iterable[any] {

	return Iterable[any]{
		Seq: func(yield func(any) bool) {
			iterable.Seq(func(item T) bool {
				return yield(selector(item))
			})
		},
	}

	/*
		linq.From([]T{...}).
			GetThese(
				func(item T) any {
					return item.ItemField
				},
			)
	*/
}

////////////////////////////////////////////////////////////////////////////////

// Get returns a new Iterable where the items are transformed by fieldName.
// If T is not a struct, or fieldName is not found, this function will panic.
func (iterable Iterable[T]) Get(fieldName string) Iterable[any] {

	return iterable.GetThese(
		getFieldNameFunc[T](fieldName),
	)

	/*
		linq.From([]T{...}).
			Get("ItemField")
	*/
}

////////////////////////////////////////////////////////////////////////////////

// AsAny converts an Iterable of any type T to an Iterable of type any.
func (iterable Iterable[T]) AsAny() Iterable[any] {

	return Iterable[any]{
		Seq: func(yield func(any) bool) {
			iterable.Seq(func(item T) bool {
				return yield(item)
			})
		},
	}

	/*
		linq.AsAnyIterable(
			linq.From([]T{...}),
		)
	*/
}
