// pipe is a library to help you write UNIX-like pipelines of operations
//
// inspired by:
//
// - http://labix.org/pipe
// - https://github.com/bitfield/script
//
// Copyright 2021-present Ganbaro Digital Ltd
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
	"os"

	ioextra "github.com/ganbarodigital/go-ioextra/v2"
)

// AttachOsStdin sets the pipe to read from your program's Stdin.
//
// You can use this both as a functional option, and/or as a
// PipeCommand.
func AttachOsStdin(p *Pipe) (int, error) {
	p.Stdin = ioextra.NewTextFile(os.Stdin)
	return StatusOkay, nil
}

// AttachOsStdout sets the pipe to write to your program's Stdout.
//
// You can use this both as a functional option, and/or as a
// PipeCommand.
func AttachOsStdout(p *Pipe) (int, error) {
	p.Stdout = ioextra.NewTextFile(os.Stdout)
	return StatusOkay, nil
}

// AttachOsStderr sets the pipe to write to your program's Stderr.
//
// You can use this both as a functional option, and/or as a
// PipeCommand.
func AttachOsStderr(p *Pipe) (int, error) {
	p.Stderr = ioextra.NewTextFile(os.Stderr)
	return StatusOkay, nil
}
