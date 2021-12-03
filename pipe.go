// pipe is a library to help you write UNIX-like pipelines of operations
//
// inspired by:
//
// - http://labix.org/pipe
// - https://github.com/bitfield/script
//
// Copyright 2019-present Ganbaro Digital Ltd
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//   * Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.
//
//   * Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in
//     the documentation and/or other materials provided with the
//     distribution.
//
//   * Neither the names of the copyright holders nor the names of his
//     contributors may be used to endorse or promote products derived
//     from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
// FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
// COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
// ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package pipe

import (
	"io"

	ioextra "github.com/ganbarodigital/go-ioextra/v2"
	envish "github.com/ganbarodigital/go_envish/v3"
)

// Pipe is our data structure. All PipeCommands read from, and/or write to
// the pipe.
type Pipe struct {
	// PipeCommands read from Stdin
	Stdin ioextra.TextReader

	// PipeCommands write to Stdout and/or Stderr
	Stdout ioextra.TextReaderWriter
	Stderr ioextra.TextReaderWriter

	// Pipe users may need to temporarily replace Stdin, Stdout and/or Stderr
	// We provide a simple stack system to support that.
	stdinStack  []ioextra.TextReader
	stdoutStack []ioextra.TextReaderWriter
	stderrStack []ioextra.TextReaderWriter

	// PipeCommands return an error. We store it here.
	err error

	// PipeCommands return a UNIX-like status code. We store it here.
	statusCode int

	// PipeCommands can have their own environment, if they want one
	Env envish.Expander

	// You can pass bitmask flags into PipeCommands. Their meaning
	// is entirely yours to interpret.
	Flags int
}

// NewPipe creates a new Pipe that's ready to use.
//
// It starts with an empty Stdin, and uses the program's environment
// by default.
//
// You can provide a list of functional options for us to call. We'll
// pass in the Pipe, for you to reconfigure.
func NewPipe(options ...func(*Pipe)) *Pipe {
	// create a pipe that's ready to go
	retval := Pipe{
		Env: envish.NewProgramEnv(),
	}
	retval.ResetBuffers()
	retval.ResetError()

	// apply any option functions we might have been given
	for _, option := range options {
		option(&retval)
	}

	// all done
	return &retval
}

// DrainStdinToStdout will copy everything that's left in the pipe's Stdin
// over to the pipe's Stdout.
func (p *Pipe) DrainStdinToStdout() {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// do we have a Stdin to drain?
	if p.Stdin == nil {
		return
	}

	// do we have a Stdout to drain to?
	if p.Stdout == nil {
		p.SetNewStdout()
	}

	// yes we do
	io.Copy(p.Stdout, p.Stdin)
}

// Error returns the error returned from the last PipeCommand
// that ran against this pipe.
func (p *Pipe) Error() error {
	// do we have a pipe to work with?
	if p == nil {
		return nil
	}

	// yes we do
	return p.err
}

// Okay confirms that the last PipeCommand run against the pipe completed
// without reporting an error.
func (p *Pipe) Okay() bool {
	// do we have a pipe to inspect?
	if p == nil {
		return true
	}

	// yes we do
	return p.err == nil
}

// ResetBuffers creates new, empty Stdin, Stdout and Stderr for the given
// pipe.
//
// It also empties the internal stacks used by PushStdin / PopStdin,
// PushStdout / PopStdout, and PushStderr / PopStderr.
func (p *Pipe) ResetBuffers() {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// set our input/output buffers
	p.SetNewStdin()
	p.SetNewStdout()
	p.SetNewStderr()

	// reset our internal stacks
	p.stdinStack = make([]ioextra.TextReader, 0)
	p.stdoutStack = make([]ioextra.TextReaderWriter, 0)
	p.stderrStack = make([]ioextra.TextReaderWriter, 0)
}

// ResetError sets the pipe's status code and error to their zero values
// of (StatusOkay, nil).
func (p *Pipe) ResetError() {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// yes we do
	p.statusCode = StatusOkay
	p.err = nil
}

// RunCommand will run a function using this pipe. The function's return
// values are stored in the pipe's StatusCode and Err fields.
func (p *Pipe) RunCommand(c PipeCommand) {
	// do we have a pipe to work with?
	if p == nil || p.Stdin == nil || p.Stdout == nil {
		return
	}

	// yes we do
	p.statusCode, p.err = c(p)

	// special case - do we have a non-zero status code, but no error?
	if p.statusCode != StatusOkay && p.err == nil {
		p.err = ErrNonZeroStatusCode{"command", p.statusCode}
	}
}

// SetNewStdin creates a new, empty Stdin buffer on this pipe.
func (p *Pipe) SetNewStdin() {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// yes we do
	p.Stdin = ioextra.NewTextBuffer()

	// all done
}

// SetStdinFromString sets the pipe's Stdin to be the given input string.
func (p *Pipe) SetStdinFromString(input string) {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// yes we do
	buf := ioextra.NewTextBuffer()
	buf.WriteString(input)

	p.Stdin = buf

	// all done
}

// PushStdin adds the pipe's existing Stdin to an internal stack,
// and then sets the pipe's Stdin to the given newStdin.
//
// You can call PopStdin to reverse this operation.
//
// This is useful for callers who need to temporarily replace the pipe's
// Stdin.
func (p *Pipe) PushStdin(newStdin ioextra.TextReader) {
	// do we have a pipe to work with?
	if p == nil {
		// no, we do not
		return
	}

	p.stdinStack = append(p.stdinStack, p.Stdin)
	p.Stdin = newStdin
}

// PopStdin sets the pipe's Stdin to its previous value.
//
// It reverses your last call to PushStdin.
//
// This is useful for callers who need to temporarily replace the pipe's
// Stdin.
func (p *Pipe) PopStdin() {
	// do we have a pipe to work with?
	if p == nil {
		// no, we do not
		return
	}

	// do we have anything to restore?
	if len(p.stdinStack) == 0 {
		return
	}

	// restore Stdin
	p.Stdin = p.stdinStack[len(p.stdinStack)-1]

	// remove the value we've just popped from the stack
	p.stdinStack = p.stdinStack[:len(p.stdinStack)-1]
}

// StdinStackLen returns the number of entries in the internal stack of
// Stdin entries.
//
// You can call PushStdin and PopStdin to add entries to & from the
// internal stack.
func (p *Pipe) StdinStackLen() int {
	// do we have a pipe to work with?
	if p == nil {
		// no, we do not
		return 0
	}

	// yes we do
	return len(p.stdinStack)
}

// SetNewStdout creates a new, empty Stdout buffer on this pipe.
func (p *Pipe) SetNewStdout() {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// yes we do
	p.Stdout = ioextra.NewTextBuffer()

	// all done
}

// PushStdout adds the pipe's existing Stdout to an internal stack,
// and then sets the pipe's Stdout to the given newStdout.
//
// You can call PopStdout to reverse this operation.
//
// This is useful for callers who need to temporarily replace the pipe's
// Stdout (for example, to redirect to /dev/null).
//
// NOTE: if p.Stdout == p.Stderr, PushStdout sets *both* p.Stdout and
// p.Stderr to the newStdout.
func (p *Pipe) PushStdout(newStdout ioextra.TextReaderWriter) {
	// do we have a pipe to work with?
	if p == nil {
		// no, we do not
		return
	}

	// special case - does the pipe's Stdout currently point at
	// the pipe's Stdin?
	if p.Stdout == p.Stderr {
		p.Stderr = newStdout
	}

	// yes we do
	p.stdoutStack = append(p.stdoutStack, p.Stdout)
	p.Stdout = newStdout
}

// PopStdout sets the pipe's Stdout to its previous value.
//
// It reverses your last call to PushStdout.
//
// This is useful for callers who need to temporarily replace the pipe's
// Stdout (for example, to redirect to /dev/null).
//
// NOTE: if p.Stdout == p.Stderr, PopStdout sets *both* p.Stdout and
// p.Stderr to the previous Stdout. Most of the time, this is the desired
// intention.
//
// Use PopStdoutOnly when you don't want to touch the pipe's Stderr at all.
func (p *Pipe) PopStdout() {
	// do we have a pipe to work with?
	if p == nil {
		// no, we do not
		return
	}

	// do we have anything to restore?
	if len(p.stdoutStack) == 0 {
		return
	}

	// fetch the old stdout from our internal stack
	oldStdout := p.stdoutStack[len(p.stdoutStack)-1]

	// remove the value we've just popped from the stack
	p.stdoutStack = p.stdoutStack[:len(p.stdoutStack)-1]

	// do Stdout and Stderr point at each other?
	if p.Stdout == p.Stderr {
		// yes, so we need to restore to both
		p.Stderr = oldStdout
	}

	// restore Stdout
	p.Stdout = oldStdout
}

// PopStdoutOnly sets the pipe's Stdout to its previous value.
//
// It reverses your last call to PushStdout.
//
// This is useful for callers who need to temporarily replace the pipe's
// Stdout (for example, to redirect to /dev/null).
//
// NOTE: even if p.Stdout == p.Stderr, PopStdoutOnly leaves p.Stderr
// untouched.
func (p *Pipe) PopStdoutOnly() {
	// do we have a pipe to work with?
	if p == nil {
		// no, we do not
		return
	}

	// do we have anything to restore?
	if len(p.stdoutStack) == 0 {
		return
	}

	// restore Stdout
	p.Stdout = p.stdoutStack[len(p.stdoutStack)-1]

	// remove the value we've just popped from the stack
	p.stdoutStack = p.stdoutStack[:len(p.stdoutStack)-1]
}

// StdoutStackLen returns the number of entries in the internal stack of
// Stdout entries.
//
// You can call PushStdout and PopStdout to add entries to & from the
// internal stack.
func (p *Pipe) StdoutStackLen() int {
	// do we have a pipe to work with?
	if p == nil {
		// no, we do not
		return 0
	}

	// yes we do
	return len(p.stdoutStack)
}

// SetNewStderr creates a new, empty Stderr buffer on this pipe.
func (p *Pipe) SetNewStderr() {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// yes we do
	p.Stderr = ioextra.NewTextBuffer()

	// all done
}

// PushStderr adds the pipe's existing Stderr to an internal stack,
// and then sets the pipe's Stderr to the given newStderr.
//
// You can call PopStderr to reverse this operation.
//
// This is useful for callers who need to temporarily replace the pipe's
// Stderr (for example, to redirect to /dev/null).
//
// NOTE: if p.Stdout == p.Stderr, PushStderr sets *both* p.Stdout and
// p.Stderr to the newStderr.
func (p *Pipe) PushStderr(newStderr ioextra.TextReaderWriter) {
	// do we have a pipe to work with?
	if p == nil {
		// no, we do not
		return
	}

	// yes we do

	// special case - does the pipe's Stdout current point at the pipe's
	// Stderr?
	if p.Stdout == p.Stderr {
		p.Stdout = newStderr
	}

	p.stderrStack = append(p.stderrStack, p.Stderr)
	p.Stderr = newStderr
}

// PopStderr sets the pipe's Stderr to its previous value.
//
// It reverses your last call to PushStderr.
//
// This is useful for callers who need to temporarily replace the pipe's
// Stderr (for example, to redirect to /dev/null).
//
// NOTE: if p.Stdout == p.Stderr, PopStderr sets *both* p.Stdout and
// p.Stderr to the previous Stderr. Most of the time, this is the desired
// intention.
//
// Use PopStderrOnly when you don't want to touch the pipe's Stdout at all.
func (p *Pipe) PopStderr() {
	// do we have a pipe to work with?
	if p == nil {
		// no, we do not
		return
	}

	// do we have anything to restore?
	if len(p.stderrStack) == 0 {
		return
	}

	// restore Stderr
	oldStderr := p.stderrStack[len(p.stderrStack)-1]

	// remove the value we've just popped from the stack
	p.stderrStack = p.stderrStack[:len(p.stderrStack)-1]

	// special case - does the pipe's Stdout currently point at
	// the pipe's Stderr?
	if p.Stdout == p.Stderr {
		p.Stdout = oldStderr
	}

	// restore Stderr
	p.Stderr = oldStderr
}

// PopStderrOnly sets the pipe's Stderr to its previous value.
//
// It reverses your last call to PushStderr.
//
// This is useful for callers who need to temporarily replace the pipe's
// Stderr (for example, to redirect to /dev/null).
//
// NOTE: even if p.Stdout == p.Stderr, PopStderrOnly leaves p.Stdout
// untouched.
func (p *Pipe) PopStderrOnly() {
	// do we have a pipe to work with?
	if p == nil {
		// no, we do not
		return
	}

	// do we have anything to restore?
	if len(p.stderrStack) == 0 {
		return
	}

	// restore Stderr
	p.Stderr = p.stderrStack[len(p.stderrStack)-1]

	// remove the value we've just popped from the stack
	p.stderrStack = p.stderrStack[:len(p.stderrStack)-1]
}

// StderrStackLen returns the number of entries in the internal stack of
// Stderr entries.
//
// You can call PushStderr and PopStderr to add entries to & from the
// internal stack.
func (p *Pipe) StderrStackLen() int {
	// do we have a pipe to work with?
	if p == nil {
		// no, we do not
		return 0
	}

	// yes we do
	return len(p.stderrStack)
}

// StatusCode returns the UNIX-like status code from the last PipeCommand
// that ran against this pipe.
func (p *Pipe) StatusCode() int {
	// do we have a pipe to work with?
	if p == nil {
		return StatusOkay
	}

	// yes we do
	return p.statusCode
}

// StatusError is a shorthand for calling p.StatusCode() and p.Error()
// to get the UNIX-like status code and the last reported Golang error.
func (p *Pipe) StatusError() (int, error) {
	// do we have a pipe to inspect?
	if p == nil {
		return StatusOkay, nil
	}

	// yes we do
	return p.statusCode, p.err
}

// SetStatusCode is a helper method, added to help us test this
// package.
//
// It is not part of our supported API. Use at your own risk!
func SetStatusCode(p *Pipe, newStatusCode int) {
	p.statusCode = newStatusCode
}
