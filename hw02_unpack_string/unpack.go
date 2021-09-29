package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const backslash = rune('\u005c')

var ErrInvalidString = errors.New("invalid string")

func resetFlag(flag *bool) {
	*flag = false
}

func setFlag(flag *bool) {
	*flag = true
}

func deleteLastRune(b *strings.Builder) {
	tmpString := b.String()
	tmpRuneSlice := []rune(tmpString)
	tmpRuneSlice = tmpRuneSlice[:len(tmpRuneSlice)-1]
	b.Reset()
	b.WriteString(string(tmpRuneSlice))
}

func Unpack(packedStr string) (unpacked string, err error) {
	var (
		buffered         rune
		backslashedFlag  bool
		multiplieredFlag bool
		b                strings.Builder
	)
	for index, currentRune := range packedStr {
		switch {
		case index == 0 && unicode.IsDigit(currentRune):
			b.Reset()
			err = ErrInvalidString
			return "", err
		case currentRune == backslash && !backslashedFlag:
			setFlag(&backslashedFlag)
			resetFlag(&multiplieredFlag)
		case unicode.IsDigit(currentRune) && !backslashedFlag:
			if multiplieredFlag {
				b.Reset()
				err = ErrInvalidString
				return "", err
			}
			multiplier, err := strconv.Atoi(string(currentRune))
			if err != nil {
				return "", err
			}
			switch multiplier {
			case 0:
				deleteLastRune(&b)
			case 1:
				// do not need multiply
			default:
				b.WriteString(strings.Repeat(string(buffered), multiplier-1))
			}
			buffered = currentRune
			resetFlag(&backslashedFlag)
			setFlag(&multiplieredFlag)
		default:
			b.WriteRune(currentRune)
			buffered = currentRune
			resetFlag(&backslashedFlag)
			resetFlag(&multiplieredFlag)
		}
	}
	unpacked = b.String()
	return unpacked, nil
}
