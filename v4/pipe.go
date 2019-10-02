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
	"os"

	envish "github.com/ganbarodigital/go_envish"
)

// Pipe is our data structure. All user-land functionality either reads from,
// and/or writes to the pipe.
type Pipe struct {
	// Pipe commands read from Stdin
	Stdin *Source

	// Pipe commands write to Stdout and/or Stderr
	Stdout *Dest
	Stderr *Dest

	// Pipe commands return an error. We store it here.
	Err error

	// Pipe commands return a UNIX-like status code. We store it here.
	StatusCode int

	// Pipe commands can have their own environment, if they want one
	Env *envish.Env
}

// NewPipe creates a new, empty Pipe.
//
// It starts with an empty Stdin.
func NewPipe(options ...func(*Pipe)) *Pipe {
	retval := Pipe{
		Stdin:      NewSourceFromString(""),
		Stdout:     new(Dest),
		Stderr:     new(Dest),
		Err:        nil,
		StatusCode: StatusOkay,
	}

	// apply any option functions we might have been given
	for _, option := range options {
		option(&retval)
	}

	// all done
	return &retval
}

// DrainStdin will copy everything that's left in the pipe's stdin
// over to the pipe's stdout
func (p *Pipe) DrainStdin() {
	// do we have a pipe to work with?
	if p == nil || p.Stdin == nil || p.Stdout == nil {
		return
	}

	// yes we do
	io.Copy(p.Stdout, p.Stdin)
}

// Error returns any error stored in the Pipe
func (p *Pipe) Error() error {
	// do we have a pipe to work with?
	if p == nil {
		return nil
	}

	// yes we do
	return p.Err
}

// Expand replaces ${var} or $var in the input string.
//
// It uses the Pipe's private environment (if the sequence has one),
// or the program's environment otherwise.
func (p *Pipe) Expand(fmt string) string {
	// do we have a pipe to work with?
	if p == nil {
		return os.Expand(fmt, os.Getenv)
	}

	// do we have an environment of our own?
	if p.Env == nil {
		return os.Expand(fmt, os.Getenv)
	}

	// yes we do
	return p.Env.Expand(fmt)
}

// Next prepares the pipe to be used by the next Command.
//
// NOTE that we DO NOT reset the StatusCode or Err here. Their value may
// be of interest to the next Command (which is why they were moved here
// in v4!)
func (p *Pipe) Next() {
	// do we have a pipe to work with?
	if p == nil || p.Stdin == nil || p.Stdout == nil || p.Stderr == nil {
		return
	}

	p.Stdin = p.Stdout.NewSource()
	p.Stdout = new(Dest)
	p.Stderr = new(Dest)
}

// Reset creates new, empty Stdin, Stdout and Stderr.
//
// It's useful for pipelines that consist of multiple lists.
//
// NOTE that we DO NOT reset the StatusCode or Err here. Their value may
// be of interest to the next Command (which is why they were moved here
// in v4!)
func (p *Pipe) Reset() {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// reset most of the things
	p.Stdin = NewSourceFromString("")
	p.Stdout = new(Dest)
	p.Stderr = new(Dest)
}

// RunCommand will run a function using this pipe. The function's return
// values are stored in the pipe's StatusCode and Err fields.
func (p *Pipe) RunCommand(c Command) {
	// do we have a pipe to work with?
	if p == nil || p.Stdin == nil || p.Stdout == nil {
		return
	}

	// yes we do
	p.StatusCode, p.Err = c(p)

	// special case - do we have a non-zero status code, but no error?
	if p.StatusCode != StatusOkay && p.Err == nil {
		p.Err = ErrNonZeroStatusCode{"command", p.StatusCode}
	}
}
