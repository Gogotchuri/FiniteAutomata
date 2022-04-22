package simulation

import (
	"fmt"
	"testing"
)

func TestCreateAutomataForLiteral(t *testing.T) {
	tests := []struct {
		input rune
		want  string
	}{
		{'a', "2 1 1\n1\n1 a 1\n0\n"},
		{'b', "2 1 1\n1\n1 b 1\n0\n"},
	}
	for _, test := range tests {
		got := createAutomataForLiteral(test.input)
		fmt.Println(got)
		if got.String() != test.want {
			t.Errorf("createAutomataForLiteral(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestBuildAutomataFromRegexForLiterals(t *testing.T) {
	tests := []struct {
		input rune
		want  string
	}{
		{'a', "2 1 1\n1\n1 a 1\n0\n"},
		{'b', "2 1 1\n1\n1 b 1\n0\n"},
		{'c', "2 1 1\n1\n1 c 1\n0\n"},
	}
	for _, test := range tests {
		got := BuildAutomataFromRegex(string(test.input))
		fmt.Println(got)
		if got.String() != test.want {
			t.Errorf("createAutomataForLiteral(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestBuildAutomataFromRegexConcat(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"a", "a"},
		{"ab1", "ab1"},
		{"a(ba)()", "aba"},
		{"ba", "ba"},
		{"baaa", "baaa"}}
	for _, test := range tests {
		got := BuildAutomataFromRegex(test.input)
		fmt.Println(test.input)
		fmt.Println(got)
		if b, _ := got.Simulate(test.want); !b {
			t.Errorf("BuildAutomataFromRegex(%q) = %q, want %q", test, got, test)
		}
	}
	//Test not matching
	for _, test := range tests {
		got := BuildAutomataFromRegex(test.input)
		if b, _ := got.Simulate("aaa"); b {
			t.Errorf("BuildAutomataFromRegex(%q) = %q, want %q", test, got, test)
		}
	}
}

func TestBuildAutomataFromRegexUnion(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"a", "a"},
		{"a|b|1", "a"},
		{"a|b|1", "b"},
		{"a|b|1", "1"},
		{"a|b", "b"},
		{"a|(ba)|()", "a"},
		{"a|(ba)|()", "ba"},
		{"a|(ba)|()", ""},
		{"((ba)|(ab|b))|a", "a"},
		{"((ba)|(ab|b))|a", "ba"},
		{"((ba)|(ab|b))|a", "ab"},
		{"((ba)|(ab|b))|a", "b"},
	}
	for _, test := range tests {
		got := BuildAutomataFromRegex(test.input)
		fmt.Println(test, "\n", got)
		if b, _ := got.Simulate(test.want); !b {
			t.Errorf("BuildAutomataFromRegex(%q) = %q, want %q", test, got, test)
		}
	}
	//Test not matching
	for _, test := range tests {
		got := BuildAutomataFromRegex(test.input)
		if b, _ := got.Simulate("aaa"); b {
			t.Errorf("BuildAutomataFromRegex(%q) = %q, want %q", test, got, test)
		}
	}
}

func TestBuildAutomataFromRegexStar(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"a*", ""},
		{"a*", "a"},
		{"a*", "aa"},
		{"a*", "aaa"},
		{"ba*", "b"},
		{"ba*", "ba"},
		{"ba*", "baaa"},
		{"a*b", "b"},
		{"a*b", "ab"},
		{"a*b", "aab"},
		{"a*|b", ""},
		{"a*|b", "b"},
		{"a*|b", "a"},
		{"a*|b", "aa"},
		{"a|b*", ""},
		{"a|b*", "a"},
		{"a|b*", "b"},
		{"a|b*", "bb"},
		{"(a|b)*", ""},
		{"(a|b)*", "a"},
		{"(a|b)*", "aa"},
		{"(a|b)*", "b"},
		{"(a|b)*", "bb"},
		{"(a|b)*", "abab"},
		{"(a|b)*", "baba"},
	}
	for _, test := range tests {
		got := BuildAutomataFromRegex(test.input)
		fmt.Println(test, "\n", got)
		if b, _ := got.Simulate(test.want); !b {
			t.Errorf("BuildAutomataFromRegex(%q) = %q, want %q", test, got, test)
		}
	}
}
