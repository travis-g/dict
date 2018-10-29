package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var (
	Version = "unknown"
	BaseURL = "https://od-api.oxforddictionaries.com/api/v1"
	Client  *http.Client

	flagRaw      = false
	flagCategory string
	flagCache    = false
	flagSimple   = false
	flagHTML     = false
)

func loadConfig() {
	viper.SetConfigName(".dictrc")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

func setupRequestCommand(c *cli.Context) error {
	// load the Viper config
	loadConfig()

	// override the default http.Client timeouts
	Client = &http.Client{
		Timeout: time.Second * 10,
	}

	BaseURL = viper.GetString("api_url")

	if flagCategory != "" {
		flagCategory = fmt.Sprintf(`#[lexicalCategory=="%s"]`, strings.Title(flagCategory))
	} else {
		flagCategory = "0"
	}
	return nil
}

func main() {
	dict := cli.NewApp()
	dict.Name = "dict"
	dict.Usage = "Oxford Dictionary CLI tool"
	dict.Version = Version

	apiRequestFlags := &[]cli.Flag{
		cli.StringFlag{
			Name:        "category, c",
			Usage:       `(unimplemented) request info for a specific lexical category`,
			Destination: &flagCategory,
		},
		cli.BoolFlag{
			Name:        "html",
			Usage:       "wrap output in an HTML document with MdMe",
			Destination: &flagHTML,
		},
		// cli.BoolFlag{
		// 	Name:        "no-cache, f",
		// 	Usage:       "force an API request, updating the cache",
		// 	Destination: &flagCache,
		// },
		cli.BoolFlag{
			Name:        "raw, r",
			Usage:       "return raw JSON from the request",
			Destination: &flagRaw,
		},
		cli.BoolFlag{
			Name:        "simple, s",
			Usage:       "output only the values, separated by newlines (excludes subsenses)",
			Destination: &flagSimple,
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "(unimplemented) increase definition verbosity",
		},
	}

	dict.Commands = []cli.Command{
		cli.Command{
			Name:    "antonyms",
			Aliases: []string{"a"},
			Usage:   "look up antonyms",
			Flags:   *apiRequestFlags,
			Before:  setupRequestCommand,
			Action: func(c *cli.Context) error {
				return AntonymCommand(c.Args().Get(0))
			},
		},
		cli.Command{
			Name:    "define",
			Aliases: []string{"d"},
			Usage:   "define a word",
			Flags:   *apiRequestFlags,
			Before:  setupRequestCommand,
			Action: func(c *cli.Context) error {
				return DefineCommand(c.Args().Get(0))
			},
		},
		cli.Command{
			Name:    "synonyms",
			Aliases: []string{"s"},
			Usage:   "look up synonyms",
			Flags:   *apiRequestFlags,
			Before:  setupRequestCommand,
			Action: func(c *cli.Context) error {
				return SynonymCommand(c.Args().Get(0))
			},
		},
	}

	err := dict.Run(os.Args)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
