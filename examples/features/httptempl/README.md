
The following example shows how to use Spreak with templates for web projects.

In `locale` you will find the files for the translations. 
The file `httptempl.pot` can be used as a template for new translations.
In `templates` you will find the template files which are delivered by the web server.


To update the `httptempl.pot` the following command can be used:
```shell
xspreak -D ./ -o locale/httptempl.pot -k ".i18n.Tr" -k "$.i18n.Tr" -t "templates/*.html"
```