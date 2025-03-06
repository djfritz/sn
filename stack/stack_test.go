package stack

import (
	"os"
	"testing"
)

func TestStack1(t *testing.T) {
	f, err := os.Open("test1.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = Decode(f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStack2(t *testing.T) {
	f, err := os.Open("test2.txt")
	if err != nil {
		t.Fatal(err)
	}

	_, err = Decode(f)
	if err != nil {
		t.Fatal(err)
	}
}
