package shared

import "encoding/hex"

// HexToBytes converts the given hex value to bytes array.
func HexToBytes(input string) []byte {
	return ignoreError(hex.DecodeString(input)).([]byte)
}

func ignoreError(val interface{}, err error) interface{} {
	return val
}
