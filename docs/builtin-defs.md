# Trex Built-In Definitions

## len

Returns the length of a given string.

Input: a string.

Parameters: none

```
--> len "example"
7
```

## count

Returns the number of values in a given list.

Input: a list.

Parameters: none

```
--> lines
one, two, three
--> count lines
3
```

## split

TODO - explanation for "split"

Input: a string.

Expected number of parameters: 1

## lines

Splits a given string into lines.

Input: a string.

Parameters: none

```
--> []
one
two
three
--> lines
one, two, three
```

## words

Splits a given string into words.

Input: a string.

Parameters: none

```
--> foo => "this is a sentence"
--> words foo
this, is, a, sentence
```

## chars

Splits a given string into a list of single characters.

Input: a string.

Parameters: none

```
--> chars 12343
1, 2, 3, 4, 3
```

## min

Finds the smallest value in a list based on a specified order.

Input: a list.

Parameters: 1
* the definition by which to order the values, which must return a value 

```
--> []
word
another
foo
--> min(#len)
foo
```

## max

Finds the largest value in a list based on a specified order.

Input: a list.

Parameters: 1
* the definition by which to order the values, which must return a value 

```
--> []
word
another
foo
--> max(#len)
another
```

## sort

Sorts a list (ascending) based on a specified order.

Input: a list.

Parameters: 1
* the definition by which to order the values, which must return a value convertible to a number.

```
--> words
one, three, four
--> sort(#len) words
one, four, three
```

## unique

Returns a list of all unique values in a given list.

Input: a list.

Parameters: none

```
--> foo => 1, 2, 3, 4, 4, 3, 2, 1, 3, 7
--> unique foo
1, 2, 3, 4, 7
```

## numOccurs

Returns the number of times a value occurs inside a given list or string.

Input: a list or string.

Parameters: 1
* the value to count occurences of

```
--> numOccurs('fo') 'foobafo'
2
```

## toUpper

Returns the input with all unicode letters mapped to their upper case.

Input: a string.

Parameters: none

```
--> toUpper "Hello World"
HELLO WORLD
```

## toLower

Returns the input with all unicode letters mapped to their lower case.

Input: a string.

Parameters: none

```
--> toLower "Hello World"
hello world
```