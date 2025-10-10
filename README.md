# gomut

A mutation testing tool for Go that helps you measure the effectiveness of your test suite.

## What is Mutation Testing?

Mutation testing introduces small changes (mutations) to your code and runs your tests to see if they catch these changes. If tests fail, the mutation is "killed" (good). If tests pass, the mutation "survived" (bad / your tests missed it).

## Installation

```bash
go install github.com/renja-g/go-mutation-testing/cmd/gomut@latest
```
### Uninstall
```bash
rm $(go env GOPATH)/bin/gomut
```

Or build from source:
```bash
git clone https://github.com/renja-g/go-mutation-testing.git
cd go-mutation-testing
go build -o gomut ./cmd/gomut
```

## Usage

Basic usage:

```bash
gomut -path ./src
```

### Options

- `-path` - Path to source directory to mutate (default: `./src`)
- `-pkg` - Go package pattern to test (default: `./...`)
- `-list` - List mutations without running tests
- `-v` - Verbose: print test output per mutation

### Examples

List all mutations without running tests:
```bash
gomut -path ./myapp -list
```

Run mutation testing on a specific package:
```bash
gomut -path ./myapp -pkg ./internal/...
```

Run with verbose output:
```bash
gomut -path ./src -v
```

## Output

The tool displays:
- Total mutations discovered
- Each mutation tested with its status (KILLED ✓ or SURVIVED ✗)
- Final score: `(killed / total) * 100%`

A higher score means your tests are more effective at catching bugs.
