# Yarn resolution logic

The way resolution of yarn lock files works is as follows:

1. Run `install --non-interactive --ignore-scripts --ignore-engines --ignore-platform --no-bin-link --production=false` in order to install all dependencies

Generated `yarn.lock` file is then uploaded together with `package.json` for scanning.
