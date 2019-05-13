# Integration tests

## Running integration tests

```console
$ go test -v ./tests
```

## Adding a new test

To add a test, drop a `test-$testname.js` file in this directory. It will be
automatically picked up and run using `jk` and its results will be compared
to the results we expect.

- The test will run with:

```console
$ jk run -o test-$testname.got test-$testname.js
```

- If the file `test-$testname.js.skip` exists, the test is skipped. This is
  useful to commit failing tests but not make them part of the test suite.

- If the file `test-$testname.js.error` exists, we'll check that `jk` exits
  with an error. Otherwise, we expect `jk` to exit with 0.

- If the file `test-$testname.js.cmd` exists, its content is used as the
  commands to excute for that tests. This allows to:

    1. Run several commands. Only the output of the jk command is compared to
       the `.expected` file.
    2. Use custom jk commands or options.
    3. Run js files that aren't in the `tests/` directory.

  `.cmd` files look like:

  ```text
  npm install
  jk run %b/test.js
  ```

  In that file, special variables can be used for convenience:

  **%f**: the name of test js file (eg. `test-foo.js`)

  **%b**: the test file base name (eg. `test-foo`)

  **%t**: the name of the test (eg. `foo`)

  **%d**: the name of the recommended output directory (eg. `test-foo.got`)


- `jk` std{out,err} will be compared to `test-$testname.js.expected`.

- If `jk` writes files to disk, they will be compared to the files in the
  `test-$testname.expected` directory.
