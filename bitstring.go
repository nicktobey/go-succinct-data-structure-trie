package bits

/*
*

	Given a string of data (eg, in BASE-64), the BitString class supports
	reading or counting a number of bits from an arbitrary position in the
	string.
*/
type BitString struct {
	data   string
	length uint
}

// W is the width of a byte in bits
const W = 8

var MaskTop = [9]uint8{
	0xff, 0x7f, 0x3f, 0x1f, 0x0f, 0x07, 0x03, 0x01, 0x00,
}

var BitsInByte = [256]uint{
	0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4, 1, 2, 2, 3, 2, 3, 3, 4, 2,
	3, 3, 4, 3, 4, 4, 5, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3,
	3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3,
	4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4,
	3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5,
	6, 6, 7, 1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5, 2, 3, 3, 4, 3, 4,
	4, 5, 3, 4, 4, 5, 4, 5, 5, 6, 2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5,
	6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 2, 3, 3, 4, 3, 4, 4, 5,
	3, 4, 4, 5, 4, 5, 5, 6, 3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 3,
	4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7, 4, 5, 5, 6, 5, 6, 6, 7, 5, 6,
	6, 7, 6, 7, 7, 8,
}

func (bs *BitString) Init(data string) {
	bs.data = data
	bs.length = uint(len(bs.data)) * 8
}

/*
*

	Returns the internal string of bytes
*/
func (bs *BitString) GetData() string {
	return bs.data
}

/*
*

	Returns a decimal number, consisting of a certain number, n, of bits
	starting at a certain position, p.
*/
func (bs *BitString) Get(p, n uint) uint {

	if n == 0 {
		return 0
	}
	// case 1: bits lie within the given byte
	if (p%W)+n <= W {
		idx := p/W | 0
		return uint((bs.data[idx] & MaskTop[p%W]) >> (W - p%W - n))

		// case 2: bits lie incompletely in the given byte
	} else {
		idx := p/W | 0
		result := uint(bs.data[idx] & MaskTop[p%W])

		l := W - p%W
		p += l
		n -= l

		for n >= W {
			idx := p/W | 0
			result = (result << W) | uint(bs.data[idx])
			p += W
			n -= W
		}

		if n > 0 {
			idx := p/W | 0
			result = (result << n) | uint(bs.data[idx]>>
				(W-n))
		}

		return result
	}
}

/*
*

	Counts the number of bits set to 1 starting at position p and
	ending at position p + n
*/
func (bs *BitString) Count(p, n uint) uint {

	var count uint = 0
	for n >= 8 {
		count += BitsInByte[bs.Get(p, 8)]
		p += 8
		n -= 8
	}

	return count + BitsInByte[bs.Get(p, n)]
}

/*
*

	Returns the number of bits set to 1 up to and including position x.
	This is the slow implementation used for testing.
*/
func (bs *BitString) Rank(x uint) uint {
	var rank uint = 0
	var i uint = 0
	for i = 0; i <= x; i++ {
		// FIXME: the above line should be the following???
		//for i = 0; i < x; i++ {
		if bs.Get(i, 1) != 0 {
			rank++
		}
	}

	return rank
}
