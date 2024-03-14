package bits

import (
	"bytes"
	"strings"
)

/*
*

	The BitWriter will create a stream of bytes, letting you write a certain
	number of bits at a time. This is part of the encoder, so it is not
	optimized for memory or speed.
*/
type BitWriter struct {
	bits []uint8
}

/*
*

	Write some data to the bit string. The number of bits must be 32 or
	fewer.
*/
func (bw *BitWriter) Write(data, numBits uint) {
	//for i := (numBits-1); i >= 0; i-- {
	// @siongui: the above commented line will cause infinite loop, why???
	// answer from @xphoenix:
	// Because i becomes uint, let's check iteration when i == 0, at the end
	// of loop, i-- takes place but as i is uint, it leads to 2^32-1 instead
	// of -1, loop condition is still true...
	for i := numBits; i > 0; i-- {
		j := i - 1
		if (data & (1 << j)) != 0 {
			bw.bits = append(bw.bits, 1)
		} else {
			bw.bits = append(bw.bits, 0)
		}
	}
}

/*
*

	Get the bitstring represented as a javascript string of bytes
*/
func (bw *BitWriter) GetData() string {
	var chars bytes.Buffer
	var b, i uint8 = 0, 0

	for j := 0; j < len(bw.bits); j++ {
		b = (b << 1) | bw.bits[j]
		i += 1
		if i == 8 {
			chars.WriteByte(b)
			i = 0
			b = 0
		}
	}

	if i != 0 {
		chars.WriteByte(b << (8 - i))
	}

	return chars.String()
}

/*
*

	Returns the bits as a human readable binary string for debugging
*/
func (bw *BitWriter) GetDebugString(group uint) string {
	var chars []string
	var i uint = 0

	for j := 0; j < len(bw.bits); j++ {
		if bw.bits[j] == 1 {
			chars = append(chars, "1")
		} else {
			chars = append(chars, "0")
		}
		i++
		if i == group {
			chars = append(chars, " ")
			i = 0
		}
	}

	return strings.Join(chars, "")
}
