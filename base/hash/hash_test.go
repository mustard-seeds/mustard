package hash

import (
	"testing"
)
func TestFingerPrint(t *testing.T) {
	_uint64 := FingerPrint("xxxxx")
	if _uint64 != 6019362431426986063 {
		t.Error("FingerPrint fail.")
	}
	_uint64str := FingerprintToString(_uint64)
	if _uint64str != "5389123d4b89804f" {
		t.Error("FingerPrintToString fail")
	}
	_uint64strint,_ := StringToFingerprint(_uint64str)
	if _uint64strint != _uint64 {
		t.Error("FingerPrintFrom String fail")
	}
}
func TestFingerPrint32(t *testing.T) {
	_uint32 := FingerPrint32("xxxxx")
	if _uint32 != 2325599039 {
		t.Error("FingerPrint32 fail")
	}
	_uint32str  := FingerprintToString(uint64(_uint32))
	if _uint32str != "8a9dd33f" {
		t.Error("FingerPrinttostring fail")
	}
	_uint32strint,_ := StringToFingerprint(_uint32str)
	if uint32(_uint32strint) != _uint32 {
		t.Error("String To FingerPrint fail")
	}
}