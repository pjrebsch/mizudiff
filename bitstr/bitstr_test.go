package bitstr_test

import(
  "testing"
  "github.com/pjrebsch/mizudiff/bitstr"
  "github.com/pjrebsch/mizudiff/bitpos"
  "math"
  "math/rand"
  "bytes"
)

var tblConstructors = []struct {
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

func TestIsEqual(t *testing.T) {
  var tbl = []struct {
    a, b []byte
    r bool
  }{
    { []byte{}, []byte{}, true },
    { []byte{0x00}, []byte{}, false },
    { []byte{}, []byte{0x00}, false },
    { []byte{0x00}, []byte{0xff}, false },
    { []byte{0xfa}, []byte{0xfa}, true },
  }
  for _, e := range tbl {
    aStr, bStr := bitstr.New(e.a), bitstr.New(e.b)

    actual := bitstr.IsEqual(aStr, bStr)
    expected := e.r

    if actual != expected {
      t.Errorf(
        "IsEqual(%02x, %02x): expected %v, got %v",
        aStr.Bytes(), bStr.Bytes(), expected, actual,
      )
    }
  }
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

func TestSetLength(t *testing.T) {
  var tblSetLength = []struct {
    byteOffset, bitOffset int64
    hasError bool
  }{
    {-1, 0, true},
    {0, -1, true},
    {0,  0, false},
    {0,  9, false},
    {10, 0, false},
  }
  for _, e := range tblSetLength {
    b := bitstr.New( []byte{ 0xff, 0xff } )
    p := bitpos.New( e.byteOffset, e.bitOffset )

    err := b.SetLength(p)
    actual := err != nil
    expected := e.hasError

    if actual != expected {
      t.Errorf(
        "SetLength(%d): expected error? %t, got %t",
        p, expected, actual,
      )
    }

    if err == nil {
      actual := b.Length()
      expected := p

      if !bitpos.IsEqual(actual, expected) {
        t.Errorf(
          "SetLength(%d): expected length to be %v, got %v",
          p, expected, actual,
        )
      }
    }
  }

  t.Run("zeros bits outside of the length", func(t *testing.T) {
    var tbl = []struct {
      l int64  // new bit length
      orig, new []byte
    }{
      {5, []byte{ 0xff }, []byte{ 0xf8 }},
      {12, []byte{ 0xff, 0xff, 0xff }, []byte{ 0xff, 0xf0 }},
      {18, []byte{ 0xff }, []byte{ 0xff, 0x00, 0x00 }},
    }
    for _, e := range tbl {
      b := bitstr.New( e.orig )
      p := bitpos.New( 0, e.l )

      b.SetLength(p)

      actual := b.Bytes()
      expected := e.new

      if !bytes.Equal(actual, expected) {
        t.Errorf(
          "SetLength(%d): expected %08b, got %08b",
          e.l, expected, actual,
        )
      }
    }
  })
}

func TestSlice(t *testing.T) {
  t.Run("length can't be less than zero", func(t *testing.T) {
    s := bitstr.New( []byte{} )
    from := bitpos.New(0,0)
    length := bitpos.New(0,-1)

    _, err := s.Slice(from, length)
    if err == nil {
      t.Errorf(
        "Slice(%d, %d): expected an error, but didn't get one",
        from, length,
      )
    }
  })

  var tbl = []struct {
    f1, f2 int64
    l1, l2 int64
    in, out []byte
  }{
    {0,0, 0,0, []byte{ 0xff }, []byte{}},
    {0,0, 0,1, []byte{ 0xff }, []byte{ 0x80 }},
    {0,0, 1,0, []byte{ 0xff }, []byte{ 0xff }},
    {0,6, 1,1, []byte{ 0xff }, []byte{ 0xc0, 0x00 }},
    {0,-1, 0,1, []byte{ 0xff }, []byte{ 0x00 }},
    {0,-1, 0,2, []byte{ 0xff }, []byte{ 0x40 }},
    {0,-1, 0,3, []byte{ 0xff }, []byte{ 0x60 }},
    {-1,0, 1,0, []byte{ 0xff }, []byte{ 0x00 }},
    {-1,0, 2,0, []byte{ 0xff }, []byte{ 0x00, 0xff }},
    {-1,0, 3,0, []byte{ 0xff }, []byte{ 0x00, 0xff, 0x00 }},
    {0,-3, 1,0, []byte{ 0x55 }, []byte{ 0x0a }},
    {0,-3, 2,0, []byte{ 0x66 }, []byte{ 0x0c, 0xc0 }},
    {1,0, 2,0, []byte{ 0xff, 0xff }, []byte{ 0xff, 0x00 }},
    {
      -2,-2,  4,7,
      []byte{ 0xd2, 0xbf, 0x78, 0xae },  // 11010010 10111111 01111000 10101110
      []byte{ 0x00, 0x00, 0x34, 0xaf, 0xde },  // 00000000 00000000 00110100 10101111 11011110
    },
  }
  for _, e := range tbl {
    s := bitstr.New(e.in)
    from := bitpos.New(e.f1,e.f2)
    length := bitpos.New(e.l1,e.l2)

    result, err := s.Slice(from, length)
    if err != nil {
      t.Fatalf(
        "Slice(%d, %d): errored: %v",
        from, length, err,
      )
    }
    if !bytes.Equal(result.Bytes(), e.out) {
      t.Errorf(
        "Slice(%d, %d): expected %08b, got %08b",
        from, length, e.out, result.Bytes(),
      )
    }
  }
}

func TestShift(t *testing.T) {
  var tbl = []struct {
    x1, x2 int64
    in, out []byte
  }{
    {0,1, []byte{ 0xff, 0xff }, []byte{ 0x7f, 0xff }},
    {0,2, []byte{ 0xff, 0xff }, []byte{ 0x3f, 0xff }},
    {0,4, []byte{ 0xff, 0xff }, []byte{ 0x0f, 0xff }},
    {1,0, []byte{ 0xff, 0xff }, []byte{ 0x00, 0xff }},
    {2,0, []byte{ 0xff, 0xff }, []byte{ 0x00, 0x00 }},
    {0,-1, []byte{ 0xff, 0xff }, []byte{ 0xff, 0xfe }},
    {0,-2, []byte{ 0xff, 0xff }, []byte{ 0xff, 0xfc }},
    {0,-4, []byte{ 0xff, 0xff }, []byte{ 0xff, 0xf0 }},
    {-1,0, []byte{ 0xff, 0xff }, []byte{ 0xff, 0x00 }},
    {-2,0, []byte{ 0xff, 0xff }, []byte{ 0x00, 0x00 }},
  }
  for _, e := range tbl {
    a := bitstr.New( e.in )
    off := bitpos.New( e.x1, e.x2 )

    b, err := a.Shift(off)
    if err != nil {
      t.Fatalf(
        "Shift(%d): did not expect an error, but got one: %v",
        off, err,
      )
    }

    actual := b.Bytes()
    expected := e.out

    if !bytes.Equal(actual, expected) {
      t.Errorf(
        "Shift(%d): expected %08b, got %08b",
        off, expected, actual,
      )
    }
  }
}

func TestXORCompress(t *testing.T) {
  t.Run("window size must be greater than zero", func(t *testing.T) {
    s := bitstr.New( []byte{} )
    adv := uint16(0)
    win := uint16(1)

    _, err := s.XORCompress(adv, win)
    if err == nil {
      t.Errorf(
        "XORCompress(%d, %d): expected an error, but didn't get one",
        adv, win,
      )
    }
  })
  t.Run("window size must be greater than zero", func(t *testing.T) {
    s := bitstr.New( []byte{} )
    adv := uint16(1)
    win := uint16(0)

    _, err := s.XORCompress(adv, win)
    if err == nil {
      t.Errorf(
        "XORCompress(%d, %d): expected an error, but didn't get one",
        adv, win,
      )
    }
  })
  t.Run("advance rate can't be greater than window size", func(t *testing.T) {
    s := bitstr.New( []byte{} )
    adv := uint16(2)
    win := uint16(1)

    _, err := s.XORCompress(adv, win)
    if err == nil {
      t.Errorf(
        "XORCompress(%d, %d): expected an error, but didn't get one",
        adv, win,
      )
    }
  })

  var tbl = []struct {
    in, out []byte
    advanceRate, windowSize uint16
  }{
    {
      []byte{},
      []byte{},
      1, 1,
    }, {
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

  for _, e := range tbl {
    s, err := bitstr.New(e.in).XORCompress( e.advanceRate, e.windowSize )
    if err != nil {
      t.Fatalf(
        "XORCompress(%d, %d): did not expect an error, but got one: %v",
        e.advanceRate, e.windowSize, err,
      )
    }

    actual := s.Bytes()
    expected := e.out

    if !bytes.Equal(actual, expected) {
      t.Errorf(
        "XORCompress(%d, %d): expected %08b, got %08b",
        e.advanceRate, e.windowSize, expected, actual,
      )
    }
  }
}

func TestDiff(t *testing.T) {
  var tbl = []struct {
    w1, w2 int64
    a, b, r []byte
  }{
    { 0,8, []byte{}, []byte{}, []byte{} },
    { 0,8, []byte{0x00}, []byte{}, []byte{} },
    { 0,8, []byte{0xf2}, []byte{0xf2}, []byte{0x00} },
    { 0,8, []byte{0x00,0x00}, []byte{0x00,0x11}, []byte{0x40} },
    {
      0,7,
      []byte{0xa8,0x1b}, // 1010100 0000110 11
      []byte{0x32},      // 0011001 0
      []byte{0xc0},      // 11
    }, {
      0,3,
      []byte{0xa8,0x1b}, // 101 010 000 001 101 1
      []byte{0xb4,0x7a}, // 101 101 000 111 101 0
      []byte{0x54}, // 010101
    },
  }
  for _, e := range tbl {
    w := bitpos.New(e.w1, e.w2)
    a := bitstr.New(e.a)
    b := bitstr.New(e.b)

    d, err := bitstr.Diff(a, b, w)
    if err != nil {
      t.Fatalf(
        "Diff(%08b, %08b, %d): did not expect an error but got one: %#v",
        a.Bytes(), b.Bytes(), w, err.Error(),
      )
    }

    actual := d.Bytes()
    expected := e.r

    if !bytes.Equal(actual, expected) {
      t.Errorf(
        "Diff(%08b, %08b, %d): expected %08b, got %08b",
        a.Bytes(), b.Bytes(), w, expected, actual,
      )
    }
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
