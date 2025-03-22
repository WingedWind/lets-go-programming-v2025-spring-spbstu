package Requester

import (
	"errors"
	"math/rand/v2"

	"github.com/quaiion/go-practice/lru-cache/internal/CacheDir"
)

type Requester struct {
        reqRange uint32
}

func CreateRequester(reqRange uint32) Requester {
        return Requester{
                reqRange: reqRange,
        }
}

var errReqProcFailed = errors.New("failed processing a request")

func (r Requester) Request(cd *CacheDir.CacheDir) (bool, error) {
        hit, err := cd.GetRequest(rand.Uint32N(r.reqRange))
        if err != nil {
                return hit, errors.Join(errReqProcFailed, err)
        }

        return hit, nil
}

var errReqSerProcFailed = errors.New("failed processing a request series")

func (r Requester) RequestN(cd *CacheDir.CacheDir, n uint32) (uint32, error) {
        var nHits uint32 = 0

        for i := uint32(0) ; i < n ; i += 1 {
                hit, err := r.Request(cd)
                if err != nil {
                        return 0, errors.Join(errReqSerProcFailed, err)
                }

                if hit {
                        nHits += 1
                }
        }

        return nHits, nil
}
