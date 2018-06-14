package main

import (
  "math"
)

// Storing a bit position itself could be a very large number, potentially
// overflowing the containing variable's bit width, so instead, we store
// the offset of the byte within which the bit position falls, then we store
// the offset that the last bit has in that byte.
type bitPosition struct {
  byte_offset uint
  bit_offset uint
}

func (self bitPosition) plus(other bitPosition) bitPosition {
  calc := self.bit_offset + other.bit_offset
  bytes := self.byte_offset + other.byte_offset + calc / BITS_IN_BYTE
  bits := calc % BITS_IN_BYTE
  return bitPosition{bytes, uint(bits)}
}

func (self bitPosition) minus(other bitPosition) bitPosition {
  calc := float64(int(self.bit_offset) - int(other.bit_offset))
  borrow := uint(math.Min(0, math.Floor(calc / BITS_IN_BYTE)))
  bytes := self.byte_offset - other.byte_offset + borrow
  bits := int(math.Abs(calc)) % BITS_IN_BYTE
  if calc < 0 && bits > 0 {
    bits = BITS_IN_BYTE - bits
  }
  return bitPosition{bytes, uint(bits)}
}

func (self *bitPosition) normalize() {
  self.byte_offset += self.bit_offset / BITS_IN_BYTE
  self.bit_offset %= BITS_IN_BYTE
}

func (self bitPosition) ceilByteOffset() uint {
  if self.byte_offset < 0 {
    return 0
  }
  if self.bit_offset > 0 {
    return self.byte_offset + 1
  }
  return self.byte_offset
}
