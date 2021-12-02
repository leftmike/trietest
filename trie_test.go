package trietest_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/leftmike/trietest"
)

type testOp int

const (
	testDelete testOp = iota
	testGet
	testHash
	testPut
	testSerialize
)

type testCase struct {
	op         testOp
	k, v, h, s []byte
	notFound   bool
}

func testTrie(t *testing.T, who string, trie trietest.Trie, cases []testCase) {
	t.Helper()

	for _, c := range cases {
		switch c.op {
		case testDelete:
			err := trie.Delete(c.k)
			if c.notFound {
				if err != trietest.ErrNotFound {
					t.Errorf("%s.Delete(%v) returned %v, expected not found", who, c.k, err)
				}
			} else if err != nil {
				t.Errorf("%s.Delete(%v) failed with %s", who, c.k, err)
			}

		case testGet:
			v, err := trie.Get(c.k)
			if c.notFound {
				if err != trietest.ErrNotFound {
					t.Errorf("%s.Get(%#v) returned %v, expected not found", who, c.k, err)
				}
			} else if err != nil {
				t.Errorf("%s.Get(%#v) failed with %s", who, c.k, err)
			} else if !bytes.Equal(c.v, v) {
				t.Errorf("%s.Get(%#v): got %#v, want %#v", who, c.k, v, c.v)
			}

		case testHash:
			h := trie.Hash()
			if !bytes.Equal(c.h, h) {
				t.Errorf("%s.Hash(): got %#v, want %#v", who, h, c.h)
			}

		case testPut:
			err := trie.Put(c.k, c.v)
			if err != nil {
				t.Errorf("%s.Put(%#v, %#v) failed with %s", who, c.k, c.v, err)
			}

		case testSerialize:
			s, ok := trie.Serialize()
			if ok {
				if !bytes.Equal(s, c.s) {
					t.Errorf("%s.Serialize(): got %#v, want %#v", who, s, c.s)
				}
			}

		default:
			panic(fmt.Sprintf("unexpected test op: %#v", c.op))
		}
	}
}

func testBasic(t *testing.T, who string, newTrie func() trietest.Trie) {
	t.Helper()

	testTrie(t, who, newTrie(),
		[]testCase{
			{
				op: testPut,
				k:  []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
				v:  []byte("0123456789012345678901234567890123456789"),
			},
			{
				op: testHash,
				h: []byte{0xa3, 0x20, 0xa, 0x22, 0x8c, 0x3f, 0xfd, 0xaf, 0x4c, 0x2a, 0x76, 0xb7,
					0x22, 0xfb, 0x81, 0x9c, 0xae, 0x1, 0x9d, 0xe2, 0xd3, 0x88, 0xac, 0x9f, 0x8f,
					0xc4, 0x4d, 0x56, 0x55, 0xba, 0x61, 0x2d},
			},
			{
				op: testSerialize,
				s: []byte{0xf8, 0x5f, 0xb5, 0x20, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68,
					0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75,
					0x76, 0x77, 0x78, 0x79, 0x7a, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48,
					0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55,
					0x56, 0x57, 0x58, 0x59, 0x5a, 0xa8, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,
					0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32,
					0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39},
			},
		})

	key1 := []byte{0x00, 0x12, 0x34}
	val1 := []byte{0x01, 0x23, 0x45}
	key2 := []byte{0xA0, 0x12, 0x34}
	val2 := []byte{0xA1, 0x23, 0x45}

	testTrie(t, who, newTrie(),
		[]testCase{
			{
				op: testHash,
				h: []byte{0x56, 0xe8, 0x1f, 0x17, 0x1b, 0xcc, 0x55, 0xa6, 0xff, 0x83, 0x45, 0xe6,
					0x92, 0xc0, 0xf8, 0x6e, 0x5b, 0x48, 0xe0, 0x1b, 0x99, 0x6c, 0xad, 0xc0, 0x1,
					0x62, 0x2f, 0xb5, 0xe3, 0x63, 0xb4, 0x21},
			},
			{op: testGet, k: key1, notFound: true},
			{op: testPut, k: key1, v: val1},
			{op: testGet, k: key1, v: val1},
			{
				op: testHash,
				h: []byte{0x72, 0x94, 0xf6, 0x2, 0x63, 0xb8, 0x94, 0xf, 0x7a, 0x61, 0x11, 0x12,
					0x25, 0x42, 0x7a, 0xfc, 0x2f, 0xe2, 0xe8, 0x11, 0x13, 0x2a, 0xaa, 0x50, 0x86,
					0xcd, 0x9a, 0xa1, 0x27, 0xa6, 0x72, 0xe},
			},
			{
				op: testSerialize,
				s:  []byte{0xc9, 0x84, 0x20, 0x0, 0x12, 0x34, 0x83, 0x1, 0x23, 0x45},
			},
			{op: testGet, k: key2, notFound: true},
			{op: testPut, k: key2, v: val2},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, v: val2},
			{
				op: testHash,
				h: []byte{0xe0, 0x17, 0x26, 0x94, 0x41, 0xb7, 0x76, 0x97, 0xa3, 0x2f, 0x0, 0x62,
					0xdc, 0x3f, 0x7a, 0xaa, 0xa9, 0x7, 0x1, 0x6c, 0x64, 0x2a, 0x8e, 0x82, 0xae,
					0x28, 0xc1, 0x72, 0x59, 0xb2, 0x3e, 0xc7},
			},
			{
				op: testSerialize,
				s: []byte{0xe1, 0xc8, 0x83, 0x30, 0x12, 0x34, 0x83, 0x1, 0x23, 0x45, 0x80, 0x80,
					0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0xc8, 0x83, 0x30, 0x12, 0x34, 0x83,
					0xa1, 0x23, 0x45, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80},
			},
		})

	key3 := []byte{0x00, 0x23, 0x45}
	val3 := []byte{0x01, 0x01, 0x01}

	testTrie(t, who, newTrie(),
		[]testCase{
			{op: testGet, k: key1, notFound: true},
			{op: testPut, k: key1, v: val1},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, notFound: true},
			{op: testPut, k: key3, v: val3},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, v: val3},
			{
				op: testHash,
				h: []byte{0x48, 0x30, 0xba, 0x5a, 0x8c, 0x0, 0xae, 0x34, 0xeb, 0x70, 0xe8, 0x90,
					0x30, 0xa4, 0x91, 0x84, 0x28, 0xea, 0x11, 0xba, 0x49, 0x88, 0xc3, 0xa8, 0x4a,
					0xf4, 0x30, 0xee, 0xaa, 0xb1, 0xca, 0xfc},
			},
		})

	key4 := []byte{0x00, 0x34, 0x56}
	val4 := []byte{0x02, 0x03, 0x04}

	testTrie(t, who, newTrie(),
		[]testCase{
			{op: testGet, k: key1, notFound: true},
			{op: testPut, k: key1, v: val1},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, notFound: true},
			{op: testPut, k: key2, v: val2},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, v: val2},
			{op: testGet, k: key4, notFound: true},
			{op: testPut, k: key4, v: val4},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key2, v: val2},
			{op: testGet, k: key4, v: val4},
			{
				op: testHash,
				h: []byte{0x30, 0x7c, 0x7f, 0x43, 0xab, 0x18, 0x12, 0xa6, 0x78, 0x57, 0x45, 0x26,
					0xf2, 0x5f, 0xb1, 0xa6, 0x11, 0xf0, 0xb3, 0x8e, 0xa5, 0xc, 0x59, 0x7d, 0x3a,
					0x44, 0x77, 0x6b, 0x8e, 0x7e, 0x60, 0x13},
			},
		})

	key5 := []byte{0x00}
	val5 := []byte{0x11, 0x22, 0x33, 0x44}

	testTrie(t, who, newTrie(),
		[]testCase{
			{op: testGet, k: key1, notFound: true},
			{op: testPut, k: key1, v: val1},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, notFound: true},
			{op: testPut, k: key3, v: val3},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, v: val3},
			{op: testGet, k: key4, notFound: true},
			{op: testPut, k: key4, v: val4},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, v: val3},
			{op: testGet, k: key4, v: val4},
			{op: testGet, k: key5, notFound: true},
			{op: testPut, k: key5, v: val5},
			{op: testGet, k: key1, v: val1},
			{op: testGet, k: key3, v: val3},
			{op: testGet, k: key4, v: val4},
			{op: testGet, k: key5, v: val5},
			{op: testGet, k: []byte{0x01}, notFound: true},
			{
				op: testHash,
				h: []byte{0xdc, 0xda, 0x5e, 0x3, 0x35, 0xa0, 0xcd, 0x7f, 0xf, 0x78, 0xe0, 0x8f,
					0x66, 0x53, 0xe7, 0x65, 0x27, 0xd, 0xf1, 0xc9, 0x81, 0xaf, 0x89, 0xf3, 0xc2,
					0x0, 0xb5, 0x43, 0x79, 0xfe, 0x2e, 0x1a},
			},
		})

	key6 := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0x00}
	val6 := []byte{0x66, 0x66}
	key7 := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0x01}
	val7 := []byte{0x77, 0x77}
	key8 := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0x10}
	val8 := []byte{0x88, 0x88}
	key9 := []byte{0x01, 0x23, 0x45, 0x00}
	val9 := []byte{0x99, 0x99}
	val9a := []byte{0x9a, 0x9A}

	testTrie(t, who, newTrie(),
		[]testCase{
			{op: testGet, k: key6, notFound: true},
			{op: testPut, k: key6, v: val6},
			{op: testGet, k: key6, v: val6},
			{op: testGet, k: key7, notFound: true},
			{op: testPut, k: key7, v: val7},
			{op: testGet, k: key6, v: val6},
			{op: testGet, k: key7, v: val7},
			{op: testGet, k: key8, notFound: true},
			{op: testPut, k: key8, v: val8},
			{op: testGet, k: key6, v: val6},
			{op: testGet, k: key7, v: val7},
			{op: testGet, k: key8, v: val8},
			{op: testGet, k: key9, notFound: true},
			{op: testPut, k: key9, v: val9},
			{op: testGet, k: key6, v: val6},
			{op: testGet, k: key7, v: val7},
			{op: testGet, k: key8, v: val8},
			{op: testGet, k: key9, v: val9},
			{op: testPut, k: key9, v: val9a},
			{op: testGet, k: key6, v: val6},
			{op: testGet, k: key7, v: val7},
			{op: testGet, k: key8, v: val8},
			{op: testGet, k: key9, v: val9a},
			{
				op: testHash,
				h: []byte{0x9d, 0xf, 0xef, 0xbb, 0xce, 0x5a, 0xe7, 0x63, 0xbe, 0x29, 0x77, 0xf0,
					0x81, 0x80, 0xf, 0x22, 0x9e, 0xa8, 0x4e, 0x4c, 0xcf, 0xe0, 0x6d, 0xb7, 0x81,
					0xfb, 0x62, 0xf0, 0x28, 0x19, 0xa4, 0x96},
			},
		})

	key10 := []byte{0x01, 0x23, 0x45, 0x67, 0x89}
	val10 := []byte{0xAA, 0xAA}
	key11 := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0x01}
	val11 := []byte{0xBB, 0xBB}
	key12 := []byte{0x01, 0x23, 0x45, 0x67}
	val12 := []byte{0xCC, 0xCC}

	testTrie(t, who, newTrie(),
		[]testCase{
			{op: testGet, k: key10, notFound: true},
			{op: testPut, k: key10, v: val10},
			{op: testGet, k: key10, v: val10},
			{op: testGet, k: key11, notFound: true},
			{op: testPut, k: key11, v: val11},
			{op: testGet, k: key10, v: val10},
			{op: testGet, k: key11, v: val11},
			{op: testGet, k: key12, notFound: true},
			{op: testPut, k: key12, v: val12},
			{op: testGet, k: key10, v: val10},
			{op: testGet, k: key11, v: val11},
			{op: testGet, k: key12, v: val12},
			{
				op: testHash,
				h: []byte{0xc, 0x1a, 0xc5, 0xc7, 0x84, 0x7d, 0xce, 0x2a, 0x26, 0x15, 0x71, 0xef,
					0x3a, 0x2a, 0xe2, 0x36, 0xa4, 0xe5, 0x7e, 0xad, 0x6b, 0xc9, 0x88, 0xa1, 0xae,
					0x3d, 0xc7, 0x35, 0x1, 0xfa, 0x28, 0x4c},
			},
		})

	testTrie(t, who, newTrie(),
		[]testCase{
			{op: testGet, k: key10, notFound: true},
			{op: testPut, k: key10, v: val10},
			{op: testGet, k: key10, v: val10},
			{op: testGet, k: key12, notFound: true},
			{op: testPut, k: key12, v: val12},
			{op: testGet, k: key10, v: val10},
			{op: testGet, k: key12, v: val12},
			{op: testGet, k: key11, notFound: true},
			{op: testPut, k: key11, v: val11},
			{op: testGet, k: key10, v: val10},
			{op: testGet, k: key11, v: val11},
			{op: testGet, k: key12, v: val12},
			{
				op: testHash,
				h: []byte{0xc, 0x1a, 0xc5, 0xc7, 0x84, 0x7d, 0xce, 0x2a, 0x26, 0x15, 0x71, 0xef,
					0x3a, 0x2a, 0xe2, 0x36, 0xa4, 0xe5, 0x7e, 0xad, 0x6b, 0xc9, 0x88, 0xa1, 0xae,
					0x3d, 0xc7, 0x35, 0x1, 0xfa, 0x28, 0x4c},
			},
		})
}

func TestBasic(t *testing.T) {
	testBasic(t, "mptrie", trietest.NewMPTrie)
	testBasic(t, "zhang", trietest.NewZhangTrie)
	testBasic(t, "eth", trietest.NewEthTrie)
}

func testDeleteTrie(t *testing.T, who string, trie trietest.Trie, k []byte) {
	t.Helper()

	err := trie.Delete(k)
	if err != nil {
		t.Errorf("%s.Get(%#v) failed with %s", who, k, err)
	}
}

func testGetTrie(t *testing.T, who string, trie trietest.Trie, k, v []byte) {
	t.Helper()

	val, err := trie.Get(k)
	if err != nil {
		t.Errorf("%s.Get(%#v) failed with %s", who, k, err)
	} else if !bytes.Equal(val, v) {
		t.Errorf("%s.Get(%#v): got %#v, want %#v", who, k, val, v)
	}
}

func testHashTrie(t *testing.T, who string, trie trietest.Trie, hash []byte) {
	t.Helper()

	h := trie.Hash()
	if !bytes.Equal(hash, h) {
		t.Errorf("%s.Hash(): got %#v, want %#v", who, h, hash)
	}
}

func testPutTrie(t *testing.T, who string, trie trietest.Trie, k, v []byte) {
	t.Helper()

	err := trie.Put(k, v)
	if err != nil {
		t.Errorf("%s.Put(%#v, %#v) failed with %s", who, k, v, err)
	}
}

func makeKey(el int, lb []byte, ll int) []byte {
	if (el+ll)%2 == 1 {
		ll += 1
	}
	nk := append(bytes.Repeat([]byte{0x01}, el), bytes.Repeat(lb, ll)...)

	k := make([]byte, 0, (el+ll)/2)
	for ni := 0; ni < el+ll; ni += 2 {
		k = append(k, (nk[ni]<<4)|nk[ni+1])
	}
	return k
}

func testEdge(t *testing.T, who string, trie trietest.Trie, el, ll1, ll2, vl int) {
	t.Helper()

	k1 := makeKey(el, []byte{0x0A}, ll1)
	v1 := bytes.Repeat([]byte{0xAA}, vl)
	testPutTrie(t, who, trie, k1, v1)

	k2 := makeKey(el, []byte{0x0B}, ll2)
	v2 := []byte{0xBB}
	testPutTrie(t, who, trie, k2, v2)

	testGetTrie(t, who, trie, k1, v1)
	testGetTrie(t, who, trie, k2, v2)
}

func trieSize(el, ll1, ll2, vl int) int {
	return el/2 + ll1/2 + ll2/2 + vl + 1
}

func TestEdge(t *testing.T) {
	for vl := 1; vl < 35; vl++ {
		for el := 0; el < 70; el++ {
			for ll1 := 1; ll1 < 70; ll1++ {
				for _, ll2 := range []int{0, ll1, ll1 + 1} {
					sz := trieSize(el, ll1, ll2, vl)
					if sz < 28 || sz > 34 {
						continue
					}

					trie := trietest.NewEthTrie()
					testEdge(t, "eth", trie, el, ll1, ll2, vl)
					hash := trie.Hash()

					trie = trietest.NewMPTrie()
					testEdge(t, "mptrie", trie, el, ll1, ll2, vl)
					testHashTrie(t, "mptrie", trie, hash)

					trie = trietest.NewZhangTrie()
					testEdge(t, "zhang", trie, el, ll1, ll2, vl)
					testHashTrie(t, "zhang", trie, hash)
				}
			}
		}
	}
}

type keyValue struct {
	k, v []byte
}

func testGetPut(t *testing.T, who string, trie trietest.Trie, seed int64, kv []keyValue) {
	t.Helper()

	r := rand.New(rand.NewSource(seed))
	for i := range r.Perm(len(kv)) {
		testPutTrie(t, who, trie, kv[i].k, kv[i].v)
	}

	for i := range r.Perm(len(kv)) {
		testGetTrie(t, who, trie, kv[i].k, kv[i].v)
	}
}

func randomBytes(r *rand.Rand, min, max int) []byte {
	bl := r.Intn(max-min+1) + min
	b := make([]byte, 0, bl)
	for bl > 0 {
		b = append(b, byte(r.Intn(256)))
		bl -= 1
	}

	return b
}

func randomKeyValues(seed int64, n, minKey, maxKey, minVal, maxVal int) []keyValue {
	var kv []keyValue
	r := rand.New(rand.NewSource(seed))
	keys := map[string]struct{}{}

	for n > 0 {
		var k []byte
		for {
			k = randomBytes(r, minKey, maxKey)
			s := hex.EncodeToString(k)
			if _, ok := keys[s]; !ok {
				keys[s] = struct{}{}
				break
			}
		}

		kv = append(kv,
			keyValue{
				k: k,
				v: randomBytes(r, minVal, maxVal),
			})

		n -= 1
	}

	return kv
}

func testRandomGetPut(t *testing.T, seed int64, n int) {
	t.Helper()

	kv := randomKeyValues(seed, n, 1, 64, 1, 128)

	trie := trietest.NewEthTrie()
	testGetPut(t, "eth", trie, seed, kv)
	hash := trie.Hash()

	trie = trietest.NewMPTrie()
	testGetPut(t, "mptrie", trie, seed, kv)
	testHashTrie(t, "mptrie", trie, hash)

	trie = trietest.NewZhangTrie()
	testGetPut(t, "zhang", trie, seed, kv)
	testHashTrie(t, "zhang", trie, hash)
}

func TestRandomGetPut(t *testing.T) {
	start := time.Now()
	for {
		for _, n := range []int{20, 200, 2000, 20000} {
			seed := time.Now().UnixNano()
			testRandomGetPut(t, seed, n)
		}

		if testing.Short() {
			break
		}

		if time.Since(start).Seconds() > 60 {
			break
		}
	}
}

func testDeleteOk(t *testing.T, who string, trie trietest.Trie, kv []keyValue, bs []bool) {
	t.Helper()

	for i := 0; i < len(kv); i += 1 {
		if bs[i] {
			testDeleteTrie(t, who, trie, kv[i].k)
		}
	}
}

func testDeleteNotFound(t *testing.T, who string, trie trietest.Trie, kv []keyValue, bs []bool) {
	t.Helper()

	for i := 0; i < len(kv); i += 1 {
		if bs[i] {
			err := trie.Delete(kv[i].k)
			if err != trietest.ErrNotFound {
				t.Errorf("%s.Delete(%v) returned %v, expected not found", who, kv[i].k, err)
			}
		}
	}
}

func testGetNotFound(t *testing.T, who string, trie trietest.Trie, kv []keyValue, bs []bool) {
	t.Helper()

	for i := 0; i < len(kv); i += 1 {
		if bs[i] {
			_, err := trie.Get(kv[i].k)
			if err != trietest.ErrNotFound {
				t.Errorf("%s.Get(%v) returned %v, expected not found", who, kv[i].k, err)
			}
		}
	}
}

func testPutOk(t *testing.T, who string, trie trietest.Trie, kv []keyValue, bs []bool) {
	t.Helper()

	for i := 0; i < len(kv); i += 1 {
		if bs[i] {
			err := trie.Put(kv[i].k, kv[i].v)
			if err != nil {
				t.Errorf("%s.Put(%#v, %#v) failed with %s", who, kv[i].k, kv[i].v, err)
			}
		}
	}
}

func randomBoolSlice(seed int64, n, t int) []bool {
	bs := make([]bool, n)
	for t > 0 {
		t -= 1
		bs[t] = true
	}

	r := rand.New(rand.NewSource(seed))
	r.Shuffle(n,
		func(i, j int) {
			bs[i], bs[j] = bs[j], bs[i]
		})
	return bs
}

func testRandomDeleteGetPut(t *testing.T, seed int64, n int) {
	t.Helper()

	kv := randomKeyValues(seed, n, 1, 64, 1, 128)
	bs := randomBoolSlice(seed, n, n/4)

	trie := trietest.NewEthTrie()
	testGetPut(t, "eth", trie, seed, kv)
	hash1 := trie.Hash()
	testDeleteOk(t, "eth", trie, kv, bs)
	hash2 := trie.Hash()
	testDeleteNotFound(t, "eth", trie, kv, bs)
	testGetNotFound(t, "eth", trie, kv, bs)
	testHashTrie(t, "eth", trie, hash2)
	testPutOk(t, "eth", trie, kv, bs)
	testHashTrie(t, "eth", trie, hash1)

	trie = trietest.NewMPTrie()
	testGetPut(t, "mptrie", trie, seed, kv)
	testHashTrie(t, "mptrie", trie, hash1)
	testDeleteOk(t, "mptrie", trie, kv, bs)
	testHashTrie(t, "mptrie", trie, hash2)
	testDeleteNotFound(t, "mptrie", trie, kv, bs)
	testGetNotFound(t, "mptrie", trie, kv, bs)
	testHashTrie(t, "mptrie", trie, hash2)
	testPutOk(t, "mptrie", trie, kv, bs)
	testHashTrie(t, "mptrie", trie, hash1)
}

func TestRandomDeleteGetPut(t *testing.T) {
	start := time.Now()
	for {
		for _, n := range []int{20, 200, 2000, 20000} {
			seed := time.Now().UnixNano()
			testRandomDeleteGetPut(t, seed, n)
		}

		if testing.Short() {
			break
		}

		if time.Since(start).Seconds() > 60 {
			break
		}
	}
}

func TestRandomUpdate(t *testing.T) {
	// XXX: test Puts which are updates
}

func TestRandom(t *testing.T) {
	// XXX: test a random sequence of operations, keeping number of keys in some range
}
