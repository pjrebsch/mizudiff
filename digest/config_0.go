package digest

import(
  "github.com/pjrebsch/mizudiff/bitpos"
  "errors"
)

type Config_0 struct {
  ByteLength uint64
  BitLength uint8
}
func (c Config_0) AdvanceRate() uint16 {
  return 1
}
func (c Config_0) WindowSize() uint16 {
  return 8
}
func (c Config_0) DataLength() (bitpos.BitPosition, error) {
  p := bitpos.New( int64(c.ByteLength), int64(c.BitLength) )

  if p.Sign() == -1 {
    return bitpos.BitPosition{},
      errors.New("digest config byte length overflowed int64")
  }
  return p, nil
}
