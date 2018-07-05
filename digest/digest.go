package digest

import(
  "github.com/pjrebsch/mizudiff/bitpos"
  "github.com/pjrebsch/mizudiff/bitstr"
  "encoding/binary"
  "errors"
)

const CurrentVersion = 0x0

// Versions defines the valid versions and the expected byte length of their
// configs.
var Versions = map[uint32]uint16 {
  0x0: 13,
}

type Digest struct {
  Version uint32
  Config interface{}
  Data bitstr.BitString
}

type Config interface {
  DataLength() (bitpos.BitPosition, error)
}

func New(s bitstr.BitString) (Digest, error) {
  c := Config_0{}
  c.AdvanceRate = 1
  c.WindowSize  = 8

  data, err := s.XORCompress(c.AdvanceRate, c.WindowSize)
  if err != nil {
    return Digest{}, err
  }

  l := data.Length()
  c.ByteLength  = uint64(l.ByteOffset())
  c.BitLength   = uint8(l.BitOffset())

  return Digest{ CurrentVersion, c, data }, nil
}

func Load(raw []byte) (Digest, error) {
  version, err := getVersion(raw)
  if err != nil {
    return Digest{}, err
  }

  // Offset is initially set to the size of the version data.
  offset := uint16(4)

  config, size, err := getConfig(version, raw[offset:])
  if err != nil {
    return Digest{}, err
  }

  offset += size

  data, err := getData(config, raw[offset:])
  if err != nil {
    return Digest{}, err
  }

  return Digest{ version, config, data }, nil
}

func getVersion(raw []byte) (uint32, error) {
  if len(raw) < 4 {
    return 0, errors.New("digest data is too short to contain version info")
  }
  return binary.BigEndian.Uint32(raw[:4]), nil
}

func getConfig(version uint32, raw []byte) (Config, uint16, error) {
  size, ok := Versions[version]
  if !ok {
    return nil, size, errors.New("digest version is not recognized")
  }
  if len(raw) < int(size) {
    return nil, size,
      errors.New("digest data is too short to contain config info")
  }

  // Slice just where the config should exist. This will cause an error if
  // and catch if the code tries to grab outside of where it should.
  s := raw[:size]

  if version == 0x0 {
    c := Config_0{}
    c.AdvanceRate = binary.BigEndian.Uint16(s[0:2])
    c.WindowSize  = binary.BigEndian.Uint16(s[2:4])
    c.ByteLength  = binary.BigEndian.Uint64(s[4:12])
    c.BitLength   = uint8(s[12])
    return c, size, nil
  }

  return nil, size, errors.New("no config was defined in the source code")
}

func getData(config Config, raw []byte) (bitstr.BitString, error) {
  l, err := config.DataLength()
  if err != nil {
    return bitstr.BitString{}, err
  }

  s := bitstr.New(raw)

  // There is a problem (possibly an attempted denial of service) if the
  // configuration's length is greater than the actual length of the data.
  if s.Length().Cmp(l.Int) == -1 {
    return bitstr.BitString{},
      errors.New("configured length is greater than the actual data's length")
  }

  s.SetLength(l)
  return s, nil
}
