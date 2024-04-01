package pow

import (
	"crypto/sha1"
	"fmt"
)

const (
	zeroByte      = 48
	ZerosCount    = 3
	Duration      = 300
	MaxIterations = 1000000
)

// HashcashData
// ver: Hashcash format version, 1 (which supersedes version 0).
// bits: Number of "partial pre-image" (zero) bits in the hashed code.
// date: The time that the message was sent, in the format YYMMDD[hhmm[ss]].
// resource: Resource data string being transmitted, e.g., an IP address or email address.
// rand: String of random characters, encoded in base-64 format.
// counter: Binary counter, encoded in base-64 format.
type HashcashData struct {
	Ver      int
	Bits     int
	Date     int64
	Resource string
	Rand     string
	Counter  int
}

// sha1Hash
func sha1Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// IsHashCorrect
func IsHashCorrect(hash string, zerosCount int) bool {
	if zerosCount > len(hash) {
		return false
	}
	for _, ch := range hash[:zerosCount] {
		if ch != zeroByte {
			return false
		}
	}
	return true
}

// ComputeHashcash
func (h HashcashData) ComputeHashcash(maxIterations int) (HashcashData, error) {
	for h.Counter <= maxIterations || maxIterations <= 0 {
		header := h.toString()
		hash := sha1Hash(header)
		if IsHashCorrect(hash, h.Bits) {
			return h, nil
		}
		h.Counter++
	}
	return h, fmt.Errorf("max iterations exceeded")
}

// toString
// example: 1:20:1303030600:anni@cypherspace.org::McMybZIhxKXu57jd:ckvi
func (h HashcashData) toString() string {
	return fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Ver, h.Bits, h.Date, h.Resource, h.Rand, h.Counter)
}
