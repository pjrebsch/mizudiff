package bitpos_test

import (
  "testing"
  "github.com/pjrebsch/mizudiff/bitpos"
  "math"
  "math/big"
)

var tblNew = []struct {
  x1 uint32; x2 uint8  // byte offset inputs
  r  int64
}{
  {0, 0,  0},
  {1, 1,  9},
  {1, 10, 18},
  {math.MaxUint32, math.MaxUint8, math.MaxUint32 * bitpos.C + math.MaxUint8},
}

var tblPlus = []struct {
  x1 uint32; x2 uint8  // byte offset inputs
  y1 uint32; y2 uint8  // bit offset inputs
  r  int64             // expectation
}{
  {0,0, 0,0,  0},
  {1,2, 3,4,  38},
  {0,8, 0,1,  9},
  {9,9, 9,10, 163},
  { math.MaxUint32, math.MaxUint8,
    math.MaxUint32, math.MaxUint8,
    math.MaxUint32 * bitpos.C * 2 + math.MaxUint8 * 2 },
}

var tblMinus = []struct {
  x1 uint32; x2 uint8  // byte offset inputs
  y1 uint32; y2 uint8  // bit offset inputs
  r  int64             // expectation
}{
  {0,0, 0,0,  0},
  {1,2, 3,4,  -18},
  {0,8, 0,1,  7},
  {9,9, 9,10, -1},
  {0,0, math.MaxUint32, math.MaxUint8, -math.MaxUint32 * bitpos.C - math.MaxUint8},
}

var tblCeilByteOffset = []struct {
  x1 uint32; x2 uint8  // byte offset inputs
  r  uint64            // expectation
}{
  {0,0, 0},
  {0,1, 1},
  {1,0, 1},
  {1,1, 2},
  {2,7, 3},
  {math.MaxUint32, bitpos.C - 1, math.MaxUint32 + 1},
}

func TestIsEqual(t *testing.T) {
  a := bitpos.New(1,1)
  b := bitpos.New(1,1)

  if !bitpos.IsEqual(a,b) {
    t.Error("expected", a, "to equal", b)
  }
}

func TestNew(t *testing.T) {
  for _, e := range tblNew {
    a := bitpos.New( e.x1, e.x2 )
    b := bitpos.BitPosition{ big.NewInt(e.r) }

    if !bitpos.IsEqual(a,b) {
      t.Errorf(
        "New(%d, %d): expected %d, got %d",
        e.x1, e.x2, b, a,
      )
    }
  }
}

func TestPlus(t *testing.T) {
  for _, e := range tblPlus {
    x := bitpos.New( e.x1, e.x2 )
    y := bitpos.New( e.y1, e.y2 )

    actual := x.Plus(y)
    expected := bitpos.BitPosition{ big.NewInt(e.r) }

    if !bitpos.IsEqual(actual, expected) {
      t.Errorf(
        "%d.Plus(%d): expected %d, got %d",
        x, y, expected, actual,
      )
    }
  }
}

func TestMinus(t *testing.T) {
  for _, e := range tblMinus {
    x := bitpos.New( e.x1, e.x2 )
    y := bitpos.New( e.y1, e.y2 )

    actual := x.Minus(y)
    expected := bitpos.BitPosition{ big.NewInt(e.r) }

    if !bitpos.IsEqual(actual, expected) {
      t.Errorf(
        "%d.Minus(%d): expected %d, got %d",
        x, y, expected, actual,
      )
    }
  }
}

func TestCeilByteOffset(t *testing.T) {
  for _, e := range tblCeilByteOffset {
    x := bitpos.New( e.x1, e.x2 )

    actual := x.CeilByteOffset()
    expected := e.r

    if actual != expected {
      t.Errorf(
        "%d.CeilByteOffset(): expected %d, got %d",
        x, expected, actual,
      )
    }
  }
}
