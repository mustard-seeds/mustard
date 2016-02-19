package hash
import (
    "fmt"
    "time"
    "strconv"
    "math/rand"
    "mustard/base/time_util"
    "strings"
)

const (
    kFingerPrintSeed = 19860413
    m64 = 0xc6a4a7935bd1e995
    r64 = 47
    m32 = 0x5bd1e995
    r32 = 24
    // murmurhash3
    c1  = 0xcc9e2d51
    c2  = 0x1b873593

    // key gen symbols
    KEY_SYMBOLS_NUM = "23456789"  // rm 0,1
    KEY_SYMBOLS_LOWERCASE_LETTER = "abcdefghijklmnopqrstuvwxyz"  // rm l
    KEY_SYMBOLS_UPPERCASE_LETTER = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"  // rm I,O
)

func FingerPrint32(s string) uint32 {
    return MurmurHash3_32([]byte(s), kFingerPrintSeed)
}
func FingerPrint(s string) uint64 {
    return MurmurHash64A([]byte(s), kFingerPrintSeed)
}
func FingerprintToString(fp uint64) string {
    return fmt.Sprintf("%x",fp)
}
func StringToFingerprint(s string) (uint64,error) {
    return strconv.ParseUint(s, 16, 64)
}

// 64-bit hash for 64-bit platforms
func MurmurHash64A(key []byte, seed uint32) uint64 {
    keyLen := len(key)
    var h uint64 = uint64(seed) ^ (uint64(keyLen) * m64)
    i := 0
    for ; i+8 <= keyLen; i += 8 {
        k := uint64(key[i]) | uint64(key[i + 1]) << 8 |
            uint64(key[i + 2]) << 16 | uint64(key[i + 3]) << 24 |
            uint64(key[i + 4]) << 32 | uint64(key[i + 5]) << 40 |
            uint64(key[i + 6]) << 48 | uint64(key[i + 7]) << 56
        k *= m64;
        k ^= k >> r64;
        k *= m64;

        h ^= k;
        h *= m64;
    }
    switch keyLen - i {
    case 7:
        h ^= uint64(key[i+6]) << 48
        fallthrough
    case 6:
        h ^= uint64(key[i+5]) << 40
        fallthrough
    case 5:
        h ^= uint64(key[i+4]) << 32
        fallthrough
    case 4:
        h ^= uint64(key[i+3]) << 24
        fallthrough
    case 3:
        h ^= uint64(key[i+2]) << 16
        fallthrough
    case 2:
        h ^= uint64(key[i+1]) << 8
        fallthrough
    case 1:
        h ^= uint64(key[i])
        h *= m64;
    }

    h ^= h >> r64;
    h *= m64;
    h ^= h >> r64;

    return h;
}
// 32-bit hash
func MurmurHash32A(key []byte, seed uint32) uint32 {
    keyLen := len(key)
    var h uint32 = seed ^ (uint32(keyLen) * m32)
    i := 0
    for ; i+4 <= keyLen; i += 4 {
        k := uint32(key[i]) | uint32(key[i + 1]) << 8 |
            uint32(key[i + 2]) << 16 | uint32(key[i + 3]) << 24
        k *= m32;
        k ^= k >> r32;
        k *= m32;

        h *= m32;
        h ^= k;
    }
    // Handle the last few bytes of the input array
    switch keyLen - i {
    case 4:
        h ^= uint32(key[i+3]) << 24
        fallthrough
    case 3:
        h ^= uint32(key[i+2]) << 16
        fallthrough
    case 2:
        h ^= uint32(key[i+1]) << 8
        fallthrough
    case 1:
        h ^= uint32(key[i])
        h *= m32;
    }
    // Do a few final mixes of the hash to ensure the last few
    // bytes are well-incorporated.

    h ^= h >> 13
    h *= m32
    h ^= h >> 15
    return h
}

func MurmurHash3_32(key []byte, seed uint32) uint32 {
    keyLen := len(key)
    var h uint32 = seed
    //----------
    // body
    i := 0
    for ; i + 4 <= keyLen; i += 4 {
        k := uint32(key[i]) | uint32(key[i + 1]) << 8 | uint32(key[i + 2]) << 16 | uint32(key[i + 3]) << 24
        k *= c1
        k = (k << 15) | (k >> (32 - 15))
        k *= c2
        h ^= k
        h = (h << 13) | (h >> (32 - 13))
        h = h * 5 + 0xe6546b64
    }
    //----------
    // tail
    var k1 uint32
    switch keyLen - i {
    case 3:
        k1 ^= uint32(key[i + 2]) << 16
        fallthrough
    case 2:
        k1 ^= uint32(key[i + 1]) << 8
        fallthrough
    case 1:
        k1 ^= uint32(key[i])
        k1 *= c1
        k1 = (k1 << 15) | (k1 >> (32 - 15))
        k1 *= c2
        h ^= k1
    }
    //----------
    // finalization

    h ^= uint32(keyLen)

    h ^= h >> 16
    h *= 0x85ebca6b
    h ^= h >> 13
    h *= 0xc2b2ae35
    h ^= h >> 16
    return h
}

func RandomIntn(n int) int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(n)
}
func GenerateKey(length int) string {
    return GenerateKeyWithGroup(length, []string{KEY_SYMBOLS_NUM,KEY_SYMBOLS_LOWERCASE_LETTER,KEY_SYMBOLS_UPPERCASE_LETTER})
}
func GenerateEasierKey(length int) string {
    return GenerateKeyWithGroup(length, []string{KEY_SYMBOLS_NUM,KEY_SYMBOLS_UPPERCASE_LETTER})
}
func GenerateKeyWithGroup(length int, group []string) string {
    /*
        Generate a random key from the supplied symbol groups.
        At least one symbol from each group will be included. Unpredictable
        results if length is less than the number of symbol groups.
        Believed to be reasonably secure (with a reasonable key length!)
    */
    key := []byte{}
    for _,g := range group {
        key = append(key, g[RandomIntn(len(g))])
    }
    Shuffle(&key)
    if len(key) > length {
        key = key[:length]
    }
    length -= len(key)
    // then fill with random characters from all symbol groups
    symbols := strings.Join(group,"")
    for i := 0;i < length; i++ {
        key = append(key, symbols[RandomIntn(len(symbols))])
    }
    // finally shuffle to ensure first x characters aren't from a predictable group
    Shuffle(&key)
    return string(key)
}

func Shuffle(slice *[]byte) {
    for i := range *slice {
        rand.Seed(time_util.GetCurrentTimeStamp())
        j := RandomIntn(i + 1)
        (*slice)[i], (*slice)[j] = (*slice)[j],(*slice)[i]
    }
}