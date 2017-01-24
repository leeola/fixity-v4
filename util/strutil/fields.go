package strutil

import "unicode"

func QuotedFields(s string) []string {
	var lastQuote rune
	f := func(c rune) bool {
		switch {
		case c == rune(0):
			lastQuote = rune(0)
			return false
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)
		}
	}

	return FieldsFunc(s, f)
}

// FieldsFunc is similar to the stdlib FieldsFunc except that "flushes" the field func
//
// This flushing allows a field func to maintain state and reset it between uses,
// since FieldsFunc iterates over s twice.
func FieldsFunc(s string, isSplit func(rune) bool) []string {
	FlushFieldFunc(isSplit)

	fCount := 0
	inField := false
	for _, c := range s {
		wasInField := inField
		inField = !isSplit(c)
		if inField && !wasInField {
			fCount++
		}
	}

	if fCount == 0 {
		return nil
	}

	// flush the isSplit func, in case it uses state tracking between the two
	// loops.
	FlushFieldFunc(isSplit)

	fields := make([]string, fCount)
	fieldIndex := 0
	fieldStart := 0
	inField = false
	for i, c := range s {
		if isSplit(c) {
			if inField {
				// If split returned true and we *were* in a field, get the range from
				// fieldStart to the last index, as that is our field.
				fields[fieldIndex] = s[fieldStart:i]
				fieldIndex++
				inField = false
			}
		} else {
			inField = true
		}

		if !inField {
			// If split returned true and we were not in a field, then we're in
			// a series of splits (spaces/etc outside of a field), so we need to
			// move the fieldStart over until we are actually in a field.
			fieldStart = i + 1
		}
	}

	// If we looped through all the runes and didn't find an end to a given field,
	// set the last field to the remaining chars.
	if inField {
		fields[fieldIndex] = s[fieldStart:]
	}

	return fields
}

func FlushFieldFunc(f func(rune) bool) {
	f(rune(0))
}
