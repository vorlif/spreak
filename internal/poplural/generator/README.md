# How to create the built-in plural functions?

1. Clone https://github.com/php-gettext/Languages
2. Execute ./bin/export-plural-rules prettyjson --output=plurals.json --parenthesis=yes --reduce=no
3. Copy the file `plurals.json` to this directory
5. Run `go run .`