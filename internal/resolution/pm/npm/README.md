# NPM resolution logic

The way resolution of NPM lock files works is as follows:

1. Run `npm install --ignore-scripts --audit=false --bin-links=false` in order to install all dependencies

Generated `package-lock.json` file is then uploaded together with `package.json` for scanning.
