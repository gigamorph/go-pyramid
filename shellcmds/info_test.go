package shellcmds

import (
	"fmt"
	"testing"
)

func TestWidth(t *testing.T) {
	fmt.Println("TestWidth")
	fpath := "../test/resources/images/ag-obj-286-0033-pub.jpg"
	want := uint(3840)
	if w, err := Width(fpath); err != nil {
		t.Error(err)
	} else if w != want {
		t.Error("width: want", want, "got ", w)
	}
}

func TestHeight(t *testing.T) {
	fpath := "../test/resources/images/ag-obj-286-0033-pub.jpg"
	want := uint(2556)
	if h, err := Height(fpath); err != nil {
		t.Error(err)
	} else if h != want {
		t.Error("height: want", want, "got ", h)
	}
}
