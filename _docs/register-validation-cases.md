# Register Validation – All Possible Error Cases

> This list is calculated largely using a [combination permutations calculator](https://www.mathsisfun.com/combinatorics/combinations-permutations-calculator.html).
> Some edge/special cases do exist (eg. if "empty" error occurs, no other error can). Some of these are thought of and the combinations are adjusted accordingly, and the rest will be dealth with while working on them as it is not practical to try to think of each and every single one of them upfront.

## Username

### 1 error
- [X] empty
- [X] too_short
- [X] too_long
- [X] invalid_char
- [X] digit_start

### 2 errors
- [X] too_short,invalid_char
- [X] too_short,digit_start
- [X] too_long,invalid_char
- [X] too_long,digit_start
- [X] invalid_char,digit_start

### 3 errors
- [X] too_short,invalid_char,digit_start
- [X] too_long,invalid_char,digit_start


## Password

### 1 error
- [X] empty 
- [X] too_short 
- [X] too_long 
- [X] no_lowercase
- [X] no_uppercase
- [X] no_digits
- [X] no_symbols
- [X] has_spaces
- [X] non_ascii

### 2 errors
- [X] too_short,no_lowercase
- [X] too_short,no_uppercase
- [X] too_short,no_digits
- [X] too_short,no_symbols
- [X] too_short,has_spaces
- [X] too_short,non_ascii
- [X] too_long,no_lowercase
- [X] too_long,no_uppercase
- [X] too_long,no_digits
- [X] too_long,no_symbols
- [X] too_long,has_spaces
- [X] too_long,non_ascii
- [X] no_lowercase,no_uppercase
- [X] no_lowercase,no_digits
- [X] no_lowercase,no_symbols
- [X] no_lowercase,has_spaces
- [X] no_lowercase,non_ascii
- [X] no_uppercase,no_digits
- [X] no_uppercase,no_symbols
- [X] no_uppercase,has_spaces
- [X] no_uppercase,non_ascii
- [X] no_digits,no_symbols
- [X] no_digits,has_spaces
- [X] no_digits,non_ascii
- [X] no_symbols,has_spaces
- [X] no_symbols,non_ascii
- [X] has_spaces,non_ascii

### 3 errors
- [X] too_short,no_lowercase,no_uppercase
- [X] too_short,no_lowercase,no_digits
- [X] too_short,no_lowercase,no_symbols
- [X] too_short,no_lowercase,has_spaces
- [X] too_short,no_lowercase,non_ascii
- [X] too_short,no_uppercase,no_digits
- [X] too_short,no_uppercase,no_symbols
- [X] too_short,no_uppercase,has_spaces
- [X] too_short,no_uppercase,non_ascii
- [X] too_short,no_digits,no_symbols
- [X] too_short,no_digits,has_spaces
- [X] too_short,no_digits,non_ascii
- [X] too_short,no_symbols,has_spaces
- [X] too_short,no_symbols,non_ascii
- [X] too_short,has_spaces,non_ascii
- [X] too_long,no_lowercase,no_uppercase
- [X] too_long,no_lowercase,no_digits
- [X] too_long,no_lowercase,no_symbols
- [X] too_long,no_lowercase,has_spaces
- [X] too_long,no_lowercase,non_ascii
- [X] too_long,no_uppercase,no_digits
- [X] too_long,no_uppercase,no_symbols
- [X] too_long,no_uppercase,has_spaces
- [X] too_long,no_uppercase,non_ascii
- [X] too_long,no_digits,no_symbols
- [X] too_long,no_digits,has_spaces
- [X] too_long,no_digits,non_ascii
- [X] too_long,no_symbols,has_spaces
- [X] too_long,no_symbols,non_ascii
- [X] too_long,has_spaces,non_ascii
- [X] no_lowercase,no_uppercase,no_digits
- [X] no_lowercase,no_uppercase,no_symbols
- [X] no_lowercase,no_uppercase,has_spaces
- [X] no_lowercase,no_uppercase,non_ascii
- [X] no_lowercase,no_digits,no_symbols
- [X] no_lowercase,no_digits,has_spaces
- [X] no_lowercase,no_digits,non_ascii
- [X] no_lowercase,no_symbols,has_spaces
- [X] no_lowercase,no_symbols,non_ascii
- [X] no_lowercase,has_spaces,non_ascii
- [X] no_uppercase,no_digits,no_symbols
- [X] no_uppercase,no_digits,has_spaces
- [X] no_uppercase,no_digits,non_ascii
- [X] no_uppercase,no_symbols,has_spaces
- [X] no_uppercase,no_symbols,non_ascii
- [X] no_uppercase,has_spaces,non_ascii
- [X] no_digits,no_symbols,has_spaces
- [X] no_digits,no_symbols,non_ascii
- [X] no_digits,has_spaces,non_ascii
- [X] no_symbols,has_spaces,non_ascii

### 4 errors
- [X] too_short,no_lowercase,no_uppercase,no_digits
- [X] too_short,no_lowercase,no_uppercase,no_symbols
- [X] too_short,no_lowercase,no_uppercase,has_spaces
- [X] too_short,no_lowercase,no_uppercase,non_ascii
- [X] too_short,no_lowercase,no_digits,no_symbols
- [X] too_short,no_lowercase,no_digits,has_spaces
- [X] too_short,no_lowercase,no_digits,non_ascii
- [X] too_short,no_lowercase,no_symbols,has_spaces
- [X] too_short,no_lowercase,no_symbols,non_ascii
- [X] too_short,no_lowercase,has_spaces,non_ascii
- [X] too_short,no_uppercase,no_digits,no_symbols
- [X] too_short,no_uppercase,no_digits,has_spaces
- [X] too_short,no_uppercase,no_digits,non_ascii
- [X] too_short,no_uppercase,no_symbols,has_spaces
- [X] too_short,no_uppercase,no_symbols,non_ascii
- [X] too_short,no_uppercase,has_spaces,non_ascii
- [X] too_short,no_digits,no_symbols,has_spaces
- [X] too_short,no_digits,no_symbols,non_ascii
- [X] too_short,no_digits,has_spaces,non_ascii
- [X] too_short,no_symbols,has_spaces,non_ascii
- [X] too_long,no_lowercase,no_uppercase,no_digits
- [X] too_long,no_lowercase,no_uppercase,no_symbols
- [X] too_long,no_lowercase,no_uppercase,has_spaces
- [X] too_long,no_lowercase,no_uppercase,non_ascii
- [X] too_long,no_lowercase,no_digits,no_symbols
- [X] too_long,no_lowercase,no_digits,has_spaces
- [X] too_long,no_lowercase,no_digits,non_ascii
- [X] too_long,no_lowercase,no_symbols,has_spaces
- [X] too_long,no_lowercase,no_symbols,non_ascii
- [X] too_long,no_lowercase,has_spaces,non_ascii
- [X] too_long,no_uppercase,no_digits,no_symbols
- [X] too_long,no_uppercase,no_digits,has_spaces
- [X] too_long,no_uppercase,no_digits,non_ascii
- [X] too_long,no_uppercase,no_symbols,has_spaces
- [X] too_long,no_uppercase,no_symbols,non_ascii
- [X] too_long,no_uppercase,has_spaces,non_ascii
- [X] too_long,no_digits,no_symbols,has_spaces
- [X] too_long,no_digits,no_symbols,non_ascii
- [X] too_long,no_digits,has_spaces,non_ascii
- [X] too_long,no_symbols,has_spaces,non_ascii
- [ ] ~~no_lowercase,no_uppercase,no_digits,no_symbols~~
  - impossible case
- [X] no_lowercase,no_uppercase,no_digits,has_spaces
- [X] no_lowercase,no_uppercase,no_digits,non_ascii
- [X] no_lowercase,no_uppercase,no_symbols,has_spaces
- [X] no_lowercase,no_uppercase,no_symbols,non_ascii
- [X] no_lowercase,no_uppercase,has_spaces,non_ascii
- [X] no_lowercase,no_digits,no_symbols,has_spaces
- [X] no_lowercase,no_digits,no_symbols,non_ascii
- [X] no_lowercase,no_digits,has_spaces,non_ascii
- [X] no_lowercase,no_symbols,has_spaces,non_ascii
- [X] no_uppercase,no_digits,no_symbols,has_spaces
- [X] no_uppercase,no_digits,no_symbols,non_ascii
- [X] no_uppercase,no_digits,has_spaces,non_ascii
- [X] no_uppercase,no_symbols,has_spaces,non_ascii
- [X] no_digits,no_symbols,has_spaces,non_ascii

### 5 errors
- [ ] too_short,no_lowercase,no_uppercase,no_digits,no_symbols
- [ ] too_short,no_lowercase,no_uppercase,no_digits,has_spaces
- [ ] too_short,no_lowercase,no_uppercase,no_digits,non_ascii
- [ ] too_short,no_lowercase,no_uppercase,no_symbols,has_spaces
- [ ] too_short,no_lowercase,no_uppercase,no_symbols,non_ascii
- [ ] too_short,no_lowercase,no_uppercase,has_spaces,non_ascii
- [ ] too_short,no_lowercase,no_digits,no_symbols,has_spaces
- [ ] too_short,no_lowercase,no_digits,no_symbols,non_ascii
- [ ] too_short,no_lowercase,no_digits,has_spaces,non_ascii
- [ ] too_short,no_lowercase,no_symbols,has_spaces,non_ascii
- [ ] too_short,no_uppercase,no_digits,no_symbols,has_spaces
- [ ] too_short,no_uppercase,no_digits,no_symbols,non_ascii
- [ ] too_short,no_uppercase,no_digits,has_spaces,non_ascii
- [ ] too_short,no_uppercase,no_symbols,has_spaces,non_ascii
- [ ] too_short,no_digits,no_symbols,has_spaces,non_ascii
- [ ] too_long,no_lowercase,no_uppercase,no_digits,no_symbols
- [ ] too_long,no_lowercase,no_uppercase,no_digits,has_spaces
- [ ] too_long,no_lowercase,no_uppercase,no_digits,non_ascii
- [ ] too_long,no_lowercase,no_uppercase,no_symbols,has_spaces
- [ ] too_long,no_lowercase,no_uppercase,no_symbols,non_ascii
- [ ] too_long,no_lowercase,no_uppercase,has_spaces,non_ascii
- [ ] too_long,no_lowercase,no_digits,no_symbols,has_spaces
- [ ] too_long,no_lowercase,no_digits,no_symbols,non_ascii
- [ ] too_long,no_lowercase,no_digits,has_spaces,non_ascii
- [ ] too_long,no_lowercase,no_symbols,has_spaces,non_ascii
- [ ] too_long,no_uppercase,no_digits,no_symbols,has_spaces
- [ ] too_long,no_uppercase,no_digits,no_symbols,non_ascii
- [ ] too_long,no_uppercase,no_digits,has_spaces,non_ascii
- [ ] too_long,no_uppercase,no_symbols,has_spaces,non_ascii
- [ ] too_long,no_digits,no_symbols,has_spaces,non_ascii
- [ ] no_lowercase,no_uppercase,no_digits,no_symbols,has_spaces
- [ ] no_lowercase,no_uppercase,no_digits,no_symbols,non_ascii
- [ ] no_lowercase,no_uppercase,no_digits,has_spaces,non_ascii
- [ ] no_lowercase,no_uppercase,no_symbols,has_spaces,non_ascii
- [ ] no_lowercase,no_digits,no_symbols,has_spaces,non_ascii
- [ ] no_uppercase,no_digits,no_symbols,has_spaces,non_ascii

### 6 errors
- [ ] too_short,no_lowercase,no_uppercase,no_digits,no_symbols,has_spaces
- [ ] too_short,no_lowercase,no_uppercase,no_digits,no_symbols,non_ascii
- [ ] too_short,no_lowercase,no_uppercase,no_digits,has_spaces,non_ascii
- [ ] too_short,no_lowercase,no_uppercase,no_symbols,has_spaces,non_ascii
- [ ] too_short,no_lowercase,no_digits,no_symbols,has_spaces,non_ascii
- [ ] too_short,no_uppercase,no_digits,no_symbols,has_spaces,non_ascii
- [ ] too_long,no_lowercase,no_uppercase,no_digits,no_symbols,has_spaces
- [ ] too_long,no_lowercase,no_uppercase,no_digits,no_symbols,non_ascii
- [ ] too_long,no_lowercase,no_uppercase,no_digits,has_spaces,non_ascii
- [ ] too_long,no_lowercase,no_uppercase,no_symbols,has_spaces,non_ascii
- [ ] too_long,no_lowercase,no_digits,no_symbols,has_spaces,non_ascii
- [ ] too_long,no_uppercase,no_digits,no_symbols,has_spaces,non_ascii
- [ ] no_lowercase,no_uppercase,no_digits,no_symbols,has_spaces,non_ascii

### 7 errors
- [ ] too_short,no_lowercase,no_uppercase,no_digits,no_symbols,has_spaces,non_ascii
- [ ] too_long,no_lowercase,no_uppercase,no_digits,no_symbols,has_spaces,non_ascii