package collection

import (
	"fmt"
	"testing"
)

func TestInvalidType(t *testing.T) {
	_, err := New("hello")
	if err == nil || err != ErrorInvalidDataType {
		t.Errorf("test failed")
	}

	// collection := MustNew([]string{"hello", "world", "", "you", "are"})
	// fmt.Println(collection.Filter2(func(item string) bool {
	// 	return item != ""
	// }).ToString())

	collection := MustNew([]string{"hello", "world", "", "you", "are"}).Filter(func(item string) bool {
		return item != ""
	})

	if collection.Count() != 4 {
		t.Error("test failed")
	}

	if fmt.Sprint(collection.All()) != "[hello world you are]" {
		t.Errorf("test failed")
	}

	collection.Append("abc", "def")
	if collection.ToString() != "[hello world you are abc def]" {
		t.Errorf("test failed")
	}

	collection.Each(func(item string, index int) {
		if item == "" || index < 0 {
			t.Errorf("test failed: %d:%s", index, item)
		}
	})

	collection.Each(func(item string) {
		if item == "" {
			t.Errorf("test failed: %s", item)
		}
	})

	if collection.IsEmpty() {
		t.Error("test failed")
	}

	if !collection.Filter(func(item string) bool {
		return false
	}).IsEmpty() {
		t.Error("test failed")
	}
}

func TestStringCollection(t *testing.T) {
	collection := MustNew([]interface{}{"hello", "world", "", "you", "are"})
	collection = collection.Filter(func(item string) bool {
		return item != ""
	}).Map(func(item string) string {
		return "<" + item + ">"
	})

	if collection.ToString() != "[<hello> <world> <you> <are>]" {
		t.Errorf("test failed: ^%s$", collection.ToString())
	}

	res := collection.Reduce(func(carry string, item string) string {
		return fmt.Sprintf("%s->%s", carry, item)
	}, "")

	if res != "-><hello>-><world>-><you>-><are>" {
		t.Errorf("test failed: ^%s$", res)
	}
}

func TestStructCollection(t *testing.T) {
	type Element struct {
		ID   int
		Name string
	}

	type Element2 struct {
		ID   int
		Name string
		Age  int
	}

	elements := []Element{
		Element{ID: 1, Name: "hello"},
		Element{ID: 2, Name: "world"},
		Element{ID: 3, Name: ""},
		Element{ID: 4, Name: "Tom"},
	}

	collection := MustNew(elements)
	collection = collection.Filter(func(item Element) bool {
		return item.Name != ""
	}).Map(func(item Element) Element2 {
		return Element2{
			ID:   item.ID,
			Name: item.Name,
			Age:  item.ID * 2,
		}
	})

	if collection.ToString() != "[{1 hello 2} {2 world 4} {4 Tom 8}]" {
		t.Errorf("test failed: ^%s$", collection.ToString())
	}

	res := collection.Reduce(func(carry string, item Element2) string {
		return fmt.Sprintf("%v\n%v", carry, item)
	}, "{0 Start}")

	expectValue := `{0 Start}
{1 hello 2}
{2 world 4}
{4 Tom 8}`

	if res != expectValue {
		t.Errorf("test failed: ^%s$", res)
	}
}
