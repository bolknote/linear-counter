package main

import (
    "math"
    "math/big"
    "os"
    "io"
    "bufio"
    "crypto/md5"
    "flag"
    "fmt"
)

const intsize = 32

func estimateCardinality(len float64, setBits float64) float64 {
    return math.Log(len / (len - setBits)) * len
}

func lineToBitNum(line string, size int64) int64 {
    h := md5.New()
    io.WriteString(h, line)

    hash := new(big.Int).SetBytes(h.Sum(nil))
    bit  := new(big.Int)

    bit.Mod(hash, big.NewInt(size))

    return bit.Int64()
}

// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
func calcSetBits(v uint32) uint32 {
    v -= (v >> 1) & 0x55555555
    v = ((v >> 2) & 0x33333333) + v & 0x33333333
    v = ((v >> 4)  + v) & 0x0F0F0F0F
    v = ((v >> 8)  + v) & 0x00FF00FF
    v = ((v >> 16) + v) & 0x0000FFFF

    return v
}

func calcSetBitsInArray(vec []uint32) uint64 {
    var sum uint64 = 0

    for i := 0; i<len(vec); i++ {
        sum += uint64(calcSetBits(vec[i]))
    }

    return sum
}

func readFile(bitsize int64) uint64 {
    stdin   := bufio.NewReader(os.Stdin)
    bufsize := int(math.Ceil(float64(bitsize) / float64(intsize)))
    vec     := make([]uint32, bufsize)

    for {
        if line, err := stdin.ReadString('\n'); err == nil {
            bit := lineToBitNum(line, bitsize)

            byte       := bit / intsize
            bitinbyte  := uint32(bit % intsize)

            vec[byte] |= 1 << bitinbyte
        } else {
            break
        }
    }

    setBitsTotal := calcSetBitsInArray(vec)

    return uint64(estimateCardinality(float64(bitsize), float64(setBitsTotal)))
}

func main() {
    size := flag.Int64("size", 10000, "size in bits of bit vector")
    flag.Parse()

    fmt.Println("Vector size:", *size, "bits")
    fmt.Println("Estimated result:", readFile(*size))
}