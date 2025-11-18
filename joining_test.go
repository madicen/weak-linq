package weaklinq

import (
	"iter"
	"testing"
)

//----------------------------------------------------------------------------//
// Joining                                                                    //
//----------------------------------------------------------------------------//

////////////////////////////////////////////////////////////////////////////////

func TestJoin(t *testing.T) {
	testItems := []testStruct{
		{Id: 1, Name: "Test 1"},
		{Id: 2, Name: "Test 2"},
	}
	testItems2 := []testStruct{
		{Id: 1, Name: "Test A"},
		{Id: 2, Name: "Test B"},
	}

	result := From(testItems).Join(From(testItems2).AsAny())
	resultIterator, _ := iter.Pull(result.itemIterable.Seq)

	if resultIterator == nil {
		t.Errorf("Expected iterator but got nil")
	}

	if item, ok := resultIterator(); ok && item != testItems[0] {
		t.Errorf("Expected first item but got %v", item)
	}

	if item, ok := resultIterator(); ok && item != testItems[1] {
		t.Errorf("Expected second item but got %v", item)
	}

	if _, ok := resultIterator(); ok {
		t.Errorf("Expected no item but got an item")
	}

	joinIterator, _ := iter.Pull(result.rightIterable.Seq)

	if joinIterator == nil {
		t.Errorf("Expected join iterator but got nil")
	}

	if item, ok := joinIterator(); ok && item != testItems2[0] {
		t.Errorf("Expected first joined item but got %v", item)
	}

	if item, ok := joinIterator(); ok && item != testItems2[1] {
		t.Errorf("Expected second joined item but got %v", item)
	}

	if _, ok := joinIterator(); ok {
		t.Errorf("Expected no joined item but got an item")
	}

}

////////////////////////////////////////////////////////////////////////////////

func TestJoinSlice(t *testing.T) {
	testItems := []testStruct{
		{Id: 1, Name: "Test 1"},
		{Id: 2, Name: "Test 2"},
	}
	testItems2 := []testStruct{
		{Id: 1, Name: "Test A"},
		{Id: 2, Name: "Test B"},
	}

	result := From(testItems).JoinSlice(testItems2)
	resultIterator, _ := iter.Pull(result.itemIterable.Seq)

	if resultIterator == nil {
		t.Errorf("Expected iterator but got nil")
	}

	if item, ok := resultIterator(); ok && item != testItems[0] {
		t.Errorf("Expected first item but got %v", item)
	}

	if item, ok := resultIterator(); ok && item != testItems[1] {
		t.Errorf("Expected second item but got %v", item)
	}

	if _, ok := resultIterator(); ok {
		t.Errorf("Expected no item but got an item")
	}

	joinIterator, _ := iter.Pull(result.rightIterable.Seq)

	if joinIterator == nil {
		t.Errorf("Expected join iterator but got nil")
	}

	if item, ok := joinIterator(); ok && item != testItems2[0] {
		t.Errorf("Expected first joined item but got %v", item)
	}

	if item, ok := joinIterator(); ok && item != testItems2[1] {
		t.Errorf("Expected second joined item but got %v", item)
	}

	if _, ok := joinIterator(); ok {
		t.Errorf("Expected no joined item but got an item")
	}

}

////////////////////////////////////////////////////////////////////////////////

func TestOnThis(t *testing.T) {
	testItems := []testStruct{
		{Id: 1, Name: "Test 1"},
		{Id: 2, Name: "Test 2"},
	}
	testItems2 := []testStruct{
		{Id: 1, Name: "Test A"},
		{Id: 2, Name: "Test B"},
	}

	result := From(testItems).
		JoinSlice(testItems2).
		OnThis(func(t testStruct) any { return t.Id })
	resultIterator, _ := iter.Pull(result.itemIterable.Seq)

	if result.keySelector(testItems[0]) != testItems[0].Id {
		t.Errorf("Expected first item but got %v", result.keySelector(testItems[0]))
	}

	if resultIterator == nil {
		t.Errorf("Expected iterator but got nil")
	}

	if item, ok := resultIterator(); ok && item != testItems[0] {
		t.Errorf("Expected first item but got %v", item)
	}

	if result.rightKeySelector(testItems2[0]) != testItems2[0].Id {
		t.Errorf("Expected first item but got %v", result.rightKeySelector(testItems2[0]))
	}
}

////////////////////////////////////////////////////////////////////////////////

func TestOn(t *testing.T) {

	//----------------------------------------------------------------------------//

	t.Run("generic", func(t *testing.T) {
		testItems := []testStruct{
			{Id: 1, Name: "Test 1"},
			{Id: 2, Name: "Test 2"},
		}
		testItems2 := []testStruct{
			{Id: 1, Name: "Test A"},
			{Id: 2, Name: "Test B"},
		}

		result := From(testItems).
			JoinSlice(testItems2).
			On("Id")
		resultIterator, _ := iter.Pull(result.itemIterable.Seq)

		if result.keySelector(testItems[0]) != testItems[0].Id {
			t.Errorf("Expected first item but got %v", result.keySelector(testItems[0]))
		}

		if resultIterator == nil {
			t.Errorf("Expected iterator but got nil")
		}

		if item, ok := resultIterator(); ok && item != testItems[0] {
			t.Errorf("Expected first item but got %v", item)
		}

		if result.rightKeySelector(testItems2[0]) != testItems2[0].Id {
			t.Errorf("Expected first item but got %v", result.rightKeySelector(testItems2[0]))
		}
	})

	//----------------------------------------------------------------------------//

	t.Run("bad field name", func(t *testing.T) {
		testItems := []testStruct{
			{Id: 1, Name: "Test 1"},
		}
		testItems2 := []testStruct{
			{Id: 1, Name: "Test A"},
		}

		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Expected panic but got %v", err)
			}
		}()

		result := From(testItems).
			JoinSlice(testItems2).
			On("BadFieldName")

		result.keySelector(testItems[0])
	})

	//----------------------------------------------------------------------------//
}

////////////////////////////////////////////////////////////////////////////////

func TestEqualsThis(t *testing.T) {
	testItems := []testStruct{
		{Id: 1, Name: "Test 1"},
		{Id: 2, Name: "Test 2"},
	}
	testItems2 := []testStruct{
		{Id: 1, Name: "Test A"},
		{Id: 2, Name: "Test B"},
	}

	result := From(testItems).
		JoinSlice(testItems2).
		On("Id").
		EqualsThis(func(item any) any {
			return item.(testStruct).Name
		})

	if result.keySelector(testItems[0]) != testItems[0].Id {
		t.Errorf("Expected first item but got %v", result.keySelector(testItems[0]))
	}

	if result.rightKeySelector(testItems2[0]) != testItems2[0].Name {
		t.Errorf("Expected first item but got %v", result.rightKeySelector(testItems2[0]))
	}
}

////////////////////////////////////////////////////////////////////////////////

func TestEquals(t *testing.T) {

	//----------------------------------------------------------------------------//

	t.Run("generic", func(t *testing.T) {
		testItems := []testStruct{
			{Id: 1, Name: "Test 1"},
			{Id: 2, Name: "Test 2"},
		}
		testItems2 := []testStruct{
			{Id: 1, Name: "Test A"},
			{Id: 2, Name: "Test B"},
		}

		result := From(testItems).
			JoinSlice(testItems2).
			On("Id").
			Equals("Name")

		if result.keySelector(testItems[0]) != testItems[0].Id {
			t.Errorf("Expected first item but got %v", result.keySelector(testItems[0]))
		}

		if result.rightKeySelector(testItems2[0]) != testItems2[0].Name {
			t.Errorf("Expected first item but got %v", result.rightKeySelector(testItems2[0]))
		}
	})

	//----------------------------------------------------------------------------//

	t.Run("bad field name", func(t *testing.T) {
		testItems := []testStruct{
			{Id: 1, Name: "Test 1"},
		}
		testItems2 := []testStruct{
			{Id: 1, Name: "Test A"},
		}

		defer func() {
			if err := recover(); err == nil {
				t.Errorf("Expected panic but got %v", err)
			}
		}()

		result := From(testItems).
			JoinSlice(testItems2).
			On("Id").
			Equals("BadFieldName")

		result.rightKeySelector(testItems2[0])
	})

	//----------------------------------------------------------------------------//
}

////////////////////////////////////////////////////////////////////////////////

func TestAsThis(t *testing.T) {
	testItems := []testStruct{
		{Id: 1, Name: "Test 1"},
		{Id: 2, Name: "Test 2"},
		{Id: 4, Name: "Test 4"},
	}
	testItems2 := []testStruct{
		{Id: 1, Name: "Test A"},
		{Id: 2, Name: "Test B"},
		{Id: 3, Name: "Test C"},
	}

	result := From(testItems).
		JoinSlice(testItems2).
		On("Id").
		AsThis(func(left testStruct, right any) any {
			return struct {
				Name      string
				RightName string
			}{
				Name:      left.Name,
				RightName: right.(testStruct).Name,
			}
		})

	resultIterator, _ := iter.Pull(result.Seq)

	if resultIterator == nil {
		t.Errorf("Expected iterator but got nil")
	}

	if item, ok := resultIterator(); ok {
		combined := item.(struct {
			Name      string
			RightName string
		})

		if combined.Name != testItems[0].Name {
			t.Errorf("Expected first left name but got %v", combined.Name)
		}

		if combined.RightName != testItems2[0].Name {
			t.Errorf("Expected first right name but got %v", combined.RightName)
		}
	}

	if item, ok := resultIterator(); ok {
		combined := item.(struct {
			Name      string
			RightName string
		})

		if combined.Name != testItems[1].Name {
			t.Errorf("Expected second left name but got %v", combined.Name)
		}

		if combined.RightName != testItems2[1].Name {
			t.Errorf("Expected second right name but got %v", combined.RightName)
		}
	}

	if _, ok := resultIterator(); ok {
		t.Errorf("Expected no item but got an item")
	}
}

////////////////////////////////////////////////////////////////////////////////

func TestAsPairs(t *testing.T) {
	testItems := []testStruct{
		{Id: 1, Name: "Test 1"},
		{Id: 2, Name: "Test 2"},
	}
	testItems2 := []testStruct{
		{Id: 1, Name: "Test A"},
		{Id: 2, Name: "Test B"},
	}

	result := make([]Pair[testStruct, any], 0)
	From(testItems).
		JoinSlice(testItems2).
		On("Id").
		AsPairs().
		AndAssignToSlice(&result)

	if result[0].Left != testItems[0] {
		t.Errorf("Expected first left item but got %v", result[0].Left)
	}

	if result[0].Right != testItems2[0] {
		t.Errorf("Expected first right item but got %v", result[0].Right)
	}

	if result[1].Left != testItems[1] {
		t.Errorf("Expected second left item but got %v", result[1].Left)
	}

	if result[1].Right != testItems2[1] {
		t.Errorf("Expected second right item but got %v", result[1].Right)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 items but got %v", len(result))
	}
}
