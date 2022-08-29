package cmd

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

/*************************** Command ***************************/
var (
	_version string
	_name    string
)

const _license = `
'dq' Copyright (C) 2022 Paulo Canilho paulo@canilho.net
`

/*** Backend ***/
var serialiser Serialiser
var internalLogger *logger.Entry

/*** Persistent flags ***/
var outputFormat string

/*** Stats ***/
var debug bool
var startTime time.Time

var rootCommand = &cobra.Command{
	Use:     _name,
	Version: fmt.Sprintf("\n(%s)\n%s", _version, _license),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Logger
		if debug {
			logger.SetLevel(logger.TraceLevel)
		}
		internalLogger = logger.NewEntry(logger.StandardLogger())

		// Serialiser
		switch strings.TrimSpace(strings.ToLower(outputFormat)) {
		case "junos":
			serialiser = newJunosSerialiser()
		case "yaml":
			serialiser = newYAMLSerialiser()
		case "json":
			serialiser = newJSONSerialiser()
		case "plain":
		default:
			logger.Fatalln("Unsupported output format supplied")
		}

		startTime = time.Now()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		differ := newDiffer(internalLogger.WithField("internal", "differ"))
		for _, fp := range args {
			v := viper.New()
			v.SetConfigFile(fp)
			internalLogger.Infof("Reading configuration file [%s]...", fp)
			if err := v.ReadInConfig(); err != nil {
				return errors.Wrapf(err, "unable to read the provided file [%s]", fp)
			}
			differ.AddSettings(fp, v.AllSettings())
		}

		differ.DebugDiff()
		diffed := fmt.Sprint(differ.Diff())
		if serialiser != nil {
			diffed = serialiser.Serialise(differ.Diff())
		}
		fmt.Println(diffed)
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		logger.Tracef("Elapsed time: %v", time.Now().Sub(startTime).String())
	},
}

func Execute() error {
	return rootCommand.Execute()
}

func init() {
	viper.AutomaticEnv()
	loggerSetup()

	// Environment
	rootCommand.PersistentFlags().BoolVarP(&debug, "debug", "d", false,
		"Display debug information when provided.")
	rootCommand.PersistentFlags().StringVarP(&outputFormat, "format", "o", "yaml",
		"The output format to be used by the application. Supported: [json, yaml, table]")
}

func loggerSetup() {
	logger.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		ShowFullLevel:   true,
		TimestampFormat: time.RFC3339,
	})
	logger.SetLevel(logger.ErrorLevel)
}
