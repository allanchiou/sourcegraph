// Code generated by github.com/efritz/go-mockgen 0.1.0; DO NOT EDIT.

package mocks

import (
	"context"
	commits "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/commits"
	"sync"
)

// MockUpdater is a mock implementation of the Updater interface (from the
// package
// github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/commits)
// used for unit testing.
type MockUpdater struct {
	// UpdateFunc is an instance of a mock function object controlling the
	// behavior of the method Update.
	UpdateFunc *UpdaterUpdateFunc
}

// NewMockUpdater creates a new mock of the Updater interface. All methods
// return zero values for all results, unless overwritten.
func NewMockUpdater() *MockUpdater {
	return &MockUpdater{
		UpdateFunc: &UpdaterUpdateFunc{
			defaultHook: func(context.Context, int, bool) error {
				return nil
			},
		},
	}
}

// NewMockUpdaterFrom creates a new mock of the MockUpdater interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockUpdaterFrom(i commits.Updater) *MockUpdater {
	return &MockUpdater{
		UpdateFunc: &UpdaterUpdateFunc{
			defaultHook: i.Update,
		},
	}
}

// UpdaterUpdateFunc describes the behavior when the Update method of the
// parent MockUpdater instance is invoked.
type UpdaterUpdateFunc struct {
	defaultHook func(context.Context, int, bool) error
	hooks       []func(context.Context, int, bool) error
	history     []UpdaterUpdateFuncCall
	mutex       sync.Mutex
}

// Update delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockUpdater) Update(v0 context.Context, v1 int, v2 bool) error {
	r0 := m.UpdateFunc.nextHook()(v0, v1, v2)
	m.UpdateFunc.appendCall(UpdaterUpdateFuncCall{v0, v1, v2, r0})
	return r0
}

// SetDefaultHook sets function that is called when the Update method of the
// parent MockUpdater instance is invoked and the hook queue is empty.
func (f *UpdaterUpdateFunc) SetDefaultHook(hook func(context.Context, int, bool) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Update method of the parent MockUpdater instance inovkes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *UpdaterUpdateFunc) PushHook(hook func(context.Context, int, bool) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *UpdaterUpdateFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(context.Context, int, bool) error {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *UpdaterUpdateFunc) PushReturn(r0 error) {
	f.PushHook(func(context.Context, int, bool) error {
		return r0
	})
}

func (f *UpdaterUpdateFunc) nextHook() func(context.Context, int, bool) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *UpdaterUpdateFunc) appendCall(r0 UpdaterUpdateFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of UpdaterUpdateFuncCall objects describing
// the invocations of this function.
func (f *UpdaterUpdateFunc) History() []UpdaterUpdateFuncCall {
	f.mutex.Lock()
	history := make([]UpdaterUpdateFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// UpdaterUpdateFuncCall is an object that describes an invocation of method
// Update on an instance of MockUpdater.
type UpdaterUpdateFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int
	// Arg2 is the value of the 3rd argument passed to this method
	// invocation.
	Arg2 bool
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c UpdaterUpdateFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1, c.Arg2}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c UpdaterUpdateFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}