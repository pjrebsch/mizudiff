package digest_test

import(
  "testing"
  "github.com/pjrebsch/mizudiff/digest"
)

var tblConstructors = []struct {
  raw []byte
  digest.Digest
}{
  { []byte{ 0x00, 0x00, 0x00, 0x00}, digest.Digest{} },
}

func TestNew(t *testing.T) {
  t.Run("version data can't be too short", func(t *testing.T) {
    raw := []byte{ 0x00, 0x00, 0x00 }
    _, err := digest.New(raw)

    expected := "digest data is too short to contain version info"
    if err.Error() != expected {
      t.Errorf(
        "New(0x%02x): expected %#v, but got %#v",
        raw, expected, err.Error(),
      )
    }
  })
  t.Run("can't have an unrecognized version", func(t *testing.T) {
    raw := []byte{ 0x10, 0x00, 0x00, 0x00 }
    _, err := digest.New(raw)

    expected := "digest version is not recognized"
    if err.Error() != expected {
      t.Errorf(
        "New(0x%02x): expected %#v, but got %#v",
        raw, expected, err.Error(),
      )
    }
  })
  t.Run("config data can't be too short", func(t *testing.T) {
    raw := []byte{ 0x00, 0x00, 0x00, 0x00 }
    _, err := digest.New(raw)

    expected := "digest data is too short to contain config info"
    if err.Error() != expected {
      t.Errorf(
        "New(0x%02x): expected %#v, but got %#v",
        raw, expected, err.Error(),
      )
    }
  })
  t.Run("config length can't be greater than the data length", func(t *testing.T) {
    raw := []byte{
      0x00, 0x00, 0x00, 0x00,  // version
      0x00, 0x01,  // advance rate
      0x00, 0x08,  // window size
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,  // byte length
      0x00,  // bit length
      0xff,  // data
    }
    _, err := digest.New(raw)

    expected := "configured length is greater than the actual data's length"
    if err.Error() != expected {
      t.Errorf(
        "New(0x%02x): expected %#v, but got %#v",
        raw, expected, err.Error(),
      )
    }
  })
}
