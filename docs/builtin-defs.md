# Trex Built-In Definitions

## Table of Contents:

1. [ascii](#ascii)
2. [bool](#bool)
3. [chars](#chars)
4. [count](#count)
5. [endswith](#endswith)
6. [fold](#fold)
7. [foldl](#foldl)
8. [foldr](#foldr)
9. [hasmatch](#hasmatch)
10. [indexby](#indexby)
11. [indexof](#indexof)
12. [isalnum](#isalnum)
13. [isalpha](#isalpha)
14. [isdigit](#isdigit)
15. [isletter](#isletter)
16. [islower](#islower)
17. [isnum](#isnum)
18. [isspace](#isspace)
19. [istitle](#istitle)
20. [isupper](#isupper)
21. [join](#join)
22. [lastindexby](#lastindexby)
23. [lastindexof](#lastindexof)
24. [len](#len)
25. [lines](#lines)
26. [matches](#matches)
27. [max](#max)
28. [min](#min)
29. [numoccurs](#numoccurs)
30. [replace](#replace)
31. [reverse](#reverse)
32. [sort](#sort)
33. [split](#split)
34. [startswith](#startswith)
35. [swapcase](#swapcase)
36. [tolower](#tolower)
37. [totitle](#totitle)
38. [toupper](#toupper)
39. [unique](#unique)
40. [words](#words)

## ascii

Returns a list of numbers, with every number representing the ASCII value of the corresponding character in the string.

Input: a string

Parameters: none

```
--> ascii 0123
48, 49, 50, 51
```

## bool

Returns 'true' if the input is true, otherwise 'false

Input: a string

Parameters: none

```
--> bool (1 = 2)
false
--> bool (12 > 4)
true
```

## chars

Splits a given string into a list of single characters.

Input: a string.

Parameters: none

```
--> chars 12343
1, 2, 3, 4, 3
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

## endswith

Checks whether a given string ends with a specified suffix.

Input: a string.

Parameters: 1

* The suffix

```
--> bool endswith('ab') 'kabab'
true
```

## fold

Applies a right fold to a list. Equivalent to 'foldr'.

Input: a list

Parameters: 1

* The definition by which to fold fold the values

```
--> fold(a,b -> a+b) (1, 2, 3, 4, 5)
15
```

## foldl

Applies a left fold to a list.

Input: a list

Parameters: 1

* The definition by which to fold fold the values

```
--> fold(a,b -> a+b) (1, 2, 3, 4, 5)
15
```

## foldr

Applies a right fold to a list.

Input: a list

Parameters: 1

* The definition by which to fold fold the values

```
--> fold(a,b -> a+b) (1, 2, 3, 4, 5)
15
```

## hasmatch

Finds whether a regular expression has a match whithin a string.

Input: a string

Parameters: 1

* The regular expression to match 

```
--> bool hasmatch('a[a-z]') "abbbjaja"
true
```

## indexby

Finds the index of the first character which satisfies the definition. Returns -1 if no character satisfies the definition.

Input: a string.

Parameters: 1

* The definition

```
--> indexby(->[] = 'a' or [] = 'b') "this is a string"
8
```

## indexof

Finds the index of the first instance of a substring. Returns -1 if the substring is not found.

Input: a string.

Parameters: 1

* The substring to find

```
--> indexof("s") "this is a string"
3
```

## isalnum

Checks whether if all characters in a string are alphanumeric and there is at least one character.

Input: a string.

Parameters: none

```
--> bool isalnum 'abc12'
true
--> bool isalnum 'ab$$1'
false
```

## isalpha

Checks if all characters in a string are alphabetic and there is at least one character.

Input: a string.

Parameters: none

```
--> bool isalpha 'abc12'
true
--> bool isalpha 'ab$$1'
false
```

## isdigit

Checks if a string is a single digit.

Input: a string

Parameters: none

```
--> bool isdigit 1
true
--> bool isdigit 'a'
false
--> bool isdigit 12
false
```

## isletter

Checks if a string is a single letter.

Input: a string

Parameters: none

```
--> bool isletter 1
false
--> bool isletter 'a'
true
--> bool isletter 'aa'
false
```

## islower

Checks if a string is comprised only of lowercase letters.

Input: a string

Parameters: none

```
--> bool islower 'A'
false
--> bool islower 'aa'
true
```

## isnum

Checks if all characters in a string are numeric and there is at least one character.

Input: a string.

Parameters: none

```
--> bool isnum 13
true
--> bool isnum 'ab'
false
```

## isspace

Checks if there are only whitespace characters in the string and there is at least one character

Input: a string.

Parameters: none

```
--> bool isspace '  '
true
```

## istitle

Checks if all words in a string begin with an uppercase letter and are otherwise are lowercase.

Input: a string.

Parameters: none

```
--> bool istitle 'Her Royal Highness'
true
```

## isupper

Checks if a string is comprised only of uppercase letters.

Input: a string

Parameters: none

```
--> bool isupper 'a'
false
--> bool isupper 'AA'
true
```

## join

Joins all elements in a list into a single string.

Input: a list

Parameters: none

```
--> join (1, 2, 3, 4, 5)
12345
```

## lastindexby

Finds the index of the last character which satisfies the definition. Returns -1 if no character satisfies the definition.

Input: a string.

Parameters: 1

* The definition

```
--> lastindexby(->[] = 'a' or [] = 'b') "kabab"
4
```

## lastindexof

Finds the index of the last instance of a substring. Returns -1 if the substring is not found.

Input: a string.

Parameters: 1

* The substring to find

```
--> lastindexof("s") "this is a string"
10
```

## len

Returns the length of a given string.

Input: a string.

Parameters: none

```
--> len "example"
7
```

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

## matches

Finds all matches of a regular expression whithin a string.

Input: a string

Parameters: 1

* The regular expression to match 

```
--> matches('a[a-z]') "abbbjaja"
ab, aj
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

## numoccurs

Returns the number of times a value occurs inside a given list or string.

Input: a list or string.

Parameters: 1

* the value to count occurences of

```
--> numoccurs('fo') 'foobafo'
2
```

## replace

Replaces all occurences of a certain string whithin a string with another string.

Input: a string

Parameters:

* the string to search for 

* the string to replace with

```
--> replace('a', 'AA') 'a bar'
AA bAAr
```

## reverse

Reverses a string or list.

Input: a string or list

Parameters: none

```
--> reverse (1, 2, 3, 4)
4, 3, 2, 1
--> reverse 1234
4321
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

## split

Splits a string into a list based on a seperator.

Input: a string.

Expected number of parameters: 1

* The seperator string

```
--> split(' ') "12 13 14 15"
12, 13, 14, 15
```

## startswith

Checks whether a given string starts with a specified prefix.

Input: a string.

Parameters: 1

* The prefix

```
--> bool startswith('tr') 'trex'
true
```

## swapcase

Swaps uppercase letters with their lowercase counterparts and vice versa. 

Input: a string

Parameters: none

```
--> swapcase "Her Royal Highness"
hER rOYAL hIGHNESS
```

## tolower

Returns the input with all unicode letters mapped to their lower case.

Input: a string.

Parameters: none

```
--> tolower "Hello World"
hello world
```

## totitle

Converts the letters at the beginning of each word to uppercase.

Input: a string

Parameters: none

```
--> totitle "her royal highness"
Her Royal Highness
```

## toupper

Returns the input with all unicode letters mapped to their upper case.

Input: a string.

Parameters: none

```
--> toupper "Hello World"
HELLO WORLD
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

## words

Splits a given string into words.

Input: a string.

Parameters: none

```
--> foo => "this is a sentence"
--> words foo
this, is, a, sentence
```

