package bitpos_test

import (
  "github.com/pjrebsch/mizudiff/bitpos"
  "testing"
)

var tblPlus = []struct {
  x1, x2, y1, y2 uint // inputs
  r1, r2         uint // results
}{
  {0,0, 0,0, 0,0},
  {1,2, 3,4, 4,6},
  {0,8, 0,0, 1,0},
  {9,9, 9,10, 20,3},
}

func TestNew(t *testing.T) {
  p := bitpos.New(0, bitpos.ByteBitCount + 2)

  if p.ByteOffset != 1 {
    t.Error("expected", 1, "got", p.ByteOffset)
  }

  if p.BitOffset != 2 {
    t.Error("expected", 2, "got", p.BitOffset)
  }
}

func TestInt(t *testing.T) {
  by := uint(10)
  bi := uint(1)

  p := bitpos.New(by, bi)
  x := p.Int()

  if !x.IsUint64() {
    t.Fatal("uint64 representation is not possible")
  }

  y := x.Uint64()
  z := uint64(by * bitpos.ByteBitCount + bi)
  if y != z {
    t.Error("expected", z, "got", y)
  }
}

func TestPlus(t *testing.T) {
  for _, e := range tblPlus {
    x := bitpos.BitPosition{ e.x1, e.x2 }
    y := bitpos.BitPosition{ e.y1, e.y2 }

    xActual := x.Plus(y)
    yActual := y.Plus(x)

    if xActual.ByteOffset != yActual.ByteOffset {
      t.Error(
        "(xActual byte offset)", xActual.ByteOffset,
        "!= (yActual byte offset)", yActual.ByteOffset,
      )
    }
    if xActual.BitOffset != yActual.BitOffset {
      t.Error(
        "(xActual bit offset)", xActual.BitOffset,
        "!= (yActual bit offset)", yActual.BitOffset,
      )
    }

    // If the commutative property of addition hasn't held, we
    // can't expect the following tests to pass reliably, so exit early.
    if t.Failed() {
      return
    }

    expected := bitpos.BitPosition{ e.r1, e.r2 }

    if xActual.ByteOffset != expected.ByteOffset {
      t.Error(
        "(actual byte offset)", xActual.ByteOffset,
        "!= (expected byte offset)", expected.ByteOffset,
      )
    }
    if xActual.BitOffset != expected.BitOffset {
      t.Error(
        "(actual bit offset)", xActual.BitOffset,
        "!= (expected bit offset)", expected.BitOffset,
      )
    }
  }
}
