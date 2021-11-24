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
  - [Composing Commands](#composing-commands)
- [Creating A Pipe](#creating-a-pipe)
  - [NewPipe()](#newpipe)
- [Using A Pipe](#using-a-pipe)
  - [Command](#command)
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
  - [Pipe.SetNewStderr()](#pipesetnewstderr)
  - [Pipe.SetNewStdout()](#pipesetnewstdout)
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

Define a [Command][Command], and then pass it into your pipe's `RunCommand()` function:

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

### Composing Commands

Pipe gives you the glue that you can use to chain [Commands][Command] together. It standardises the behaviour of input, output and error handling, so that you can build your own solution for composing standardised [command functions][Command].

```golang
// this example code is based on the PipelineController from Scriptish
func runCommands(p *pipe.Pipe, steps ...Pipe.Command) {
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

### Command

`Command` is the signature of any function that will work with our Pipe.

```golang
type Command = func(*pipe.Pipe) (int, error)
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

## Pipe

```golang
// Pipe is our data structure. All Commands read from, and/or write to
// the pipe.
type Pipe struct {
	// Pipe commands read from Stdin
	Stdin ioextra.TextReader

	// Pipe commands write to Stdout and/or Stderr
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
// Error returns the error returned from the last Command
// that ran against this pipe
func (p *Pipe) Error() error
```

### Pipe.Okay()

```golang
// Okay confirms that the last Command run against the pipe completed
// without reporting an error
func (p *Pipe) Okay() bool
```

### Pipe.ResetBuffers()

```golang
// ResetBuffers creates new, empty buffers for the pipe
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
func (p *Pipe) RunCommand(c Command)
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

### Pipe.SetNewStderr()

```golang
// SetNewStderr creates a new, empty Stderr buffer on this pipe
func (p *Pipe) SetNewStderr()
```

### Pipe.SetNewStdout()

```golang
// SetNewStdout creates a new, empty Stdout buffer on this pipe
func (p *Pipe) SetNewStdout()
```

### Pipe.StatusCode()

```golang
// StatusCode returns the UNIX-like status code from the last Command
// that ran against this pipe
func (p *Pipe) StatusCode() int
```

### Pipe.StatusError()

```golang
// StatusError is a shorthand for calling p.StatusCode() and p.Error()
// to get the UNIX-like status code and the last reported Golang error
func (p *Pipe) StatusError() (int, error)
```

[Command]: #command
[RunCommand]: #piperuncommand
[Envish]: https://github.com/ganbarodigital/go_envish
[Scriptish]: https://github.com/ganbarodigital/go_scriptish