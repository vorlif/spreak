
Here you will find a number of helpful examples of how you can use spreak. 

## Setup

Follow these setup to run the `helloworld` example:

1. Get the code

  ```shell
  $ git clone https://github.com/vorlif/spreak
  ```

2. Run the code

  ```shell
  cd spreak/examples/helloworld
  go run main.go
  ```

## What is included?

Unless otherwise required for the demonstration, all `.pot` files and the associated translations are
stored in the folder `locale`.

* [helloworld](./helloworld/main.go)
  * Simple example to demonstrate how spreak can be used
  * xspreak command `xspreak -e -D helloworld/ -p locale/ -d helloworld`

* [dayinfo](./dayinfo/main.go)
  * Advanced example to demonstrate how spreak can be used with xspreak
  * xspreak command `xspreak -e -D dayinfo/ -p locale/ -d dayinfo`

* [httptempl](features/httptempl): Example of how to use spreak with a web server and `http/templates`

* [jhttptempl](features/jhttptempl): Example of how to use spreak with a web server, json and `http/templates`

* [embed](./features/embed/main.go): Example how spreak can be used with the embed api

* [errors](./features/errors/main.go): Example how errors can be translated with spreak and xspreak

* [loaders](./features/loaders/main.go): Example how to load PO files from other sources

* [resolver](./features/resolver/main.go): Example how to resolve the path to a file with translations.

* [printer](./features/printer/main.go): Example how to use your own printer

* [decoder](./features/decoder): Example of implementation of a decoder and a catalog for importing JSON files for translation.
