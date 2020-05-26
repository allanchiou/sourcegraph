// Code generated by github.com/efritz/go-mockgen 0.1.0; DO NOT EDIT.

package mocks

import (
	"context"
	"io"
	"sync"

	db "github.com/sourcegraph/sourcegraph/internal/codeintel/db"
	gitserver "github.com/sourcegraph/sourcegraph/internal/codeintel/gitserver"
)

// MockClient is a mock impelementation of the Client interface (from the
// package github.com/sourcegraph/sourcegraph/internal/codeintel/gitserver)
// used for unit testing.
type MockClient struct {
	// ArchiveFunc is an instance of a mock function object controlling the
	// behavior of the method Archive.
	ArchiveFunc *ClientArchiveFunc
	// CommitsNearFunc is an instance of a mock function object controlling
	// the behavior of the method CommitsNear.
	CommitsNearFunc *ClientCommitsNearFunc
	// DirectoryChildrenFunc is an instance of a mock function object
	// controlling the behavior of the method DirectoryChildren.
	DirectoryChildrenFunc *ClientDirectoryChildrenFunc
	// FileExistsFunc is an instance of a mock function object controlling
	// the behavior of the method FileExists.
	FileExistsFunc *ClientFileExistsFunc
	// HeadFunc is an instance of a mock function object controlling the
	// behavior of the method Head.
	HeadFunc *ClientHeadFunc
}

// NewMockClient creates a new mock of the Client interface. All methods
// return zero values for all results, unless overwritten.
func NewMockClient() *MockClient {
	return &MockClient{
		ArchiveFunc: &ClientArchiveFunc{
			defaultHook: func(context.Context, db.DB, int, string) (io.Reader, error) {
				return nil, nil
			},
		},
		CommitsNearFunc: &ClientCommitsNearFunc{
			defaultHook: func(context.Context, db.DB, int, string) (map[string][]string, error) {
				return nil, nil
			},
		},
		DirectoryChildrenFunc: &ClientDirectoryChildrenFunc{
			defaultHook: func(context.Context, db.DB, int, string, []string) (map[string][]string, error) {
				return nil, nil
			},
		},
		FileExistsFunc: &ClientFileExistsFunc{
			defaultHook: func(context.Context, db.DB, int, string, string) (bool, error) {
				return false, nil
			},
		},
		HeadFunc: &ClientHeadFunc{
			defaultHook: func(context.Context, db.DB, int) (string, error) {
				return "", nil
			},
		},
	}
}

// NewMockClientFrom creates a new mock of the MockClient interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockClientFrom(i gitserver.Client) *MockClient {
	return &MockClient{
		ArchiveFunc: &ClientArchiveFunc{
			defaultHook: i.Archive,
		},
		CommitsNearFunc: &ClientCommitsNearFunc{
			defaultHook: i.CommitsNear,
		},
		DirectoryChildrenFunc: &ClientDirectoryChildrenFunc{
			defaultHook: i.DirectoryChildren,
		},
		FileExistsFunc: &ClientFileExistsFunc{
			defaultHook: i.FileExists,
		},
		HeadFunc: &ClientHeadFunc{
			defaultHook: i.Head,
		},
	}
}

// ClientArchiveFunc describes the behavior when the Archive method of the
// parent MockClient instance is invoked.
type ClientArchiveFunc struct {
	defaultHook func(context.Context, db.DB, int, string) (io.Reader, error)
	hooks       []func(context.Context, db.DB, int, string) (io.Reader, error)
	history     []ClientArchiveFuncCall
	mutex       sync.Mutex
}

// Archive delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockClient) Archive(v0 context.Context, v1 db.DB, v2 int, v3 string) (io.Reader, error) {
	r0, r1 := m.ArchiveFunc.nextHook()(v0, v1, v2, v3)
	m.ArchiveFunc.appendCall(ClientArchiveFuncCall{v0, v1, v2, v3, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Archive method of
// the parent MockClient instance is invoked and the hook queue is empty.
func (f *ClientArchiveFunc) SetDefaultHook(hook func(context.Context, db.DB, int, string) (io.Reader, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Archive method of the parent MockClient instance inovkes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *ClientArchiveFunc) PushHook(hook func(context.Context, db.DB, int, string) (io.Reader, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ClientArchiveFunc) SetDefaultReturn(r0 io.Reader, r1 error) {
	f.SetDefaultHook(func(context.Context, db.DB, int, string) (io.Reader, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ClientArchiveFunc) PushReturn(r0 io.Reader, r1 error) {
	f.PushHook(func(context.Context, db.DB, int, string) (io.Reader, error) {
		return r0, r1
	})
}

func (f *ClientArchiveFunc) nextHook() func(context.Context, db.DB, int, string) (io.Reader, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ClientArchiveFunc) appendCall(r0 ClientArchiveFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ClientArchiveFuncCall objects describing
// the invocations of this function.
func (f *ClientArchiveFunc) History() []ClientArchiveFuncCall {
	f.mutex.Lock()
	history := make([]ClientArchiveFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ClientArchiveFuncCall is an object that describes an invocation of method
// Archive on an instance of MockClient.
type ClientArchiveFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 db.DB
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 int
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 io.Reader
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ClientArchiveFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ClientArchiveFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ClientCommitsNearFunc describes the behavior when the CommitsNear method
// of the parent MockClient instance is invoked.
type ClientCommitsNearFunc struct {
	defaultHook func(context.Context, db.DB, int, string) (map[string][]string, error)
	hooks       []func(context.Context, db.DB, int, string) (map[string][]string, error)
	history     []ClientCommitsNearFuncCall
	mutex       sync.Mutex
}

// CommitsNear delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockClient) CommitsNear(v0 context.Context, v1 db.DB, v2 int, v3 string) (map[string][]string, error) {
	r0, r1 := m.CommitsNearFunc.nextHook()(v0, v1, v2, v3)
	m.CommitsNearFunc.appendCall(ClientCommitsNearFuncCall{v0, v1, v2, v3, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the CommitsNear method
// of the parent MockClient instance is invoked and the hook queue is empty.
func (f *ClientCommitsNearFunc) SetDefaultHook(hook func(context.Context, db.DB, int, string) (map[string][]string, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// CommitsNear method of the parent MockClient instance inovkes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ClientCommitsNearFunc) PushHook(hook func(context.Context, db.DB, int, string) (map[string][]string, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ClientCommitsNearFunc) SetDefaultReturn(r0 map[string][]string, r1 error) {
	f.SetDefaultHook(func(context.Context, db.DB, int, string) (map[string][]string, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ClientCommitsNearFunc) PushReturn(r0 map[string][]string, r1 error) {
	f.PushHook(func(context.Context, db.DB, int, string) (map[string][]string, error) {
		return r0, r1
	})
}

func (f *ClientCommitsNearFunc) nextHook() func(context.Context, db.DB, int, string) (map[string][]string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ClientCommitsNearFunc) appendCall(r0 ClientCommitsNearFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ClientCommitsNearFuncCall objects
// describing the invocations of this function.
func (f *ClientCommitsNearFunc) History() []ClientCommitsNearFuncCall {
	f.mutex.Lock()
	history := make([]ClientCommitsNearFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ClientCommitsNearFuncCall is an object that describes an invocation of
// method CommitsNear on an instance of MockClient.
type ClientCommitsNearFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 db.DB
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 int
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 map[string][]string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ClientCommitsNearFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ClientCommitsNearFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ClientDirectoryChildrenFunc describes the behavior when the
// DirectoryChildren method of the parent MockClient instance is invoked.
type ClientDirectoryChildrenFunc struct {
	defaultHook func(context.Context, db.DB, int, string, []string) (map[string][]string, error)
	hooks       []func(context.Context, db.DB, int, string, []string) (map[string][]string, error)
	history     []ClientDirectoryChildrenFuncCall
	mutex       sync.Mutex
}

// DirectoryChildren delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockClient) DirectoryChildren(v0 context.Context, v1 db.DB, v2 int, v3 string, v4 []string) (map[string][]string, error) {
	r0, r1 := m.DirectoryChildrenFunc.nextHook()(v0, v1, v2, v3, v4)
	m.DirectoryChildrenFunc.appendCall(ClientDirectoryChildrenFuncCall{v0, v1, v2, v3, v4, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the DirectoryChildren
// method of the parent MockClient instance is invoked and the hook queue is
// empty.
func (f *ClientDirectoryChildrenFunc) SetDefaultHook(hook func(context.Context, db.DB, int, string, []string) (map[string][]string, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// DirectoryChildren method of the parent MockClient instance inovkes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *ClientDirectoryChildrenFunc) PushHook(hook func(context.Context, db.DB, int, string, []string) (map[string][]string, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ClientDirectoryChildrenFunc) SetDefaultReturn(r0 map[string][]string, r1 error) {
	f.SetDefaultHook(func(context.Context, db.DB, int, string, []string) (map[string][]string, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ClientDirectoryChildrenFunc) PushReturn(r0 map[string][]string, r1 error) {
	f.PushHook(func(context.Context, db.DB, int, string, []string) (map[string][]string, error) {
		return r0, r1
	})
}

func (f *ClientDirectoryChildrenFunc) nextHook() func(context.Context, db.DB, int, string, []string) (map[string][]string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ClientDirectoryChildrenFunc) appendCall(r0 ClientDirectoryChildrenFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ClientDirectoryChildrenFuncCall objects
// describing the invocations of this function.
func (f *ClientDirectoryChildrenFunc) History() []ClientDirectoryChildrenFuncCall {
	f.mutex.Lock()
	history := make([]ClientDirectoryChildrenFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ClientDirectoryChildrenFuncCall is an object that describes an invocation
// of method DirectoryChildren on an instance of MockClient.
type ClientDirectoryChildrenFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 db.DB
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 int
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 string
	// Arg4 is the value of the 5th argument passed to this method
	// invocation.
	Arg4 []string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 map[string][]string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ClientDirectoryChildrenFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3, c.Arg4}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ClientDirectoryChildrenFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ClientFileExistsFunc describes the behavior when the FileExists method of
// the parent MockClient instance is invoked.
type ClientFileExistsFunc struct {
	defaultHook func(context.Context, db.DB, int, string, string) (bool, error)
	hooks       []func(context.Context, db.DB, int, string, string) (bool, error)
	history     []ClientFileExistsFuncCall
	mutex       sync.Mutex
}

// FileExists delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockClient) FileExists(v0 context.Context, v1 db.DB, v2 int, v3 string, v4 string) (bool, error) {
	r0, r1 := m.FileExistsFunc.nextHook()(v0, v1, v2, v3, v4)
	m.FileExistsFunc.appendCall(ClientFileExistsFuncCall{v0, v1, v2, v3, v4, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the FileExists method of
// the parent MockClient instance is invoked and the hook queue is empty.
func (f *ClientFileExistsFunc) SetDefaultHook(hook func(context.Context, db.DB, int, string, string) (bool, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// FileExists method of the parent MockClient instance inovkes the hook at
// the front of the queue and discards it. After the queue is empty, the
// default hook function is invoked for any future action.
func (f *ClientFileExistsFunc) PushHook(hook func(context.Context, db.DB, int, string, string) (bool, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ClientFileExistsFunc) SetDefaultReturn(r0 bool, r1 error) {
	f.SetDefaultHook(func(context.Context, db.DB, int, string, string) (bool, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ClientFileExistsFunc) PushReturn(r0 bool, r1 error) {
	f.PushHook(func(context.Context, db.DB, int, string, string) (bool, error) {
		return r0, r1
	})
}

func (f *ClientFileExistsFunc) nextHook() func(context.Context, db.DB, int, string, string) (bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ClientFileExistsFunc) appendCall(r0 ClientFileExistsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ClientFileExistsFuncCall objects describing
// the invocations of this function.
func (f *ClientFileExistsFunc) History() []ClientFileExistsFuncCall {
	f.mutex.Lock()
	history := make([]ClientFileExistsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ClientFileExistsFuncCall is an object that describes an invocation of
// method FileExists on an instance of MockClient.
type ClientFileExistsFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 db.DB
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 int
	// Arg3 is the value of the 4th argument passed to this method
	// invocation.
	Arg3 string
	// Arg4 is the value of the 5th argument passed to this method
	// invocation.
	Arg4 string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 bool
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ClientFileExistsFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2, c.Arg3, c.Arg4}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ClientFileExistsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}

// ClientHeadFunc describes the behavior when the Head method of the parent
// MockClient instance is invoked.
type ClientHeadFunc struct {
	defaultHook func(context.Context, db.DB, int) (string, error)
	hooks       []func(context.Context, db.DB, int) (string, error)
	history     []ClientHeadFuncCall
	mutex       sync.Mutex
}

// Head delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockClient) Head(v0 context.Context, v1 db.DB, v2 int) (string, error) {
	r0, r1 := m.HeadFunc.nextHook()(v0, v1, v2)
	m.HeadFunc.appendCall(ClientHeadFuncCall{v0, v1, v2, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the Head method of the
// parent MockClient instance is invoked and the hook queue is empty.
func (f *ClientHeadFunc) SetDefaultHook(hook func(context.Context, db.DB, int) (string, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Head method of the parent MockClient instance inovkes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *ClientHeadFunc) PushHook(hook func(context.Context, db.DB, int) (string, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *ClientHeadFunc) SetDefaultReturn(r0 string, r1 error) {
	f.SetDefaultHook(func(context.Context, db.DB, int) (string, error) {
		return r0, r1
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *ClientHeadFunc) PushReturn(r0 string, r1 error) {
	f.PushHook(func(context.Context, db.DB, int) (string, error) {
		return r0, r1
	})
}

func (f *ClientHeadFunc) nextHook() func(context.Context, db.DB, int) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *ClientHeadFunc) appendCall(r0 ClientHeadFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of ClientHeadFuncCall objects describing the
// invocations of this function.
func (f *ClientHeadFunc) History() []ClientHeadFuncCall {
	f.mutex.Lock()
	history := make([]ClientHeadFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// ClientHeadFuncCall is an object that describes an invocation of method
// Head on an instance of MockClient.
type ClientHeadFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 db.DB
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 int
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c ClientHeadFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c ClientHeadFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}
