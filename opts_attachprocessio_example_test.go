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

package pipe_test

import (
	"fmt"

	pipe "github.com/ganbarodigital/go_pipe/v7"
)

func ExampleAttachOsStdin_asFunctionalOption() {
	// the pipe will now read from os.Stdin
	p := pipe.NewPipe(pipe.AttachOsStdin)

	// prove that the option did not error out
	statusCode, err := p.StatusError()
	fmt.Printf("statusCode is: %d\n", statusCode)
	fmt.Printf("err is: %v\n", err)
	// Output:
	// statusCode is: 0
	// err is: <nil>
}

func ExampleAttachOsStdin_asPipeCommand() {
	// create a new pipe
	p := pipe.NewPipe()

	// the pipe will now read from os.Stdin
	p.RunCommand(pipe.AttachOsStdin)
}

func ExampleAttachOsStdout_asFunctionalOption() {
	// the pipe will now read from os.Stdout
	p := pipe.NewPipe(pipe.AttachOsStdout)

	// prove that the option did not error out
	statusCode, err := p.StatusError()
	fmt.Printf("statusCode is: %d\n", statusCode)
	fmt.Printf("err is: %v\n", err)
	// Output:
	// statusCode is: 0
	// err is: <nil>
}

func ExampleAttachOsStdout_asPipeCommand() {
	// create a new pipe
	p := pipe.NewPipe()

	// the pipe will now read from os.Stdout
	p.RunCommand(pipe.AttachOsStdout)
}

func ExampleAttachOsStderr_asFunctionalOption() {
	// the pipe will now read from os.Stdin
	p := pipe.NewPipe(pipe.AttachOsStderr)

	// prove that the option did not error out
	statusCode, err := p.StatusError()
	fmt.Printf("statusCode is: %d\n", statusCode)
	fmt.Printf("err is: %v\n", err)
	// Output:
	// statusCode is: 0
	// err is: <nil>
}

func ExampleAttachOsStderr_asPipeCommand() {
	// create a new pipe
	p := pipe.NewPipe()

	// the pipe will now read from os.Stderr
	p.RunCommand(pipe.AttachOsStderr)
}
