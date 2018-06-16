package bitstr

import (
  "fmt"
  "math"
  "strings"
  "math/bits"
  "github.com/pjrebsch/mizudiff/bitpos"
)

type BitString struct {
  Bytes []byte
  Length bitpos.BitPosition
}

func (s BitString) Debug() {
  str := fmt.Sprintf("%08b", s.Bytes)
  str = strings.Replace(str, " ", "", -1)
  str = strings.Trim(str, "[]")

  bytestr := make([]byte, len(str))
  copy(bytestr, str)

  bits := s.Length.BitOffset

  end := len(str) + int(bits - bitpos.ByteBitCount) % bitpos.ByteBitCount
  bytestr = bytestr[:end]

  fmt.Printf("%s\n", bytestr)
}

func New(bytes []byte, length bitpos.BitPosition) BitString {
  bs := BitString{ bytes, length }
  bs.zeroExtraBits()
  return bs
}

func (s BitString) SplitBy(window_size uint16) []BitString {
  var list []BitString

  x := bitpos.BitPosition{}

  for i := 0; x.ByteOffset < uint(len(s.Bytes)); i += 1 {
    y := x.Plus(bitpos.BitPosition{ 0, uint(window_size) })

    s := s.slice(x, y)
    list = append(list, s)

    x = bitpos.BitPosition{ y.ByteOffset, y.BitOffset }
  }

  return list
}

func (s BitString) ShiftRight(offset uint) BitString {
  new_length := s.Length.Plus(bitpos.BitPosition{ 0, offset })

  buf := make([]byte, new_length.CeilByteOffset())

  // We may need to add additional "0" bytes to the beginning of the new string.
  new_lead_length := new_length.Minus(s.Length)

  for i := 0; uint(i) < new_lead_length.ByteOffset; i += 1 {
    buf[i] = byte(0x00)
  }

  var saved byte
  var shift_by uint = offset % bitpos.ByteBitCount

  for i := new_lead_length.ByteOffset; i < new_length.CeilByteOffset(); i += 1 {
    j := i - new_lead_length.ByteOffset

    var b byte
    if j < uint(len(s.Bytes)) {
      b = s.Bytes[j]
    }

    buf[i] = b >> shift_by | saved
    saved = b << (bitpos.ByteBitCount - shift_by)
  }

  return New(buf, new_length)
}

func (s BitString) ShiftLeft(offset uint) BitString {
  new_length := s.Length.Minus(bitpos.BitPosition{ 0, offset })
  if new_length.CeilByteOffset() <= 0 {
    return BitString{ []byte(""), bitpos.BitPosition{} }
  }

  buf := make([]byte, len(s.Bytes))

  for i, b := range s.Bytes {
    buf[len(buf) - i - 1] = byte(bits.Reverse8(b))
  }

  bs := BitString{ buf, bitpos.BitPosition{ uint(len(buf)), 0 } }
  bs = bs.ShiftRight(offset)
  // debug("SHIFT : %08b\n", bs.Bytes)

  chopped := bs.Bytes[:new_length.CeilByteOffset()]
  buf = make([]byte, len(chopped))
  // debug("CHOP  : %08b\n", chopped)

  for i, b := range chopped {
    buf[len(buf) - i - 1] = byte(bits.Reverse8(b))
  }

  return New(buf, new_length)



  //
  // buf := make([]byte, new_length.CeilByteOffset())
  //
  // var saved byte
  // var shift_by uint = offset % bitpos.ByteBitCount
  // var start_at uint = offset / bitpos.ByteBitCount
  //
  // for i := 0; uint(i) < new_length.CeilByteOffset(); i += 1 {
  //   j := uint(i) + start_at
  //   b := s.Bytes[j]
  //
  //   buf[i] = b << shift_by | saved
  //   saved = b >> (bitpos.ByteBitCount - shift_by)
  // }
  //
  // return BitString{ buf, new_length }
}

func (s BitString) slice(from, to bitpos.BitPosition) BitString {
  top := int(math.Max(0, float64(from.ByteOffset)))
  bot := int(math.Min(float64(len(s.Bytes)), float64(to.CeilByteOffset())))

  real_length := to.Minus(from)

  var out []byte

  buf := s.Bytes[top:bot]

  if from.BitOffset != 0 {
    // We have the bytes which contain the window that we are looking for,
    // but they are not offset properly.

    out = make([]byte, real_length.CeilByteOffset())

    for i := 0; i < len(buf); i += 2 {
      tmp := buf[i] << from.BitOffset

      if i + 1 < len(buf) {
        tmp |= buf[i+1] >> (bitpos.ByteBitCount - from.BitOffset)
      }

      out[i/2] = tmp
    }
  } else {
    out = buf
  }

  // We should zero any bits in the bytes that are outside of the length.
  return New(out, real_length)
}

func (s *BitString) zeroExtraBits() {
  res := make([]byte, len(s.Bytes))
  n := copy(res, s.Bytes)
  if s.Length.BitOffset > 0 {
    off := bitpos.ByteBitCount - s.Length.BitOffset
    res[n-1] = res[n-1] >> off << off
  }
  s.Bytes = res
}
