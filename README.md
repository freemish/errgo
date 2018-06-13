freemish/errgo
================
This package is a fork of go-errors/go and will change it in the following ways, and possibly others:

- Ability to easily access root error message
- Make stacktrace string (more) configurable by users, and change the default stacktrace
- Rather than encouraging the use of New, then Wrap with some number, just let users call Wrap and not worry about passing in a skip number, if they just want a quick stacktrace solution
- Some other stuff

This package is licensed under the MIT license, see LICENSE.MIT for details.
