# An API for managing Docker containers

This is a simplified and stripped down Go API for managing Docker containers.
It provides basic functionality for creating, inspecting, and removing containers, as well as managing images and networks.

The API is mostly compatible with the official [Docker API], but has been simplified somewhat.
That is, some of the API calls have slightly different method signatures.
Moreover, many methods, options and fields have been removed.

We will add features as the need arises.
However, the API is not intended to be a complete implementation of the Docker API.
Moreover, we may also change the API without warning.
In particular, we may make the API more compatible with the official Docker API in the future.

Feel free to open an issue or a pull request if you find something missing or have suggestions for improvements.

## Why not just use the official Docker API?

The official [Docker API] brings in a lot of dependencies, which we do not want.
This API is a simplified version that only includes the client functionality we needed for two of our projects.

## Credits

Some of the struct definitions and method signatures have been copied (and modified) from the official Docker API types.

Hein Meling wrote the rest of the code in this package.

[Docker API]: https://pkg.go.dev/github.com/docker/docker
