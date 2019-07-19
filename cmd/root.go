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
	rootCmd.Flags().String()
}

func initConfig() {
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func run(cmd *cobra.Command, args []string) {
	// Method URL Headers Cookies Params
	client := &http.Client{}

	postParam := url.Values{
		"a": {"1"},
		"b": {"2"},
	}
	postParam.Set("c", "3")

	req, err := http.NewRequest("GET", "http://10.96.110.34:8086/dcs", strings.NewReader(postParam.Encode()))
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

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
