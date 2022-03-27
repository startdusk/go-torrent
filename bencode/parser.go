package bencode

import (
	"bufio"
	"io"
)

func Parse(r io.Reader) (*BObject, error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	//recursive descent parsing
	b, err := br.Peek(1)
	if err != nil {
		return nil, err
	}
	var ret BObject
	switch {
	case b[0] >= '0' && b[0] <= '9':
		// parse string
		val, err := DecodeString(br)
		if err != nil {
			return nil, err
		}
		ret.typ = BSTR
		ret.val = val
	case b[0] == 'i':
		// parse int
		val, err := DecodeInt(br)
		if err != nil {
			return nil, err
		}
		ret.typ = BINT
		ret.val = val
	case b[0] == 'l':
		// parse list
		br.ReadByte()
		var list []*BObject
		for {
			if p, _ := br.Peek(1); p[0] == 'e' {
				br.ReadByte()
				break
			}
			elem, err := Parse(br)
			if err != nil {
				return nil, err
			}
			list = append(list, elem)
		}
		ret.typ = BLIST
		ret.val = list
	case b[0] == 'd':
		// parse map
		br.ReadByte()
		dict := make(map[string]*BObject)
		for {
			if p, _ := br.Peek(1); p[0] == 'e' {
				br.ReadByte()
				break
			}
			key, err := DecodeString(br)
			if err != nil {
				return nil, err
			}
			val, err := Parse(br)
			if err != nil {
				return nil, err
			}
			dict[key] = val
		}
		ret.typ = BDICT
		ret.val = dict
	default:
		return nil, ErrIvd
	}
	return &ret, nil
}
