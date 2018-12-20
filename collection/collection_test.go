package collection_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/mylxsw/go-toolkit/collection"
)

var testMapData = map[string]string{
	"k1": "v1",
	"k2": "v2",
	"k3": "v3",
	"":   "v4",
	"k5": "",
	"xx": "yy",
}

type Element struct {
	ID   int
	Name string
}

type Element2 struct {
	ID   int
	Name string
	Age  int
}

func TestInvalidTypeForMap(t *testing.T) {
	cc := collection.MustNew(testMapData)

	collectionWithoutEmpty := cc.Filter(func(value string) bool {
		return value != ""
	}).Filter(func(value string, key string) bool {
		return key != ""
	})

	collectionWithoutEmpty.Each(func(value, key string) {
		if value == "" || key == "" {
			t.Errorf("test failed: %s=>%s", key, value)
		}
	})

	if collectionWithoutEmpty.Size() != 4 || cc.Size() != 6 {
		t.Error("test failed")
	}

	if _, ok := collectionWithoutEmpty.All().(map[interface{}]interface{}); !ok {
		t.Error("test failed")
	}

	if _, ok := cc.All().([]string); ok {
		t.Error("test failed")
	}

	if cc.IsEmpty() {
		t.Error("test failed")
	}

	if !cc.Filter(func(value string) bool {
		return false
	}).IsEmpty() {
		t.Error("test failed")
	}
}

func TestInvalidTypeForArray(t *testing.T) {
	_, err := collection.New("hello")
	if err == nil || err != collection.ErrorInvalidDataType {
		t.Errorf("test failed")
	}

	// collection := MustNew([]string{"hello", "world", "", "you", "are"})
	// fmt.Println(collection.Filter2(func(item string) bool {
	// 	return item != ""
	// }).ToString())

	collection := collection.MustNew([]string{"hello", "world", "", "you", "are"}).Filter(func(item string) bool {
		return item != ""
	})

	if collection.Size() != 4 {
		t.Error("test failed")
	}

	if fmt.Sprint(collection.All()) != "[hello world you are]" {
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

func TestStringMapCollection(t *testing.T) {
	collection := collection.MustNew(testMapData)
	collection = collection.Filter(func(_, key string) bool {
		return key != ""
	}).Filter(func(value string) bool {
		return value != ""
	}).Map(func(value, key string) string {
		return fmt.Sprintf("<%s(%s)>", value, key)
	})

	collection.Each(func(value, key string) {
		if !regexp.MustCompile(fmt.Sprintf(`^<\w+\(%s\)>$`, key)).MatchString(value) {
			t.Error("test failed")
		}
	})

	joinedValue := collection.Reduce(func(carry string, value, key string) string {
		if collection.MapIndex(key).(string) != value {
			t.Error("test failed")
		}
		return carry + " " + collection.MapIndex(key).(string)
	}, "value: ").(string)

	if len(joinedValue) <= len("value: ") {
		t.Error("test failed")
	}

	collection.Map(func(value string, key string) (string, string) {
		return value, key + "(modified)"
	}).Each(func(_, key string) {
		if !strings.HasSuffix(key, "(modified)") {
			t.Error("test failed")
		}
	})
}

func TestStringCollection(t *testing.T) {
	collection := collection.MustNew([]interface{}{"hello", "world", "", "you", "are"})
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

func TestComplexMapCollection(t *testing.T) {
	elements := map[string]Element{
		"one":   {ID: 1, Name: "hello"},
		"two":   {ID: 2, Name: "world"},
		"three": {ID: 3, Name: ""},
		"four":  {ID: 4, Name: "Tom"},
	}

	collection := collection.MustNew(elements)
	collection = collection.Filter(func(value Element, key string) bool {
		return value.Name != ""
	}).Map(func(value Element) Element2 {
		return Element2{
			ID:   value.ID,
			Name: value.Name,
			Age:  value.ID * 2,
		}
	})

	if collection.Size() != 3 {
		t.Errorf("test failed")
	}

	if !collection.MapHasIndex("one") {
		t.Error("test failed")
	}

	if collection.MapHasIndex("three") {
		t.Error("test failed")
	}

	collection.Each(func(value Element2, key string) {
		if key == "" {
			t.Error("test failed")
		}
	})
}

func TestComplexCollection(t *testing.T) {

	elements := []Element{
		{ID: 1, Name: "hello"},
		{ID: 2, Name: "world"},
		{ID: 3, Name: ""},
		{ID: 4, Name: "Tom"},
	}

	collection := collection.MustNew(elements)
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
