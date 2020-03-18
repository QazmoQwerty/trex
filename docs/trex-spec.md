# Trex Language Specification

## Notation
The syntax is specified using [Extended Backus-Naur Form](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form) (EBNF):

**Usage**       | **Notation**
------------    | -------------
definition      |   =
concatenation   |   ,
termination     |   ;
optional        |   [ ... ]
repetition      |   { ... }
grouping        |   ( ... )
terminal string |   " ... "
comment         |   (\* ... \*)
special sequence|   ? ... ?


## Lexical Elements

```EBNF
all_characters  = ? all visible characters ? ;
terminator      = '|' | '\n';
white_space     = ? whitespace characters ? ;
letter          = 'a'...'z' | 'A'...'Z' ;
digit           = '0'...'9';
```


### Whitespace

Whitespace has a purpose in the language's syntax and as such will *not* be discarded.

###  Terminators

The character ‘|’ is used as a terminator (similar to ‘;’ in other languages).

### Newlines

A newline is converted into a terminator is the token before it is one of the following:

* an *identifier*
* a *literal*
* one of the operators: ')', ']', or '}'

### Comments

There are two types of comments: 

* Single-line comments start with ‘//’ and end at the end of the line.
* Multi-line comments start with ‘/\*’ and end with ‘\*/’.

### Identifiers

```EBNF 
Identifier = Letter | "_", { Letter | Digit | "_" };
```

Identifiers name definitions. An identifier is a sequence of one or more letters and digits. The first character in an identifier must be a letter. 

```
name  
_a12  
IdentA_13
```

### Operators

The following character sequences are turned into operators:

```EBNF
+   -   *   /   %       (   )
#   ,   :   =>  <<  **  {   }
=   !=  <=  <   >   >=  [   ]
not and or  if  for in  else
```

### Literals

```EBNF
literal = string_literal | number_literal | character_literal;
```
```EBNF
string_literal =    ('"' { all_characters } '"') | ("'" { all_characters } "'");
number_literal =    digit { digit };
character_literal = '\t' | '\n' | '\r' (*TODO - other escaped chars*);
```

All literals are treated as strings. There are 3 types of literals:

1. String literals
```
"hello"
'abc'
"string literals can 
span multiple lines"
```
2. Number literals
```
123
```

3. Character literals
```
\n
\t
```

## Programs

```EBNF
Program     = ProgramLine { Terminator ProgramLine };
ProgramLine = Statement | Expression;
```

A program consists of *definitions* and *expressions*. Each program gets an input string (called the *argument*) and outputs a *Value*.  

Entering a top-level expression will cause the program to output it's value. Entering another top-level expression will cause the program to output a newline + the expression (note: if the value is not a string it will be converted into one).

```
>>> a => 10 + 20
>>> a
30
>>> a << 'abc'
30abc
```

## Definitions

```EBNF
Definition      = identifier [Parameters] DefinitionBody;
DefinitionBody  = ( "=>" Expression ) | ( '{' Program '}' );
Parameters      = '(' IdentifierList ')' ;
IdentifierList  = identifier { ',' identifier };
```

Definitions bind a *program* to an identifier. In addition to the program's argument, a definition may also specify a list of *parameters* to be passed when the definition is called.

```
>>> bar(a) => a ** 2
>>> foo(a, b) {
...     bar(a)
...     b ** 3
... }
>>> foo('a', 'b')
aa
bbb
```

## Subscripts

```EBNF
Subscript = [Expression] '[' [Expression] [':' [ Expression ] ] [':' Expression] ']';
```

Subscript expressions construct a substring or list from a string or a list.

The expression upon which the subscript acts may be omitted, in which case it will default to the argument value.

A subscript can get 0-3 indices, which must all evaluate to a string which is convertible to a number.

**0 indices**

Does nothing. This is used to express the value of the argument string.

```
>>> foo => []
>>> foo "abc"
abc
```

**1 index:**

Will fetch the value at said index:

```
>>> '0123456'[4]
4
>>> a => 'aa', 'bb', 'cc', 'dd'
>>> a[2]
cc
```

**2 indices:**

Will fetch the values from the first index to the last index:

```
>>> "0123456789"[3:7]
3456
```

**3 indices:**

The third index specifies to only treat the value every *n* values. If the third index is negative, the string/list will be flipped.

```
>>> "0123456789"[::2]
02468
>>> "01234"[::-1]
43210
>>> "0123456789"[::-2]
97531
```

Any of the 3 indices may be ommitted. A missing index is equivalent to the value "". A missing first index defaults to 0, a missing second index defaults to the length of the value, and a missing third index defaults to 1.

A negative index (as one of the first two indices) counts from the back of the string/list:

```
>>> "01234"[-1]
4
>>> "01234"[:-1]
0123
```

## Lists

```EBNF
ExpressionList = Expression { ',' ExpressionList };
```

Lists are constructed using the ',' operator.

```
>>> 1, 2, 3, 4
[1, 2, 3, 4]
>>> (1, 2, 3, 4)[0]
1
```

## Calls

```EBNF
Call    = Callee [ Params ] [ WS Expression ];
Callee  = Expression;
Params  = '(' Expression ')';
Arg     = Expression;
```

Calls evaluate the program specified in a *definition*. The *callee* is the definition to be called, which must be a definition value.

Calls may specify *parameters* to be bound to the identifiers specified in the definition. In addition Calls pass an *argument* to the program (if the argument is omitted it defaults to the current argument value).

```
// Input string is "foobar"
>>> a(n) => [:n]
>>> a(3)
foo
>>> a(2) "abcd"
ab
```

## Conditionals

```EBNF
Conditional = Expression "if" Expression "else" Expression;
```

First the middle expression (the condition) is evaluated. If it evaluates to a true value (any non-empty string) then the the left-most expression is evaluated and returned by the conditional. Otherwise if the *condition* evaluated to *false* (and empty string - ""), the right-most value is evaluated instead.

```
>>> 'a' if 12 < 13 else 'b'
a
>>> 1 if '' else 1 + 1
2
```

## Recursion

Programs can call themselves:

```
>>> factorial(n) => 1 if n = 0 else n * fact(n - 1)
>>> factorial(4)
24
```

## Binary Operators

The following are binary operators:

```
+   -   *   /   =   !=  or  and
<<  **  %   <   >   >=  <=
```

**Operator**    |   **Usage**
--------------  |   -------
\+              |   numeric addition
\-              |   numeric subtraction
\*              |   numeric multiplication
/               |   numeric division
%               |   numeric remainder
<<              |   string addition
**              |   string multiplication
<               |   smaller than
<=              |   smaller or equal to 
>               |   greater than
>=              |   greater or equal to
=               |   equal
!=              |   not equal
and             |   logical and
or              |   logical or

The operands of the numeric operators, as well as the right operand of the string multiplication operator, MUST be convertible to a number.

The comparision operators will compare the strings based on lexical order.

NOTE: all binary operations only accept strings as their operands.

## Unary Operators

The following are unary operators:

```
+   -   #   not
```

**Operator**    |   **Usage**               
--------------  |   ----------------------  
not             |   logical not        
\#              |   definition reference (accepts identifiers)     
\+              |   numeric unary +  
\-              |   numeric unary -

The '#' operator is used to pass definitions without calling them.

```
>>> foo(a) => a ** 2
>>> bar(f) => f('aa')
>>> bar(#foo)
aaaa
>>> bar(foo)  // foo will be called with no parameters causing an error
```