package socket

import "net"

const (
	Reason_None                   = iota
	Reason_Timeout
	Reason_PackageTagNotMatch
	Reason_PackageDataSizeInvalid
	Reason_PackageTooBig
)

func errToReason(err error) int {

	switch n := err.(type) {
	case net.Error:
		if n.Timeout() {
			return Reason_Timeout
		}
	}

	switch err {
	case ErrPackageTagNotMatch:
		return Reason_PackageTagNotMatch
	case ErrPackageDataSizeInvalid:
		return Reason_PackageDataSizeInvalid
	case ErrPackageTooBig:
		return Reason_PackageTooBig
	}

	return Reason_None

}
