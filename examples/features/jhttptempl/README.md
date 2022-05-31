The following example shows how to use Spreak with the translation managed in key-value JSON files and `html/template`
for a web project.

In the `locale/` folder you will find the files for the translations.
The file `locale/httptempl.json` can be used as a template for new translations.
In `templates/` you will find the template files which are delivered by the web server.

To start the web server just use `go run .`.

To update the `httptempl.json` the following command can be used:

```shell
xspreak -D ./ -o locale/httptempl.json -k ".i18n.Tr" -k "$.i18n.Tr" -k ".i18n.TrN:1,1" -t "templates/*.html" -f json --template-use-kv
```