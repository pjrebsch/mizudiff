package bitstr_test

import(
  "testing"
  "github.com/pjrebsch/mizudiff/bitstr"
  "github.com/pjrebsch/mizudiff/bitpos"
  "math"
  "math/rand"
  "bytes"
)

var tblNew = []struct {
  init []byte           // initial byte data
  l bitpos.BitPosition  // bit string length
  expect []byte         // stored byte data
}{
  {
    []byte{ 0x15, 0xc3, 0x84, 0xb4, 0x7f },
    bitpos.New(0, 0),
    []byte{},
  }, {
    []byte{ 0x15, 0xc3, 0x84, 0xb4, 0x7f },
    bitpos.New(3, 4),
    []byte{ 0x15, 0xc3, 0x84, 0xb0 },
  },
}

func TestNew(t *testing.T) {
  for _, e := range tblNew {
    s := bitstr.New(e.init, e.l)

    t.Run("returns correct bit length", func(t *testing.T) {
      actual := s.Length().Uint64()
      expected := e.l.Uint64()

      if actual != expected {
        t.Errorf(
          "New(%02x, %d): expected bit length %d, got %d",
          e.init, e.l, expected, actual,
        )
      }
    })

    t.Run("returns correct byte length", func(t *testing.T) {
      actual := len(s.Bytes())
      expected := int(e.l.CeilByteOffset())

      if actual != expected {
        t.Errorf(
          "New(%02x, %d): expected byte length %d, got %d",
          e.init, e.l, expected, actual,
        )
      }
    })

    t.Run("refers to correct byte data", func(t *testing.T) {
      actual := s.Bytes()
      expected := e.expect

      if !bytes.Equal(actual, expected) {
        t.Errorf(
          "New(%02x, %d): expected %02x, got %02x",
          e.init, e.l, expected, actual,
        )
      }
    })
  }

  t.Run("uses a copy of the original slice", func(t *testing.T) {
    b := []byte{ 0xff, 0xff }
    l := bitpos.New(2,0)
    s := bitstr.New(b, l)

    // Change the slice originally used to create the bit string.
    // The changes should not be seen in the bit string's bytes.
    b[0] = byte(0x00)

    actual := s.Bytes()[0]
    expected := byte(0xff)

    if actual != expected {
      t.Errorf(
        "New(%02x, %d): expected %02x, got %02x",
        b, l, expected, actual,
      )
    }
  })
}

// func TestBytes(t *testing.T) {
//   t.Run("returns the initializing bytes", func(t *testing.T) {
//     s, b, p := makeBitString(15, byte(0xff), 15, 0)
//
//     actual := s.Bytes()
//     expected := b
//
//     if !bytes.IsEqual(actual, expected) {
//       t.Errorf("Bytes(): expected %02x, got %02x", expected, actual)
//     }
//   })
//
//   t.Run("returns separate slices", func(t *testing.T) {
//     s, b, _ := makeBitString(15, byte(0xff), 15, 0)
//
//     x := s.Bytes()
//     x[0] = byte(0x00)
//     y := s.Bytes()
//
//     actual := y[0]
//     expected := b[0]
//
//     if actual != expected {
//       t.Errorf("Bytes(): expected %02x, got %02x", expected, actual)
//     }
//   })
// }

// func makeBitString(
//   bl int, b byte, p bitpos.BitPosition,
// ) (
//   bitstr.BitString, []byte, bitpos.BitPosition,
// ) {
//   bytes := make([]byte, bl)
//   for i := range bytes {
//     bytes[i] = b
//   }
//
//   p := bitpos.New(x1, x2)
//   s := bitstr.New(bytes, p)
//
//   return s, bytes, p
// }

func deterministicBytes(length int, seed int64) []byte {
  src := rand.NewSource(seed)
  r := rand.New(src)

  str := make([]byte, length, length)

  for i := 0; i < length; i++ {
    str[i] = byte(r.Intn(math.MaxUint8))
  }

  return str
}
