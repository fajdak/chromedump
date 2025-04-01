package crypto

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	dllcrypt32  = syscall.NewLazyDLL("Crypt32.dll")
	dllkernel32 = syscall.NewLazyDLL("Kernel32.dll")

	procDecryptData = dllcrypt32.NewProc("CryptUnprotectData")
	procLocalFree   = dllkernel32.NewProc("LocalFree")
)

type dataBlob struct {
	cbData uint32
	pbData *byte
}

func newBlob(d []byte) *dataBlob {
	if len(d) == 0 {
		return &dataBlob{}
	}
	return &dataBlob{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (b *dataBlob) bytes() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}

// 将函数名首字母改为大写，使其可导出
func DecryptWithDPAPI(data []byte) ([]byte, error) {
	var outBlob dataBlob
	r, _, err := procDecryptData.Call(
		uintptr(unsafe.Pointer(newBlob(data))),
		0,
		0,
		0,
		0,
		0,
		uintptr(unsafe.Pointer(&outBlob)))

	if r == 0 {
		return nil, fmt.Errorf("DPAPI解密失败: %v", err)
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outBlob.pbData)))

	return outBlob.bytes(), nil
}
