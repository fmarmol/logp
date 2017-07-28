# Documentation

A simple program that logs inputs from a pipe to the standard output

## Install

```
go get github.com/fmarmol/logp
```

## Example

```bash
python -c "print('start of the function');sum(range(1000));print('end of the function')" | logp
# output:
# 2017/07/28 10:10:17 start of the function
# 2017/07/28 10:10:17 end of the function. Time elapsed since last message: 3.96µs
```
