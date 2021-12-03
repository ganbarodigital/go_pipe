/*
pipe is a Golang library. It helps you write code that wants to work with input sources and output destinations.

It is released under the 3-clause New BSD license. See ./LICENSE.md for details.

What Does Pipe Do

We've built Pipe to simulate the UNIX-style process execution environment, complete with: stdin, stdout, stderr, a status code and a Go error.

Oh, and we also threw in an emulation of environment variables too, using our
Envish https://github.com/ganbarodigital/go_envish package.

Instead of running external processes, Pipe runs Golang functions that conform to the PipeCommand type.

It is inspired by:

- http://labix.org/pipe

- https://github.com/bitfield/script

If you want to see Pipe in use, take a look at our Scriptish library https://github.com/ganbarodigital/go_scriptish.


Why Did We Build Pipe

Pipe is one of the reusable packages that we extracted out of Scriptish.

We created Scriptish to make it easier to replace UNIX shell scripts with
compiled Golang binaries. UNIX shell scripts work by piping data out of one
command and into the next. We built Pipe to represent that boundary between
commands in a pipeline.

Consider a classic UNIX terminal command sequence, that people around the
world use all the time:

  # show me the files and folders in my home directory
  # but one page at a time!
  ls $HOME | less

When we describe what these commands are doing, we might say "`ls` my home
directory, and pipe the results into `less`".


Getting Started

Import Pipe into your Golang code:

  import pipe "github.com/ganbarodigital/go_pipe/v7"

Create a new pipe, and provide a source to read from:

  p := pipe.NewPipe()
  p.Stdin = NewSourceFromReader(os.Stdin)

Define a function that's compatible with our PipeCommand type, and then pass
it into your pipe's `RunCommand()` function:

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

Once your command has run, you can get its status code and Golang error from
the pipe:

  statusCode, err := p.StatusError()


Composing PipeCommands

Pipe gives you the glue that you can use to chain PipeCommands together. It
standardises the behaviour of input, output and error handling, so that you
can build your own solution for composing standardised command functions.

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


Creating A Pipe

Call `NewPipe()` when you want to create a new `pipe.Pipe`:

  p := NewPipe()

The new pipe:

* has an empty `p.Stdin`

* has an empty `p.Stdout`

* has an empty `p.Stderr`

* has a `p.StatusCode()` of `StatusOkay` (ie, `0`)

* has a `p.Error()` of `nil`

* has a `p.Env` that works directly with your program's environment variables

Using A Pipe

`PipeCommand` is the signature of any function that will work with our Pipe.

  type PipeCommand = func(*pipe.Pipe) (int, error)

It takes a `pipe.Pipe` as its input parameter, and it returns:

* a UNIX-style status code (0 for success, 1 or higher for an error), and

* a Golang-style error to provide more details about what went wrong

Once you have defined your command function, pass it into your Pipe's RunCommand() function:

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


Using The Stdin, Stdout And Stderr Stacks

Sometimes, you may want to temporarily replace the pipe's Stdin, Stdout or
Stderr (for example, to simulate redirecting to `/dev/null`).

You can use the pipe's push & pop functions for this, so that you don't have
to keep track of it yourself:

  Input/Output   | Push Method      | Pop Method
  ---------------|------------------|-----------
  `p.Stdin`      | `p.PushStdin()`  | `p.PopStdin()`
  `p.Stdout`     | `p.PushStdout()` | `p.PopStdout()`, `p.PopStdoutOnly()`
  `p.Stderr`     | `p.PushStderr()` | `p.PopStderr()`, `p.PopStderrOnly()`

Here's an example of how to use the stacks:

  pipe := NewPipe()

  // temporarily redirect to /dev/null
  p.PushStderr(ioextra.NewTextDevNull())

  // run a command
  p.RunCommand(myCommand)

  // restore the previous stderr
  p.PopStderr()
*/
package pipe
