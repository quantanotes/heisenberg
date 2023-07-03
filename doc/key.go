package doc

type key []string

func (k key) Compress() string {
	bytes := make([]byte, 0)
	for _, s := range k {
		bytes = append(bytes, byte(len(s)))
		bytes = append(bytes, []byte(s)...)
	}
	return string(bytes)
}

func UncompressKey(k string) key {
	bytes := []byte(k)
	keys := make(key, 0)
	i := 0
	for i < len(bytes) {
		strLen := int(bytes[i])
		i++
		bs := bytes[i : i+strLen]
		s := string(bs)
		keys = append(keys, s)
		i += strLen
	}
	return keys
}
