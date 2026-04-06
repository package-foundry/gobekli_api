package apikey

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var base58Decodemap = make(map[byte]int)

func init() {
	for i, b := range []byte(base58Alphabet) {
		base58Decodemap[b] = i
	}
}

func base58Encode(src []byte) string {
	if len(src) == 0 {
		return ""
	}

	count := 0
	for _, b := range src {
		if b == 0 {
			count++
		} else {
			break
		}
	}

	result := make([]byte, 0, len(src)*2)

	num := newInteger(src)
	for num.GreaterThan(zero) {
		div, rem := num.DivMod(58)
		result = append(result, base58Alphabet[rem])
		num = div
	}

	for range src[:count] {
		result = append(result, base58Alphabet[0])
	}

	reverse(result)
	return string(result)
}

type integer struct {
	data []byte
}

func newInteger(src []byte) *integer {
	dst := make([]byte, len(src))
	copy(dst, src)
	return &integer{data: dst}
}

func (i *integer) GreaterThan(o *integer) bool {
	if len(i.data) != len(o.data) {
		return len(i.data) > len(o.data)
	}
	for j := 0; j < len(i.data); j++ {
		if i.data[j] != o.data[j] {
			return i.data[j] > o.data[j]
		}
	}
	return false
}

func (i *integer) DivMod(divisor int) (*integer, int) {
	quotient := make([]byte, 0, len(i.data))
	remainder := 0

	for _, b := range i.data {
		current := remainder*256 + int(b)
		quot := current / divisor
		remainder = current % divisor

		if len(quotient) > 0 || quot > 0 {
			quotient = append(quotient, byte(quot))
		}
	}

	if len(quotient) == 0 {
		quotient = []byte{0}
	}

	return &integer{data: quotient}, remainder
}

var zero = &integer{data: []byte{0}}

func reverse(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
