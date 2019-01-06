/**
* Copyright 2018 Comcast Cable Communications Management, LLC
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	// Command-line flags
	cfConfig      = "config"
	cfVersion     = "version"
	cfLogLevel    = "log-level"
	cfInstanceID  = "instance-id"
	cfOrigin      = "origin"
	cfProxyPort   = "proxy-port"
	cfMetricsPort = "metrics-port"

	// Environment variables
	evOrigin      = "TRK_ORIGIN"
	evProxyPort   = "TRK_PROXY_PORT"
	evMetricsPort = "TRK_METRICS_PORT"
	evLogLevel    = "TRK_LOG_LEVEL"
)

// loadConfiguration reads the config path from Flags,
// Loads the configs (w/ default values where missing)
// and then evaluates any provided flags as overrides
func loadConfiguration(c *Config, arguments []string) error {
	args := loadFlags(arguments)

	// Display version information then exit the program
	if args.version {
		fmt.Println(applicationVersion)
		os.Exit(3)
	}

	err := loadConfigurationFile(c, args)
	if err != nil {
		return err
	}

	// Set the configuration loaded from the arguments
	args.setConfig(c)

	// Load from Environment Variables
	loadEnvVars(c)

	return nil
}

func loadEnvVars(c *Config) {
	// Origin
	if x := os.Getenv(evOrigin); x != "" {
		c.DefaultOriginURL = x
	}

	// Proxy Port
	if x := os.Getenv(evProxyPort); x != "" {
		if y, err := strconv.ParseInt(x, 10, 64); err == nil {
			c.ProxyServer.ListenPort = int(y)
		}
	}

	// Metrics Port
	if x := os.Getenv(evMetricsPort); x != "" {
		if y, err := strconv.ParseInt(x, 10, 64); err == nil {
			c.Metrics.ListenPort = int(y)
		}
	}

	// LogLevel
	if x := os.Getenv(evLogLevel); x != "" {
		c.Logging.LogLevel = x
	}

}

type args struct {
	path              string
	origin            string
	proxyListenPort   int
	metricsListenPort int
	logLevel          string
	instanceID        int
	version           bool
}

func (a args) setConfig(c *Config) {
	if len(a.origin) > 0 {
		c.DefaultOriginURL = a.origin
	}
	if a.proxyListenPort > 0 {
		c.ProxyServer.ListenPort = a.proxyListenPort
	}
	if a.metricsListenPort > 0 {
		c.Metrics.ListenPort = a.metricsListenPort
	}
	c.Logging.LogLevel = a.logLevel
	c.Main.InstanceID = a.instanceID

}

// loadFlags loads command line flags.
func loadFlags(arguments []string) args {
	var args args

	f := flag.NewFlagSet(applicationName, flag.ExitOnError)
	f.StringVar(&args.logLevel, cfLogLevel, "INFO", "Level of Logging to use (debug, info, warn, error)")
	f.IntVar(&args.instanceID, cfInstanceID, 0, "Instance ID for when running multiple processes")
	f.StringVar(&args.origin, cfOrigin, "", "URL to the Prometheus Origin. Enter it like you would in grafana, e.g., http://prometheus:9090")
	f.IntVar(&args.proxyListenPort, cfProxyPort, 0, "Port that the Proxy server will listen on.")
	f.IntVar(&args.metricsListenPort, cfMetricsPort, 0, "Port that the /metrics endpoint will listen on.")
	f.StringVar(&args.path, cfConfig, "", "Path to Trickster Config File")

	f.Parse(arguments)

	return args
}

// loadConfigurationFile gets the arguments passed to the executable and reads
// the basic flags necessary to identify the configuration file. It then loads
// the configuration file.
func loadConfigurationFile(c *Config, a args) error {
	// If the config file is not specified on the cmdline then try the default
	// location to load the config file.  If the default config does not exist
	// then move on, no big deal.
	if a.path != "" {
		if err := c.LoadFile(a.path); err != nil {
			return err
		}
	} else {
		_, err := os.Open(c.Main.ConfigFile)
		if err == nil {
			if err := c.LoadFile(c.Main.ConfigFile); err != nil {
				return err
			}
		}
	}

	return nil
}
