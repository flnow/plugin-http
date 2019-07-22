package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "plugin-http",
	Short: "A plugin for Flnow to send HTTP request.",
	Long:  `A plugin for Flnow to send HTTP request.`,
	Run:   run,
}

// Execute command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	testF := pflag.Flag{}
	testF.Name = "testaaa"
	// testF.Shorthand = "tt"
	testF.Usage = "this is test usage"
	testF.DefValue = "1"
	rootCmd.Flags().AddFlag(&testF)

	rootCmd.Flags().StringP("method", "M", "GET", "HTTP Method, default to GET")
	rootCmd.Flags().StringP("url", "U", "", "Call URL, start with http:// or https://")
	rootCmd.Flags().StringSlice("headers", []string{}, "Request Headers, e.g. color=black")
	rootCmd.Flags().StringSliceP("cookies", "C", []string{}, "Request with cookies, e.g. name=jack")
	rootCmd.Flags().StringSliceP("params", "P", []string{}, "Reuqest with params, form/url params/post data")
	rootCmd.MarkFlagRequired("url")
}

func initConfig() {
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func run(cmd *cobra.Command, args []string) {
	// required flags
	method := cmd.Flag("method")
	callURL := cmd.Flag("url")
	postParam := url.Values{}
	params, _ := cmd.Flags().GetStringSlice("params")

	// Posted data
	for _, param := range params {
		kvParam := strings.SplitN(param, "=", 2)
		if len(kvParam) == 2 {
			postParam.Set(kvParam[0], kvParam[1])
		} else if len(kvParam) == 1 {
			postParam.Set(kvParam[0], "")
		}
	}

	// Method URL Headers Cookies Params
	client := &http.Client{}
	req, err := http.NewRequest(method.Value.String(), callURL.Value.String(), strings.NewReader(postParam.Encode()))

	// Cookies setting
	cookies, _ := cmd.Flags().GetStringSlice("cookies")
	for _, cookie := range cookies {
		kvCookie := strings.SplitN(cookie, "=", 2)
		if len(kvCookie) == 2 {
			req.AddCookie(&http.Cookie{
				Name:  kvCookie[0],
				Value: kvCookie[1],
			})
		}
	}

	if err != nil {
		log.Println(err)
		return
	}

	// Headers setting
	headers, _ := cmd.Flags().GetStringSlice("headers")

	for _, headerStr := range headers {
		kvSlice := strings.SplitN(headerStr, "=", 2)
		if len(kvSlice) == 2 {
			req.Header.Set(kvSlice[0], kvSlice[1])
		} else if len(kvSlice) == 1 {
			req.Header.Set(kvSlice[0], "")
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	// fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}
