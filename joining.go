package weaklinq

import "reflect"

//----------------------------------------------------------------------------//
// Joining                                                                    //
//----------------------------------------------------------------------------//

////////////////////////////////////////////////////////////////////////////////

// JoinIterable is a specialized iterable that includes a right iterable, a key
// selector for the left iterable, and a key selector for the right iterable.
// These selectors are stored until the collection is iterated, and then
// applied to the items.
type JoinIterable[T any] struct {
	itemIterable     Iterable[T]
	rightIterable    Iterable[any]
	keySelector      func(T) any
	rightKeySelector func(any) any
}

/////////////////////////////////////////////////////////////////////////////////

// DeferredJoinIterable is a JoinIterable where the key selectors have not yet
// been set. Designed to be used in tandem with the On and Equals functions.
// Has very little use outside of that.
type DeferredJoinIterable[T any] JoinIterable[T]

/////////////////////////////////////////////////////////////////////////////////

// Pair is a struct that holds a left and right item for use in join operations.
type Pair[TLeft any, TRight any] struct {
	Left  TLeft
	Right TRight
}

////////////////////////////////////////////////////////////////////////////////

func defaultJoinIterable[T any](iterable Iterable[T], joinIterable Iterable[any]) JoinIterable[T] {

	return JoinIterable[T]{
		itemIterable:     iterable,
		rightIterable:    joinIterable,
		keySelector:      identitySelector[T],
		rightKeySelector: identitySelector[any],
	}
}

////////////////////////////////////////////////////////////////////////////////

// Join returns a new DeferredJoinIterable that will join the items of the given iterable
func (iterable Iterable[T]) Join(joinIterable Iterable[any]) DeferredJoinIterable[T] {

	return DeferredJoinIterable[T](defaultJoinIterable(iterable, joinIterable))

	/*
		linq.From([]T{...}).
			Join(linq.From([]TRight{...}).AsAny())
	*/
}

////////////////////////////////////////////////////////////////////////////////

// JoinSlice returns a new DeferredJoinIterable that will join the items of the given slice
func (iterable Iterable[T]) JoinSlice(joinSlice any) DeferredJoinIterable[T] {

	if reflect.TypeOf(joinSlice).Kind() != reflect.Slice {
		panic("'joinSlice' must be a slice")
	}

	joinIterable := Iterable[any]{
		Seq: func(yield func(any) bool) {
			s := reflect.ValueOf(joinSlice)
			for i := 0; i < s.Len(); i++ {
				if !yield(s.Index(i).Interface()) {
					return
				}
			}
		},
	}

	return iterable.Join(joinIterable)

	/*
		linq.From([]T{...}).
			JoinSlice([]TRight{...})
	*/
}

////////////////////////////////////////////////////////////////////////////////

// OnThis sets the key selector functions for both left and right iterables.
func (iterable DeferredJoinIterable[T]) OnThis(keySelector func(T) any) DeferredJoinIterable[T] {

	iterable.keySelector = keySelector
	iterable.rightKeySelector = func(item any) any {
		return keySelector(item.(T))
	}

	return iterable

	/*
		linq.From([]T{...}).
			Join(linq.From([]TRight{...}).AsAny()).
			OnThis(
				func(left T) any {
					return ...
				},
			)
	*/
}

////////////////////////////////////////////////////////////////////////////////

// On sets the key selector functions for both left and right iterables using
// the given field name.
func (iterable DeferredJoinIterable[T]) On(fieldName string) DeferredJoinIterable[T] {

	return iterable.OnThis(
		getFieldNameFunc[T](fieldName),
	)

	/*
		linq.From([]T{...}).
			Join(linq.From([]TRight{...}).AsAny()).
			On("KeyField")
	*/
}

////////////////////////////////////////////////////////////////////////////////

// EqualsThis sets the right key selector function.
func (iterable DeferredJoinIterable[T]) EqualsThis(rightKeySelector func(any) any) DeferredJoinIterable[T] {

	iterable.rightKeySelector = rightKeySelector

	return iterable

	/*
		linq.From([]T{...}).
			Join(linq.From([]TRight{...}).AsAny()).
			On("LeftKeyField").
			EqualsThis(
				func(right any) any {
					return ...
				},
			)
	*/
}

////////////////////////////////////////////////////////////////////////////////

// Equals sets the right key selector function using the given field name.
func (iterable DeferredJoinIterable[T]) Equals(fieldName string) DeferredJoinIterable[T] {

	return iterable.EqualsThis(
		getFieldNameFunc[any](fieldName),
	)

	/*
		linq.From([]T{...}).
			Join(linq.From([]TRight{...}).AsAny()).
			On("LeftKeyField").
			Equals("RightKeyField")
	*/
}

////////////////////////////////////////////////////////////////////////////////

// AsThis projects the joined items using the given joinSelector function.
func (iterable DeferredJoinIterable[T]) AsThis(joinSelector func(T, any) any) Iterable[any] {

	rightKeysToRightItems := make(map[any][]any)
	for rightItem := range iterable.rightIterable.Seq {
		rightKey := iterable.rightKeySelector(rightItem)
		rightKeysToRightItems[rightKey] = append(rightKeysToRightItems[rightKey], rightItem)
	}

	return Iterable[any]{
		Seq: func(yield func(any) bool) {
			iterable.itemIterable.Seq(func(leftItem T) bool {

				leftKey := iterable.keySelector(leftItem)

				if rightItems, ok := rightKeysToRightItems[leftKey]; ok {
					for _, rightItem := range rightItems {
						result := joinSelector(leftItem, rightItem)
						if !yield(result) {
							return false
						}
					}
				}

				return true
			})
		},
	}

	/*
		linq.From([]T{...}).
			Join(linq.From([]TRight{...}).AsAny()).
			On("LeftKeyField").
			Equals("RightKeyField").
			AsThis(
				func(left T, right any) any {
					return ...
				},
			)
	*/

}

////////////////////////////////////////////////////////////////////////////////

// AsPairs projects the joined items as Pair structs.
func (iterable DeferredJoinIterable[T]) AsPairs() Iterable[any] {

	return iterable.AsThis(
		func(left T, right any) any {
			return Pair[T, any]{
				Left:  left,
				Right: right,
			}
		},
	)

	/*
		linq.From([]T{...}).
			Join(linq.From([]TRight{...}).AsAny()).
			On("LeftKeyField").
			Equals("RightKeyField").
			AsPairs()
	*/
}
