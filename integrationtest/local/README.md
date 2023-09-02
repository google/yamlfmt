# Local Integration Tests

These are tests that can be run directly on the host machine. 

Each test runs by:
* Creating a temporary directory
* Copying everything from `before` in the testdata folder for the given test into the temp directory
* Run the specified command for the given test with the temp directory as the working directory
* Compare goldens for command output and state of the directory
    - If running with a `-update` flag, simply overwrite all golden files
    - If running normally, compare the golden files to ensure all the files are the same and the content of each file matches
