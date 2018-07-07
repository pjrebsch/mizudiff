package bitstr

import (
  "github.com/pjrebsch/mizudiff/bitpos"
  "errors"
  "math"
  "bytes"
)

type BitString struct {
  bytes []byte  // raw data
  length bitpos.BitPosition  // bit length of the string
}

func IsEqual(a, b BitString) bool {
  return bytes.Equal(a.bytes, b.bytes)
}

func New(bytes []byte) BitString {
  l := bitpos.New(int64(len(bytes)), 0)
  b := make([]byte, uint64(l.ByteOffset()))
  copy(b, bytes)
  return BitString{ b, l }
}

func (s BitString) Bytes() []byte {
  b := make([]byte, len(s.bytes))
  copy(b, s.bytes)
  return b
}

func (s BitString) Length() bitpos.BitPosition {
  return s.length
}

func (s *BitString) SetLength(p bitpos.BitPosition) error {
  if p.Sign() == -1 {
    return errors.New("length cannot be negative")
  }
  s.length = p

  err := s.updateDataSize()
  if err != nil {
    return err
  }

  s.zeroExtraBits()
  return nil
}

func (s BitString) Slice(from, length bitpos.BitPosition) (BitString, error) {
  if length.Sign() == -1 {
    return BitString{}, errors.New("length can't be less than zero")
  }

  l, err := length.CeilByteOffset()
  if err != nil {
    return BitString{}, err
  }

  buf := make([]byte, l)
  bufOff := bitpos.Zero()

  // If the starting position is negative, then we need to make the buffer
  // start with zero-bits for the offset of `from`.
  if from.Sign() == -1 {
    bufOff.Abs(from.Int)
  }

  bytesLen := uint64(len(s.bytes))

  fromAbs := bitpos.Zero()
  fromAbs.Abs(from.Int)
  bitOff := uint8(fromAbs.BitOffset())

  byteOff := uint64(0)
  if from.Sign() == 1 {
    byteOff = uint64(from.ByteOffset())
  }

  for ; bufOff.Cmp(length.Int) == -1; byteOff += 1 {
    thisPart, savedPart := byte(0x00), byte(0x00)

    if j := byteOff; j >= 0 && j < bytesLen {
      thisPart = s.bytes[j]

      if from.Sign() == -1 {
        thisPart >>= bitOff
      } else {
        thisPart <<= bitOff
      }
    }

    if from.Sign() == -1 {
      if j := byteOff - 1; j >= 0 && j < bytesLen {
        savedPart = s.bytes[j] << (bitpos.C - bitOff)
      }
    } else {
      if j := byteOff + 1; j >= 0 && j < bytesLen {
        savedPart = s.bytes[j] >> (bitpos.C - bitOff)
      }
    }

    buf[bufOff.ByteOffset()] = thisPart | savedPart

    bufOff.Add(bufOff.Int, bitpos.New(1,0).Int)
  }

  out := New(buf)
  out.SetLength(length)
  return out, nil
}

// Shift performs a bitwise shift on the bit string.
// A positive offset shifts right, negative offset shifts left.
func (s BitString) Shift(offset bitpos.BitPosition) (BitString, error) {
  from := bitpos.Zero()
  from.Neg(offset.Int)
  return s.Slice(from, s.Length())
}

func (s BitString) XORCompress(adv, win uint16) (BitString, error) {
  if adv == 0 {
    return BitString{}, errors.New("advance rate must be greater than zero")
  }
  if win == 0 {
    return BitString{}, errors.New("window size must be greater than zero")
  }
  if adv > win {
    err := errors.New("advance rate can't be greater than window size")
    return BitString{}, err
  }

  if len(s.bytes) == 0 {
    return BitString{}, nil
  }

  advRate := bitpos.New(0, int64(adv))
  winSize := bitpos.New(0, int64(win))

  windows := s.Length().CeilDividedBy(winSize)

  // This isn't true growth, because the first window never grows.
  growth := windows.MultipliedBy(advRate)

  // Correct the growth to get the real length.
  length := growth.Minus(advRate).Plus(winSize)

  l, err := length.CeilByteOffset()
  if err != nil {
    return BitString{}, err
  }

  out := make([]byte, l)

  // Bit index for `out`.
  i := bitpos.Zero()

  // Bit index for `s.bytes`.
  j := bitpos.Zero()

  for j.Cmp(s.Length().Int) == -1 {
    slice, err := s.Slice(j, winSize)
    if err != nil {
      return BitString{}, err
    }

    // Add a byte to the end of the buffer so that shifting right preserves
    // the latter bits.
    buf := append(slice.bytes, byte(0x00))

    // Only shift the bytes by the bit offset for `out`. The byte offset
    // will be taken care of later.
    bitOff := bitpos.New(0, i.BitOffset())
    shifted, err := New(buf).Shift(bitOff)
    if err != nil {
      return BitString{}, err
    }
    buf = shifted.bytes

    // Byte index for `buf`.
    m := int64(0)

    // Byte index for `out`.
    n := m + i.ByteOffset()

    for m < int64(len(buf)) && n < int64(len(out)) {
      out[n] ^= buf[m]
      m++
      n++
    }

    i = i.Plus(advRate)
    j = j.Plus(winSize)
  }

  r := New(out)
  r.SetLength(length)
  return r, nil
}

// Diff produces a bit string from two given bit strings which represents
// the difference between them. Each output bit represents the difference
// of each window compared between the bit strings. A 0 bit represents
// equality, whereas a 1 bit represents inequality.
//
// When given strings of different lengths, it will only compare as much as
// the length of the shortest string.
func Diff(a, b BitString, w bitpos.BitPosition) (BitString, error) {
  s := BitString{}

  if w.Sign() < 1 {
    return s, errors.New("window size must be greater than zero")
  }

  minLength := bitpos.Min(a.Length(), b.Length())
  outLength := minLength.CeilDividedBy(w)

  l, err := outLength.CeilByteOffset()
  if err != nil {
    return s, err
  }
  out := make([]byte, l)

  // Bit index for `out`.
  i := bitpos.Zero()

  // Bit index for `a` and `b`.
  j := bitpos.Zero()

  for i.Cmp(outLength.Int) < 0 {
    aWin, err := a.Slice(j, w)
    if err != nil {
      return s, err
    }
    bWin, err := b.Slice(j, w)
    if err != nil {
      return s, err
    }

    if !IsEqual(aWin, bWin) {
      out[i.ByteOffset()] |= 0x1 << (bitpos.C - uint8(i.BitOffset()) - 1)
    }

    i = i.Plus(bitpos.New(0,1))
    j = j.Plus(w)
  }

  s = New(out)
  s.SetLength(outLength)
  return s, nil
}

// func (s BitString) Debug() {
//   str := fmt.Sprintf("%08b", s.Bytes)
//   str = strings.Replace(str, " ", "", -1)
//   str = strings.Trim(str, "[]")
//
//   bytestr := make([]byte, len(str))
//   copy(bytestr, str)
//
//   bits := s.Length.BitOffset
//
//   end := len(str) + int(bits - bitpos.C) % bitpos.C
//   bytestr = bytestr[:end]
//
//   fmt.Printf("%s\n", bytestr)
// }

func (s *BitString) updateDataSize() error {
  n := int64(len(s.bytes))
  l, err := s.length.CeilByteOffset()
  if err != nil {
    return err
  }

  if l != n {
    least := uint64(math.Min(float64(l), float64(n)))
    buf := make([]byte, l)
    copy(buf, s.bytes[:least])
    s.bytes = buf
  }
  return nil
}

func (s *BitString) zeroExtraBits() {
  n := len(s.bytes)
  if n == 0 {
    return
  }
  if bits := s.length.BitOffset(); bits > 0 {
    off := uint64(bitpos.C - bits)
    s.bytes[n-1] = s.bytes[n-1] >> off << off
  }
}
