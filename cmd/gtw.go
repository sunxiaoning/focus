/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"focus/app"
	"focus/cfg"
	resourceservice "focus/serv/resource"
	"github.com/spf13/cobra"
)

// gtwCmd represents the gtw command
var gtwCmd = &cobra.Command{
	Use:   "gtw",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunGtw()
	},
}

func RunGtw() error {
	fmt.Printf("runtime Env: %s", runtimeConfig.Env)
	var err error
	if err = app.InitCfg(runtimeConfig); err != nil {
		return err
	}
	app.InitLog()
	if err = app.InitDB(); err != nil {
		return err
	}
	defer cfg.FocusCtx.DB.Close()
	if err = resourceservice.InitServiceResource(); err != nil {
		return err
	}
	app.InitServer(cfg.FocusCtx.Cfg.Server.ListenPort+1, app.Gtw, app.GtwFilters)
	go app.InitTask()
	app.StartServer()
	return nil
}

func init() {
	rootCmd.AddCommand(gtwCmd)
	gtwCmd.Flags().StringVar(&runtimeConfig.Env, "env", "alpha", "server env")
	gtwCmd.Flags().StringVar(&runtimeConfig.SecretKeyPath, "", "", "server aeskey")
}
