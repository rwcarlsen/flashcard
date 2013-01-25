
package colwidth

import (
	"unicode"
)

type interval struct {
	First rune
	Last rune
}

/* The following two functions define the column width of an ISO 10646
 * character as follows:
 *
 *    - The null character (U+0000) has a column width of 0.
 *
 *    - Other C0/C1 control characters and DEL will lead to a return
 *      value of -1.
 *
 *    - Non-spacing and enclosing combining characters (general
 *      category code Mn or Me in the Unicode database) have a
 *      column width of 0.
 *
 *    - SOFT HYPHEN (U+00AD) has a column width of 1.
 *
 *    - Other format characters (general category code Cf in the Unicode
 *      database) and ZERO WIDTH SPACE (U+200B) have a column width of 0.
 *
 *    - Hangul Jamo medial vowels and final consonants (U+1160-U+11FF)
 *      have a column width of 0.
 *
 *    - Spacing characters in the East Asian Wide (W) or East Asian
 *      Full-width (F) category as defined in Unicode Technical
 *      Report #11 have a column width of 2.
 *
 *    - All remaining characters (including all printable
 *      ISO 8859-1 and WGL4 characters, Unicode control characters,
 *      etc.) have a column width of 1.
 *
 * This implementation assumes that wchar_t characters are encoded
 * in ISO 10646.
 */

var nonSpacing = []*unicode.RangeTable{
	unicode.Mn,
	unicode.Me,
	unicode.Cf,
}

var fullAndWideSpacing = []*unicode.RangeTable{
	unicode.Hangul,
}

func Char(r rune) int {
	if r == 0 {
		return 0 // null character
	} else if unicode.Is(unicode.C, r) { // is r control or special character
		return -1
	/* binary search in table of non-spacing characters */
	} else if unicode.IsOneOf(nonSpacing, r) {
		return 0
	} else if unicode.IsOneOf(fullAndWideSpacing, r) {
		return 2
	}
	return 1
}

func String(s string) int {
	width := 0
	for _, char := range s {
		if w := Char(char); w < 0 {
			return -1;
		} else {
			width += w;
		}
	}
	return width;
}

