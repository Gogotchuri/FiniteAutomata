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
		{'a', "2 1 1\n1\n0 a 1\n"},
		{'b', "2 1 1\n1\n0 b 1\n"},
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
		{'a', "2 1 1\n1\n0 a 1\n"},
		{'b', "2 1 1\n1\n0 b 1\n"},
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

func TestBuild11(t *testing.T) {
	in := `(hc|((n|ecv|nycu)(((7b)*|(4x5)*|d)*)(pjt))*|m|((qjk|xw0|9)*)(235w)(ll|u68s|(oxk)*|(2)*)(1xh5|j|(o7j)*|r))*`
	out := `51 13 209
			0 2 19 20 31 33 37 40 41 45 46 49 50
			9 h 1 n 3 n 7 e 4 m 20 q 21 x 24 9 27 2 28 
			1 c 2 
			9 h 1 n 3 n 7 e 4 m 20 q 21 x 24 9 27 2 28 
			4 7 11 4 13 d 16 p 17 
			1 c 5 
			1 v 6 
			4 7 11 4 13 d 16 p 17 
			1 y 8 
			1 c 9 
			1 u 10 
			4 7 11 4 13 d 16 p 17 
			1 b 12 
			4 7 11 4 13 d 16 p 17 
			1 x 14 
			1 5 15 
			4 4 13 7 11 d 16 p 17 
			4 7 11 4 13 d 16 p 17 
			1 j 18 
			1 t 19 
			9 n 3 n 7 e 4 h 1 m 20 q 21 x 24 9 27 2 28 
			9 h 1 n 3 n 7 e 4 m 20 q 21 x 24 9 27 2 28 
			1 j 22 
			1 k 23 
			4 q 21 x 24 9 27 2 28 
			1 w 25 
			1 0 26 
			4 q 21 x 24 9 27 2 28 
			4 q 21 x 24 9 27 2 28 
			1 3 29 
			1 5 30 
			1 w 31 
			17 l 32 u 34 o 38 o 47 2 41 2 28 1 42 j 46 r 50 h 1 n 3 n 7 e 4 m 20 q 21 x 24 9 27 
			1 l 33 
			13 1 42 j 46 o 47 r 50 h 1 n 3 n 7 e 4 m 20 q 21 x 24 9 27 2 28 
			1 6 35 
			1 8 36 
			1 s 37 
			13 1 42 j 46 o 47 r 50 h 1 n 3 n 7 e 4 m 20 q 21 x 24 9 27 2 28 
			1 x 39 
			1 k 40 
			14 o 38 o 47 1 42 j 46 r 50 h 1 n 3 n 7 e 4 m 20 q 21 x 24 9 27 2 28 
			14 2 41 2 28 1 42 j 46 o 47 r 50 h 1 n 3 n 7 e 4 m 20 q 21 x 24 9 27 
			1 x 43 
			1 h 44 
			1 5 45 
			9 h 1 n 3 n 7 e 4 m 20 q 21 x 24 9 27 2 28 
			9 h 1 n 3 n 7 e 4 m 20 q 21 x 24 9 27 2 28 
			1 7 48 
			1 j 49 
			10 o 47 h 1 n 3 n 7 e 4 m 20 q 21 x 24 9 27 2 28 
			9 h 1 n 3 n 7 e 4 m 20 q 21 x 24 9 27 2 28`

	got := BuildAutomataFromRegex(in)
	fmt.Println(got)
	if got.String() != out {
		t.Errorf("Expected: %s\n Got: %s\n", out, got.String())
	}
}
