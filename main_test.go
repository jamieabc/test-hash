package hashing

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/bitmark-inc/bitmarkd/blockdigest"
)

const (
	strSize   = 100
	iteration = 1000
)

func BenchmarkHashing(b *testing.B) {
	nCPU := runtime.NumCPU()
	var wg sync.WaitGroup

	fmt.Printf("%d cpus\n", nCPU)

	s := make([]byte, strSize)
	for i := 0; i < strSize; i++ {
		s[i] = byte(i)
	}

	var total uint32
	total = 0

	for cpu := 0; cpu < nCPU; cpu++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			counter := atomic.LoadUint32(&total)

			for counter < iteration {
				blockdigest.NewDigest(s)
				atomic.AddUint32(&total, 1)
				counter = atomic.LoadUint32(&total)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("%d hashes generated.", total)
}
