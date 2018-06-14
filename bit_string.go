package main

import (
  "fmt"
  "math"
  "strings"
  "math/bits"
)

type bitString struct {
  bytes []byte
  length bitPosition
}

func (self bitString) debug() {
  str := fmt.Sprintf("%08b", self.bytes)
  str = strings.Replace(str, " ", "", -1)
  str = strings.Trim(str, "[]")

  bytestr := make([]byte, len(str))
  copy(bytestr, str)

  bit_offset := self.length.bit_offset

  bytestr = bytestr[:len(str) + int(bit_offset - BITS_IN_BYTE) % BITS_IN_BYTE]

  fmt.Printf("%s\n", bytestr)
}

func (self bitString) slice(from, to bitPosition) bitString {
  top := int(math.Max(0, float64(from.byte_offset)))
  bot := int(math.Min(float64(len(self.bytes)), float64(to.ceilByteOffset())))

  real_length := to.minus(from)

  var out []byte

  buf := self.bytes[top:bot]

  if from.bit_offset != 0 {
    // We have the bytes which contain the window that we are looking for,
    // but they are not offset properly.

    out = make([]byte, real_length.ceilByteOffset())

    for i := 0; i < len(buf); i += 2 {
      tmp := buf[i] << from.bit_offset

      if i + 1 < len(buf) {
        tmp |= buf[i+1] >> (BITS_IN_BYTE - from.bit_offset)
      }

      out[i/2] = tmp
    }
  } else {
    out = buf
  }

  // We should zero any bits in the bytes that are outside of the length.
  res := bitString{ out, real_length }
  res.zeroExtraBits()
  return res
}

func (self *bitString) zeroExtraBits() {
  res := make([]byte, len(self.bytes))
  n := copy(res, self.bytes)
  if self.length.bit_offset > 0 {
    off := BITS_IN_BYTE - self.length.bit_offset
    res[n-1] = res[n-1] >> off << off
  }
  self.bytes = res
}

func (self bitString) splitBy(window_size uint16) []bitString {
  var list []bitString

  x := bitPosition{}

  for i := 0; x.byte_offset < uint(len(self.bytes)); i += 1 {
    y := x.plus(bitPosition{ 0, uint(window_size) })

    s := self.slice(x, y)
    list = append(list, s)

    x = bitPosition{ y.byte_offset, y.bit_offset }
  }

  return list
}

func (self bitString) shiftRight(offset uint) bitString {
  new_length := self.length.plus(bitPosition{ 0, offset })

  buf := make([]byte, new_length.ceilByteOffset())

  // We may need to add additional "0" bytes to the beginning of the new string.
  new_lead_length := new_length.minus(self.length)

  for i := 0; uint(i) < new_lead_length.byte_offset; i += 1 {
    buf[i] = byte(0x00)
  }

  var saved byte
  var shift_by uint = offset % BITS_IN_BYTE

  for i := new_lead_length.byte_offset; i < new_length.ceilByteOffset(); i += 1 {
    j := i - new_lead_length.byte_offset

    var b byte
    if j < uint(len(self.bytes)) {
      b = self.bytes[j]
    }

    buf[i] = b >> shift_by | saved
    saved = b << (BITS_IN_BYTE - shift_by)
  }

  bs := bitString{ buf, new_length }
  bs.zeroExtraBits()
  return bs
}

func (self bitString) shiftLeft(offset uint) bitString {
  new_length := self.length.minus(bitPosition{ 0, offset })
  if new_length.ceilByteOffset() <= 0 {
    return bitString{ []byte(""), bitPosition{} }
  }

  buf := make([]byte, len(self.bytes))

  for i, b := range self.bytes {
    buf[len(buf) - i - 1] = byte(bits.Reverse8(b))
  }

  bs := bitString{ buf, bitPosition{ uint(len(buf)), 0 } }
  bs = bs.shiftRight(offset)
  debug("SHIFT : %08b\n", bs.bytes)

  chopped := bs.bytes[:new_length.ceilByteOffset()]
  buf = make([]byte, len(chopped))
  debug("CHOP  : %08b\n", chopped)

  for i, b := range chopped {
    buf[len(buf) - i - 1] = byte(bits.Reverse8(b))
  }

  bs = bitString{ buf, new_length }
  bs.zeroExtraBits()
  return bs



  //
  // buf := make([]byte, new_length.ceilByteOffset())
  //
  // var saved byte
  // var shift_by uint = offset % BITS_IN_BYTE
  // var start_at uint = offset / BITS_IN_BYTE
  //
  // for i := 0; uint(i) < new_length.ceilByteOffset(); i += 1 {
  //   j := uint(i) + start_at
  //   b := self.bytes[j]
  //
  //   buf[i] = b << shift_by | saved
  //   saved = b >> (BITS_IN_BYTE - shift_by)
  // }
  //
  // return bitString{ buf, new_length }
}
