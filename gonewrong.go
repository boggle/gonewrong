package gonewrong

/*
#include <errno.h>
#include <stdint.h>
int64_t get_errno() { return errno; }
*/
import "C"

import rt "runtime"


// ******* Thunks ************************************************************

// A simple argumentless function with no return value
type Thunk func()

// Wrap thunk in calls that lock the executing go routine to some OSThread
func (p Thunk) WithOSThread() Thunk {
    return Thunk(func() {
        rt.LockOSThread()
        defer rt.UnlockOSThread()

        p()
    })
}

// Helper for calling thunk within a separate go routine bound to a
// fixed OSThread
func (p Thunk) RunInOSThread() {
    go (p.WithOSThread())()
}

// Wrap thunk such that it sends msg on ch after finishing
// (May discard errors!)
func (p Thunk) SendOnFinish(ch chan interface{}, msg interface{}) Thunk {
    return Thunk(func() {
        defer func() { ch <- msg }()
        p()
    })
}


// ******* Error Handling ****************************************************

type ErrKnow interface {
    GetError() error

    OkIf(cond bool) error
    ErrorIf(cond bool) error
}

// If cond is true returns nil, error otherwise
func OkIf(cond bool, error error) error {
    return ErrorIf(!cond, error)
}

// If cond is true returns error, nil otherwise
func ErrorIf(cond bool, error error) error {
    if cond {
        return error
    }
    return nil
}

// Panics if err is != nil
func PanicUnlessNil(err error) {
    if err != nil {
        panic(err)
    }
}

// Returns true if ptr is C NULL (== 0)
// (Spec doesnt define go nil to be == NULL)
func IsCNullPtr(ptr uintptr) bool {
    return ptr == uintptr(0)
}

// Returns errno from C for the current thread
//
// Safe use may require locking this go routine to the underlying OSThread
func GetCErrno() int64 {
    return int64(C.get_errno())
}
