package bitstr

import (
  "github.com/pjrebsch/mizudiff/bitpos"
  "errors"
  "math"
  // "fmt"
  // "strings"
  // "math/bits"
)

type BitString struct {
  bytes []byte  // raw data
  length bitpos.BitPosition  // bit length of the string
}

func New(bytes []byte) BitString {
  l := len(bytes)
  b := make([]byte, l)
  copy(b, bytes)

  s := BitString{ b, bitpos.New(int64(l), 0) }
  return s
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

func (s BitString) XORCompress(adv, win uint16) (BitString, error) {
  if adv > win {
    err := errors.New("advanceRate can't be greater than windowSize")
    return BitString{}, err
  }

  return BitString{}, nil
}

// func (s BitString) SplitBy(window_size uint16) []BitString {
//   var list []BitString
//
//   x := bitpos.BitPosition{}
//
//   for i := 0; x.ByteOffset < uint(len(s.Bytes)); i += 1 {
//     y := x.Plus(bitpos.BitPosition{ 0, uint(window_size) })
//
//     s := s.slice(x, y)
//     list = append(list, s)
//
//     x = bitpos.BitPosition{ y.ByteOffset, y.BitOffset }
//   }
//
//   return list
// }
//
// func (s BitString) ShiftRight(offset uint) BitString {
//   new_length := s.Length.Plus(bitpos.BitPosition{ 0, offset })
//
//   buf := make([]byte, new_length.CeilByteOffset())
//
//   // We may need to add additional "0" bytes to the beginning of the new string.
//   new_lead_length := new_length.Minus(s.Length)
//
//   for i := 0; uint(i) < new_lead_length.ByteOffset; i += 1 {
//     buf[i] = byte(0x00)
//   }
//
//   var saved byte
//   var shift_by uint = offset % bitpos.C
//
//   for i := new_lead_length.ByteOffset; i < new_length.CeilByteOffset(); i += 1 {
//     j := i - new_lead_length.ByteOffset
//
//     var b byte
//     if j < uint(len(s.Bytes)) {
//       b = s.Bytes[j]
//     }
//
//     buf[i] = b >> shift_by | saved
//     saved = b << (bitpos.C - shift_by)
//   }
//
//   return New(buf, new_length)
// }
//
// func (s BitString) ShiftLeft(offset uint) BitString {
//   new_length := s.Length.Minus(bitpos.BitPosition{ 0, offset })
//   if new_length.CeilByteOffset() <= 0 {
//     return BitString{ []byte(""), bitpos.BitPosition{} }
//   }
//
//   buf := make([]byte, len(s.Bytes))
//
//   for i, b := range s.Bytes {
//     buf[len(buf) - i - 1] = byte(bits.Reverse8(b))
//   }
//
//   bs := BitString{ buf, bitpos.BitPosition{ uint(len(buf)), 0 } }
//   bs = bs.ShiftRight(offset)
//   // debug("SHIFT : %08b\n", bs.Bytes)
//
//   chopped := bs.Bytes[:new_length.CeilByteOffset()]
//   buf = make([]byte, len(chopped))
//   // debug("CHOP  : %08b\n", chopped)
//
//   for i, b := range chopped {
//     buf[len(buf) - i - 1] = byte(bits.Reverse8(b))
//   }
//
//   return New(buf, new_length)
//
//
//
//   //
//   // buf := make([]byte, new_length.CeilByteOffset())
//   //
//   // var saved byte
//   // var shift_by uint = offset % bitpos.C
//   // var start_at uint = offset / bitpos.C
//   //
//   // for i := 0; uint(i) < new_length.CeilByteOffset(); i += 1 {
//   //   j := uint(i) + start_at
//   //   b := s.Bytes[j]
//   //
//   //   buf[i] = b << shift_by | saved
//   //   saved = b >> (bitpos.C - shift_by)
//   // }
//   //
//   // return BitString{ buf, new_length }
// }

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

  if from.BitOffset() == 0 {
    // If there is no bit offset, the bytes can simply be copied.
    copy(buf[bufOff.ByteOffset():], s.bytes)
  } else {
    bytesLen := uint64(len(s.bytes))

    fromAbs := bitpos.Zero()
    fromAbs.Abs(from.Int)
    bitOff := uint8(fromAbs.BitOffset())

    for byteOff := uint64(0); bufOff.Cmp(length.Int) == -1; byteOff += 1 {
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
  }

  out := New(buf)
  out.SetLength(length)
  return out, nil
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
