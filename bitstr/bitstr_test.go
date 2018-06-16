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
  bl int               // byte slice length
  x1 uint32; x2 uint8  // bit string length
}{
  {0, 0, 0},
  {15, 3, 4},
}

func makeBitString(
  bl int, b byte, x1 uint32, x2 uint8,
) (
  bitstr.BitString, []byte, bitpos.BitPosition,
) {
  bytes := make([]byte, bl)
  for i := range bytes {
    bytes[i] = b
  }

  p := bitpos.New(x1, x2)
  s := bitstr.New(bytes, p)

  return s, bytes, p
}

func TestNew(t *testing.T) {
  for _, e := range tblNew {
    s, b, p := makeBitString(e.bl, byte(0xff), e.x1, e.x2)

    t.Run("bit length", func(t *testing.T) {
      actual := s.Length().Uint64()
      expected := p.Uint64()

      if actual != expected {
        t.Errorf(
          "New(%02x, %d): expected bit length %d, got %d",
          b, p, expected, actual,
        )
      }
    })

    t.Run("byte length", func(t *testing.T) {
      actual := len(s.Bytes())
      expected := int(p.CeilByteOffset())

      if actual != expected {
        t.Errorf(
          "New(%02x, %d): expected byte length %d, got %d",
          b, p, expected, actual,
        )
      }
    })
  }

  s, b, p := makeBitString(10, byte(0xff), 3, 4)

  t.Run("byte equality", func(t *testing.T) {
    actual := s.Bytes()
    expected := []byte{ 0xff, 0xff, 0xff, 0xf0 }

    if !bytes.Equal(actual, expected) {
      t.Errorf(
        "New(%02x, %d): expected %02x, got %02x",
        b, p, expected, actual,
      )
    }
  })

  t.Run("use different slices", func(t *testing.T) {
    // Change the slice originally used to create the bit string.
    // The changes should not be seen in the bit string's bytes.
    b[0] = byte(0x00)

    actual := s.Bytes()[0]
    expected := byte(0xff)

    if actual != expected {
      t.Errorf(
        "New(%02x, %d): expected %02x, got %02x",
        b, p, expected, actual,
      )
    }
  })
}

func TestBytes(t *testing.T) {
  s, _, _ := makeBitString(15, byte(0xff), 15, 0)

  // A change to the return value of Bytes should not affect the bit string'
  // bytes.
  t.Run("use different slices", func(t *testing.T) {
    a := s.Bytes()
    a[0] = byte(0x00)
    b := s.Bytes()

    actual := b[0]
    expected := byte(0xff)

    if actual != expected {
      t.Errorf(
        "Bytes(): expected %02x, got %02x",
        expected, actual,
      )
    }
  })
}

// str := deterministicBytes(20, 2023114929827)
// t.Logf("%02x\n", str)
//
func deterministicBytes(length int, seed int64) []byte {
  src := rand.NewSource(seed)
  r := rand.New(src)

  str := make([]byte, length, length)

  for i := 0; i < length; i++ {
    str[i] = byte(r.Intn(math.MaxUint8))
  }

  return str
}
