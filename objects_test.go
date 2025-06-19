package flowdata

import "testing"

func TestMergeObjects(t *testing.T) {
	a := Object{Type: "test", Name: "test"}
	b := Object{Type: "test", Name: "test-2", Data: map[string]interface{}{"brightness": 123}}

	a.Merge(&b)

	if a.Name != b.Name {
		t.Error("failed to merge name")
	}

	v, ok := a.Data["brightness"]
	if !ok || v != 123 {
		t.Error("failed to merge data")
	}
}

func TestMergeObjectsWithChilds(t *testing.T) {
	a := Object{Type: "test", Name: "test",
		Objects: []*Object{{Type: "sub-object", Name: "obj1"}, {Type: "sub-object", Name: "obj2"}},
	}

	b := Object{Type: "test", Name: "test",
		Objects: []*Object{{Type: "sub-object", Name: "obj1"}, {Type: "sub-object", Name: "obj3"}},
		Items:   []*Object{{Type: "sub-item", Name: "item1"}},
	}

	a.Merge(&b)

	t.Logf("%+v", a.Objects)

	if len(a.Objects) != 4 {
		t.Errorf("should be 4 child items, has %d", len(a.Objects))
	}
}
