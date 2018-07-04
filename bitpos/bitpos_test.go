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
  r int64       // expectation
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
  {2, -10, 6},
  {-2, 10, -6},
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

var tblOffsets = []struct {
  x1, x2 int64  // inputs
  r1, r2 int64  // expectations
}{
  {0,0, 0,0},
  {0,5, 0,5},
  {1,0, 1,0},
  {0,10, 1,2},
  {3,17, 5,1},
}
func TestByteOffset(t *testing.T) {
  for _, e := range tblOffsets {
    p := bitpos.New( e.x1, e.x2 )

    actual := p.ByteOffset()
    expected := e.r1

    if actual != expected {
      t.Errorf(
        "ByteOffset(): expected %d, got %d",
        expected, actual,
      )
    }
  }
}
func TestBitOffset(t *testing.T) {
  for _, e := range tblOffsets {
    p := bitpos.New( e.x1, e.x2 )

    actual := p.BitOffset()
    expected := e.r2

    if actual != expected {
      t.Errorf(
        "BitOffset(): expected %d, got %d",
        expected, actual,
      )
    }
  }
}

var tblPlus = []struct {
  x1, x2 int64  // byte offset inputs
  y1, y2 int64  // byte offset inputs
  r int64       // expectation
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
  r int64       // expectation
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

func TestDividedBy(t *testing.T) {
  var tbl = []struct {
    x1, x2 int64  // byte offset inputs
    y1, y2 int64  // bit offset inputs
    r1, r2 int64  // expectation
  }{
    {0,1, 0,1, 0,1},
    {1,0, 0,1, 1,0},
    {0,3, 0,2, 0,1},
    {0,1, 0,2, 0,0},
    {1,0, 0,2, 0,4},
    {-1,0, 0,1, -1,0},
    {-2,1, 0,2, -1,0},
    {-3,7, -1,0, 0,3},
  }
  for _, e := range tbl {
    x := bitpos.New( e.x1, e.x2 )
    y := bitpos.New( e.y1, e.y2 )

    actual := x.DividedBy(y)
    expected := bitpos.New( e.r1, e.r2 )

    if !bitpos.IsEqual(actual, expected) {
      t.Errorf(
        "%d.DividedBy(%d): expected %d, got %d",
        x, y, expected, actual,
      )
    }
  }
}

func TestMultipliedBy(t *testing.T) {
  var tbl = []struct {
    x1, x2 int64  // byte offset inputs
    y1, y2 int64  // bit offset inputs
    r1, r2 int64  // expectation
  }{
    {0,1, 0,1, 0,1},
    {1,0, 0,1, 1,0},
    {0,3, 0,2, 0,6},
    {0,1, 0,2, 0,2},
    {1,0, 0,2, 2,0},
    {-1,0, 0,1, -1,0},
    {-2,1, 0,2, -4,2},
    {-3,7, -1,0, 17,0},
  }
  for _, e := range tbl {
    x := bitpos.New( e.x1, e.x2 )
    y := bitpos.New( e.y1, e.y2 )

    actual := x.MultipliedBy(y)
    expected := bitpos.New( e.r1, e.r2 )

    if !bitpos.IsEqual(actual, expected) {
      t.Errorf(
        "%d.MultipliedBy(%d): expected %d, got %d",
        x, y, expected, actual,
      )
    }
  }
}

func TestCeilByteOffset(t *testing.T) {
  t.Run("bit position can't overflow int64", func(t *testing.T) {
    var tbl = []struct {
      x1, x2 int64  // inputs
    }{
      {math.MaxInt64, bitpos.C - 1},
      {math.MinInt64, -bitpos.C + 1},
      {math.MaxInt64 / bitpos.C + 1, bitpos.C - 1},
      {math.MinInt64 / bitpos.C - 1, -bitpos.C + 1},
    }
    for _, e := range tbl {
      x := bitpos.New(e.x1, e.x2)
      _, err := x.CeilByteOffset()

      if err == nil {
        t.Errorf(
          "%d.CeilByteOffset(): expected an error, but didn't get one",
          x,
        )
      }
    }
  })

  var tbl = []struct {
    x1, x2 int64  // byte offset inputs
    r int64       // expectation
  }{
    {0,0, 0},
    {0,1, 1},
    {1,0, 1},
    {1,1, 2},
    {2,7, 3},
    {0,-1, 0},
    {-1,-1, -1},
  }
  for _, e := range tbl {
    x := bitpos.New( e.x1, e.x2 )

    actual, err := x.CeilByteOffset()
    expected := e.r

    if err != nil {
      t.Fatalf(
        "%d.CeilByteOffset(): didn't expect an error, but got one: %v",
        x, err,
      )
    }
    if actual != expected {
      t.Errorf(
        "%d.CeilByteOffset(): expected %d, got %d",
        x, expected, actual,
      )
    }
  }
}
