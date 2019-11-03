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

	envish "github.com/ganbarodigital/go_envish/v3"
)

// Pipe is our data structure. All Commands read from, and/or write to
// the pipe.
type Pipe struct {
	// Pipe commands read from Stdin
	Stdin *Source

	// Pipe commands write to Stdout and/or Stderr
	Stdout *Dest
	Stderr *Dest

	// Pipe commands return an error. We store it here.
	err error

	// Pipe commands return a UNIX-like status code. We store it here.
	statusCode int

	// Pipe commands can have their own environment, if they want one
	Env envish.Expander

	// You can pass bitmask flags into pipe commands. Their meaning
	// is entirely yours to interpret.
	Flags int
}

// NewPipe creates a new Pipe that's ready to use.
//
// It starts with an empty Stdin, and uses the program's environment
// by default.
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
// over to the pipe's Stdout
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

// Error returns the error returned from the last Command
// that ran against this pipe
func (p *Pipe) Error() error {
	// do we have a pipe to work with?
	if p == nil {
		return nil
	}

	// yes we do
	return p.err
}

// Okay confirms that the last Command run against the pipe completed
// without reporting an error
func (p *Pipe) Okay() bool {
	// do we have a pipe to inspect?
	if p == nil {
		return true
	}

	// yes we do
	return p.err == nil
}

// ResetBuffers creates new, empty buffers for the pipe
func (p *Pipe) ResetBuffers() {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// set our input/output buffers
	p.SetNewStdin()
	p.SetNewStdout()
	p.SetNewStderr()
}

// ResetError sets the pipe's status code and error to their zero values
// of (StatusOkay, nil)
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
func (p *Pipe) RunCommand(c Command) {
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

// SetNewStdin creates a new, empty Stdin buffer on this pipe
func (p *Pipe) SetNewStdin() {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// yes we do
	p.Stdin = NewSourceFromString("")

	// all done
}

// SetStdinFromString sets the pipe's Stdin to be the given input string
func (p *Pipe) SetStdinFromString(input string) {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// yes we do
	p.Stdin = NewSourceFromString(input)

	// all done
}

// SetNewStdout creates a new, empty Stdout buffer on this pipe
func (p *Pipe) SetNewStdout() {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// yes we do
	p.Stdout = new(Dest)

	// all done
}

// SetNewStderr creates a new, empty Stderr buffer on this pipe
func (p *Pipe) SetNewStderr() {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// yes we do
	p.Stderr = new(Dest)

	// all done
}

// StatusCode returns the UNIX-like status code from the last Command
// that ran against this pipe
func (p *Pipe) StatusCode() int {
	// do we have a pipe to work with?
	if p == nil {
		return StatusOkay
	}

	// yes we do
	return p.statusCode
}

// StatusError is a shorthand for calling p.StatusCode() and p.Error()
// to get the UNIX-like status code and the last reported Golang error
func (p *Pipe) StatusError() (int, error) {
	// do we have a pipe to inspect?
	if p == nil {
		return StatusOkay, nil
	}

	// yes we do
	return p.statusCode, p.err
}
