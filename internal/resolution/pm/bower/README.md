# Bower resolution logic

The way resolution of bower lock files works is as follows:

1. Run `bower install --save --save-dev --save-exact --allow-root` in order to install all dependencies
2. Run `bower list` to get installed dependencies tree

The result of `bower list` command is then being written into the lock file.
