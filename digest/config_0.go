package digest

import(
  "github.com/pjrebsch/mizudiff/bitpos"
  "errors"
)

type config_0 struct {
  advanceRate uint16
  windowSize uint16
  byteLength uint64
  bitLength uint8
}
func (c config_0) DataLength() (bitpos.BitPosition, error) {
  p := bitpos.New( int64(c.byteLength), int64(c.bitLength) )

  if p.Sign() == -1 {
    return bitpos.BitPosition{},
      errors.New("digest config byte length overflowed int64")
  }
  return p, nil
}
