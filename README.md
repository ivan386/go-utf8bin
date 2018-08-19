# go-utf8bin
convert binary data to valid utf8 string

```
package main

import (
	"bytes"
	"errors"
	"fmt"
	"unicode/utf8"
	"encoding/json"
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

func main() {
	// This binary data that contain invalid UTF-8 byte sequences at the end
	binary_data := []byte("\u0080\u0081\u0082\u0083 ðñòÿ test тест 12345 \x01\x02\x03\x04\x05\x06 \x7F\x80\x81\x82\x83 \xF0\xF1\xF2\xFF")
	fmt.Printf("binary_data: %q\n", binary_data)
	
	// If you try to convert it to json
	binary_to_json, _ := json.Marshal(string(binary_data))
	
	// invalid utf-8 byte sequences will be replaced by \ufffd rune
	fmt.Printf("binary_to_json: %q\n", binary_to_json)
	
	// If you convert binary data to valid UTF-8 string
	utf8_string := BinToUTF8(binary_data)
	
	// all invalid UTF-8 byte sequences vill be converted to valid
	fmt.Printf("utf8_string: %q\n", utf8_string)
	
	// convert it to json
	utf8_in_json, _ := json.Marshal(utf8_string)
	
	// all data is save
	fmt.Printf("utf8_in_json: %q\n", utf8_in_json)
	
	// on other end unmarshal it
	utf8_from_json := ""
	err := json.Unmarshal(utf8_in_json, &utf8_from_json)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("(utf8_string == utf8_from_json) is %t\n", utf8_string == utf8_from_json)
	fmt.Printf("utf8_from_json: %q\n", utf8_from_json)
	
	// and decode it to binary data
	decoded_binary, _ := UTF8ToBin(utf8_from_json)
	
	// all data is saved
	fmt.Printf("(binary_data == decoded_binary) is %t\n", bytes.Equal(binary_data, decoded_binary))
	fmt.Printf("decoded_binary: %q\n", decoded_binary)
}
```

Result:
```
binary_data: "\u0080\u0081\u0082\u0083 ðñòÿ test тест 12345 \x01\x02\x03\x04\x05\x06 \u007f\x80\x81\x82\x83 \xf0\xf1\xf2\xff"
binary_to_json: "\"\u0080\u0081\u0082\u0083 ðñòÿ test тест 12345 \\u0001\\u0002\\u0003\\u0004\\u0005\\u0006 \u007f\\ufffd\\ufffd\\ufffd\\ufffd \\ufffd\\ufffd\\ufffd\\ufffd\""
utf8_string: "Â\u0080Â\u0081Â\u0082Â\u0083 Ã°Ã±Ã²Ã¿ test тест 12345 \x01\x02\x03\x04\x05\x06 \u007f\u0080\u0081\u0082\u0083 ðñòÿ"
utf8_in_json: "\"Â\u0080Â\u0081Â\u0082Â\u0083 Ã°Ã±Ã²Ã¿ test тест 12345 \\u0001\\u0002\\u0003\\u0004\\u0005\\u0006 \u007f\u0080\u0081\u0082\u0083 ðñòÿ\""
(utf8_string == utf8_from_json) is true
utf8_from_json: "Â\u0080Â\u0081Â\u0082Â\u0083 Ã°Ã±Ã²Ã¿ test тест 12345 \x01\x02\x03\x04\x05\x06 \u007f\u0080\u0081\u0082\u0083 ðñòÿ"
(binary_data == decoded_binary) is true
decoded_binary: "\u0080\u0081\u0082\u0083 ðñòÿ test тест 12345 \x01\x02\x03\x04\x05\x06 \u007f\x80\x81\x82\x83 \xf0\xf1\xf2\xff"

Program exited.
```

https://play.golang.org/p/j25UvMAxeab
