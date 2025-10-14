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

## Mutators

### Arithmetic
| Name | Original | Mutated |
| --- | --- | --- |
| Bitwise AND (`Arithmetic_AND`) | `a & b` | `a \| b` |
| Bitwise OR (`Arithmetic_OR`) | `a \| b` | `a & b` |
| Bitwise XOR (`Arithmetic_XOR`) | `a ^ b` | `a & b` |
| Bitwise NOT (`Arithmetic_NOT`) | `^a` | `a` |
| Addition (`Arithmetic_ADD`) | `a + b` | `a - b` |
| Addition Assign (`Arithmetic_ADD_ASSIGN`) | `a += b` | `a -= b` |
| Subtraction (`Arithmetic_SUB`) | `a - b` | `a + b` |
| Subtraction Assign (`Arithmetic_SUB_ASSIGN`) | `a -= b` | `a += b` |
| Multiplication (`Arithmetic_QUO`) | `a * b` | `a / b` |
| Multiplication Assign (`Arithmetic_MUL_ASSIGN`) | `a *= b` | `a /= b` |
| Division (`Arithmetic_QUO`) | `a / b` | `a * b` |
| Division Assign (`Arithmetic_QUO_ASSIGN`) | `a /= b` | `a *= b` |
| Modulus (`Arithmetic_REM`) | `a % b` | `a * b` |
| Modulus Assign (`Arithmetic_REM_ASSIGN`) | `a %= b` | `a *= b` |
| Shift Left (`Arithmetic_SHL`) | `a << b` | `a >> b` |
| Shift Right (`Arithmetic_SHR`) | `a >> b` | `a << b` |

> Note: The addition mutator implements a type-aware guard to avoid mutating string concatenations.

### Boolean
| Name | Original | Mutated |
| --- | --- | --- |
| Boolean TRUE (`Boolean_TRUE`) | `true` | `false` |
| Boolean FALSE (`Boolean_FALSE`) | `false` | `true` |

### Conditional Boundary
| Name | Original | Mutated |
| --- | --- | --- |
| Greater Than (`ConditionalBoundary_GTR_GEQ`) | `a > b` | `a >= b` |
| Greater Than Or Equal (`ConditionalBoundary_GEQ_GTR`) | `a >= b` | `a > b` |
| Less Than (`ConditionalBoundary_LSS_LEQ`) | `a < b` | `a <= b` |
| Less Than Or Equal (`ConditionalBoundary_LEQ_LSS`) | `a <= b` | `a < b` |

### Logical
| Name | Original | Mutated |
| --- | --- | --- |
| Logical AND (`Logical_AND`) | `a && b` | `a \|\| b` |
| Logical OR (`Logical_OR`) | `a \|\| b` | `a && b` |

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
