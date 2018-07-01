package bitpos_test

import (
  "testing"
  "github.com/pjrebsch/mizudiff/bitpos"
  "math"
  "math/big"
)

func TestIsEqual(t *testing.T) {
  a := bitpos.New(1,1)
  b := bitpos.New(1,1)

  if !bitpos.IsEqual(a,b) {
    t.Error("expected", a, "to equal", b)
  }
}

var tblNew = []struct {
  x1, x2 int64  // byte offset inputs
  r  int64      // expectation
}{
  {0, 0,  0},
  {0, 1,  1},
  {1, 1,  9},
  {1, 10, 18},
  {math.MaxInt32, math.MaxInt8, math.MaxInt32 * bitpos.C + math.MaxInt8},
  {-1, 0, -8},
  {0, -1, -1},
  {-1, -1, -9},
  {-1, -10, -18},
  {math.MinInt32, math.MinInt8, math.MinInt32 * bitpos.C + math.MinInt8},
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

var tblPlus = []struct {
  x1, x2 int64  // byte offset inputs
  y1, y2 int64  // byte offset inputs
  r  int64      // expectation
}{
  {0,0, 0,0,  0},
  {1,2, 3,4,  38},
  {0,8, 0,1,  9},
  {9,9, 9,10, 163},
  { math.MaxInt32, math.MaxInt8,
    math.MaxInt32, math.MaxInt8,
    math.MaxInt32 * 2 * bitpos.C + math.MaxInt8 * 2 },
  {0,0, 0,-1, -1},
  {1,2, -3,-4, -18},
  {9,9, -9,-9, 0},
  { math.MinInt32, math.MinInt8,
    math.MinInt32, math.MinInt8,
    math.MinInt32 * 2 * bitpos.C + math.MinInt8 * 2 },
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

var tblMinus = []struct {
  x1, x2 int64  // byte offset inputs
  y1, y2 int64  // bit offset inputs
  r  int64      // expectation
}{
  {0,0, 0,0,  0},
  {1,2, 3,4,  -18},
  {0,8, 0,1,  7},
  {9,9, 9,10, -1},
  {0,0, math.MaxInt32, math.MaxInt8, -math.MaxInt32 * bitpos.C - math.MaxInt8},
  {0,0, 0,1, -1},
  {1,2, -3,-4, 38},
  {0,0, math.MinInt32, math.MinInt8, -math.MinInt32 * bitpos.C - math.MinInt8},
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

var tblCeilByteOffset = []struct {
  x1, x2 int64  // byte offset inputs
  r  uint64     // expectation
}{
  {0,0, 0},
  {0,1, 1},
  {1,0, 1},
  {1,1, 2},
  {2,7, 3},
  {math.MaxInt32,  bitpos.C - 1, math.MaxInt32 + 1},
  {0,-1, 1},
  {-1,-1, 2},
  {math.MinInt32, -bitpos.C + 1, math.MaxInt32 + 2},
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
