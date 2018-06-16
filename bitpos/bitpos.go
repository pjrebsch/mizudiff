package bitpos

import (
  "math"
  "math/big"
)

const ByteBitCount = 8

// Storing a bit position itself could be a very large number, potentially
// overflowing the containing variable's bit width, so instead, we store
// the offset of the byte within which the bit position falls, then we store
// the offset that the last bit has in that byte.
type BitPosition struct {
  ByteOffset uint
  BitOffset uint
}

func New(ByteOffset uint, BitOffset uint) BitPosition {
  bp := BitPosition{ ByteOffset, BitOffset }
  bp.normalize()
  return bp
}

func (p BitPosition) Int() *big.Int {
  res := big.NewInt(ByteBitCount)
  res.Mul(res, big.NewInt(int64(p.ByteOffset)))
  res.Add(res, big.NewInt(int64(p.BitOffset)))
  return res
}

func (p BitPosition) Plus(other BitPosition) BitPosition {
  calc := p.BitOffset + other.BitOffset
  bytes := p.ByteOffset + other.ByteOffset + calc / ByteBitCount
  bits := calc % ByteBitCount
  return BitPosition{bytes, uint(bits)}
}

func (p BitPosition) Minus(other BitPosition) BitPosition {
  calc := float64(int(p.BitOffset) - int(other.BitOffset))
  borrow := uint(math.Min(0, math.Floor(calc / ByteBitCount)))
  bytes := p.ByteOffset - other.ByteOffset + borrow
  bits := int(math.Abs(calc)) % ByteBitCount
  if calc < 0 && bits > 0 {
    bits = ByteBitCount - bits
  }
  return BitPosition{bytes, uint(bits)}
}

func (p BitPosition) CeilByteOffset() uint {
  if p.ByteOffset < 0 {
    return 0
  }
  if p.BitOffset > 0 {
    return p.ByteOffset + 1
  }
  return p.ByteOffset
}

func (p *BitPosition) normalize() {
  p.ByteOffset += p.BitOffset / ByteBitCount
  p.BitOffset %= ByteBitCount
}
