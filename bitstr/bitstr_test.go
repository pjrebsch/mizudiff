package bitstr_test

import(
  "testing"
  "github.com/pjrebsch/mizudiff/bitstr"
  "github.com/pjrebsch/mizudiff/bitpos"
  "math"
  "math/rand"
  "bytes"
)

var tblConstructors = []struct{
  byteLen int
  strSeed int64
}{
  { 0, 80192744 },
  { 5, 2803843 },
  { 127, 199332 },
  { 256, 245532 },
  { 5000, 78229892 },
  { 10000, 3094199 },
}

func TestNew(t *testing.T) {
  t.Run("uses a copy of the original slice", func(t *testing.T) {
    b := []byte{ 0xff }
    s := bitstr.New(b)

    // Change the slice originally used to create the bit string.
    // The changes should not be seen in the bit string's bytes.
    b[0] = byte(0x00)

    actual := s.Bytes()[0]
    expected := byte(0xff)

    if actual != expected {
      t.Errorf(
        "New(%02x): expected %02x, got %02x",
        b, expected, actual,
      )
    }
  })
}

func TestBytes(t *testing.T) {
  for _, e := range tblConstructors {
    b := deterministicBytes(e.byteLen, e.strSeed)
    s := bitstr.New(b)

    t.Run("returns the initializing bytes", func(t *testing.T) {
      actual := s.Bytes()
      expected := b

      if !bytes.Equal(actual, expected) {
        t.Errorf("Bytes(): expected %02x, got %02x", expected, actual)
      }
    })
  }

  t.Run("returns separate slices", func(t *testing.T) {
    b := []byte{ 0xff }
    s := bitstr.New(b)

    x := s.Bytes()
    x[0] = byte(0x00)
    y := s.Bytes()

    actual := y[0]
    expected := b[0]

    if actual != expected {
      t.Errorf("Bytes(): expected %02x, got %02x", expected, actual)
    }
  })
}

func TestLength(t *testing.T) {
  for _, e := range tblConstructors {
    b := deterministicBytes(e.byteLen, e.strSeed)
    s := bitstr.New(b)

    t.Run("returns correct bit length", func(t *testing.T) {
      actual := s.Length().Uint64()
      expected := uint64(e.byteLen) * bitpos.C

      if actual != expected {
        t.Errorf(
          "New(%02x): expected bit length %d, got %d",
          b, expected, actual,
        )
      }
    })
  }

  t.Run("returns separate length structs", func(t *testing.T) {
    b := []byte{ 0xff }
    s := bitstr.New(b)

    x := s.Length()
    xPtr := &x
    y := s.Length()
    yPtr := &y

    if xPtr == yPtr {
      t.Errorf("Length(): expected address %p to differ from %p", xPtr, yPtr)
    }
  })
}

var tblXORCompress = []struct{
  in, out []byte
  advanceRate, windowSize uint8
}{
  {
    []byte{0xf8},
    []byte{0xf8},
    1, 8,
  }, {
    // 11111000  (0xf8)
    //  10000000  (0x80)
    // 101110000
    []byte{0xf8, 0x80},
    []byte{0xb8, 0x00},
    1, 8,
  }, {
    // 10011000  (0x98)
    //  01000000  (0x40)
    // 101110000
    []byte{0xf8, 0x80},
    []byte{0xb8, 0x00},
    1, 8,
  }, {
    // 11111000  (0xf8)
    //  10101100  (0xac)
    //   01001000  (0x48)
    //    01101110  (0x6e)
    //     00001111  (0x0f)
    //      11011010  (0xda)
    //       10011000  (0x98)
    //        01101001  (0x69)
    //         00111100  (0x3c)
    //          00110101  (0x35)
    // 10110101011101001
    []byte{0xf8, 0xac, 0x48, 0x6e, 0x0f, 0xda, 0x98, 0x69, 0x3c, 0x35},
    []byte{0xb5, 0x74, 0x80},
    1, 8,
  },
}

func TestXORCompress(t *testing.T) {
  // TODO: test empty byte slice as a special case where window size and
  // advance rate are irrelevant.

  for _, e := range tblXORCompress {
    t.Logf("%08b", e.in)
    t.Logf("%08b", e.out)
  }
}

func deterministicBytes(length int, seed int64) []byte {
  src := rand.NewSource(seed)
  r := rand.New(src)

  str := make([]byte, length, length)

  for i := 0; i < length; i++ {
    str[i] = byte(r.Intn(math.MaxUint8))
  }

  return str
}
