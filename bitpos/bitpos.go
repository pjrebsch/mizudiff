package bitpos

import (
  "math/big"
)

// C represents the number of bits in a byte.
const C = 8

type BitPosition struct {
  *big.Int
}

func IsEqual(a, b BitPosition) bool {
  return a.Cmp(b.Int) == 0
}

// New allocates and returns a new BitPosition.
func New(byteOffset int64, bitOffset int8) BitPosition {
  p := big.NewInt(C)
  p.Mul(p, big.NewInt(byteOffset))
  p.Add(p, big.NewInt(int64(bitOffset)))
  return BitPosition{ p }
}

func (p BitPosition) Plus(other BitPosition) BitPosition {
  return BitPosition{ initInt().Add(p.Int, other.Int) }
}

func (p BitPosition) Minus(other BitPosition) BitPosition {
  return BitPosition{ initInt().Sub(p.Int, other.Int) }
}

func (p BitPosition) ByteOffset() int64 {
  r := initInt().Div(p.Int, big.NewInt(C))
  return r.Int64()
}

func (p BitPosition) BitOffset() int64 {
  r := initInt().Mod(p.Int, big.NewInt(C))
  return r.Int64()
}

// CeilByteOffset returns the byte index that the bit position corresponds
// to, or 1 greater.
//
// When the bit position divided by the "byte bit count" is still greater than
// a 64-bit value, overflow could occur, and this should be dealt with.
func (p BitPosition) CeilByteOffset() uint64 {
  r := p.Abs(p.Int)
  r.Add(p.Int, big.NewInt(C - 1))
  r.Div(r,     big.NewInt(C))
  return r.Uint64()
}

func initInt() *big.Int {
  return big.NewInt(0)
}
