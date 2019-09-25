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

import "io"

// Pipe is our data structure. All user-land functionality either reads from,
// and/or writes to the pipe.
type Pipe struct {
	// Pipe operations read from Stdin
	Stdin *Source

	// Pipe operations write to Stdout and/or Stderr
	Stdout *Dest
	Stderr *Dest
}

// NewPipe creates a new, empty Pipe.
//
// It starts with an empty Stdin.
func NewPipe() *Pipe {
	return &Pipe{
		Stdin:  NewSourceFromString(""),
		Stdout: new(Dest),
		Stderr: new(Dest),
	}
}

// Next prepares the pipe to be used by the next PipeOperation
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
func (p *Pipe) Reset() {
	// do we have a pipe to work with?
	if p == nil {
		return
	}

	// reset all the things
	p.Stdin = NewSourceFromString("")
	p.Stdout = new(Dest)
	p.Stderr = new(Dest)
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
