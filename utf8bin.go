package utf8bin

import (
	"errors"
	"unicode/utf8"
)

func BinToUTF8(data []byte) string {
	size := len(data)
	runes := make([]rune, 0)

	for i := 0; i < size; i++ {
		r, l := utf8.DecodeRune(data[i:])
		if l > 1 {
			if (r > 0x7F && r <= 0xFF) || (r > 0xD800 && r <= 0xDFFF) {
				for ; l > 0; l-- {
					runes = append(runes, rune(data[i]))
					i++
				}
			} else {
				runes = append(runes, r)
			}
			i += l - 1
		} else {
			runes = append(runes, rune(data[i]))
		}
	}
	return string(runes)
}

func UTF8ToBin(str string) ([]byte, error) {
	data_in := []byte(str)
	size := len(data_in)
	data_out := make([]byte, 0)

	for i := 0; i < size; i++ {
		if data_in[i] >= 0xC2 && data_in[i] <= 0xC3 {
			r, l := utf8.DecodeRune(data_in[i:])
			if l > 0 && utf8.RuneError != r {
				data_out = append(data_out, byte(r))
				i++
			} else {
				return nil, errors.New("not UTF-8 string")
			}
		} else {
			data_out = append(data_out, data_in[i])
		}
	}
	return data_out, nil
}
