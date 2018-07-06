package digest_test

import(
  "bytes"
  "testing"
  "github.com/pjrebsch/mizudiff/digest"
  "github.com/pjrebsch/mizudiff/bitstr"
)

func TestNew(t *testing.T) {
  var tbl = []struct {
    bytes []byte
  }{
    {
      []byte{0xf8, 0xac, 0x48, 0x6e, 0x0f, 0xda, 0x98, 0x69, 0x3c, 0x35},
    },
  }
  for _, e := range tbl {
    s := bitstr.New( e.bytes )

    adv := uint16(1)
    win := uint16(8)

    data, err := s.XORCompress(adv, win)
    if err != nil {
      t.Fatalf(
        "New(0x%02x): XORCompress(): did not expect an error, but got one: %#v",
        s.Bytes(), err.Error(),
      )
    }

    l := data.Length()

    d, err := digest.New(s)
    if err != nil {
      t.Fatalf(
        "New(0x%02x): did not expect an error, but got one: %#v",
        s.Bytes(), err.Error(),
      )
    }

    if d.Version != digest.CurrentVersion {
      t.Errorf(
        "New(0x%02x): expected version %v, got %v",
        s.Bytes(), digest.CurrentVersion, d.Version,
      )
    }

    c := d.Config.(digest.Config_0)

    if c.ByteLength != uint64(l.ByteOffset()) {
      t.Errorf(
        "New(0x%02x): expected config byte length %v, got %v",
        s.Bytes(), l.ByteOffset(), c.ByteLength,
      )
    }

    if c.BitLength != uint8(l.BitOffset()) {
      t.Errorf(
        "New(0x%02x): expected config bit length %v, got %v",
        s.Bytes(), l.BitOffset(), c.BitLength,
      )
    }

    if !bytes.Equal(d.Data.Bytes(), data.Bytes()) {
      t.Errorf(
        "New(0x%02x): expected data %02x, got %02x",
        s.Bytes(), data.Bytes(), d.Data.Bytes(),
      )
    }
  }
}

func TestLoad(t *testing.T) {
  t.Run("version data can't be too short", func(t *testing.T) {
    raw := []byte{ 0x00, 0x00, 0x00 }
    _, err := digest.Load(raw)

    expected := "digest data is too short to contain version info"
    if err.Error() != expected {
      t.Errorf(
        "Load(0x%02x): expected %#v, but got %#v",
        raw, expected, err.Error(),
      )
    }
  })
  t.Run("can't have an unrecognized version", func(t *testing.T) {
    raw := []byte{ 0x10, 0x00, 0x00, 0x00 }
    _, err := digest.Load(raw)

    expected := "digest version is not recognized"
    if err.Error() != expected {
      t.Errorf(
        "Load(0x%02x): expected %#v, but got %#v",
        raw, expected, err.Error(),
      )
    }
  })
  t.Run("config data can't be too short", func(t *testing.T) {
    raw := []byte{ 0x00, 0x00, 0x00, 0x00 }
    _, err := digest.Load(raw)

    expected := "digest data is too short to contain config info"
    if err.Error() != expected {
      t.Errorf(
        "Load(0x%02x): expected %#v, but got %#v",
        raw, expected, err.Error(),
      )
    }
  })
  t.Run("config length can't be greater than the data length", func(t *testing.T) {
    raw := []byte{
      0x00, 0x00, 0x00, 0x00,  // version
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,  // byte length
      0x00,  // bit length
      0xff,  // data
    }
    _, err := digest.Load(raw)

    expected := "configured length is greater than the actual data's length"
    if err.Error() != expected {
      t.Errorf(
        "Load(0x%02x): expected %#v, but got %#v",
        raw, expected, err.Error(),
      )
    }
  })

  // var tbl = []struct {
  //   raw []byte
  //   digest.Digest
  // }{
  //   { []byte{ 0x00, 0x00, 0x00, 0x00 }, digest.Digest{} },
  // }
}
