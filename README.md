# Welcome to Pipe!

pipe is a Golang library. It helps you write code that wants to work with input sources and output destinations.

It is released under the 3-clause New BSD license. See [./LICENSE.md](./LICENSE.md) for details.

- [What Does Pipe Do](#what-does-pipe-do)
- [Why Did We Build Pipe](#why-did-we-build-pipe)
- [Getting Started](#getting-started)
- [Composing PipeCommands](#composing-pipecommands)
- [Creating A Pipe](#creating-a-pipe)
- [Using A Pipe](#using-a-pipe)
- [Using The Stdin, Stdout And Stderr Stacks](#using-the-stdin-stdout-and-stderr-stacks)
- [Constants](#constants)
- [Functions](#functions)
  - [func AttachOsStderr](#func-attachosstderr)
  - [AsFunctionalOption](#asfunctionaloption)
  - [AsPipeCommand](#aspipecommand)
  - [func AttachOsStdin](#func-attachosstdin)
  - [AsFunctionalOption](#asfunctionaloption-1)
  - [AsPipeCommand](#aspipecommand-1)
  - [func AttachOsStdout](#func-attachosstdout)
  - [AsFunctionalOption](#asfunctionaloption-2)
  - [AsPipeCommand](#aspipecommand-2)
  - [func SetStatusCode](#func-setstatuscode)
- [Types](#types)
  - [type ErrNonZeroStatusCode](#type-errnonzerostatuscode)
  - [type Pipe](#type-pipe)
  - [type PipeCommand](#type-pipecommand)
  - [type PipeOption](#type-pipeoption)

## What Does Pipe Do

We've built Pipe to simulate the UNIX-style process execution environment, complete with: stdin, stdout, stderr, a status code and a Go error.

Oh, and we also threw in an emulation of environment variables too, using our
Envish [https://github.com/ganbarodigital/go_envish](https://github.com/ganbarodigital/go_envish) package.

Instead of running external processes, Pipe runs Golang functions that conform to the PipeCommand type.

It is inspired by:

- [http://labix.org/pipe](http://labix.org/pipe)

- [https://github.com/bitfield/script](https://github.com/bitfield/script)

If you want to see Pipe in use, take a look at our Scriptish library [https://github.com/ganbarodigital/go_scriptish](https://github.com/ganbarodigital/go_scriptish).

## Why Did We Build Pipe

Pipe is one of the reusable packages that we extracted out of Scriptish.

We created Scriptish to make it easier to replace UNIX shell scripts with
compiled Golang binaries. UNIX shell scripts work by piping data out of one
command and into the next. We built Pipe to represent that boundary between
commands in a pipeline.

Consider a classic UNIX terminal command sequence, that people around the
world use all the time:

```go
# show me the files and folders in my home directory
# but one page at a time!
ls $HOME | less
```

When we describe what these commands are doing, we might say "`ls` my home
directory, and pipe the results into `less`".

## Getting Started

Import Pipe into your Golang code:

```go
import pipe "github.com/ganbarodigital/go_pipe/v7"
```

Create a new pipe, and provide a source to read from:

```go
p := pipe.NewPipe()
p.Stdin = NewSourceFromReader(os.Stdin)
```

Define a function that's compatible with our PipeCommand type, and then pass
it into your pipe's `RunCommand()` function:

```go
// this example is abridged from the Scriptish source code
//
// it reads from the pipe's Stdin, sorts the data, and then writes
// out the sorted lines to the pipe's Stdout
func Sort(p *Pipe) (int, error) {
  lines := p.Stdin.Strings()
  sort.Strings(lines)

  for _, line := range lines {
      p.Stdout.WriteString(line)
      p.Stdout.WriteRune('\n')
  }

  return StatusOkay, nil
}

p.RunCommand(Sort)
```

Once your command has run, you can get its status code and Golang error from
the pipe:

```go
statusCode, err := p.StatusError()
```

## Composing PipeCommands

Pipe gives you the glue that you can use to chain PipeCommands together. It
standardises the behaviour of input, output and error handling, so that you
can build your own solution for composing standardised command functions.

```go
// this example code is based on the PipelineController from Scriptish
func runCommands(p *pipe.Pipe, steps ...Pipe.PipeCommand) {
	// do we have a pipeline to play with?
	if p == nil {
		return
	}

	// execute everything in our pipeline
	for _, step := range steps {
		// at this point, stdout needs to become the next
		// stdin
		preparePipeForNextCommand(p)

		// run the next step
		p.RunCommand(step)

		// we stop executing the moment something goes wrong
		err := p.Error()
		if err != nil {
			// we cannot continue
			return
		}
	}
}

// helper function
func preparePipeForNextCommand(p *Pipe) {
	// the output from our previous command becomes the input to the next
	p.SetStdinFromString(p.Stdout.String())

	// the next command starts with no output
	p.SetNewStdout()

	// we throw away any errors that have been written here
	p.SetNewStderr()
}
```

## Creating A Pipe

Call `NewPipe()` when you want to create a new `pipe.Pipe`:

```go
p := NewPipe()
```

The new pipe:

* has an empty `p.Stdin`

* has an empty `p.Stdout`

* has an empty `p.Stderr`

* has a `p.StatusCode()` of `StatusOkay` (ie, `0`)

* has a `p.Error()` of `nil`

* has a `p.Env` that works directly with your program's environment variables

## Using A Pipe

`PipeCommand` is the signature of any function that will work with our Pipe.

```go
type PipeCommand = func(*pipe.Pipe) (int, error)
```

It takes a `pipe.Pipe` as its input parameter, and it returns:

* a UNIX-style status code (0 for success, 1 or higher for an error), and

* a Golang-style error to provide more details about what went wrong

Once you have defined your command function, pass it into your Pipe's RunCommand() function:

```go
// this example is abridged from the Scriptish source code
func Sort(p *Pipe) (int, error) {
    lines := p.Stdin.Strings()
    sort.Strings(lines)

    for _, line := range lines {
        p.Stdout.WriteString(line)
        p.Stdout.WriteRune('\n')
    }

    return StatusOkay, nil
}

p.RunCommand(Sort)
```

## Using The Stdin, Stdout And Stderr Stacks

Sometimes, you may want to temporarily replace the pipe's Stdin, Stdout or
Stderr (for example, to simulate redirecting to `/dev/null`).

You can use the pipe's push & pop functions for this, so that you don't have
to keep track of it yourself:

```go
Input/Output   | Push Method      | Pop Method
---------------|------------------|-----------
`p.Stdin`      | `p.PushStdin()`  | `p.PopStdin()`
`p.Stdout`     | `p.PushStdout()` | `p.PopStdout()`, `p.PopStdoutOnly()`
`p.Stderr`     | `p.PushStderr()` | `p.PopStderr()`, `p.PopStderrOnly()`
```

Here's an example of how to use the stacks:

```go
pipe := NewPipe()

// temporarily redirect to /dev/null
p.PushStderr(ioextra.NewTextDevNull())

// run a command
p.RunCommand(myCommand)

// restore the previous stderr
p.PopStderr()
```

## Constants

```golang
const (
    // StatusOkay is what a PipeCommand returns when everything worked.
    StatusOkay = iota

    // StatusNotOkay is what a PipeCommand returns when it did not work.
    StatusNotOkay
)
```

## Functions

### func [AttachOsStderr](/opts_attachprocessio.go#L70)

`func AttachOsStderr(p *Pipe) (int, error)`

AttachOsStderr sets the pipe to write to your program's Stderr.

You can use this both as a functional option, and/or as a
PipeCommand.

### AsFunctionalOption

```golang
package main

import (
	"fmt"
	pipe "github.com/ganbarodigital/go_pipe/v6"
)

func main() {
	// the pipe will now read from os.Stdin
	p := pipe.NewPipe(pipe.AttachOsStderr)

	// prove that the option did not error out
	statusCode, err := p.StatusError()
	fmt.Printf("statusCode is: %d\n", statusCode)
	fmt.Printf("err is: %v\n", err)
}

```

 Output:

```
statusCode is: 0
err is: <nil>
```

### AsPipeCommand

```golang
package main

import (
	pipe "github.com/ganbarodigital/go_pipe/v6"
)

func main() {
	// create a new pipe
	p := pipe.NewPipe()

	// the pipe will now read from os.Stderr
	p.RunCommand(pipe.AttachOsStderr)
}

```

### func [AttachOsStdin](/opts_attachprocessio.go#L52)

`func AttachOsStdin(p *Pipe) (int, error)`

AttachOsStdin sets the pipe to read from your program's Stdin.

You can use this both as a functional option, and/or as a
PipeCommand.

### AsFunctionalOption

```golang
package main

import (
	"fmt"
	pipe "github.com/ganbarodigital/go_pipe/v6"
)

func main() {
	// the pipe will now read from os.Stdin
	p := pipe.NewPipe(pipe.AttachOsStdin)

	// prove that the option did not error out
	statusCode, err := p.StatusError()
	fmt.Printf("statusCode is: %d\n", statusCode)
	fmt.Printf("err is: %v\n", err)
}

```

 Output:

```
statusCode is: 0
err is: <nil>
```

### AsPipeCommand

```golang
package main

import (
	pipe "github.com/ganbarodigital/go_pipe/v6"
)

func main() {
	// create a new pipe
	p := pipe.NewPipe()

	// the pipe will now read from os.Stdin
	p.RunCommand(pipe.AttachOsStdin)
}

```

### func [AttachOsStdout](/opts_attachprocessio.go#L61)

`func AttachOsStdout(p *Pipe) (int, error)`

AttachOsStdout sets the pipe to write to your program's Stdout.

You can use this both as a functional option, and/or as a
PipeCommand.

### AsFunctionalOption

```golang
package main

import (
	"fmt"
	pipe "github.com/ganbarodigital/go_pipe/v6"
)

func main() {
	// the pipe will now read from os.Stdout
	p := pipe.NewPipe(pipe.AttachOsStdout)

	// prove that the option did not error out
	statusCode, err := p.StatusError()
	fmt.Printf("statusCode is: %d\n", statusCode)
	fmt.Printf("err is: %v\n", err)
}

```

 Output:

```
statusCode is: 0
err is: <nil>
```

### AsPipeCommand

```golang
package main

import (
	pipe "github.com/ganbarodigital/go_pipe/v6"
)

func main() {
	// create a new pipe
	p := pipe.NewPipe()

	// the pipe will now read from os.Stdout
	p.RunCommand(pipe.AttachOsStdout)
}

```

### func [SetStatusCode](/pipe.go#L581)

`func SetStatusCode(p *Pipe, newStatusCode int)`

SetStatusCode is a helper method, added to help us test this
package.

It is not part of our supported API. Use at your own risk!

## Types

### type [ErrNonZeroStatusCode](/errors.go#L47)

`type ErrNonZeroStatusCode struct { ... }`

ErrNonZeroStatusCode is the error returned by Pipe.RunCommand when
a PipeCommand has finished with a non-zero status code, and no error
of its own.

#### func (ErrNonZeroStatusCode) [Error](/errors.go#L52)

`func (e ErrNonZeroStatusCode) Error() string`

### type [Pipe](/pipe.go#L51)

`type Pipe struct { ... }`

Pipe is our data structure. All PipeCommands read from, and/or write to
the pipe.

#### func [NewPipe](/pipe.go#L88)

`func NewPipe(options ...PipeOption) *Pipe`

NewPipe creates a new Pipe that's ready to use.

It starts with an empty Stdin, and uses the program's environment
by default.

You can provide a list of functional options for us to call. We'll
pass in the Pipe, for you to reconfigure. If your functional option
returns an error, we'll store that in the Pipe, and stop processing
any further functional options.

```golang
package main

import (
	"fmt"
	pipe "github.com/ganbarodigital/go_pipe/v6"
)

func main() {
	// create a new pipe
	p := pipe.NewPipe()

	// it starts with no error set
	statusCode, err := p.StatusError()

	fmt.Printf("statusCode is: %d\n", statusCode)
	fmt.Printf("err is: %v\n", err)
}

```

 Output:

```
statusCode is: 0
err is: <nil>
```

#### func (*Pipe) [DrainStdinToStdout](/pipe.go#L118)

`func (p *Pipe) DrainStdinToStdout()`

DrainStdinToStdout will copy everything that's left in the pipe's Stdin
over to the pipe's Stdout.

#### func (*Pipe) [Error](/pipe.go#L140)

`func (p *Pipe) Error() error`

Error returns the error returned from the last PipeCommand
that ran against this pipe.

#### func (*Pipe) [Okay](/pipe.go#L152)

`func (p *Pipe) Okay() bool`

Okay confirms that the last PipeCommand run against the pipe completed
without reporting an error.

#### func (*Pipe) [PopStderr](/pipe.go#L481)

`func (p *Pipe) PopStderr()`

PopStderr sets the pipe's Stderr to its previous value.

It reverses your last call to PushStderr.

This is useful for callers who need to temporarily replace the pipe's
Stderr (for example, to redirect to /dev/null).

NOTE: if p.Stdout == p.Stderr, PopStderr sets *both* p.Stdout and
p.Stderr to the previous Stderr. Most of the time, this is the desired
intention.

Use PopStderrOnly when you don't want to touch the pipe's Stdout at all.

#### func (*Pipe) [PopStderrOnly](/pipe.go#L518)

`func (p *Pipe) PopStderrOnly()`

PopStderrOnly sets the pipe's Stderr to its previous value.

It reverses your last call to PushStderr.

This is useful for callers who need to temporarily replace the pipe's
Stderr (for example, to redirect to /dev/null).

NOTE: even if p.Stdout == p.Stderr, PopStderrOnly leaves p.Stdout
untouched.

#### func (*Pipe) [PopStdin](/pipe.go#L267)

`func (p *Pipe) PopStdin()`

PopStdin sets the pipe's Stdin to its previous value.

It reverses your last call to PushStdin.

This is useful for callers who need to temporarily replace the pipe's
Stdin.

#### func (*Pipe) [PopStdout](/pipe.go#L355)

`func (p *Pipe) PopStdout()`

PopStdout sets the pipe's Stdout to its previous value.

It reverses your last call to PushStdout.

This is useful for callers who need to temporarily replace the pipe's
Stdout (for example, to redirect to /dev/null).

NOTE: if p.Stdout == p.Stderr, PopStdout sets *both* p.Stdout and
p.Stderr to the previous Stdout. Most of the time, this is the desired
intention.

Use PopStdoutOnly when you don't want to touch the pipe's Stderr at all.

#### func (*Pipe) [PopStdoutOnly](/pipe.go#L392)

`func (p *Pipe) PopStdoutOnly()`

PopStdoutOnly sets the pipe's Stdout to its previous value.

It reverses your last call to PushStdout.

This is useful for callers who need to temporarily replace the pipe's
Stdout (for example, to redirect to /dev/null).

NOTE: even if p.Stdout == p.Stderr, PopStdoutOnly leaves p.Stderr
untouched.

#### func (*Pipe) [PushStderr](/pipe.go#L450)

`func (p *Pipe) PushStderr(newStderr ioextra.TextReaderWriter)`

PushStderr adds the pipe's existing Stderr to an internal stack,
and then sets the pipe's Stderr to the given newStderr.

You can call PopStderr to reverse this operation.

This is useful for callers who need to temporarily replace the pipe's
Stderr (for example, to redirect to /dev/null).

NOTE: if p.Stdout == p.Stderr, PushStderr sets *both* p.Stdout and
p.Stderr to the newStderr.

#### func (*Pipe) [PushStdin](/pipe.go#L250)

`func (p *Pipe) PushStdin(newStdin ioextra.TextReader)`

PushStdin adds the pipe's existing Stdin to an internal stack,
and then sets the pipe's Stdin to the given newStdin.

You can call PopStdin to reverse this operation.

This is useful for callers who need to temporarily replace the pipe's
Stdin.

#### func (*Pipe) [PushStdout](/pipe.go#L325)

`func (p *Pipe) PushStdout(newStdout ioextra.TextReaderWriter)`

PushStdout adds the pipe's existing Stdout to an internal stack,
and then sets the pipe's Stdout to the given newStdout.

You can call PopStdout to reverse this operation.

This is useful for callers who need to temporarily replace the pipe's
Stdout (for example, to redirect to /dev/null).

NOTE: if p.Stdout == p.Stderr, PushStdout sets *both* p.Stdout and
p.Stderr to the newStdout.

#### func (*Pipe) [ResetBuffers](/pipe.go#L167)

`func (p *Pipe) ResetBuffers()`

ResetBuffers creates new, empty Stdin, Stdout and Stderr for the given
pipe.

It also empties the internal stacks used by PushStdin / PopStdin,
PushStdout / PopStdout, and PushStderr / PopStderr.

#### func (*Pipe) [ResetError](/pipe.go#L186)

`func (p *Pipe) ResetError()`

ResetError sets the pipe's status code and error to their zero values
of (StatusOkay, nil).

#### func (*Pipe) [RunCommand](/pipe.go#L199)

`func (p *Pipe) RunCommand(c PipeCommand)`

RunCommand will run a function using this pipe. The function's return
values are stored in the pipe's StatusCode and Err fields.

#### func (*Pipe) [SetNewStderr](/pipe.go#L428)

`func (p *Pipe) SetNewStderr()`

SetNewStderr creates a new, empty Stderr buffer on this pipe.

#### func (*Pipe) [SetNewStdin](/pipe.go#L215)

`func (p *Pipe) SetNewStdin()`

SetNewStdin creates a new, empty Stdin buffer on this pipe.

#### func (*Pipe) [SetNewStdout](/pipe.go#L303)

`func (p *Pipe) SetNewStdout()`

SetNewStdout creates a new, empty Stdout buffer on this pipe.

#### func (*Pipe) [SetStdinFromString](/pipe.go#L228)

`func (p *Pipe) SetStdinFromString(input string)`

SetStdinFromString sets the pipe's Stdin to be the given input string.

#### func (*Pipe) [StatusCode](/pipe.go#L555)

`func (p *Pipe) StatusCode() int`

StatusCode returns the UNIX-like status code from the last PipeCommand
that ran against this pipe.

#### func (*Pipe) [StatusError](/pipe.go#L567)

`func (p *Pipe) StatusError() (int, error)`

StatusError is a shorthand for calling p.StatusCode() and p.Error()
to get the UNIX-like status code and the last reported Golang error.

#### func (*Pipe) [StderrStackLen](/pipe.go#L542)

`func (p *Pipe) StderrStackLen() int`

StderrStackLen returns the number of entries in the internal stack of
Stderr entries.

You can call PushStderr and PopStderr to add entries to & from the
internal stack.

#### func (*Pipe) [StdinStackLen](/pipe.go#L291)

`func (p *Pipe) StdinStackLen() int`

StdinStackLen returns the number of entries in the internal stack of
Stdin entries.

You can call PushStdin and PopStdin to add entries to & from the
internal stack.

#### func (*Pipe) [StdoutStackLen](/pipe.go#L416)

`func (p *Pipe) StdoutStackLen() int`

StdoutStackLen returns the number of entries in the internal stack of
Stdout entries.

You can call PushStdout and PopStdout to add entries to & from the
internal stack.

### type [PipeCommand](/pipecommand.go#L44)

`type PipeCommand = func(*Pipe) (int, error)`

PipeCommand is the signature of any function that will work with
our Pipe.

### type [PipeOption](/pipeoption.go#L44)

`type PipeOption = PipeCommand`

PipeOption describes functional options that you can pass into
NewPipe.

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
