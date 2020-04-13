# Trex

A toy language for quick & easy string manipulation/analysis.

For more info read the [language specification](docs/trex-spec.md)

## Install

Download the relevant binary from here: https://gitlab.com/QazmoQwerty/trex/releases

Or alternatively if you have go installed use

```
go get gitlab.com/QazmoQwerty/trex
```

and this will build the binary in $GOPATH/bin.

## Objectives

Trex runs in a CLI, and is meant to be used for simple tasks involving string manipulation and analysis, especially concerning plaintext files.

## Examples

```
factorial(n) => 1 if n = 0 else n * factorial(n - 1)
```

```
max(#len) words // longest word in a string
```

```
// For input "aabdbg" output would be (a, 2), (b, 2), (d, 1), (g, 1)
c, numoccurs(c) for c in unique chars
```


```
// primes(n) returns all prime numbers from 0 to n
isprime(n) => count (i from 2..n if n % i = 0) = 0
primes(n) => i from 0..n if isprime(i)
```


```
sum => fold(a,b -> a+b) // sum of numbers in list
sum (1, 2, 3, 4, 5, 6) // will output 21
```


## Status

The project is currently in a fairly usable state. There are a few issues and other than that teh main thing left to add is documentation/tutorials for how to use the language and the terminal application.

## Issues

* The program currently can occasionally have issues when being stress-tested.
* The semantics of the range '..' operator are still undecided.