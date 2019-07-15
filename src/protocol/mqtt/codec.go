package mqtt

import "io"

func encodeBool(b bool) (i int) {
	if b {
		i = 1
	}
	return
}
func boolToByte(b bool) byte {
	switch b {
	case true:
		return 1
	default:
		return 0
	}
}

func decodeBool(i int) (b bool) {
	if i > 0 {
		b = true
	}
	return
}

func encodeString(str string) []byte {
	buf := []byte(str)
	size := len(buf)

	return append(encodeInt16(size), buf...)
}
func encodeBinary(b []byte) []byte {
	size := len(b)
	return append(encodeInt16(size), b...)
}
func encodeVariable(size uint64) []byte {
	ret := []byte{}
	for size > 0 {
		digit := size % 0x80
		size /= 0x80
		if size > 0 {
			digit |= 0x80
		}
		ret = append(ret, byte(digit))
	}
	return ret
}

func decodeString(buffer []byte, start int) (string, int) {
	size := (int(buffer[start]) << 8) | int(buffer[start+1])
	start += 2
	data := string(buffer[start:(start + size)])
	return data, start + size
}
func decodeBinary(buffer []byte, start int) ([]byte, int) {
	size := (int(buffer[start]) << 8) | int(buffer[start+1])
	start += 2
	data := buffer[start:(start + size)]
	return data, start + size
}
func decodeVariable(buffer []byte, start int) (uint64, int) {
	var (
		b    int
		size uint64
		mul  uint64 = 1
	)
	for {
		b, start = decodeInt(buffer, start)
		size += uint64(b&0x7F) * mul
		mul *= 0x80
		if b&0x80 == 0 {
			break
		}
	}
	return size, start
}

func decodeLength(r io.Reader) int {
	var rLength uint32
	var multiplier uint32
	b := make([]byte, 1)
	for multiplier < 27 { //fix: Infinite '(digit & 128) == 1' will cause the dead loop
		io.ReadFull(r, b)
		digit := b[0]
		rLength |= uint32(digit&127) << multiplier
		if (digit & 128) == 0 {
			break
		}
		multiplier += 7
	}
	return int(rLength)
}

func encodeInt(v int) byte {
	return byte(v)
}
func encodeInt16(v int) []byte {
	return append([]byte{}, byte(v>>8), byte(v&0xFF))
}
func encodeInt32(v int) []byte {
	return append([]byte{}, byte(v>>24), byte(v>>16), byte(v>>8), byte(v&0xFF))
}
func encodeUint(v uint8) byte {
	return byte(v)
}
func encodeUint16(v uint16) []byte {
	iv := int(v)
	return append([]byte{}, byte(iv>>8), byte(iv&0xFF))
}
func encodeUint32(v uint32) []byte {
	iv := int(v)
	return append([]byte{}, byte(iv>>24), byte(iv>>16), byte(iv>>8), byte(iv&0xFF))
}

func decodeInt(buffer []byte, start int) (int, int) {
	return int(buffer[start]), start + 1
}
func decodeInt16(buffer []byte, start int) (int, int) {
	return (int(buffer[start]) << 8) | int(buffer[start+1]), start + 2
}
func decodeInt32(buffer []byte, start int) (int, int) {
	return (int(buffer[start]) << 24) | (int(buffer[start+1]) << 16) | (int(buffer[start+2]) << 8) | int(buffer[start+3]), start + 4
}
func decodeUint(buffer []byte, start int) (uint8, int) {
	return uint8(buffer[start]), start + 1
}
func decodeUint16(buffer []byte, start int) (uint16, int) {
	return uint16((int(buffer[start]) << 8) | int(buffer[start+1])), start + 2
}
func decodeUint32(buffer []byte, start int) (uint32, int) {
	return uint32((int(buffer[start]) << 24) | (int(buffer[start+1]) << 16) | (int(buffer[start+2]) << 8) | int(buffer[start+3])), start + 4
}


