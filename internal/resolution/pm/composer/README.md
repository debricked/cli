# Composer resolution logic

The way resolution of composer lock files works is as follows:

1. Run `composer update --no-interaction --no-scripts --ignore-platform-reqs --no-autoloader --no-install --no-plugins --no-audit` in order to install all dependencies

Generated `composer.lock` file is then uploaded together with `composer.json` for scanning.
