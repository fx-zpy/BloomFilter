package BloomFilter

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

var dict = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	a := make([]rune, n)
	for i := range a {
		a[i] = dict[rand.Intn(len(dict))]
	}
	return string(a)
}

func TestAdd(t *testing.T) {
	filter := New(1024, 3, false)
	filter.Add([]byte("bloom")).
		AddString("filter").
		AddUint16(uint16(1)).
		AddUint32(uint32(2)).
		AddUint64(uint64(4)).
		AddUint16Batch([]uint16{17, 21, 38}).
		AddUint32Batch([]uint32{22, 31, 109}).
		AddUint64Batch([]uint64{35, 29, 91})

	t.Logf("bloom exist:%t", filter.Test([]byte("bloom")))
	t.Logf("filter exist:%t", filter.TestString("filter"))
	t.Logf("uint16(1) exist:%t", filter.TestUint16(uint16(1)))
	t.Logf("uint16(17) exist:%t", filter.TestUint16(uint16(17)))
	t.Logf("uint32(2) exist:%t", filter.TestUint32(uint32(2)))
	t.Logf("uint32(22) exist:%t", filter.TestUint32(uint32(22)))
	t.Logf("uint32(4) exist:%t", filter.TestUint64(uint64(4)))
	t.Logf("uint32(35) exist:%t", filter.TestUint64(uint64(35)))

	t.Logf("blllm exist:%t", filter.Test([]byte("blllm")))
	t.Logf("filtrr exist:%t", filter.TestString("filtrr"))
	t.Logf("uint16(2) exist:%t", filter.TestUint16(uint16(2)))
	t.Logf("uint16(21) exist:%t", filter.TestUint16(uint16(21)))
	t.Logf("uint32(67) exist:%t", filter.TestUint32(uint32(67)))
	t.Logf("uint32(31) exist:%t", filter.TestUint32(uint32(31)))
	t.Logf("uint32(3) exist:%t", filter.TestUint64(uint64(3)))
	t.Logf("uint32(91) exist:%t", filter.TestUint64(uint64(91)))

}

func TestAddAndGet(t *testing.T) {
	dataSize := 1000000
	dataMap := make(map[string]struct{}, dataSize)
	stringLen := 30
	filter := New(uint64(dataSize), 3, false)
	for i := 0; i < dataSize; i++ {
		randStr := RandStringRunes(stringLen)
		// add unique random string
		if _, ok := dataMap[randStr]; !ok {
			dataMap[randStr] = struct{}{}
			filter.Add([]byte(randStr))
		}
	}
	for k := range dataMap {
		exist := filter.Test([]byte(k))
		if !exist {
			t.Fatalf("key %s not exist", k)
		}
	}
}

func TestSync(t *testing.T) {
	sizeData := 100000
	stringLen := 30
	parts := 10

	filter := New(uint64(sizeData), 3, true)
	// concurrent write and read
	fn := func(size int, wg *sync.WaitGroup) {
		defer wg.Done()
		m := make(map[string]struct{}, size)
		for i := 0; i < size; i++ {
			randStr := RandStringRunes(stringLen)
			// add unique random string
			if _, ok := m[randStr]; !ok {
				m[randStr] = struct{}{}
				// write
				filter.AddString(randStr)
				// read
				exist := filter.TestString(randStr)
				if !exist {
					t.Errorf("key %s not exist", randStr)
				}
			}
		}
	}
	var waitGroup sync.WaitGroup
	for i := 0; i < parts; i++ {
		waitGroup.Add(1)
		go fn(sizeData/parts, &waitGroup)
	}
	waitGroup.Wait()
}

func TestFalsePositive(t *testing.T) {
	dataSize := 1000000
	dataNoSize := 100000
	dataMap := make(map[string]struct{}, dataSize)
	dataNoMap := make(map[string]struct{}, dataNoSize)
	stringLen := 30
	filter := New(uint64(dataSize), 3, false)

	for i := 0; i < dataSize; i++ {
		randStr := RandStringRunes(stringLen)
		if _, ok := dataMap[randStr]; !ok {
			dataMap[randStr] = struct{}{}
			filter.AddString(randStr)
		}
	}
	for i := 0; i < dataNoSize; i++ {
		randStr := RandStringRunes(stringLen)
		// add unique random string
		_, ok := dataMap[randStr]
		if !ok {
			dataNoMap[randStr] = struct{}{}
		}
	}
	falsePositiveCount := 0
	for k := range dataNoMap {
		exist := filter.TestString(k)
		if exist {
			falsePositiveCount++
		}
	}
	falsePositiveRatio := float64(falsePositiveCount) / float64(dataNoSize)
	t.Logf("false positive count:%d,false positive ratio:%f", falsePositiveCount, falsePositiveRatio)
}

func BenchmarkFilter_Add(b *testing.B) {
	b.StopTimer()
	dataTestSize := 100000
	dataTestMap := make(map[string]struct{}, dataTestSize)
	dataTestArr := make([]string, dataTestSize)
	stringLen := 100
	for i := 0; i < dataTestSize; i++ {
		randStr := RandStringRunes(stringLen)
		dataTestMap[randStr] = struct{}{}
		dataTestArr = append(dataTestArr, randStr)
	}
	filter := New(uint64(dataTestSize), 3, false)
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < len(dataTestArr); i++ {
		filter.Add([]byte(dataTestArr[i]))
	}

}
