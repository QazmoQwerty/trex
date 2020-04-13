# Trex

A toy language for quick & easy string manipulation.

For more info read the [language specification](docs/trex-spec.md)

## Install

Download the relevant binary from here: gitlab.com/QazmoQwerty/trex/releases

Or alternatively if you have go installed use

```
go get gitlab.com/QazmoQwerty/trex
```

and this will build the binary in $GOPATH/bin.

## Examples

```
factorial(n) => 1 if n = 0 else n * factorial(n - 1)
```

```
// returns longest word in a string
longestWord => max(#len) words
```

```
// countChars "aabdbg" = (a, 2), (b, 2), (d, 1), (g, 1) 
countChars => ch, numOccurs(ch) for ch in unique chars
```


```
// primes(n) returns all prime numbers from 0 to n
isPrime(n) => count (i from 2..n if n % i = 0) = 0
primes(n) => i from 0..n if isPrime(i)
```


```
// (fold(a,b -> a+b) chars 12345) = 15
fold(f) => f([0], fold(#f) [1:]) if len > 1 else [0] if len > 0 else ()
```

