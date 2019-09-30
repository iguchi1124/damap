package damap

import (
	"reflect"
	"testing"
)

func TestExactMatchSearch(t *testing.T) {
	keys := []string{"Lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit", "sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore", "magna", "aliqua", "Ut", "enim", "ad", "minim", "veniam", "quis", "nostrud", "exercitation", "ullamco", "laboris", "nisi", "ut", "aliquip", "ex", "ea", "commodo", "consequat", "Duis", "aute", "irure", "dolor", "in", "reprehenderit", "in", "voluptate", "velit", "esse", "cillum", "dolore", "eu", "fugiat", "nulla", "pariatur", "Excepteur", "sint", "occaecat", "cupidatat", "non", "proident", "sunt", "in", "culpa", "qui", "officia", "deserunt", "mollit", "anim", "id", "est", "laborum"}
	tests := []struct {
		in  string
		out bool
	}{
		{"Lorem", true},
		{"ipsum", true},
		{"dolor", true},
		{"sit", true},
		{"amet", true},
		{"consectetur", true},
		{"adipiscing", true},
		{"elit", true},
		{"sed", true},
		{"do", true},
		{"eiusmod", true},
		{"tempor", true},
		{"incididunt", true},
		{"ut", true},
		{"labore", true},
		{"et", true},
		{"dolore", true},
		{"magna", true},
		{"aliqua", true},
		{"Lore", false},
		{"lorem", false},
		{"ipsu", false},
		{"olor", false},
		{"i", false},
	}

	d := New()
	for _, key := range keys {
		d.Write(key, nil)
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			result := d.ExactMatchSearch(test.in)

			if result != test.out {
				t.Errorf("got %t, want %t", result, test.out)
			}
		})
	}
}

func TestCommonPrefixSearch(t *testing.T) {
	d := New()
	d.Write("pine", "foo")
	d.Write("apple", "bar")
	d.Write("pineapple", "foobar")

	test := struct {
		in  string
		out CommonPrefixSearchResult
	}{
		"I have a pineapple.",
		[]struct {
			Pos   int
			Key   string
			Value interface{}
		}{
			{9, "pine", "foo"},
			{9, "pineapple", "foobar"},
			{13, "apple", "bar"},
		},
	}

	t.Run(test.in, func(t *testing.T) {
		result := d.CommonPrefixSearch(test.in)

		if !reflect.DeepEqual(result, test.out) {
			t.Errorf("got %v, want %v", result, test.out)
		}
	})
}

func TestMultiByteTextCommonPrefixSearch(t *testing.T) {
	d := New()
	d.Write("こんにちは", "hello")
	d.Write("さようなら", "bye")

	test := struct {
		in  string
		out CommonPrefixSearchResult
	}{
		"おはようこんにちはさようなら",
		[]struct {
			Pos   int
			Key   string
			Value interface{}
		}{
			{4, "こんにちは", "hello"},
			{9, "さようなら", "bye"},
		},
	}

	t.Run(test.in, func(t *testing.T) {
		result := d.CommonPrefixSearch(test.in)

		if !reflect.DeepEqual(result, test.out) {
			t.Errorf("got %v, want %v", result, test.out)
		}
	})
}
