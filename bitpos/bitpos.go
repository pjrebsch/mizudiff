package bitpos

import (
  "errors"
  "math"
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

func Zero() BitPosition {
  return BitPosition{ big.NewInt(0) }
}

func Min(x, y BitPosition) BitPosition {
  if x.Cmp(y.Int) > 0 {
    return y
  }
  return x
}

func Max(x, y BitPosition) BitPosition {
  if x.Cmp(y.Int) < 0 {
    return y
  }
  return x
}

// New allocates and returns a new BitPosition.
func New(byteOffset, bitOffset int64) BitPosition {
  p := big.NewInt(C)
  p.Mul(p, big.NewInt(byteOffset))
  p.Add(p, big.NewInt(bitOffset))
  return BitPosition{ p }
}

func (p BitPosition) ByteOffset() int64 {
  r := Zero().Div(p.Int, big.NewInt(C))
  return r.Int64()
}

func (p BitPosition) BitOffset() int64 {
  r := Zero().Mod(p.Int, big.NewInt(C))
  return r.Int64()
}

func (p BitPosition) Plus(other BitPosition) BitPosition {
  return BitPosition{ Zero().Add(p.Int, other.Int) }
}

func (p BitPosition) Minus(other BitPosition) BitPosition {
  return BitPosition{ Zero().Sub(p.Int, other.Int) }
}

func (p BitPosition) DividedBy(other BitPosition) BitPosition {
  return BitPosition{ Zero().Div(p.Int, other.Int) }
}

func (p BitPosition) CeilDividedBy(other BitPosition) BitPosition {
  result, remainder := Zero().DivMod(p.Int, other.Int, Zero().Int)

  if remainder.Cmp(big.NewInt(0)) != 0 {
    result.Add(result, big.NewInt(int64(remainder.Sign())))
  }

  return BitPosition{ result }
}

func (p BitPosition) MultipliedBy(other BitPosition) BitPosition {
  return BitPosition{ Zero().Mul(p.Int, other.Int) }
}

// CeilByteOffset takes the absolute value bit index and returns the
// the ceiling byte offset that it would correspond to. This is primarily
// used for determining the correct byte slice size for a given bit string.
func (p BitPosition) CeilByteOffset() (int64, error) {
  if p.Cmp(big.NewInt(math.MaxInt64)) >= 0 {
    err := errors.New("reciever is greater than or equal to the max possible byte offset")
    return 0, err
  }
  if p.Cmp(big.NewInt(math.MinInt64)) <= 0 {
    err := errors.New("reciever is less than or equal to the min possible byte offset")
    return 0, err
  }

  r := Zero().Add(p.Int, big.NewInt(C - 1))
  r.Div(r, big.NewInt(C))
  return r.Int64(), nil
}
