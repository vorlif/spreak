### Spreak + json + `http/template`

The following example shows how to use Spreak with the translation managed in key-value JSON files and `html/template`
for a web project.
For an example with po files, see the [httptempl](../httptempl) example.

In the `locale/` folder you will find the files for the translations.
The file `locale/httptempl.json` can be used as a template for new translations.
In `templates/` you will find the template files which are delivered by the web server.

To start the web server just use `go run .`.

1. To update the `httptempl.json` the following command can be used: 
    ```shell
    xspreak -D ./ -o locale/httptempl.json -k ".i18n.Tr" -k "$.i18n.Tr" -k ".i18n.TrN:1,1" -t "templates/*.html" -f json --template-use-kv
    ```
   
2. To create a translation file for a new language use
    ```shell
    xspreak merge -i locale/httptempl.json -o locale/de.json -l de
    # alternative with copy of english translations
    xspreak merge -i locale/en.json -o locale/de.json -l de
    ```
   
3. To add new keys for an existing language use:
   ```shell
    xspreak merge -i locale/httptempl.json -o locale/de.json -l de
    # alternative with copy of english translations
    xspreak merge -i locale/en.json -o locale/de.json -l de
    ```