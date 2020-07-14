//
//  Copyright (c) 2020 Sugesh Chandran
//

package errors

import (
	"fmt"
)

var (
	//ErrNil : No error
	ErrNil = fmt.Errorf("No error in Agent")
	// ErrInvalidInput : Invalid input error
	ErrInvalidInput  = fmt.Errorf("Invalid Input, Cannot process in App")
	ErrInputTooShort = fmt.Errorf("Invalid input, input length is too short")
	// ErrInvalidOp : Invalid operation request.
	ErrInvalidOp = fmt.Errorf("Invalid operation request, Cannot process in App")
	// ErrInvalidState : System is at invalid state
	ErrInvalidState = fmt.Errorf("Invalid state of application")
	//ErrWriteFail : Failed to write to a stream(can be file, stdout, channel etc)
	ErrWriteFail = fmt.Errorf("Failed to write to the stream")
	//ErrOpTimeout : Failed to complete the operation in specific time
	ErrOpTimeout = fmt.Errorf("Operation timeout")
	//ErrInvalidAddr : Invalid protocol address
	ErrInvalidAddr = fmt.Errorf("Invalid protocol address")
	//ErrSkipOp : Skipping the request for specific operation
	ErrSkipOp = fmt.Errorf("Skipping the request for operation")
	//ErrRetry
	ErrRetry = fmt.Errorf("Operation failed, try again")
	//ErrNotExists : Error for file/object not exists in the system
	ErrNotExists = fmt.Errorf("Not exists in the system")
)