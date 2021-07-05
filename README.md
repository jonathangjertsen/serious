# Serious serial utility

A terminal emulator like Tera Term and PuTTY, but made specifically for interfacing with
serial/COM ports, placing important configuration like port and baud rate front and center.
This makes it suited for more specialized tasks and scripting.

## How to get serious

There is currently no release for `serious` because it doesn't do anything yet.
If you have `git` and `go`, you can build it yourself (should work on all common platforms):

```sh
git clone https://github.com/jonathangjertsen/serious.git
cd serious
go install github.com/magefile/mage
go install .
mage build
```

The executable can then be found under `bin/<platform>-<arch>/serious`.
