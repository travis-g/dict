# dict

[![GitHub](https://img.shields.io/github/license/travis-g/dict.svg)](https://github.com/travis-g/dict/blob/master/LICENSE)

A CLI for querying [Oxford Dictionaries'](https://www.oxforddictionaries.com/) API. Functions and features:

- Look up a word's definition, synonyms, and antonyms,
- Manpage-like or simplified output formats,
- Retrieve raw JSON API responses,
- Optional HTML page generation using [Blackfriday](https://github.com/russross/blackfriday), a Markdown renderer.

## V2 API Changes

V2 of Oxford Dictionaries' API removes free access to the `thesaurus` endpoint. Synonym and antonym lookups will fail with a Prototype tier plan or lower.

## Installation

With a proper [Go environment](https://golang.org/doc/install) set up:

```console
$ go get -u -v github.com/travis-g/dict
github.com/travis-g/dict (download)
...
```

Before using the CLI you'll need to retrieve an API key from the [Oxford Dictionaries developer portal](https://developer.oxforddictionaries.com/). You will also need to create a `.dictrc.yaml` file in your home directory (other filetypes work but are not documented). The CLI will fall back to checking the current working directory for a config file.

The `.dictrc.yaml` config file should have the following properties:

```yaml
app_id: <APP_ID>
app_key: <APP_KEY>
api_url: https://od-api.oxforddictionaries.com/api/v2
region: en-gb # or "en-us" for the New Oxford American Dictionary
```

Note that certain spellings are used in only certain regions (try "skillful" vs. "skilful" across regions).

## Usage

```console
$ dict -h
NAME:
   dict - Oxford Dictionary CLI tool

USAGE:
   dict [global options] command [command options] [arguments...]

VERSION:
   ...

COMMANDS:
     antonyms, a  look up antonyms
     define, d    define a word
     synonyms, s  look up synonyms
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

For example, to look up the definition of "console":

```console
$ dict define "console"
console
=======

...
```

The raw JSON response object can be output with the `--raw` or `-r` flag:

```console
$ dict define --raw console
{"metadata":{"provider":"Oxford University Press"},"results":[...]}
```

The subcommands have their own usage information.

## Notes

- If you're concerned about burning an API call on a word that may not exist, try `grep`ing through local wordlist files, possibly in [`/usr/share/dict/`](https://en.wikipedia.org/wiki/Words_(Unix)), and then make the request.
- See also:
  - [dictionary network protocol](https://en.wikipedia.org/wiki/DICT)
  - [words (Unix)](https://en.wikipedia.org/wiki/Words_(Unix))

## License

Source code is available under the [MIT license](/LICENSE).

Please be sure to read Oxford Dictionaries' [Terms and Conditions](https://developer.oxforddictionaries.com/api-terms-and-conditions) and their [FAQ](https://developer.oxforddictionaries.com/faq).

<img alt="Powered by OXFORD" src="https://developer.oxforddictionaries.com/images/PBO_black.png" height=50/>
