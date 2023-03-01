package BloomFilter

import (
	"errors"
	"fmt"
	"github.com/spaolacci/murmur3"
	"math"
	"sync"
)

type Filter struct {
	lock       *sync.RWMutex
	concurrent bool
	m          uint64
	n          uint64
	log2m      uint64
	k          uint64
	keys       []byte
}

const (
	mod7       = 1<<3 - 1
	bitPerByte = 8
)

/**
创建一个filter
*/
func New(size uint64, k uint64, race bool) *Filter {
	log2 := uint64(math.Ceil(math.Log2(float64(size))))
	filter := &Filter{
		m:          1 << log2,
		log2m:      log2,
		k:          k,
		keys:       make([]byte, 1<<log2),
		concurrent: race,
	}
	if filter.concurrent {
		filter.lock = &sync.RWMutex{}
	}
	return filter
}

/**
使用murmur3 hash
*/
func baseHash(data []byte) []uint64 {
	a := []byte{1}
	hasher := murmur3.New128()
	hasher.Write(data)
	v1, v2 := hasher.Sum128()
	hasher.Write(a)
	v3, v4 := hasher.Sum128()
	return []uint64{
		v1, v2, v3, v4,
	}
}

/**
location 返回字节数组中的位位置
*/
func (f *Filter) location(loc uint64) (uint64, uint64) {
	slot := (loc / bitPerByte) & (f.m - 1)
	mod := loc & mod7
	return slot, mod
}

/**
location 使用四个基本哈希值返回第 i 个哈希位置
*/
func location(h []uint64, i uint64) uint64 {
	return h[i&1] + i*h[2+(((i+(i&1))&3)/2)]
}

/**
将字节数组添加到布隆过滤器
*/
func (f *Filter) Add(data []byte) *Filter {
	if f.concurrent {
		f.lock.Lock()
		defer f.lock.Unlock()
	}
	h := baseHash(data)
	for i := uint64(0); i < f.k; i++ {
		loc := location(h, i)
		slot, mod := f.location(loc)
		f.keys[slot] |= 1 << mod
	}
	f.n++
	return f
}

/**
测试布隆过滤器中是否存在字节数组
*/
func (f *Filter) Test(data []byte) bool {
	if f.concurrent {
		f.lock.Lock()
		defer f.lock.Unlock()
	}
	h := baseHash(data)
	for i := uint64(0); i < f.k; i++ {
		loc := location(h, i)
		slot, mod := f.location(loc)
		if f.keys[slot]&(1<<mod) == 0 {
			return false
		}
	}
	return true
}

/**
将字符串添加到过滤器
*/
func (f *Filter) AddString(s string) *Filter {
	data := StrToBytes(s)
	return f.Add(data)
}

/**
测试字符串是否加入过滤器
*/
func (f *Filter) TestString(s string) bool {
	data := StrToBytes(s)
	return f.Test(data)
}

/**
将uint16添加到过滤器
*/
func (f *Filter) AddUint16(num uint16) *Filter {
	data := Uint16ToBytes(num)
	return f.Add(data)
}

/**
测试uint16是否加入过滤器
*/
func (f *Filter) TestUint16(num uint16) bool {
	data := Uint16ToBytes(num)
	return f.Test(data)
}

/**
将uint32添加到过滤器
*/
func (f *Filter) AddUint32(num uint32) *Filter {
	data := Uint32ToBytes(num)
	return f.Add(data)
}

/**
测试uint32是否加入过滤器
*/
func (f *Filter) TestUint32(num uint32) bool {
	data := Uint32ToBytes(num)
	return f.Test(data)
}

/**
将uint64添加到过滤器
*/
func (f *Filter) AddUint64(num uint64) *Filter {
	data := Uint64ToBytes(num)
	return f.Add(data)
}

/**
测试uint64是否加入过滤器
*/
func (f *Filter) TestUint64(num uint64) bool {
	data := Uint64ToBytes(num)
	return f.Test(data)
}

/**
批量添加字符数组
*/
func (f *Filter) AddBatch(dataarr [][]byte) *Filter {
	if f.concurrent {
		f.lock.Lock()
		defer f.lock.Unlock()
	}
	for i := 0; i < len(dataarr); i++ {
		data := dataarr[i]
		h := baseHash(data)
		for j := uint64(0); j < f.k; j++ {
			loc := location(h, j)
			slot, mod := f.location(loc)
			f.keys[slot] |= 1 << mod
		}
		f.n++
	}
	return f
}

/**
批量添加uint16数组
*/
func (f *Filter) AddUint16Batch(numarr []uint16) *Filter {
	data := make([][]byte, 0, len(numarr))
	for i := 0; i < len(numarr); i++ {
		byteArr := Uint16ToBytes(numarr[i])
		data = append(data, byteArr)
	}
	return f.AddBatch(data)
}

/**
批量添加uint32数组
*/
func (f *Filter) AddUint32Batch(numarr []uint32) *Filter {
	data := make([][]byte, 0, len(numarr))
	for i := 0; i < len(numarr); i++ {
		byteArr := Uint32ToBytes(numarr[i])
		data = append(data, byteArr)
	}
	return f.AddBatch(data)
}

/**
批量添加uint64数组
*/
func (f *Filter) AddUint64Batch(numarr []uint64) *Filter {
	data := make([][]byte, 0, len(numarr))
	for i := 0; i < len(numarr); i++ {
		byteArr := Uint64ToBytes(numarr[i])
		data = append(data, byteArr)
	}
	return f.AddBatch(data)
}

/**
将过滤器中使用的位重置为零
*/
func (f *Filter) Reset() {
	if f.concurrent {
		f.lock.Lock()
		defer f.lock.Unlock()
	}
	for i := 0; i < len(f.keys); i++ {
		f.keys[i] &= 0
	}
	f.n = 0

}

/**
将另一个过滤器合并到当前过滤器中
*/
func (f *Filter) MergeInPlace(g *Filter) error {
	if f.m != g.m {
		return fmt.Errorf("m's don't match: %d != %d", f.m, g.m)
	}

	if f.k != g.k {
		return fmt.Errorf("k's don't match: %d != %d", f.m, g.m)
	}
	if g.concurrent {
		return errors.New("merging concurrent filter is not support")
	}

	if f.concurrent {
		f.lock.Lock()
		defer f.lock.Unlock()
	}
	for i := 0; i < len(f.keys); i++ {
		f.keys[i] |= g.keys[i]
	}
	return nil
}

/**
获取位数组的容量
*/
func (f *Filter) Cap() uint64 {
	if f.concurrent {
		f.lock.RLock()
		defer f.lock.RUnlock()
	}
	return f.m
}

/**
获取插入元素个数
*/
func (f *Filter) Size() uint64 {
	if f.concurrent {
		f.lock.RLock()
		defer f.lock.RUnlock()
	}
	return f.n
}

/**
求假阳性率
*/
func (f *Filter) FalsePositiveRate() float64 {
	if f.concurrent {
		f.lock.RLock()
		defer f.lock.RUnlock()
	}
	expoInner := -(float64)(f.k*f.n) / float64(f.m)
	rate := math.Pow(1-math.Pow(math.E, expoInner), float64(f.k))
	return rate
}
