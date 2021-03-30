package emoji

import (
	"fmt"
	"testing"
)

func TestEmojiExists(t *testing.T) {
	x := reverseEmojiMap["❤️"]
	if ":heart:" != x {
		t.Fatal("not found emoji")
	}

	msg := "country flag alias 🇺🇸"
	m := Deparse(msg)
	fmt.Printf("m: %s\n", m)
	msg = "country flag alias 🇬🇧"
	m = Deparse(msg)
	fmt.Printf("m: %s\n", m)
	msg = Parse(":flag-gb:")
	fmt.Printf("msg: %s\n", msg)

}
