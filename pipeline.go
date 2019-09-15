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

import "io/ioutil"

// Pipeline is a set of pipe operations to be executed
type Pipeline struct {
	Pipe *Pipe
	// keep track of the steps that belong to this pipeline
	Steps []PipelineOperation

	// If anything goes wrong, we track the error here
	Err error

	// The UNIX-like status code from the last executed step
	StatusCode int
}

// NewPipeline creates a pipeline that's ready to run
func NewPipeline(steps ...PipelineOperation) *Pipeline {
	pipe := NewPipe()
	pipeline := Pipeline{
		pipe,
		steps,
		nil,
		0,
	}

	return &pipeline
}

// Bytes returns the contents of the pipeline's stdout as a byte slice
func (pl *Pipeline) Bytes() ([]byte, error) {
	// do we have a pipeline?
	if pl == nil {
		return []byte{}, nil
	}

	// was the pipeline initialised correctly?
	if pl.Pipe == nil {
		return []byte{}, pl.Err
	}

	// did an error occur?
	if pl.Err != nil {
		retval, _ := ioutil.ReadAll(pl.Pipe.Stderr.NewReader())
		return retval, pl.Err
	}

	// if we get here, then all is well
	retval, _ := ioutil.ReadAll(pl.Pipe.Stdout.NewReader())
	return retval, pl.Err
}

// Exec executes a pipeline
func (pl *Pipeline) Exec() {
	// do we have a pipeline to play with?
	if pl == nil {
		return
	}

	// is the pipeline fit to use?
	if pl.Pipe == nil {
		return
	}

	for _, step := range pl.Steps {
		// at this point, stdout needs to become the next
		// stdin
		pl.Pipe.Next()

		// run the next step
		pl.StatusCode, pl.Err = step(pl.Pipe)

		// we stop executing the moment something goes wrong
		if pl.Err != nil {
			return
		}
	}

	// all done
}

// String returns the pipeline's stdout as a single string
func (pl *Pipeline) String() (string, error) {
	// do we have a pipeline to play with?
	if pl == nil {
		return "", nil
	}

	// was the pipeline correctly initialised?
	if pl.Pipe == nil {
		return "", pl.Err
	}

	// did an error occur?
	if pl.Err != nil {
		return pl.Pipe.Stderr.String(), pl.Err
	}

	// if we get here, then all is well
	return pl.Pipe.Stdout.String(), nil
}

// Strings returns the pipeline's stdout, one string per line
func (pl *Pipeline) Strings() ([]string, error) {
	// do we have a pipeline to play with?
	if pl == nil {
		return []string{}, nil
	}

	// was the pipeline correctly initialised?
	if pl.Pipe == nil {
		return []string{}, pl.Err
	}

	// did an error occur?
	if pl.Err != nil {
		return pl.Pipe.Stderr.Strings(), pl.Err
	}

	// if we get here, then all is well
	return pl.Pipe.Stdout.Strings(), nil
}
