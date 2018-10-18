# dict

A CLI for querying [Oxford Dictionaries'](https://www.oxforddictionaries.com/) API. Functions and features:

- Look up a word's definition, synonyms, and antonyms,
- Manpage-like or simplified output formats,
- Retrive raw JSON API responses,
- Optional HTML page generation using [MdMe](https://github.com/susam/mdme), a JavaScript library for self-rendering Markdown documents.

## Installation

With a proper [Go environment](https://golang.org/doc/install) set up:

```console
$ go get -u github.com/travis-g/dict
```

Before using the CLI you'll need to retrieve an API key from the [Oxford Dictionaries developer portal](https://developer.oxforddictionaries.com/). You will also need to create a `.dictrc.yaml` file in your home directory (other filetypes work but are not documented). The CLI will fall back to checking the current working directory for a config file.

The `.dictrc.yaml` config file should have the following properties:

```yaml
app_id: <APP_ID>
app_key: <APP_KEY>
api_url: https://od-api.oxforddictionaries.com/api/v1
region: gb # or "us" for the New Oxford American Dictionary
```

Note that certain spellings are used in only certain regions (try "skillful" vs. "skilful" across regions).

## Usage

```console
$ dict -h
NAME:
   dict - Oxford Dictionary CLI tool

USAGE:
   dict [global options] command [command options] [arguments...]

COMMANDS:
     antonyms, a  look up antonyms
     define, d    define a word
     synonyms, s  look up synonyms
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

For example, to look up the definition of "command line":

```console
$ dict define "command line"
```

The subcommands have their own usage information.

## Notes

- This tool will probably be split into two components,one that will become a Golang library designed solely to implement the Oxford Dictionaries API, and a separate CLI with the fetching functionality.

## License

Source code is available under the [MIT license](/LICENSE).

Please be sure to read Oxford Dictionaries' [Terms and Conditions](https://developer.oxforddictionaries.com/api-terms-and-conditions) and their [FAQ](https://developer.oxforddictionaries.com/faq).

<img alt="Powered by OXFORD" src="https://developer.oxforddictionaries.com/images/PBO_black.png" height=50/>
