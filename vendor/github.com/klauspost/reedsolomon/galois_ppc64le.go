//go:build !noasm && !appengine && !gccgo
// +build !noasm,!appengine,!gccgo

// Copyright 2015, Klaus Post, see LICENSE for details.
// Copyright 2018, Minio, Inc.

package reedsolomon

//go:noescape
func galMulPpc(low, high, in, out []byte)

//go:noescape
func galMulPpcXor(low, high, in, out []byte)

// This is what the assembler routines do in blocks of 16 bytes:
/*
func galMulPpc(low, high, in, out []byte) {
	for n, input := range in {
		l := input & 0xf
		h := input >> 4
		out[n] = low[l] ^ high[h]
	}
}
func galMulPpcXor(low, high, in, out []byte) {
	for n, input := range in {
		l := input & 0xf
		h := input >> 4
		out[n] ^= low[l] ^ high[h]
	}
}
*/

func galMulSlice(c byte, in, out []byte, o *options) {
	if c == 1 {
		copy(out, in)
		return
	}
	done := (len(in) >> 4) << 4
	if done > 0 {
		galMulPpc(mulTableLow[c][:], mulTableHigh[c][:], in[:done], out)
	}
	remain := len(in) - done
	if remain > 0 {
		mt := mulTable[c][:256]
		for i := done; i < len(in); i++ {
			out[i] = mt[in[i]]
		}
	}
}

func galMulSliceXor(c byte, in, out []byte, o *options) {
	if c == 1 {
		sliceXor(in, out, o)
		return
	}
	done := (len(in) >> 4) << 4
	if done > 0 {
		galMulPpcXor(mulTableLow[c][:], mulTableHigh[c][:], in[:done], out)
	}
	remain := len(in) - done
	if remain > 0 {
		mt := mulTable[c][:256]
		for i := done; i < len(in); i++ {
			out[i] ^= mt[in[i]]
		}
	}
}

// slice galois add
func sliceXor(in, out []byte, o *options) {
	for n, input := range in {
		out[n] ^= input
	}
}

// 4-way butterfly
func ifftDIT4(work [][]byte, dist int, log_m01, log_m23, log_m02 ffe, o *options) {
	ifftDIT4Ref(work, dist, log_m01, log_m23, log_m02, o)
}

// 4-way butterfly
func fftDIT4(work [][]byte, dist int, log_m01, log_m23, log_m02 ffe, o *options) {
	fftDIT4Ref(work, dist, log_m01, log_m23, log_m02, o)
}

// 2-way butterfly forward
func fftDIT2(x, y []byte, log_m ffe, o *options) {
	// Reference version:
	refMulAdd(x, y, log_m)
	sliceXor(x, y, o)
}

// 2-way butterfly inverse
func ifftDIT2(x, y []byte, log_m ffe, o *options) {
	// Reference version:
	sliceXor(x, y, o)
	refMulAdd(x, y, log_m)
}

func mulgf16(x, y []byte, log_m ffe, o *options) {
	refMul(x, y, log_m)
}
