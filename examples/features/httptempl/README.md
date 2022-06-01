### Spreak + po + `http/template`

The following example shows how to use Spreak with the translation managed in `po` files and `html/template`
for a web project.
For an example with json files, see the [jhttptempl](../jhttptempl) example.

In the `locale/` folder you will find the files for the translations. 
The file `locale/httptempl.pot` can be used as a template for new translations.
In `templates/` you will find the template files which are delivered by the web server.

To start the web server just use `go run .`.

To update the `httptempl.pot` the following command can be used:
```shell
xspreak -D ./ -o locale/httptempl.pot -k ".i18n.Tr" -k "$.i18n.Tr" -t "templates/*.html"
```