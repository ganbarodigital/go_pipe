# Welcome to pipe

## Introduction

`pipe` is a Golang library to help you write UNIX-like processing pipelines.

It is inspired by:

- http://labix.org/pipe
- https://github.com/bitfield/script

If you want to see `pipe` in use, take a look at our [scriptish library](https://github.com/ganbarodigital/scriptish).

## What Does It Do?

`pipe` gives you three types:

* `Pipeline` - a list of operations to execute, and
* `Pipe` - an input source/output buffer that's passed into each operation in turn
* `PipeOperation` - a function that does some work

`PipeOperation`s read from the `Pipe.Stdin`, and write to the `Pipe.Stdout` and/or the `Pipe.Stderr`.

When the pipeline has finished executing, there's some helper functions for you to get the final value of `Pipe.Stdout` back into your regular Golang code.

Together, these provide the primitive building blocks needed to create higher-level UNIX-like processing pipelines.

