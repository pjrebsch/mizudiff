package main

import (
  "fmt"
  "log"
  "math"
  // "encoding/hex"
  "io/ioutil"
  // "bytes"
  "github.com/pjrebsch/mizudiff/bitstr"
  "github.com/pjrebsch/mizudiff/bitpos"
)

func loadSources() ([]byte, []byte) {
  a_raw, err := ioutil.ReadFile("testdata/git-2.9.5.tar.gz.txt")
  if err != nil {
    log.Fatal(err)
  }
  b_raw, err := ioutil.ReadFile("testdata/git-2.9.5.tar.gz.txt-2")
  if err != nil {
    log.Fatal(err)
  }

  return a_raw, b_raw

  // a_raw = bytes.Trim(a_raw, "\n")
  // b_raw = bytes.Trim(b_raw, "\n")
  //
  // a := make([]byte, hex.DecodedLen(len(a_raw)))
  // b := make([]byte, hex.DecodedLen(len(b_raw)))

  // _, err = hex.Decode(a, a_raw)
  // if err != nil {
  //   log.Fatal(err)
  // }
  // _, err = hex.Decode(b, b_raw)
  // if err != nil {
  //   log.Fatal(err)
  // }

  // return a, b
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

func main() {
  log.Println("Go-ing...")
  a, b := loadSources()
  // a := []byte("abcd")
  // b := []byte("abcdefghijklmnopqrstuvwxyz...............abcdefghijklmnopqrstuvwxyz")
  // a := []byte("abcd")
  // b := []byte("0abcd")

  sa, sb := bitstr.New(a), bitstr.New(b)

  log.Println("Digesting...")

  da, err := sa.XORCompress(1,8)
  if err != nil {
    log.Fatalln(err)
  }
  db, err := sb.XORCompress(1,8)
  if err != nil {
    log.Fatalln(err)
  }

  // db, _ = db.Shift(bitpos.New(0,-1))

  // debug("da: %02x\n", da.Bytes())
  // debug("db: %02x\n", db.Bytes())

  la, _ := da.Length().CeilByteOffset()
  lb, _ := db.Length().CeilByteOffset()

  debug("len(A): %#v | len(digest_A): %#v | Ratio: %f%%\n", len(a), la, float64(la)/float64(len(a))*100)
  debug("len(B): %#v | len(digest_B): %#v | Ratio: %f%%\n", len(b), lb, float64(lb)/float64(len(b))*100)

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
    // res := diffCompare(sigA, sigB)
    // prettyDiffComparison(res)
  // }

  log.Println("Diffing...")

  diff, err := bitstr.Diff(da, db, bitpos.New(1,0))
  if err != nil {
    log.Fatalln(err)
  }

  debug("diff: %08b\n", diff.Bytes())
}
