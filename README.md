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
runner. You will also need [goimports](golang.org/x/tools/cmd/goimports) and
[golint](github.com/golang/lint). To install via `go get`:

```console
go get -u github.com/tj/robo
go get -u golang.org/x/tools/cmd/goimports
go get -u github.com/golang/lint
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

### useful commands

`run` compiles and runs the project on the fly.

```console
robo run
```

`check` asserts code quality. It vets, formats, and lints the project and runs
the unit tests.

Run this whenever you edit the project.

```console
robo check
```

`all` safely builds the project by running `check` and then building it only
if all checks pass.

Run this before pushing changes to the repo.

```console
robo all
```

If arguments are given to the `check` and `all` tasks, the shell builtin `set`
will be called with those arguments. For example, the `+e` option can be used
to disable the behavior which exits early if a command fails:

```console
robo check +e
robo all +e
```

A `pre-push` git hook that runs is included in `tools/hooks/`. To use it, just
copy it to `.git/hooks/`:

```console
cp tools/hooks/pre-push .git/hooks 
 ```

Alternatively, just run the `install-hooks` command to symlink all bundled git
hooks to the `.git/hooks/` directory:

```console
robo install-hooks
```
