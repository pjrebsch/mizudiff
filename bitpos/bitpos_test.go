package bitpos_test

import (
  "testing"
  "github.com/pjrebsch/mizudiff/bitpos"
  "math/big"
)

var tblPlus = []struct {
  x1 uint32; x2 uint8  // byte offset inputs
  y1 uint32; y2 uint8  // bit offset inputs
  r1 uint32; r2 uint8  // expectation
}{
  {0,0, 0,0, 0,0},
  {1,2, 3,4, 4,6},
  {0,8, 0,1, 1,1},
  {9,9, 9,10, 20,3},
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
}

func TestIsEqual(t *testing.T) {
  a := bitpos.New(1,1)
  b := bitpos.New(1,1)

  if !bitpos.IsEqual(a,b) {
    t.Error("expected", a, "to equal", b)
  }
}

func TestNew(t *testing.T) {
  a := bitpos.New(1, bitpos.ByteBitCount + 2)
  b := bitpos.BitPosition{ big.NewInt(18) }

  if !bitpos.IsEqual(a,b) {
    t.Error("expected", a, "to equal", b)
  }
}

func TestPlus(t *testing.T) {
  for _, e := range tblPlus {
    x := bitpos.New( e.x1, e.x2 )
    y := bitpos.New( e.y1, e.y2 )

    actual := x.Plus(y)
    expected := bitpos.New( e.r1, e.r2 )

    if !bitpos.IsEqual(actual, expected) {
      t.Error("expected", actual, "to equal", expected)
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
      t.Error("expected", actual, "to equal", expected)
    }
  }
}

func TestCeilByteOffset(t *testing.T) {
  for _, e := range tblCeilByteOffset {
    x := bitpos.New( e.x1, e.x2 )

    actual := x.CeilByteOffset()
    expected := e.r

    if actual != expected {
      t.Error("expected", expected, "got", actual)
    }
  }
}
