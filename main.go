package main

import (
  "fmt"
  "log"
  "math"
  "encoding/hex"
  "io/ioutil"
  "bytes"
  "github.com/pjrebsch/mizudiff/bitstr"
  "github.com/pjrebsch/mizudiff/bitpos"
)

func loadSources() ([]byte, []byte) {
  a_raw, err := ioutil.ReadFile("tmp/srcC.hex.txt")
  if err != nil {
    log.Fatal(err)
  }
  b_raw, err := ioutil.ReadFile("tmp/srcD.hex.txt")
  if err != nil {
    log.Fatal(err)
  }

  a_raw = bytes.Trim(a_raw, "\n")
  b_raw = bytes.Trim(b_raw, "\n")

  a := make([]byte, hex.DecodedLen(len(a_raw)))
  b := make([]byte, hex.DecodedLen(len(b_raw)))

  _, err = hex.Decode(a, a_raw)
  if err != nil {
    log.Fatal(err)
  }
  _, err = hex.Decode(b, b_raw)
  if err != nil {
    log.Fatal(err)
  }

  return a, b
}

// See https://dave.cheney.net/2014/09/28/using-build-to-switch-between-debug-and-release
// for future implementation.
func debug(str string, args ...interface{}) {
	fmt.Printf(str, args...)
}

// SRC_BLOCK_SIZE must be even.
const SRC_BLOCK_SIZE = 4 // bytes

// SIG_ADVANCE_RATE must be greater than 0 but less than SRC_BLOCK_SIZE.
const SIG_ADVANCE_RATE = 1 // bytes per block

// SIG_ABSORPTION_LENGTH is the amount of overlap that two source blocks
// share in a signature.
const SIG_ABSORPTION_LENGTH = SRC_BLOCK_SIZE - SIG_ADVANCE_RATE

func calculateSourceSignatureLength(src_length int) int {
  if src_length <= SRC_BLOCK_SIZE {
    if src_length < 0 {
      return 0
    }
    return src_length
  }

  full_block_count := src_length / SRC_BLOCK_SIZE
  remainder_length := src_length % SRC_BLOCK_SIZE

  growth_due_to_full_blocks := (full_block_count - 1) * SIG_ADVANCE_RATE
  growth_due_to_remainder := int(math.Max(0, float64(remainder_length - SIG_ABSORPTION_LENGTH)))

  return SRC_BLOCK_SIZE + growth_due_to_full_blocks + growth_due_to_remainder
}

func translateSourceToSignaturePosition(pos int) int {
  full_block_count := pos / SRC_BLOCK_SIZE
  return pos % SRC_BLOCK_SIZE + SIG_ADVANCE_RATE * full_block_count
}

func generateSourceSignature(src []byte) []byte {
  result := make([]byte, calculateSourceSignatureLength(len(src)))

  prevtranspos := 0

  for i, b := range src {
    trans_pos := translateSourceToSignaturePosition(i)

    if trans_pos < prevtranspos {
      prefix := make([]byte, i / SRC_BLOCK_SIZE * (SIG_ADVANCE_RATE * 2))
      for i := range prefix {
        prefix[i] = ' '
      }
      fmt.Printf("\n%s", prefix)
    }
    prevtranspos = trans_pos
    fmt.Printf("%02x", b)

    result[trans_pos] ^= b
  }
  fmt.Print("\n")

  return result
}

func diffCompare(a, b []byte) ([]byte) {
  min_len := math.Min( float64(len(a)), float64(len(b)) )
  max_len := math.Max( float64(len(a)), float64(len(b)) )

  // Dividing by 2 will short the slice by 1, so we'll Ceil to fix that case.
  result_len := int(math.Ceil(max_len))

  result := make([]byte, result_len)

  for i := 0; i < int(max_len); i += 1 {
    abs_diff := byte(0xff)

    if i < int(min_len) {
      // Calculates the absolute difference between the individual bytes from
      // the two sources.
      abs_diff = a[i] ^ b[i]
    }

    // The absolute difference can only ever be as great as half of a byte,
    // so we can compress the resulting slice by 50% by adding consecutive
    // differences results.
    result[i] ^= abs_diff
  }

  return result
}

func prettyDiffComparison(diff []byte) {
  for _, val := range diff {
    if val == 0 {
      fmt.Print("██")
    } else {
      fmt.Print("--")
    }
  }
  fmt.Print("\n\n")
}

// Assumes that bit_pos is zero-based.
func bytePositionToNearestBitPosition(src []byte, bit_pos int) int {
  return bit_pos / (bitpos.ByteBitCount - 1)
}

func calculateNewLength(element_count, element_length, advance_rate uint) uint {
	return element_count * element_length - (element_length - advance_rate) * (element_count - 1)
}

func xorCompress(in []byte) []byte {
  var advance_rate uint = 3
  var window_size uint16 = 8

  root := bitstr.New(in, bitpos.New(uint(len(in)), 0))
  slices := root.SplitBy(window_size)

  new_bit_count := calculateNewLength(uint(len(slices)), uint(window_size), advance_rate)
  new_length := bitpos.New(0, new_bit_count)

  out := make([]byte, new_length.CeilByteOffset())

  for i := 0; i < len(slices); i += 1 {
    pos := bitpos.New(0, uint(i) * advance_rate)

    prefix := make([]byte, uint(i) * advance_rate)
    for d := range prefix {
      prefix[d] = ' '
    }
    fmt.Printf("%s", prefix)
    slices[i].Debug()

    for j, b := range slices[i].Bytes {
      k := pos.ByteOffset + uint(j)

      out[k] ^= b >> pos.BitOffset

      if pos.BitOffset > 0 {
        out[k+1] |= b << (bitpos.ByteBitCount - pos.BitOffset)
      }
    }
  }

  return out
}

func main() {
  log.Println("Go-ing...")
  a, b := loadSources()
  // a := []byte("abcd")
  // b := []byte("abcdefghijklmnopqrstuvwxyz...............abcdefghijklmnopqrstuvwxyz")
  // a := []byte("abcd")
  // b := []byte("0abcd")

  // sigA := xorCompress(a)[:]
  // sigB := xorCompress(b)[:]
  sigA := a
  sigB := b

  debug("BEFORE: %08b\n", sigB)
  bs := bitstr.New(sigB, bitpos.New(uint(len(sigB)), 0))
  bs = bs.ShiftLeft(7)
  sigB = bs.Bytes
  debug("AFTER : %08b\n", sigB)

  debug("len(A): %#v | len(sig_A): %#v | Ratio: %f%%\n", len(a), len(sigA), float64(len(sigA))/float64(len(a))*100)
  debug("len(B): %#v | len(sig_B): %#v | Ratio: %f%%\n", len(b), len(sigB), float64(len(sigB))/float64(len(b))*100)

  // for i := 0; i < 5; i += 1 {
  //   new_block := make([]byte, SRC_BLOCK_SIZE)
  //
  //   for j := 0; i != 0 && j < SRC_BLOCK_SIZE; j += 1 {
  //     new_block[j] = byte(i) & byte(j) ^ byte(0xff)
  //   }
  //
  //   a = append(a, new_block...)
  //   b = append(b, new_block...)
  //
  //   sigA, sigB := generateSourceSignature(a)[:], generateSourceSignature(b)[1:]
  //
  //   debug("New: %02x\n", new_block)
    // debug("%08b\n", sigA)
    // debug("%08b\n", sigB)
  //
    res := diffCompare(sigA, sigB)
    prettyDiffComparison(res)
  // }
}
