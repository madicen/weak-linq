package weaklinq

import (
	"fmt"
	"iter"
	"reflect"
	"slices"
)

//----------------------------------------------------------------------------//
// Common                                                                     //
//----------------------------------------------------------------------------//

////////////////////////////////////////////////////////////////////////////////

// Iterable is the base structure for an iterable. It allows for the lazy
// iteration of a collection of items and exists to allow functions to be
// called on the collection.
type Iterable[T any] struct {
	iter.Seq[T]
}

////////////////////////////////////////////////////////////////////////////////

// identitySelector is a function that returns the item passed to it. Used as
// the default selector for many functions.
func identitySelector[T any](item T) any {
	return item
}

////////////////////////////////////////////////////////////////////////////////

// getFieldNameFunc returns a function that returns the value of the given
// field name. If T is not a struct or pointer to struct, or fieldName is
// not found, this function will panic.
func getFieldNameFunc[T any](fieldName string) func(T) any {

	return func(item T) any {

		res := reflect.ValueOf(item)
		if res.Kind() != reflect.Struct {

			if res.Kind() == reflect.Pointer {
				if reflect.Indirect(res).Kind() != reflect.Struct {
					panic(fmt.Sprintf("item is a pointer, but not to a struct: %T", item))
				}

				res = res.Elem()

			} else {
				panic(fmt.Sprintf("item is not a struct or pointer to struct: %T", item))
			}
		}

		field := res.FieldByName(fieldName)
		if !field.IsValid() {
			panic(fmt.Sprintf("field name '%s' not found in struct %T", fieldName, item))
		}

		return field.Interface()
	}
}

//----------------------------------------------------------------------------//
// Constructors                                                               //
//----------------------------------------------------------------------------//

////////////////////////////////////////////////////////////////////////////////

// From creates a new Iterable from a slice of items.
func From[T any](items []T) Iterable[T] {

	return Iterable[T]{
		Seq: slices.Values(items),
	}

	/*
		linq.From([]T{...})
	*/
}
