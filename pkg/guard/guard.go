package guard

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
	"wow/pkg/pow"
)

type Guard struct {
	cache Cache
	rgen  *rand.Rand
}

// NewGuard
func NewGuard(cache Cache) *Guard {
	return &Guard{
		cache: cache,
		rgen:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Protect
func (p *Guard) Protect(_ context.Context, inputHeader int, inputPayload string, clientInfo string) (bool, int, string, error) {

	//
	switch inputHeader {
	case Quit:
		return false, Quit, "", fmt.Errorf("client %s close connection", clientInfo)
	case RequestChallenge:
		log.Printf("client %s challenge", clientInfo)

		//
		randValue := p.rgen.Intn(100000)
		err := p.cache.Add(randValue, pow.Duration)
		if err != nil {
			return false, Quit, "", fmt.Errorf("add rand to cache error: %w", err)
		}

		hashcash := pow.HashcashData{
			Ver:      1,
			Bits:     pow.ZerosCount,
			Date:     time.Now().Unix(),
			Resource: clientInfo,
			Rand:     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", randValue))),
			Counter:  0,
		}
		hashcashMarshaled, err := json.Marshal(hashcash)
		if err != nil {
			return false, Quit, "", fmt.Errorf("marshal hashcash error: %v", err)
		}

		//
		return false, ResponseChallenge, string(hashcashMarshaled), nil
	case RequestResource:
		log.Printf("client %s resource, payload %s", clientInfo, inputPayload)

		//
		var hashcash pow.HashcashData
		err := json.Unmarshal([]byte(inputPayload), &hashcash)
		if err != nil {
			return false, Quit, "", fmt.Errorf("unmarshal hashcash error: %w", err)
		}
		//
		if hashcash.Resource != clientInfo {
			return false, Quit, "", fmt.Errorf("invalid hashcash resource")
		}

		//
		randValueBytes, err := base64.StdEncoding.DecodeString(hashcash.Rand)
		if err != nil {
			return false, Quit, "", fmt.Errorf("decode rand error: %w", err)
		}
		randValue, err := strconv.Atoi(string(randValueBytes))
		if err != nil {
			return false, Quit, "", fmt.Errorf("decode rand error: %w", err)
		}

		//
		exists, err := p.cache.Get(randValue)
		if err != nil {
			return false, Quit, "", fmt.Errorf("get rand from cache error: %w", err)
		}
		if !exists {
			return false, Quit, "", fmt.Errorf("challenge expired or not sent")
		}

		//
		if time.Now().Unix()-hashcash.Date > pow.Duration {
			return false, Quit, "", fmt.Errorf("challenge expired")
		}

		//
		maxIter := hashcash.Counter
		if maxIter == 0 {
			maxIter = 1
		}
		_, err = hashcash.ComputeHashcash(maxIter)
		if err != nil {
			return false, Quit, "", fmt.Errorf("invalid hashcash")
		}

		//
		log.Printf("client %s computed hashcash %s", clientInfo, inputPayload)
		p.cache.Delete(randValue)
		return true, ResponseResource, "", nil
	default:
		return false, Quit, "", fmt.Errorf("header was undefined")
	}
}
