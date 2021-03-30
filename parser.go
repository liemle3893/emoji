package emoji

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	flagRegex = regexp.MustCompile(`^:flag-([a-zA-Z]{2}):$`)
)

// Parse replaces emoji aliases (:pizza:) with unicode representation.
func Parse(input string) string {
	var matched strings.Builder
	var output strings.Builder

	for _, r := range input {
		// when it's not `:`, it might be inner or outer of the emoji alias
		if r != ':' {
			// if matched is empty, it's the outer of the emoji alias
			if matched.Len() == 0 {
				output.WriteRune(r)
				continue
			}

			matched.WriteRune(r)

			// if it's space, the alias's not valid.
			// reset matched for breaking the emoji alias
			if unicode.IsSpace(r) {
				output.WriteString(matched.String())
				matched.Reset()
			}
			continue
		}

		// r is `:` now
		// if matched is empty, it's the beginning of the emoji alias
		if matched.Len() == 0 {
			matched.WriteRune(r)
			continue
		}

		// it's the end of the emoji alias
		match := matched.String()
		alias := match + ":"

		// check for emoji alias
		if code, ok := Find(alias); ok {
			output.WriteString(code)
			matched.Reset()
			continue
		}

		// not found any emoji
		output.WriteString(match)
		// it might be the beginning of the another emoji alias
		matched.Reset()
		matched.WriteRune(r)

	}

	// if matched not empty, add it to output
	if matched.Len() != 0 {
		output.WriteString(matched.String())
		matched.Reset()
	}

	return output.String()
}

// Deparse return string which replacing all emoji with :alias:
func Deparse(msg string) string {
	var cRunes []rune
	var output strings.Builder

	for len(msg) > 0 {
		r, size := utf8.DecodeRuneInString(msg)
		cRunes = append(cRunes, r)
		c := fmt.Sprintf("%s", string(cRunes))
		if alias, ok := reverseEmojiMap[c]; ok {
			// Found alias
			normalizedStr := normalizedString(msg)
			lge, s := longestEmoji(normalizedStr)
			if lge != "" {
				output.WriteString(lge)
				size = s
			} else {
				output.WriteString(alias)
			}
			// Reset current rune
			cRunes = nil
		}
		if s := RunesToHexKey([]rune{r}); len(s) >= 4 {
			msg = msg[size:]
			continue
		}
		// Flush cRunes if any
		if len(cRunes) > 0 {
			output.WriteString(string(cRunes))
			cRunes = nil
		}
		msg = msg[size:]
	}
	return output.String()
}

// Map returns the emojis map.
// Key is the alias of the emoji.
// Value is the code of the emoji.
func Map() map[string]string {
	return emojiMap
}

// ReversedMap returns the reversed emoji map of aliases
// Key is the code of the emoji
// Value is the alias
func ReversedMap() map[string]string {
	return reverseEmojiMap
}

// AppendAlias adds new emoji pair to the emojis map.
func AppendAlias(alias, code string) error {
	if c, ok := emojiMap[alias]; ok {
		return fmt.Errorf("emoji already exist: %q => %+q", alias, c)
	}

	for _, r := range alias {
		if unicode.IsSpace(r) {
			return fmt.Errorf("emoji alias is not valid: %q", alias)
		}
	}

	emojiMap[alias] = code

	return nil
}

// Exist checks existence of the emoji by alias.
func Exist(alias string) bool {
	_, ok := Find(alias)

	return ok
}

// Find returns the emoji code by alias.
func Find(alias string) (string, bool) {
	if code, ok := emojiMap[alias]; ok {
		return code, true
	}

	if flag := checkFlag(alias); len(flag) > 0 {
		return flag, true
	}

	return "", false
}

func FindReverse(unicode string) (string, bool) {
	if alias, ok := reverseEmojiMap[unicode]; ok {
		return alias, true
	}
	return "", false
}

// checkFlag finds flag emoji for `flag-[CODE]` pattern
func checkFlag(alias string) string {
	if matches := flagRegex.FindStringSubmatch(alias); len(matches) == 2 {
		flag, _ := CountryFlag(matches[1])

		return flag.String()
	}

	return ""
}

// RunesToHexKey - Convert a slice of runes to hex string representation of their Unicode Code Point value
func RunesToHexKey(runes []rune) (output string) {
	// Build a slice of hex representations of each rune
	hexParts := []string{}
	for _, rune := range runes {
		hexParts = append(hexParts, fmt.Sprintf("%X", rune))
	}

	// Join the hex strings with a hypen - this is the key used in the emojis map
	output = strings.Join(hexParts, "-")
	return
}

func normalizedString(input string) string {
	runes := []rune(input)
	for idx, r := range runes {
		if hk := RunesToHexKey([]rune{r}); len(hk) < 4 {
			return string(runes[:idx+1])
		}
	}
	return input
}

func longestEmoji(normalizedStr string) (string, int) {
	runes := []rune(normalizedStr)
	size := 0
	for len(runes) > 0 {
		emoji := fmt.Sprintf("%s", string(runes))
		if alias, ok := reverseEmojiMap[emoji]; ok {
			for _, r := range runes {
				size += utf8.RuneLen(r)
			}
			return alias, size
		}
		runes = runes[:len(runes)-1]
	}
	return "", 0
}
