package test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestContains(t *testing.T) {
	fmt.Println(strings.Contains("seafood", "foo"))
	fmt.Println(strings.Contains("seafood", "bar"))
	fmt.Println(strings.Contains("seafood", ""))
	fmt.Println(strings.Contains("", ""))
}

func TestJoin(t *testing.T) {
	fmt.Println(strings.Join([]string{"foo", "bar", "baz"}, "---"))
}

func TestIndex(t *testing.T) {
	fmt.Println(strings.Index("chicken", "ken"))
	fmt.Println(strings.Index("chicken", "dmr"))
}

func TestRepeat(t *testing.T) {
	fmt.Println("ba" + strings.Repeat("na", 2))
}

func TestReplace(t *testing.T) {
	fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))
	fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1))
}

func TestSplit(t *testing.T) {
	fmt.Printf("%q\n", strings.Split("a,b,c", ","))
	fmt.Printf("%q\n", strings.Split(",a,b,c,", ","))
	fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a "))
	fmt.Printf("%q\n", strings.Split(" xyz ", ""))
	fmt.Printf("%q\n", strings.Split("", "Bernardo O'Higgins"))
}

func TestTrim(t *testing.T) {
	fmt.Printf("[%q]", strings.Trim(" !!! Achtung !!! ", " !A"))
}

func TestFields(t *testing.T) {
	fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   "))
}

func TestAppend(t *testing.T) {
	str := make([]byte, 0, 100)
	str = strconv.AppendInt(str, 4567, 10)
	str = strconv.AppendBool(str, false)
	str = strconv.AppendQuote(str, "abcdefg")
	str = strconv.AppendQuoteRune(str, '单')
	fmt.Println(string(str))
}

func TestFormat(t *testing.T) {
	a := strconv.FormatBool(false)
	b := strconv.FormatFloat(123.23, 'g', 12, 64)
	c := strconv.FormatInt(1234, 16)
	d := strconv.FormatUint(12345, 16)
	e := strconv.Itoa(1023)
	fmt.Println(a, b, c, d, e)
}

func TestParse(t *testing.T) {
	a, err := strconv.ParseBool("false")
	checkError(err)
	b, err := strconv.ParseFloat("123.23", 64)
	checkError(err)
	c, err := strconv.ParseInt("1234", 10, 64)
	checkError(err)
	d, err := strconv.ParseUint("12345", 10, 64)
	checkError(err)
	e, err := strconv.Atoi("1023")
	checkError(err)
	fmt.Println(a, b, c, d, e)
}
