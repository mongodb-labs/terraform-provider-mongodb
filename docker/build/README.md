# Build tooling

This is a docker container which can be used to build the project for all supported platforms, without needing to
install Golang on your current system.

To run, simply execute `make build` in this directory and once finished check the 
[out](../../out) directory in the project's root.

If this is the first time you are running these tools, execute `make build-image` first.
