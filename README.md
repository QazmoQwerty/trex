# Trex

A toy language for quick & easy string manipulation.

For more info read the [language specification](docs/trex-spec,md)

# Code Examples


```
factorial(n) => 1 if n = 0 else n * factorial(n - 1)
```

```
// returns longest word in a string
longestWord => max(#len) words
```

```
// countChars "aabdbg" = (a, 2), (b, 2), (d, 1), (g, 1) 
countChars => (ch, count(ch)) for ch in unique chars
```


```
// primes(n) returns all prime numbers from 0 to n
isPrime(n) => count (i for i in range(2, n) where n % i = 0) = 0
primes(n) => i for i in range(100) where isPrime(n)
```
