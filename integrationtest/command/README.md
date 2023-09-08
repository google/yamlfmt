# Command Integration Tests

These are tests that run a yamlfmt binary with different combos of commands in a temp directory set up by the test data. 

Each test runs by:
* Accepting the absolute path to a binary in the `YAMLFMT_BIN` environment variable
* Creating a temporary directory
* Copying everything from `before` in the testdata folder for the given test into the temp directory
* Run the specified command for the given test with the temp directory as the working directory
* Compare goldens for command output and state of the directory
    - If running with a `-update` flag, simply overwrite all golden files
    - If running normally, compare the golden files to ensure all the files are the same and the content of each file matches

You can run the tests by running `make integrationtest` which will build the binary and run the tests with it.