package weaklinq

import (
	"iter"
	"testing"
)

//----------------------------------------------------------------------------//
// Flatten                                                                    //
//----------------------------------------------------------------------------//

////////////////////////////////////////////////////////////////////////////////

func TestGetFlattenedThese(t *testing.T) {

	//----------------------------------------------------------------------------//

	t.Run("generic", func(t *testing.T) {

		testItems := []testStruct{
			{Id: 1, Name: "Test 1"},
			{Id: 2, Name: "Test 2"},
		}

		// Flatten by creating an iterable of Ids and Names for each item
		result := From(testItems).FlattenThese(
			func(item testStruct) Iterable[any] {
				return From([]any{item.Id, item.Name})
			},
		)
		resultIterator, _ := iter.Pull(result.Seq)

		if resultIterator == nil {
			t.Errorf("Expected iterator but got nil")
		}

		// First item's Id and Name
		if item, ok := resultIterator(); ok && item != testItems[0].Id {
			t.Errorf("Expected first item Id but got %v", item)
		}

		if item, ok := resultIterator(); ok && item != testItems[0].Name {
			t.Errorf("Expected first item Name but got %v", item)
		}

		// Second item's Id and Name
		if item, ok := resultIterator(); ok && item != testItems[1].Id {
			t.Errorf("Expected second item Id but got %v", item)
		}

		if item, ok := resultIterator(); ok && item != testItems[1].Name {
			t.Errorf("Expected second item Name but got %v", item)
		}

		if item, ok := resultIterator(); ok {
			t.Errorf("Expected no item but got %v", item)
		}
	})

	//----------------------------------------------------------------------------//

	t.Run("empty source", func(t *testing.T) {

		var testItems []testStruct

		result := From(testItems).FlattenThese(
			func(item testStruct) Iterable[any] {
				return From([]any{item.Id, item.Name})
			},
		)
		resultIterator, _ := iter.Pull(result.Seq)

		if resultIterator == nil {
			t.Errorf("Expected iterator but got nil")
		}

		if item, ok := resultIterator(); ok {
			t.Errorf("Expected no item but got %v", item)
		}
	})

	//----------------------------------------------------------------------------//

	t.Run("empty selector results", func(t *testing.T) {

		testItems := []testStruct{
			{Id: 1, Name: "Test 1"},
			{Id: 2, Name: "Test 2"},
		}

		// Selector returns empty iterable for each item
		result := From(testItems).FlattenThese(
			func(item testStruct) Iterable[any] {
				return From([]any{})
			},
		)
		resultIterator, _ := iter.Pull(result.Seq)

		if resultIterator == nil {
			t.Errorf("Expected iterator but got nil")
		}

		if item, ok := resultIterator(); ok {
			t.Errorf("Expected no item but got %v", item)
		}
	})

	//----------------------------------------------------------------------------//

	t.Run("varying selector results", func(t *testing.T) {

		testItems := []testStruct{
			{Id: 1, Name: "Test 1"},
			{Id: 2, Name: "Test 2"},
			{Id: 3, Name: "Test 3"},
		}

		// Selector returns different number of items based on Id
		result := From(testItems).FlattenThese(
			func(item testStruct) Iterable[any] {
				items := make([]any, item.Id)
				for i := 0; i < item.Id; i++ {
					items[i] = item.Name
				}
				return From(items)
			},
		)
		resultIterator, _ := iter.Pull(result.Seq)

		if resultIterator == nil {
			t.Errorf("Expected iterator but got nil")
		}

		// First item: 1 occurrence of "Test 1"
		if item, ok := resultIterator(); ok && item != "Test 1" {
			t.Errorf("Expected 'Test 1' but got %v", item)
		}

		// Second item: 2 occurrences of "Test 2"
		if item, ok := resultIterator(); ok && item != "Test 2" {
			t.Errorf("Expected 'Test 2' but got %v", item)
		}

		if item, ok := resultIterator(); ok && item != "Test 2" {
			t.Errorf("Expected 'Test 2' but got %v", item)
		}

		// Third item: 3 occurrences of "Test 3"
		if item, ok := resultIterator(); ok && item != "Test 3" {
			t.Errorf("Expected 'Test 3' but got %v", item)
		}

		if item, ok := resultIterator(); ok && item != "Test 3" {
			t.Errorf("Expected 'Test 3' but got %v", item)
		}

		if item, ok := resultIterator(); ok && item != "Test 3" {
			t.Errorf("Expected 'Test 3' but got %v", item)
		}

		if item, ok := resultIterator(); ok {
			t.Errorf("Expected no item but got %v", item)
		}
	})

	//----------------------------------------------------------------------------//

	t.Run("nested collections", func(t *testing.T) {

		type parent struct {
			Id       int
			Children []int
		}

		testItems := []parent{
			{Id: 1, Children: []int{10, 11}},
			{Id: 2, Children: []int{20, 21, 22}},
		}

		result := From(testItems).FlattenThese(
			func(item parent) Iterable[any] {
				return From(item.Children).AsAny()
			},
		)
		resultIterator, _ := iter.Pull(result.Seq)

		if resultIterator == nil {
			t.Errorf("Expected iterator but got nil")
		}

		// First parent's children
		if item, ok := resultIterator(); ok && item != 10 {
			t.Errorf("Expected 10 but got %v", item)
		}

		if item, ok := resultIterator(); ok && item != 11 {
			t.Errorf("Expected 11 but got %v", item)
		}

		// Second parent's children
		if item, ok := resultIterator(); ok && item != 20 {
			t.Errorf("Expected 20 but got %v", item)
		}

		if item, ok := resultIterator(); ok && item != 21 {
			t.Errorf("Expected 21 but got %v", item)
		}

		if item, ok := resultIterator(); ok && item != 22 {
			t.Errorf("Expected 22 but got %v", item)
		}

		if item, ok := resultIterator(); ok {
			t.Errorf("Expected no item but got %v", item)
		}
	})

	//----------------------------------------------------------------------------//
}

////////////////////////////////////////////////////////////////////////////////

func TestGetFlattened(t *testing.T) {

	//----------------------------------------------------------------------------//

	t.Run("generic", func(t *testing.T) {

		testItems := []testStruct{
			{Id: 1, Name: "Test 1"},
			{Id: 2, Name: "Test 2"},
		}

		// Note: GetFlattened calls iterable.Get(fieldName) for each source item,
		// which returns an iterable of all field values from the source.
		// This results in [Name1, Name2] for item1, then [Name1, Name2] for item2
		result := From(testItems).Flatten("Name")
		resultIterator, _ := iter.Pull(result.Seq)

		if resultIterator == nil {
			t.Errorf("Expected iterator but got nil")
		}

		// First iteration: yields all Names (Test 1, Test 2)
		if item, ok := resultIterator(); ok && item != testItems[0].Name {
			t.Errorf("Expected %v but got %v", testItems[0].Name, item)
		}

		if item, ok := resultIterator(); ok && item != testItems[1].Name {
			t.Errorf("Expected %v but got %v", testItems[1].Name, item)
		}

		// Second iteration: yields all Names again (Test 1, Test 2)
		if item, ok := resultIterator(); ok && item != testItems[0].Name {
			t.Errorf("Expected %v but got %v", testItems[0].Name, item)
		}

		if item, ok := resultIterator(); ok && item != testItems[1].Name {
			t.Errorf("Expected %v but got %v", testItems[1].Name, item)
		}

		if item, ok := resultIterator(); ok {
			t.Errorf("Expected no item but got %v", item)
		}
	})

	//----------------------------------------------------------------------------//

	t.Run("nested structs", func(t *testing.T) {

		type testStruct struct {
			Id   int
			Ints []int
		}

		testItems := []testStruct{
			{Id: 1, Ints: []int{1, 2}},
			{Id: 2, Ints: []int{3, 4}},
		}

		result := From(testItems).Flatten("Ints")
		resultIterator, _ := iter.Pull(result.Seq)

		if resultIterator == nil {
			t.Errorf("Expected iterator but got nil")
		}

		if item, ok := resultIterator(); ok && item != 1 {
			t.Errorf("Expected 1 but got %v", item)
		}

		if item, ok := resultIterator(); ok && item != 2 {
			t.Errorf("Expected 2 but got %v", item)
		}

		if item, ok := resultIterator(); ok && item != 3 {
			t.Errorf("Expected 3 but got %v", item)
		}

		if item, ok := resultIterator(); ok && item != 4 {
			t.Errorf("Expected 4 but got %v", item)
		}

		if item, ok := resultIterator(); ok {
			t.Errorf("Expected no item but got %v", item)
		}
	})

	//----------------------------------------------------------------------------//

	t.Run("bad field name", func(t *testing.T) {

		testItems := []testStruct{
			{Id: 1, Name: "Test 1"},
		}

		result := From(testItems).Flatten("BadFieldName")
		resultIterator, _ := iter.Pull(result.Seq)

		if resultIterator == nil {
			t.Errorf("Expected iterator but got nil")
		}

		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Expected panic but got %v", err)
			}
		}()

		resultIterator()
	})

	//----------------------------------------------------------------------------//

	t.Run("bad item type", func(t *testing.T) {

		testItems := []int{1, 2, 3}

		result := From(testItems).Flatten("Name")
		resultIterator, _ := iter.Pull(result.Seq)

		if resultIterator == nil {
			t.Errorf("Expected iterator but got nil")
		}

		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Expected panic but got %v", err)
			}
		}()

		resultIterator()
	})

	//----------------------------------------------------------------------------//
}

////////////////////////////////////////////////////////////////////////////////
