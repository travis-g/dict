package main

import (
	"bytes"
	"encoding/json"
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
		"Accept":  []string{"application/json"},
		"app_id":  []string{appID},
		"app_key": []string{appKey},
	}

	res, err := Client.Do(req)

	// error making request => nil response, so res can't be returned.
	if err != nil && res == nil {
		fmt.Fprintf(os.Stderr, "HTTP request error: %s\n", err)
		os.Exit(1)
	}

	if res.StatusCode != 200 {
		switch res.StatusCode {
		// TODO(travis-g): add more cases for human-friendly HTTP error notifications
		default:
			return res, fmt.Errorf("HTTP response code %d received", res.StatusCode)
		}
	}

	if err != nil {
		return res, err
	}

	return res, nil
}

func DefineCommand(word string) error {
	var (
		urlPattern = "%s/entries/%s/%s?fields=definitions,domains,examples,pronunciations,registers,etymologies"
	)

	url := fmt.Sprintf(urlPattern, BaseURL, viper.GetString("region"), strings.ToLower(word))

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

	data := result
	var buf bytes.Buffer
	if flagSimple {
		err = Templates["definition-short"].Execute(&buf, data)
	} else {
		err = Templates["definition"].Execute(&buf, data)
	}
	if err != nil {
		panic(err)
	}

	if flagHTML {
		page := HTMLOutput{
			Title:   data.Results[0].Word,
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
		urlPattern = "%s/thesaurus/%s/%s?fields=synonyms,antonyms"
	)

	url := fmt.Sprintf(urlPattern, BaseURL, viper.GetString("region"), strings.ToLower(word))

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
		urlPattern = "%s/thesaurus/%s/%s?fields=synonyms%%2Cantonyms"
	)

	url := fmt.Sprintf(urlPattern, BaseURL, viper.GetString("region"), strings.ToLower(word))

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
