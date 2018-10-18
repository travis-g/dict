package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

func makeRequest(url string) (*http.Response, error) {
	appID := viper.GetString("app_id")
	appKey := viper.GetString("app_key")

	if appID == "" {
		return nil, fmt.Errorf("Application ID not provided")
	}
	if appKey == "" {
		return nil, fmt.Errorf("API key not provided")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header = map[string][]string{
		"app_id":  []string{appID},
		"app_key": []string{appKey},
	}

	res, err := Client.Do(req)

	if res.StatusCode != 200 {
		return res, errors.New(fmt.Sprintf("HTTP response code %d received", res.StatusCode))
	}

	if err != nil {
		return res, err
	}

	return res, nil
}

func DefineCommand(word string) error {
	var (
		urlPattern = "%s/entries/en/%s/regions=%s"
	)

	url := fmt.Sprintf(urlPattern, BaseURL, strings.ToLower(word), viper.GetString("region"))

	res, err := makeRequest(url)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if res.StatusCode != http.StatusOK {
		return cli.NewExitError("No entry found.", 1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	result := Results{}
	json.Unmarshal(body, &result)
	if flagRaw {
		bytes, _ := json.Marshal(result)
		fmt.Printf("%s\n", string(bytes))
		return nil
	}

	if len(result.Results) == 0 {
		return cli.NewExitError("No entry found.", 1)
	}

	// value := gjson.Get(string(body), key)
	// s.FinalMSG = fmt.Sprintf("%s is %s\n", word, value)

	// result.SplitResultsByHomograph()

	data := result.Results[0]
	var buf bytes.Buffer
	if flagSimple {
		err = Templates["definition-simple"].Execute(&buf, data)
	} else {
		err = Templates["definition"].Execute(&buf, data)
	}
	if err != nil {
		panic(err)
	}

	if flagHTML {
		page := HTMLOutput{
			Title:   data.Word,
			Content: buf.String(),
		}
		// Reset the buffer to clear the original output
		buf.Reset()
		err = Templates["webpage"].Execute(&buf, page)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(os.Stdout, buf.String())
	} else {
		fmt.Fprintf(os.Stdout, buf.String())
	}

	return nil
}

func SynonymCommand(word string) error {
	var (
		urlPattern = "%s/entries/en/%s/synonyms;antonyms"
	)

	url := fmt.Sprintf(urlPattern, BaseURL, strings.ToLower(word))

	res, err := makeRequest(url)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	result := Results{}
	json.Unmarshal(body, &result)
	if flagRaw {
		bytes, _ := json.Marshal(result)
		fmt.Printf("%s\n", string(bytes))
		return nil
	}

	if len(result.Results) == 0 {
		return cli.NewExitError("No entry found.", 1)
	}

	// values := gjson.Get(string(body), key).Array()
	// results := make([]string, len(values))
	// for i, result := range values {
	// 	results[i] = result.String()
	// }

	// s.FinalMSG = fmt.Sprintf("%s\n", strings.Join(results, "\n"))

	// result.SplitResultsByHomograph()

	data := result.Results[0]
	var buf bytes.Buffer
	if flagSimple {
		err = Templates["synonyms-simple"].Execute(&buf, data)
	} else {
		err = Templates["synonyms"].Execute(&buf, data)
	}
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stdout, buf.String())

	return nil
}

func AntonymCommand(word string) error {
	var (
		urlPattern = "%s/entries/en/%s/synonyms;antonyms"
	)

	url := fmt.Sprintf(urlPattern, BaseURL, strings.ToLower(word))

	res, err := makeRequest(url)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	result := Results{}
	json.Unmarshal(body, &result)

	if flagRaw {
		bytes, _ := json.Marshal(result)
		fmt.Printf("%s\n", string(bytes))
		return nil
	}

	if len(result.Results) == 0 {
		return cli.NewExitError("No entry found.", 1)
	}

	data := result.Results[0]
	var buf bytes.Buffer
	if flagSimple {
		err = Templates["antonyms-simple"].Execute(&buf, data)
	} else {
		err = Templates["antonyms"].Execute(&buf, data)
	}
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, buf.String())

	return nil
}