# ecoscript

A majorly hackable ecosystem simulator.

## installation

Via `go get`:

```console
go get github.com/dustinrohde/ecoscript
```

## development

This section assumes you're in the project directory, in a terminal. If you
followed the instructions in [installation](installation), the source code
should be in `$GOPATH/src/github.com/dustinrohde/ecoscript/`.

Task management is handled with [robo](https://github.com/tj/robo), a Go task
runner. To install via `go get`:

```console
go get github.com/tj/robo
```

### list available tasks

To list all available tasks, run `robo` with no arguments:

```console
robo
```

### get help on a specific task

To get help on a specific task, use `robo help` with the task name:

```console
robo help TASK
```

### best practices

The two most useful tasks are `check` and `all`.

`check` asserts code quality. It vets, formats, and lints the project and runs
the unit tests.

Run this whenever you edit the project.

```console
robo check
```

`all` safely builds the project by running `check` and then building it.

Run this before pushing changes to the repo.

```console
robo all
```

 There is `pre-push` hook in `tools/` that does this. To use it, just copy it
 to `.git/hooks/`:

 ```console
cp tools/pre-push .git/hooks 
 ```
