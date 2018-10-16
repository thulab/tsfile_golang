package utils

import (
	_ "log"
	"bytes"
)

// get one bit in input byte. the offset is from low to high and start with 0
// e.g.<br>
// data:16(00010000), if offset is 4, return 1(000 "1" 0000) if offset is 7, return 0("0" 0010000)
func GetByteN(data byte, offset int) int {
	offset %= 8

	if (data & (1 << uint(7-offset))) != 0 {
		return 1
	} else {
		return 0
	}
}

/**
 * set one bit in input byte. the offset is from low to high and start with
 * index 0<br>
 * e.g.<br>
 * data:16(00010000), if offset is 4, value is 0, return 0({000 "0" 0000})
 * if offset is 1, value is 1, return 18({00010010}) if offset is 0, value
 * is 0, return 16(no change)
 *
 * @param data   input byte variable
 * @param offset bit offset
 * @param value  value to set
 * @return byte variable
 */
func SetByteN(data byte, offset int, value int) byte {
	offset %= 8

	if value == 1 {
		return (byte)(data | (1 << uint32(7-offset)))
	} else {
		return (byte)(data & ^(1 << uint32(7-offset)))
	}
}

/**
 * get one bit in input integer. the offset is from low to high and start
 * with 0<br>
 * e.g.<br>
 * data:1000(00000000 00000000 00000011 11101000), if offset is 4, return
 * 0(111 "0" 1000) if offset is 9, return 1(00000 "1" 1 11101000)
 *
 * @param data   input int variable
 * @param offset bit offset
 * @return 0 or 1
 */
func GetIntN(data int32, offset int) int32 {
	offset %= 32

	if (data & (1 << uint32(offset))) != 0 {
		return 1
	} else {
		return 0
	}
}

// set one bit in input integer. the offset is from low to high and start with index 0
// e.g.<br>
// data:1000({00000000 00000000 00000011 11101000}),
// if offset is 4, value is 1, return 1016({00000000 00000000 00000011 111 "1" 1000})
// if offset is 9, value is 0 return 488({00000000 00000000 000000 "0" 1 11101000})
// if offset is 0, value is 0 return 1000(no change)
func SetIntN(data int32, offset int, value int) int32 {
	offset %= 32

	if value == 1 {
		return (data | (1 << uint32(offset)))
	} else {
		return (data & ^(1 << uint32(offset)))
	}
}

/**
 * get one bit in input long. the offset is from low to high and start with
 * 0<br>
 *
 * @param data   input long variable
 * @param offset bit offset
 * @return 0/1
 */
func GetLongN(data int64, offset int) int32 {
	offset %= 64

	if (data & (int64(1) << uint32(offset))) != 0 {
		return 1
	} else {
		return 0
	}
}

// set one bit in input long. the offset is from low to high and start with index 0
func SetLongN(data int64, offset int, value int) int64 {
	offset %= 64

	if value == 1 {
		return (data | (1 << uint32(offset)))
	} else {
		return (data & ^(1 << uint32(offset)))
	}
}

// given a byte array, read width bits from specified position bits and convert it to an integer
func BytesToInt(result []byte, pos int, width int) int32 {
	var value int32 = 0
	var temp int = 0

	for i := 0; i < width; i++ {
		temp = (pos + width - 1 - i) / 8
		value = SetIntN(value, i, GetByteN(result[temp], pos+width-1-i))
	}
	return value
}

// given a byte array, read width bits from specified pos bits and convert it to an long
func BytesToLong(data []byte, pos int, width int) int64 {
	var value int64 = 0
	var temp int = 0

	for i := 0; i < width; i++ {
		temp = (pos + width - 1 - i) / 8
		value = SetLongN(value, i, GetByteN(data[temp], pos+width-1-i))
	}

	return value
}

/**
 * convert an integer to a byte array which length is width, then copy this
 * array to the parameter result from pos
 *
 * @param srcNum input integer variable
 * @param result byte array to convert
 * @param pos    start position
 * @param width  bit-width
 */
func IntToBytes(srcNum int32, result []byte, pos int, width int) {
	var temp int32 = 0
	for i := 0; i < width; i++ {
		temp = int32(pos+width-1-i) / 8
		result[temp] = SetByteN(result[temp], pos+width-1-i, int(GetIntN(srcNum, i)))
	}
}

/**
 * convert an long to a byte array which length is width, then copy this
 * array to the parameter result from pos
 *
 * @param srcNum input long variable
 * @param result byte array to convert
 * @param pos    start position
 * @param width  bit-width
 */
func LongToBytes(srcNum int64, result []byte, pos int, width int) {
	temp := 0
	for i := 0; i < width; i++ {
		temp = (pos + width - 1 - i) / 8
		result[temp] = SetByteN(result[temp], pos+width-1-i, int(GetLongN(srcNum, i)))
	}
}

func NumberOfLeadingZeros(i int32) int32 {
	if i == 0 {
		return 32
	}

	var n int32 = 1
	if i>>16 == 0 {
		n += 16
		i <<= 16
	}
	if i>>24 == 0 {
		n += 8
		i <<= 8
	}
	if i>>28 == 0 {
		n += 4
		i <<= 4
	}
	if i>>30 == 0 {
		n += 2
		i <<= 2
	}
	n -= int32(uint32(i) >> 31)

	return n
}

func NumberOfTrailingZeros(i int32) int32 {
	if i == 0 {
		return 32
	}

	var y int32
	var n int32 = 31
	y = i << 16
	if y != 0 {
		n = n - 16
		i = y
	}
	y = i << 8
	if y != 0 {
		n = n - 8
		i = y
	}
	y = i << 4
	if y != 0 {
		n = n - 4
		i = y
	}
	y = i << 2
	if y != 0 {
		n = n - 2
		i = y
	}

	return n - int32(uint32(i<<1)>>31)
}

func NumberOfLeadingZerosLong(i int64) int32 {
	if i == 0 {
		return 64
	}

	var n int32 = 1
	var x int32 = int32(uint64(i) >> 32)

	if x == 0 {
		n += 32
		x = int32(i)
	}
	if uint32(x)>>16 == 0 {
		n += 16
		x <<= 16
	}
	if uint32(x)>>24 == 0 {
		n += 8
		x <<= 8
	}
	if uint32(x)>>28 == 0 {
		n += 4
		x <<= 4
	}
	if uint32(x)>>30 == 0 {
		n += 2
		x <<= 2
	}
	n -= int32(uint32(x) >> 31)

	return n
}

func NumberOfTrailingZerosLong(i int64) int32 {
	if i == 0 {
		return 64
	}

	var x, y int32
	var n int32 = 63
	y = int32(i)

	if y != 0 {
		n = n - 32
		x = y
	} else {
		x = (int32)(uint64(i) >> 32)
	}
	y = x << 16
	if y != 0 {
		n = n - 16
		x = y
	}
	y = x << 8
	if y != 0 {
		n = n - 8
		x = y
	}
	y = x << 4
	if y != 0 {
		n = n - 4
		x = y
	}
	y = x << 2
	if y != 0 {
		n = n - 2
		x = y
	}

	return n - int32(uint32(x<<1)>>31)
}

/**
	* write a value to stream using unsigned var int format.
	* for example,
	* int 123456789 has its binary format 111010-1101111-0011010-0010101, 
	* function WriteUnsignedVarInt() will split every seven bits and write them to stream from low bit to high bit like:
	* 1-0010101 1-0011010 1-1101111 0-0111010
	* 1 represents has next byte to write, 0 represents number end.
	*/
func WriteUnsignedVarInt(value int32, buffer *bytes.Buffer) {
	var position int32 = 1

	for (value & 0x7FFFFF80) != 0 {
		buffer.WriteByte(byte((value & 0x7F) | 0x80))
		value = int32(uint32(value) >> 7)
		position++
	}

	buffer.WriteByte(byte(value & 0x7F))
}