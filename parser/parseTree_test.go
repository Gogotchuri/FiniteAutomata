package parser

import (
	"fmt"
	"testing"
)

func TestInsertConcatenation(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"a", "a"},
		{"ab", "a.b"},
		{"ab()", "a.b._"},
		{"(ab)a", "(a.b).a"},
		{"(ab)a()", "(a.b).a._"},
		{"(ab())b(ba)", "(a.b._).b.(b.a)"},
		{"a(ab)a(ba)b", "a.(a.b).a.(b.a).b"},
		{"(a|b)a*(ba)|b", "(a|b).a*.(b.a)|b"},
	}
	for _, test := range tests {
		got := InsertConcatenation(test.input)
		if got != test.want {
			t.Errorf("InsertConcatenation(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestNewParseTreeWOParens(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"a", "a"},
		{"ab", `(a/.\b)`},
		{"a|b", `(a/|\b)`},
		{"a*b", `((a/*)/.\b)`},
		{"a*|b", `((a/*)/|\b)`},
		{"a*|b*", `((a/*)/|\(b/*))`},
		{"a*b*a*|a|bab*", `(((a/*)/.\((b/*)/.\(a/*)))/|\(a/|\(b/.\(a/.\(b/*)))))`},
	}
	for _, test := range tests {
		got := NewParseTree(test.input)
		fmt.Println(got)
		if got.String() != test.want {
			t.Errorf("NewParseTree(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestNewParseTreeWParens(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"(a)", "a"},
		{"((a)(b))", `(a/.\b)`},
		{"(a|b)", `(a/|\b)`},
		{"(a|(b))", `(a/|\b)`},
		{"(a*b)", `((a/*)/.\b)`},
		{"((a)*b)", `((a/*)/.\b)`},
		{"((a*)b)", `((a/*)/.\b)`},
		{"(a*(b))", `((a/*)/.\b)`},
		{"((a*)(b))", `((a/*)/.\b)`},
		{"(a)*|b", `((a/*)/|\b)`},
		{"(a)*|b", `((a/*)/|\b)`},
		{"((a)*)|b", `((a/*)/|\b)`},
		{"(((a|b)))|b", `((a/|\b)/|\b)`},
		{"((((a|b))*)|(b))", `(((a/|\b)/*)/|\b)`},
		{"((((a|b))*)|(b)*)", `(((a/|\b)/*)/|\(b/*))`},
		{"((((a|b))*)|(b*))", `(((a/|\b)/*)/|\(b/*))`},
		{"((((a|b))*)|(b*))*", `((((a/|\b)/*)/|\(b/*))/*)`},
	}
	for _, test := range tests {
		got := NewParseTree(test.input)
		fmt.Println(got)
		if got.String() != test.want {
			t.Errorf("NewParseTree(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestNewParseTreeWEmptyString(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"a", "a"},
		{"(a)", "a"},
		{"(a)()", "(a/.\\_)"},
		{"((a)()(b))", `(a/.\(_/.\b))`},
		{"(a)()*|b", `((a/.\(_/*))/|\b)`},
		{"((ba)|(ab|b))|a", `(((b/.\a)/|\((a/.\b)/|\b))/|\a)`},
	}
	for _, test := range tests {
		got := NewParseTree(test.input)
		fmt.Println(got)
		if got.String() != test.want {
			t.Errorf("NewParseTree(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}
