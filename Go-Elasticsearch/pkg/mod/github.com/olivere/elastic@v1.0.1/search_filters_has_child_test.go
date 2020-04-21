package elastic

import (
	"encoding/json"
	"testing"
)

func TestHasChildFilter(t *testing.T) {
	f := NewHasChildFilter("blog_tag")
	f = f.Query(NewTermQuery("tag", "something"))
	data, err := json.Marshal(f.Source())
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"has_child":{"query":{"term":{"tag":"something"}},"type":"blog_tag"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
