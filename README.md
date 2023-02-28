## mnms -- minimalist NMS

basic scan, config and multiple node cluster support with cli and api

### Do not directly check into main branches

must provide PR from a branch and get others to review and approve before merging to main

### run make to build

keep it simple. just use Makefile

### run go test

write as many  tests as go test. If you have added code without go tests then you have added more debt.  Do not add problems. Add solutions. Solutions must be focused, simple and have extensive tests.

Be sure to run

```
sudo go test -v -p 1
```

and all tests pass.

On Windows, run `go test -v -p 1` in a command shell or powershell in administrator mode.



### do not check in binary files or large files

don't make git slow down because of large binary files

### keep everything as simple and small as possible

don't drag in any dependency if not absolutely required. Remember: less is more.  Keep It Simple and Stupid.


### keep your worklogs up to date

https://github.com/Atop-NMS-team/Worklogs

keep a window open with editor so updating worklog daily is as easy as possible.
Commit your worklog daily.
