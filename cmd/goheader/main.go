// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

import (
	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/gss"
)

import (
	"github.com/spatialcurrent/go-header/pkg/goheader"
)

var gitTag string
var gitBranch string
var gitCommit string

func mustCompilePatterns(expressions map[string]string) map[string]*regexp.Regexp {
	patterns := map[string]*regexp.Regexp{}
	for name, expression := range expressions {
		patterns[name] = regexp.MustCompile(expression)
	}
	return patterns
}

var patterns = mustCompilePatterns(map[string]string{
	"go": "(?s)^(.*?)package",
})

func main() {
	rootCommand := &cobra.Command{
		Use:   "goheader",
		Short: "goheader",
		Long:  `goheader is a simple program for inspecting and updating the headers in files`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "print version information to stdout",
		Long:  "print version information to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(gitTag) > 0 {
				fmt.Println("Tag: " + gitTag)
			}
			if len(gitBranch) > 0 {
				fmt.Println("Branch: " + gitBranch)
			}
			if len(gitCommit) > 0 {
				fmt.Println("Commit: " + gitCommit)
			}
			return nil
		},
	}

	dumpCommand := &cobra.Command{
		Use:   "dump",
		Short: "dump headers",
		Long:  "dump headers",
		RunE: func(cmd *cobra.Command, args []string) error {

			v := viper.New()
			err := v.BindPFlags(cmd.Flags())
			if err != nil {
				panic(err)
			}
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()

			inputDir := v.GetString("input-dir")
			inputSkip := v.GetString("input-skip")
			var inputSkipPattern *regexp.Regexp
			if len(inputSkip) > 0 {
				p, err := regexp.Compile(inputSkip)
				if err != nil {
					return err
				}
				inputSkipPattern = p
			}
			inputExtension := v.GetString("input-extesion")
			outputFormat := v.GetString("output-format")

			outputObjects := make([]map[string]interface{}, 0)

			err = godirwalk.Walk(inputDir, &godirwalk.Options{
				Callback: func(p string, de *godirwalk.Dirent) error {
					if de.ModeType().IsDir() {
						if inputSkipPattern != nil && inputSkipPattern.MatchString(p) {
							return filepath.SkipDir
						}
					} else if de.ModeType().IsRegular() {
						if inputSkipPattern != nil && inputSkipPattern.MatchString(p) {
							return nil
						}
						ext := filepath.Ext(p)
						if len(ext) > 0 {
							ext = ext[1:]
						}
						if len(inputExtension) == 0 || ext == inputExtension {
							pattern, ok := patterns[ext]
							if ok {
								b, err := ioutil.ReadFile(p) // #nosec
								if err != nil {
									return err
								}
								outputMap := map[string]interface{}{
									"path": p,
								}
								if headerMatches := pattern.FindStringSubmatch(string(b)); len(headerMatches) > 0 {
									h, err := goheader.ParseHeader(goheader.RemoveEmptyLines(goheader.RemoveLinePrefix(headerMatches[1], "//")))
									if err != nil {
										return err
									}
									for k, v := range h.Map() {
										outputMap[k] = v
									}
								}
								outputObjects = append(outputObjects, outputMap)
							}
						}
					}
					return nil
				},
				ErrorCallback: func(p string, err error) godirwalk.ErrorAction {
					fmt.Println(errors.Wrap(err, "error walking "+p))
					return godirwalk.SkipNode
				},
			})
			if err != nil {
				return err
			}

			outputString, err := gss.SerializeString(outputObjects, outputFormat, []string{}, gss.NoLimit)
			if err != nil {
				return err
			}
			fmt.Println(outputString)
			return nil
		},
	}

	flags := dumpCommand.Flags()
	flags.StringP("input-dir", "i", "", "The input directory")
	flags.StringP("input-skip", "s", "", "Skip these directories and files.  A regular expression that is matched against the path.")
	flags.StringP("input-extension", "e", "", "Filter all files by this extesion")
	flags.StringP("output-format", "f", "json", "Output format")

	rootCommand.AddCommand(versionCommand, dumpCommand)

	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
