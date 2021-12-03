# Welcome to pipe!

## Introduction

Pipe is a Golang library. It helps you write code that wants to work with input sources and output destinations.

It is released under the 3-clause New BSD license. See [LICENSE.md](LICENSE.md) for details.

## Table of Contents <!-- omit in toc -->

- [Introduction](#introduction)
- [Why Use Pipe?](#why-use-pipe)
  - [What Does Pipe Do?](#what-does-pipe-do)
  - [Why Did We Build Pipe?](#why-did-we-build-pipe)
- [How Does It Work?](#how-does-it-work)
  - [Getting Started](#getting-started)
  - [Composing PipeCommands](#composing-pipecommands)
- [Creating A Pipe](#creating-a-pipe)
  - [NewPipe()](#newpipe)
- [Using A Pipe](#using-a-pipe)
  - [PipeCommand](#pipecommand)
  - [Using The Stdin, Stdout And Stderr Stacks](#using-the-stdin-stdout-and-stderr-stacks)
- [Pipe](#pipe)
  - [NewPipe()](#newpipe-1)
  - [Pipe Functional Options](#pipe-functional-options)
  - [Pipe.DrainStdinToStdout()](#pipedrainstdintostdout)
  - [Pipe.Error()](#pipeerror)
  - [Pipe.Okay()](#pipeokay)
  - [Pipe.ResetBuffers()](#piperesetbuffers)
  - [Pipe.ResetError()](#pipereseterror)
  - [Pipe.RunCommand()](#piperuncommand)
  - [Pipe.SetNewStdin()](#pipesetnewstdin)
  - [Pipe.SetStdinFromString()](#pipesetstdinfromstring)
  - [Pipe.PushStdin()](#pipepushstdin)
  - [Pipe.PopStdin()](#pipepopstdin)
  - [Pipe.StdinStackLen()](#pipestdinstacklen)
  - [Pipe.SetNewStderr()](#pipesetnewstderr)
  - [Pipe.PushStderr()](#pipepushstderr)
  - [Pipe.PopStderr()](#pipepopstderr)
  - [Pipe.PopStderrOnly()](#pipepopstderronly)
  - [Pipe.StderrStackLen()](#pipestderrstacklen)
  - [Pipe.SetNewStdout()](#pipesetnewstdout)
  - [Pipe.PushStdout()](#pipepushstdout)
  - [Pipe.PopStdout()](#pipepopstdout)
  - [Pipe.PopStdoutOnly()](#pipepopstdoutonly)
  - [Pipe.StdoutStackLen()](#pipestdoutstacklen)
  - [Pipe.StatusCode()](#pipestatuscode)
  - [Pipe.StatusError()](#pipestatuserror)

## Why Use Pipe?

### What Does Pipe Do?

We've built Pipe to simulate the UNIX-style process execution environment, complete with:

* stdin,
* stdout,
* stderr,
* a status code
* and a Go error

Oh, and we also threw in an emulation of environment variables too, using our [Envish][Envish] package.

Instead of running external processes, Pipe runs Golang functions that conform to the [Command][Command] interface.

It is inspired by:

- http://labix.org/pipe
- https://github.com/bitfield/script

If you want to see `pipe` in use, take a look at our [Scriptish library](https://github.com/ganbarodigital/go_scriptish).

### Why Did We Build Pipe?

Pipe is one of the reusable packages that we extracted out of [Scriptish][Scriptish].

We created [Scriptish][Scriptish] to make it easier to replace UNIX shell scripts with compiled Golang binaries. UNIX shell scripts work by piping data out of one command and into the next. We built Pipe to represent that boundary between commands in a pipeline.

Consider a classic UNIX terminal command sequence, that people around the world use all the time:

```shell
# show me the files and folders in my home directory
# but one page at a time!
ls $HOME | less
```

When we describe what these commands are doing, we might say "`ls` my home directory, and pipe the results into `less`".

## How Does It Work?

### Getting Started

Import Pipe into your Golang code:

```golang
import pipe "github.com/ganbarodigital/go_pipe/v5"
```

Create a new pipe, and provide a source to read from:

```golang
p := pipe.NewPipe()
p.Stdin = NewSourceFromReader(os.Stdin)
```

Define a [PipeCommand][PipeCommand], and then pass it into your pipe's `RunCommand()` function:

```golang
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

Once your command has run, you can get its status code and Golang error:

```golang
statusCode, err := p.StatusError()
```

### Composing PipeCommands

Pipe gives you the glue that you can use to chain [PipeCommands][PipeCommand] together. It standardises the behaviour of input, output and error handling, so that you can build your own solution for composing standardised [command functions][PipeCommand].

```golang
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

### NewPipe()

Call `NewPipe()` when you want to create a new `pipe.Pipe`:

```golang
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

### PipeCommand

`PipeCommand` is the signature of any function that will work with our Pipe.

```golang
type PipeCommand = func(*pipe.Pipe) (int, error)
```

It takes a `pipe.Pipe` as its input parameter, and it returns:

* a UNIX-style status code (0 for success, 1 or higher for an error), and
* a Golang-style error to provide more details about what went wrong

Once you have defined your command function, pass it into your Pipe's [RunCommand()][RunCommand] function:

```golang
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

### Using The Stdin, Stdout And Stderr Stacks

Sometimes, you may want to temporarily replace the pipe's Stdin, Stdout or Stderr (for example, to simulate redirecting to `/dev/null`).

You can use the pipe's push & pop functions for this, so that you don't have to keep track of it yourself:

Input/Output   | Push Method      | Pop Method
---------------|------------------|-----------
`p.Stdin`      | `p.PushStdin()`  | `p.PopStdin()`
`p.Stdout`     | `p.PushStdout()` | `p.PopStdout()`, `p.PopStdoutOnly()`
`p.Stderr`     | `p.PushStderr()` | `p.PopStderr()`, `p.PopStderrOnly()`

```golang
pipe := NewPipe()

// temporarily redirect to /dev/null
p.PushStderr(ioextra.NewTextDevNull())

// run a command
p.RunCommand(myCommand)

// restore the previous stderr
p.PopStderr()
```

## Pipe

```golang
// Pipe is our data structure. All PipeCommands read from, and/or write to
// the pipe.
type Pipe struct {
	// PipeCommands read from Stdin
	Stdin ioextra.TextReader

	// PipeCommands write to Stdout and/or Stderr
	Stdout ioextra.TextReaderWriter
	Stderr ioextra.TextReaderWriter

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
```

### NewPipe()

```golang
func NewPipe(options ...func(*Pipe)) *Pipe
```

`NewPipe()` creates a new Pipe that's ready to use.

It starts with an empty Stdin, and uses the program's environment by default.

You can provide a list of functional options for us to call. We'll pass in the Pipe, for you to reconfigure.

```golang
// create a new pipe
p := pipe.NewPipe()
```

```golang
// create a new pipe that reads from / writes to the program's
// stdin, stdout and stderr

p := pipe.NewPipe(
    pipe.AttachOsStdin,
    pipe.AttachOsStdout,
    pipe.AttachOsStderr,
)
```

### Pipe Functional Options

```golang
// AttachOsStdin sets the pipe to read from your program's stdin
func AttachOsStdin(p *Pipe)

// AttachOsStdout sets the pipe to write to your program's stdout
func AttachOsStdout(p *Pipe)

// AttachOsStderr sets the pipe to write to your program's stderr
func AttachOsStderr(p *Pipe)
```

### Pipe.DrainStdinToStdout()

```golang
// DrainStdinToStdout will copy everything that's left in the pipe's Stdin
// over to the pipe's Stdout
func (p *Pipe) DrainStdinToStdout()
```

### Pipe.Error()

```golang
// Error returns the error returned from the last PipeCommand
// that ran against this pipe
func (p *Pipe) Error() error
```

### Pipe.Okay()

```golang
// Okay confirms that the last PipeCommand run against the pipe completed
// without reporting an error
func (p *Pipe) Okay() bool
```

### Pipe.ResetBuffers()

```golang
// ResetBuffers creates new, empty Stdin, Stdout and Stderr for the given
// pipe.
//
// It also empties the internal stacks used by PushStdin / PopStdin,
// PushStdout / PopStdout, and PushStderr / PopStderr.
func (p *Pipe) ResetBuffers()
```

### Pipe.ResetError()

```golang
// ResetError sets the pipe's status code and error to their zero values
// of (StatusOkay, nil)
func (p *Pipe) ResetError()
```

### Pipe.RunCommand()

```golang
// RunCommand will run a function using this pipe. The function's return
// values are stored in the pipe's StatusCode and Err fields.
func (p *Pipe) RunCommand(c PipeCommand)
```

### Pipe.SetNewStdin()

```golang
// SetNewStdin creates a new, empty Stdin buffer on this pipe
func (p *Pipe) SetNewStdin()
```

### Pipe.SetStdinFromString()

```golang
// SetStdinFromString sets the pipe's Stdin to be the given input string
func (p *Pipe) SetStdinFromString(input string) {
```

### Pipe.PushStdin()

```golang
// PushStdin adds the pipe's existing Stdin to an internal stack,
// and then sets the pipe's Stdin to the given newStdin.
//
// You can call PopStdin to reverse this operation.
//
// This is useful for callers who need to temporarily replace the pipe's
// Stdin.
func (p *Pipe) PushStdin(newStdin ioextra.TextReader)
```

### Pipe.PopStdin()

```golang
// PopStdin sets the pipe's Stdin to its previous value.
//
// It reverses your last call to PushStdin.
//
// This is useful for callers who need to temporarily replace the pipe's
// Stdin.
func (p *Pipe) PopStdin()
```

### Pipe.StdinStackLen()

```golang
// StdinStackLen returns the number of entries in the internal stack of
// Stdin entries.
//
// You can call PushStdin and PopStdin to add entries to & from the
// internal stack.
func (p *Pipe) StdinStackLen()
```

### Pipe.SetNewStderr()

```golang
// SetNewStderr creates a new, empty Stderr buffer on this pipe
func (p *Pipe) SetNewStderr()
```

### Pipe.PushStderr()

```golang
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
func (p *Pipe) PushStderr(newStderr ioextra.TextReaderWriter)
```

### Pipe.PopStderr()

```golang
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
func (p *Pipe) PopStderr()
```

### Pipe.PopStderrOnly()

```golang
// PopStderrOnly sets the pipe's Stderr to its previous value.
//
// It reverses your last call to PushStderr.
//
// This is useful for callers who need to temporarily replace the pipe's
// Stderr (for example, to redirect to /dev/null).
//
// NOTE: even if p.Stdout == p.Stderr, PopStderrOnly leaves p.Stdout
// untouched.
func (p *Pipe) PopStderrOnly()
```

### Pipe.StderrStackLen()

```golang
// StderrStackLen returns the number of entries in the internal stack of
// Stderr entries.
//
// You can call PushStderr and PopStderr to add entries to & from the
// internal stack.
func (p *Pipe) StderrStackLen()
```

### Pipe.SetNewStdout()

```golang
// SetNewStdout creates a new, empty Stdout buffer on this pipe
func (p *Pipe) SetNewStdout()
```

### Pipe.PushStdout()

```golang
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
func (p *Pipe) PushStdout(newStdout ioextra.TextReaderWriter)
```

### Pipe.PopStdout()

```golang
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
func (p *Pipe) PopStdout()
```

### Pipe.PopStdoutOnly()

```golang
// PopStdoutOnly sets the pipe's Stdout to its previous value.
//
// It reverses your last call to PushStdout.
//
// This is useful for callers who need to temporarily replace the pipe's
// Stdout (for example, to redirect to /dev/null).
//
// NOTE: even if p.Stdout == p.Stderr, PopStdoutOnly leaves p.Stderr
// untouched.
func (p *Pipe) PopStdoutOnly()
```

### Pipe.StdoutStackLen()

```golang
// StdoutStackLen returns the number of entries in the internal stack of
// Stdout entries.
//
// You can call PushStdout and PopStdout to add entries to & from the
// internal stack.
func (p *Pipe) StdoutStackLen()
```

### Pipe.StatusCode()

```golang
// StatusCode returns the UNIX-like status code from the last PipeCommand
// that ran against this pipe
func (p *Pipe) StatusCode() int
```

### Pipe.StatusError()

```golang
// StatusError is a shorthand for calling p.StatusCode() and p.Error()
// to get the UNIX-like status code and the last reported Golang error
func (p *Pipe) StatusError() (int, error)
```

[PipeCommand]: #pipecommand
[RunCommand]: #piperuncommand
[Envish]: https://github.com/ganbarodigital/go_envish
[Scriptish]: https://github.com/ganbarodigital/go_scriptish