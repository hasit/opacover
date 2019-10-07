# opacover

Generate HTML representation of OPA test coverage.

## Install

Make sure you have `Go` installed on your machine.

```shell
$ go version
go version go1.13 darwin/amd64
```

Get `opacover`. The following command will download the source code and place a runnable binary inside `$GOPATH/bin`.

```shell
$ go get -u https://github.com/hasit/opacover
```

## Use

Whenever in doubt, `--help` flag will guide you.

```shell
$ go run main.go --help
generate HTML representation of OPA test coverage

Usage:
  opacover [flags]

Flags:
  -h, --help             help for opacover
  -i, --input string     input file in JSON
  -m, --modules string   directory path of Rego modules
```

To generate an HTML coverage report of OPA's test coverage, you will need a JSON formatted coverage file.

```shell
$ opa test <path-to-rego-files>/ -v --coverage > coverage.json
```

We will use this JSON file as input to `opacover`.

```shell
$ $GOPATH/bin/opacover --input coverage.json --modules ./ > coverage.html
```

Now, open the newly generated `coverage.html` file in your default browser to view it.

## Contribute

Feel free to open issues and pull requests to improve `opacover`.

## Inspirations

1. [OPA](https://www.openpolicyagent.org/)
2. [Go's Cover Tool](https://golang.org/cmd/cover/)
