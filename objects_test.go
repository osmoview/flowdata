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
