package cmd

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
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

// RawData of Command and Form UI
func RawData() []Field {
	raw := []Field{
		Field{
			Display:     "HTTP Method",
			Name:        "method",
			EnvName:     "PLUGIN_METHOD",
			ShortName:   "M",
			Default:     "GET",
			Description: "HTTP Method, default to GET",
			Required:    false,
			Type:        "string",
			Array:       false,
		},
		Field{
			Display:     "URL",
			Name:        "url",
			ShortName:   "U",
			EnvName:     "PLUGIN_URL",
			Description: "Call URL, start with http:// or https://",
			Required:    true,
			Type:        "string",
			Array:       false,
		},
		Field{
			Display:     "Headers",
			Name:        "headers",
			EnvName:     "PLUGIN_HEADERS",
			Description: "Request Headers, e.g. color=black",
			Required:    false,
			Type:        "string",
			Array:       true,
		},
		Field{
			Display:     "Cookies",
			Name:        "cookies",
			ShortName:   "C",
			EnvName:     "PLUGIN_COOKIES",
			Description: "Request with cookies, e.g. name=jack",
			Required:    false,
			Type:        "string",
			Array:       true,
		},
		Field{
			Display:     "Parameters",
			Name:        "params",
			ShortName:   "P",
			EnvName:     "PLUGIN_PARAMS",
			Description: "Reuqest with params, form/url params/post data",
			Required:    false,
			Type:        "string",
			Array:       true,
		},
	}
	return raw
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringP("method", "M", "GET", "HTTP Method, default to GET")
	rootCmd.Flags().StringP("url", "U", "", "Call URL, start with http:// or https://")
	rootCmd.Flags().StringSlice("headers", []string{}, "Request Headers, e.g. color=black")
	rootCmd.Flags().StringSliceP("cookies", "C", []string{}, "Request with cookies, e.g. name=jack")
	rootCmd.Flags().StringSliceP("params", "P", []string{}, "Reuqest with params, form/url params/post data")
	// rootCmd.MarkFlagRequired("url")
}

func initConfig() {
	// load environment variables
	viper.SetEnvPrefix("PLUGIN")
	viper.AutomaticEnv()
}

func run(cmd *cobra.Command, args []string) {
	// required flags
	method := viper.GetString("METHOD")
	if len(method) == 0 {
		method, _ = cmd.Flags().GetString("method")
	}

	callURL := viper.GetString("URL")
	if len(callURL) == 0 {
		callURL, _ = cmd.Flags().GetString("url")
	}

	postParam := url.Values{}
	paramStr := viper.GetString("PARAMS")
	params := strings.Split(paramStr, "&")
	if len(paramStr) == 0 || len(params) == 0 {
		params, _ = cmd.Flags().GetStringSlice("params")
	}

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
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 10,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				MaxVersion:         tls.VersionTLS11,
			},
		},
	}
	req, err := http.NewRequest(method, callURL, strings.NewReader(postParam.Encode()))

	// Cookies setting
	cookieStr := viper.GetString("COOKIES")
	cookies := strings.Split(cookieStr, "&")
	if len(cookieStr) == 0 || len(cookies) == 0 {
		cookies, _ = cmd.Flags().GetStringSlice("cookies")
	}

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
	headerStr := viper.GetString("HEADERS")
	headers := strings.Split(headerStr, "&")
	if len(headerStr) == 0 || len(headers) == 0 {
		headers, _ = cmd.Flags().GetStringSlice("headers")
	}

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
	fmt.Println(string(body))
}
