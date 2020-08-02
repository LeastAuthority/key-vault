package shared

import "encoding/hex"

func ignoreError(val interface{}, err error) interface{} {
	return val
}

func HexToBytes(input string) []byte {
	return ignoreError(hex.DecodeString(input)).([]byte)
}
