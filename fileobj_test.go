package fileobj

import "os"
import "testing"

func TestIt(t *testing.T) {
	os.Remove("test")
	obj := map[string]string{
		"foo": "FOO",
		"bar": "BAR",
	}
	f, err := New("test", &obj)
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	err = f.Save()
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	obj["baz"] = "BAZ"
	err = f.Save()
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	obj = map[string]string{}
	f, err = New("test", &obj)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if obj["foo"] != "FOO" || obj["bar"] != "BAR" || obj["baz"] != "BAZ" {
		t.Fatal("data error")
	}
}
