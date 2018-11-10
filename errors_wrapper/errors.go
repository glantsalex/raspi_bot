package errors_wrapper

import (
	"bytes"
	"fmt"
	"runtime"
)

type FacilityErrCode uint16

const (
	Success            FacilityErrCode = 0
	LoginErrInvalidUid FacilityErrCode = 1
	LoginErrInvalidAid FacilityErrCode = 2

	DbErrConnectionFailed FacilityErrCode = 99
	DbErr                 FacilityErrCode = 100
	DbErrNoRowsFetched    FacilityErrCode = 101
	DbErrNoRowsAffected   FacilityErrCode = 102
	DbErrDuplicateKey     FacilityErrCode = 103

	RedisErr                 FacilityErrCode = 200
	RedisErrConnectionFailed FacilityErrCode = 201

	AwsKenesisErr            FacilityErrCode = 300
	AwsS3NotInitialized      FacilityErrCode = 400
	AwsCreateSessionErr      FacilityErrCode = 401
	AwsS3Err                 FacilityErrCode = 410
	AwsS3NoNameProvidedErr   FacilityErrCode = 411
	AwsS3NoRegionProvidedErr FacilityErrCode = 412
	AwsS3CredentialsErr      FacilityErrCode = 413
	JsonErr                  FacilityErrCode = 500
	Base64ECodecErr          FacilityErrCode = 510
	ZlibErr                  FacilityErrCode = 520
	IoErr                    FacilityErrCode = 530
	UrlParseErr              FacilityErrCode = 540

	OsFileSystemErr         FacilityErrCode = 600
	UnknownStorageSpecified FacilityErrCode = 601

	SecurityTokenError   FacilityErrCode = 700
	NoSecurityTokenFound FacilityErrCode = 701
)

type IFacilityError interface {
	Error() string
	StackTrace() string
	ErrorCode() FacilityErrCode
}

type FacilityError struct {
	ErrCode      FacilityErrCode
	WrappedError error
	stackTrace   bytes.Buffer
}

func NewFacilityError(ec FacilityErrCode, we error) IFacilityError {
	fe := &FacilityError{ErrCode: ec, WrappedError: we}

	pc := make([]uintptr, 15)
	n := runtime.Callers(0, pc)

	frames := runtime.CallersFrames(pc[2:n])
	for {
		frame, more := frames.Next()
		fe.addToStackTrace(fmt.Sprintf("%s,:%d %s\n", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	return fe
}

func FacilityNilError() IFacilityError {
	return nil
}

func (fe *FacilityError) addToStackTrace(step string) {
	fe.stackTrace.WriteString(step)
}

func (fe *FacilityError) ErrorCode() FacilityErrCode {
	return fe.ErrCode
}

func (fe *FacilityError) StackTrace() string {
	return fe.stackTrace.String()
}

func (fe *FacilityError) Error() string {
	return fmt.Sprintf(" Code: %d, %s\nAt: %s", fe.ErrCode, fe.WrappedError, fe.StackTrace())
}
