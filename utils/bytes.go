package utils

func SliceTo32Bytes(in []byte) [32]byte {
	var res [32]byte
	copy(res[:], in)
	return res
}
