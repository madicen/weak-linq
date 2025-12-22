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

////////////////////////////////////////////////////////////////////////////////

func TestLeftJoin(t *testing.T) {
	t.Run("all items match", func(t *testing.T) {
		left := []testStruct{
			{Id: 1, Name: "Left 1"},
			{Id: 2, Name: "Left 2"},
		}
		right := []testStruct{
			{Id: 1, Name: "Right 1"},
			{Id: 2, Name: "Right 2"},
		}

		result := make([]Pair[testStruct, any], 0)
		From(left).
			LeftJoinSlice(right).
			On("Id").
			AsPairs().
			AndAssignToSlice(&result)

		if len(result) != 2 {
			t.Errorf("Expected 2 items but got %v", len(result))
		}

		if result[0].Right == nil {
			t.Errorf("Expected Right 1 but got nil")
		}
	})

	t.Run("left has unmatched items", func(t *testing.T) {
		left := []testStruct{
			{Id: 1, Name: "Left 1"},
			{Id: 2, Name: "Left 2"},
			{Id: 4, Name: "Left 4"},
		}
		right := []testStruct{
			{Id: 1, Name: "Right 1"},
			{Id: 2, Name: "Right 2"},
			{Id: 3, Name: "Right 3"},
		}

		result := make([]Pair[testStruct, any], 0)
		From(left).
			LeftJoinSlice(right).
			On("Id").
			AsPairs().
			AndAssignToSlice(&result)

		// All 3 left items should be present
		if len(result) != 3 {
			t.Errorf("Expected 3 items but got %v", len(result))
		}

		// First two should have matches
		if result[0].Left.Name != "Left 1" || result[0].Right.(testStruct).Name != "Right 1" {
			t.Errorf("Expected Left 1 and Right 1 but got %v and %v", result[0].Left.Name, result[0].Right)
		}

		if result[1].Left.Name != "Left 2" || result[1].Right.(testStruct).Name != "Right 2" {
			t.Errorf("Expected Left 2 and Right 2 but got %v and %v", result[1].Left.Name, result[1].Right)
		}

		// Third should have nil right
		if result[2].Left.Name != "Left 4" || result[2].Right != nil {
			t.Errorf("Expected Left 4 and nil but got %v and %v", result[2].Left.Name, result[2].Right)
		}
	})

	t.Run("no matches", func(t *testing.T) {
		left := []testStruct{
			{Id: 1, Name: "Left 1"},
			{Id: 2, Name: "Left 2"},
		}
		right := []testStruct{
			{Id: 3, Name: "Right 3"},
			{Id: 4, Name: "Right 4"},
		}

		result := make([]Pair[testStruct, any], 0)
		From(left).
			LeftJoinSlice(right).
			On("Id").
			AsPairs().
			AndAssignToSlice(&result)

		// All left items should be present with nil right
		if len(result) != 2 {
			t.Errorf("Expected 2 items but got %v", len(result))
		}

		if result[0].Left.Name != "Left 1" || result[0].Right != nil {
			t.Errorf("Expected Left 1 and nil but got %v and %v", result[0].Left.Name, result[0].Right)
		}

		if result[1].Left.Name != "Left 2" || result[1].Right != nil {
			t.Errorf("Expected Left 2 and nil but got %v and %v", result[1].Left.Name, result[1].Right)
		}
	})
}

////////////////////////////////////////////////////////////////////////////////

func TestRightJoin(t *testing.T) {
	t.Run("all items match", func(t *testing.T) {
		left := []testStruct{
			{Id: 1, Name: "Left 1"},
			{Id: 2, Name: "Left 2"},
		}
		right := []testStruct{
			{Id: 1, Name: "Right 1"},
			{Id: 2, Name: "Right 2"},
		}

		result := make([]Pair[testStruct, any], 0)
		From(left).
			RightJoinSlice(right).
			On("Id").
			AsPairs().
			AndAssignToSlice(&result)

		if len(result) != 2 {
			t.Errorf("Expected 2 items but got %v", len(result))
		}

		if result[0].Left.Id == 0 {
			t.Errorf("Expected Left 1 but got zero value")
		}
	})

	t.Run("right has unmatched items", func(t *testing.T) {
		left := []testStruct{
			{Id: 1, Name: "Left 1"},
			{Id: 2, Name: "Left 2"},
			{Id: 4, Name: "Left 4"},
		}
		right := []testStruct{
			{Id: 1, Name: "Right 1"},
			{Id: 2, Name: "Right 2"},
			{Id: 3, Name: "Right 3"},
		}

		result := make([]Pair[testStruct, any], 0)
		From(left).
			RightJoinSlice(right).
			On("Id").
			AsPairs().
			AndAssignToSlice(&result)

		// All 3 right items should be present
		if len(result) != 3 {
			t.Errorf("Expected 3 items but got %v", len(result))
		}

		// Find the matched and unmatched items
		matchedCount := 0
		unmatchedCount := 0
		for _, pair := range result {
			if pair.Left.Id == 0 {
				unmatchedCount++
				// Should be Right 3 (unmatched)
				if pair.Right.(testStruct).Id != 3 {
					t.Errorf("Expected unmatched right item to have Id 3 but got %v", pair.Right.(testStruct).Id)
				}
			} else {
				matchedCount++
			}
		}

		if matchedCount != 2 {
			t.Errorf("Expected 2 matched items but got %v", matchedCount)
		}

		if unmatchedCount != 1 {
			t.Errorf("Expected 1 unmatched item but got %v", unmatchedCount)
		}
	})

	t.Run("no matches", func(t *testing.T) {
		left := []testStruct{
			{Id: 1, Name: "Left 1"},
			{Id: 2, Name: "Left 2"},
		}
		right := []testStruct{
			{Id: 3, Name: "Right 3"},
			{Id: 4, Name: "Right 4"},
		}

		result := make([]Pair[testStruct, any], 0)
		From(left).
			RightJoinSlice(right).
			On("Id").
			AsPairs().
			AndAssignToSlice(&result)

		// All right items should be present with zero-value left
		if len(result) != 2 {
			t.Errorf("Expected 2 items but got %v", len(result))
		}

		if result[0].Left.Id != 0 || result[0].Right.(testStruct).Name != "Right 3" {
			t.Errorf("Expected zero-value left and Right 3 but got %v and %v", result[0].Left, result[0].Right)
		}

		if result[1].Left.Id != 0 || result[1].Right.(testStruct).Name != "Right 4" {
			t.Errorf("Expected zero-value left and Right 4 but got %v and %v", result[1].Left, result[1].Right)
		}
	})
}

////////////////////////////////////////////////////////////////////////////////

func TestFullOuterJoin(t *testing.T) {
	t.Run("all items match", func(t *testing.T) {
		left := []testStruct{
			{Id: 1, Name: "Left 1"},
			{Id: 2, Name: "Left 2"},
		}
		right := []testStruct{
			{Id: 1, Name: "Right 1"},
			{Id: 2, Name: "Right 2"},
		}

		result := make([]Pair[testStruct, any], 0)
		From(left).
			FullOuterJoinSlice(right).
			On("Id").
			AsPairs().
			AndAssignToSlice(&result)

		if len(result) != 2 {
			t.Errorf("Expected 2 items but got %v", len(result))
		}

		// All should be matched pairs
		for i, pair := range result {
			if pair.Left.Id == 0 || pair.Right == nil {
				t.Errorf("Expected matched pair at index %v but got %v and %v", i, pair.Left, pair.Right)
			}
		}
	})

	t.Run("both sides have unmatched items", func(t *testing.T) {
		left := []testStruct{
			{Id: 1, Name: "Left 1"},
			{Id: 2, Name: "Left 2"},
			{Id: 4, Name: "Left 4"},
		}
		right := []testStruct{
			{Id: 1, Name: "Right 1"},
			{Id: 2, Name: "Right 2"},
			{Id: 3, Name: "Right 3"},
		}

		result := make([]Pair[testStruct, any], 0)
		From(left).
			FullOuterJoinSlice(right).
			On("Id").
			AsPairs().
			AndAssignToSlice(&result)

		// Should have 2 matched + 1 left-only + 1 right-only = 4 items
		if len(result) != 4 {
			t.Errorf("Expected 4 items but got %v", len(result))
		}

		matchedCount := 0
		leftOnlyCount := 0
		rightOnlyCount := 0

		for _, pair := range result {
			if pair.Left.Id != 0 && pair.Right != nil {
				matchedCount++
			} else if pair.Left.Id != 0 && pair.Right == nil {
				leftOnlyCount++
				if pair.Left.Id != 4 {
					t.Errorf("Expected left-only item to have Id 4 but got %v", pair.Left.Id)
				}
			} else if pair.Left.Id == 0 && pair.Right != nil {
				rightOnlyCount++
				if pair.Right.(testStruct).Id != 3 {
					t.Errorf("Expected right-only item to have Id 3 but got %v", pair.Right.(testStruct).Id)
				}
			}
		}

		if matchedCount != 2 {
			t.Errorf("Expected 2 matched items but got %v", matchedCount)
		}

		if leftOnlyCount != 1 {
			t.Errorf("Expected 1 left-only item but got %v", leftOnlyCount)
		}

		if rightOnlyCount != 1 {
			t.Errorf("Expected 1 right-only item but got %v", rightOnlyCount)
		}
	})

	t.Run("no matches", func(t *testing.T) {
		left := []testStruct{
			{Id: 1, Name: "Left 1"},
			{Id: 2, Name: "Left 2"},
		}
		right := []testStruct{
			{Id: 3, Name: "Right 3"},
			{Id: 4, Name: "Right 4"},
		}

		result := make([]Pair[testStruct, any], 0)
		From(left).
			FullOuterJoinSlice(right).
			On("Id").
			AsPairs().
			AndAssignToSlice(&result)

		// All items from both sides should be present
		if len(result) != 4 {
			t.Errorf("Expected 4 items but got %v", len(result))
		}

		leftOnlyCount := 0
		rightOnlyCount := 0

		for _, pair := range result {
			if pair.Left.Id != 0 && pair.Right == nil {
				leftOnlyCount++
			} else if pair.Left.Id == 0 && pair.Right != nil {
				rightOnlyCount++
			}
		}

		if leftOnlyCount != 2 {
			t.Errorf("Expected 2 left-only items but got %v", leftOnlyCount)
		}

		if rightOnlyCount != 2 {
			t.Errorf("Expected 2 right-only items but got %v", rightOnlyCount)
		}
	})

	t.Run("empty left side", func(t *testing.T) {
		left := []testStruct{}
		right := []testStruct{
			{Id: 1, Name: "Right 1"},
			{Id: 2, Name: "Right 2"},
		}

		result := make([]Pair[testStruct, any], 0)
		From(left).
			FullOuterJoinSlice(right).
			On("Id").
			AsPairs().
			AndAssignToSlice(&result)

		// Should have all right items with zero-value left
		if len(result) != 2 {
			t.Errorf("Expected 2 items but got %v", len(result))
		}

		for _, pair := range result {
			if pair.Left.Id != 0 {
				t.Errorf("Expected zero-value left but got %v", pair.Left)
			}
		}
	})

	t.Run("empty right side", func(t *testing.T) {
		left := []testStruct{
			{Id: 1, Name: "Left 1"},
			{Id: 2, Name: "Left 2"},
		}
		right := []testStruct{}

		result := make([]Pair[testStruct, any], 0)
		From(left).
			FullOuterJoinSlice(right).
			On("Id").
			AsPairs().
			AndAssignToSlice(&result)

		// Should have all left items with nil right
		if len(result) != 2 {
			t.Errorf("Expected 2 items but got %v", len(result))
		}

		for _, pair := range result {
			if pair.Right != nil {
				t.Errorf("Expected nil right but got %v", pair.Right)
			}
		}
	})
}
