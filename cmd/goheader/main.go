// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
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
	"go": "(?s)^(.*?)(package (?:\\w+)(?:.*))$",
})

func createRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goheader",
		Short: "goheader",
		Long:  `goheader is a simple program for inspecting and updating the headers in files`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	flags := cmd.PersistentFlags()
	flags.BoolP("verbose", "v", false, "Verbose output")
	return cmd
}

func createVersionCommand() *cobra.Command {
	return &cobra.Command{
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
}

func createDumpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump",
		Short: "dump headers",
		Long:  "dump headers",
		RunE: func(cmd *cobra.Command, args []string) error {

			v := viper.New()
			err := v.BindPFlags(cmd.Flags())
			if err != nil {
				return err
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
								outputFile := goheader.File{
									Path:      p,
									Name:      filepath.Base(p),
									Extension: ext,
								}
								if matches := pattern.FindStringSubmatch(string(b)); len(matches) > 0 {
									h, err := goheader.ParseHeader(matches[1])
									if err != nil {
										return err
									}
									outputFile.Header = h
								}
								outputObjects = append(outputObjects, outputFile.Map())
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

			outputString, err := gss.SerializeString(&gss.SerializeInput{
				Object: outputObjects,
				Format: outputFormat,
				Header: gss.NoHeader,
				Limit:  gss.NoLimit,
			})
			if err != nil {
				return err
			}
			fmt.Println(outputString)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringP("input-dir", "i", "", "The input directory")
	flags.StringP("input-skip", "s", "", "Skip these directories and files.  A regular expression that is matched against the path.")
	flags.StringP("input-extension", "e", "", "Filter all files by this extesion")
	flags.StringP("output-format", "f", "json", "Output format")

	return cmd
}

func createFixCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fix",
		Short: "fix headers",
		Long:  "fix headers",
		RunE: func(cmd *cobra.Command, args []string) error {

			v := viper.New()
			err := v.BindPFlags(cmd.Flags())
			if err != nil {
				return err
			}
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()

			verbose := v.GetBool("verbose")

			exitCodeOnChanges := v.GetInt("exit-code-on-changes")

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
			fixYear := v.GetInt("fix-year")

			changes := 0
			err = godirwalk.Walk(inputDir, &godirwalk.Options{
				Callback: func(p string, de *godirwalk.Dirent) error {
					modeType := de.ModeType()
					if modeType.IsDir() {
						if inputSkipPattern != nil && inputSkipPattern.MatchString(p) {
							return filepath.SkipDir
						}
					} else if modeType.IsRegular() {
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
								fileContents, err := ioutil.ReadFile(p) // #nosec
								if err != nil {
									return err
								}
								before := string(fileContents)
								if matches := pattern.FindStringSubmatch(before); len(matches) > 0 {
									after := goheader.FixHeader(matches[1], fixYear) + matches[2]
									if before != after {
										changes++
										err := ioutil.WriteFile(p, []byte(after), modeType)
										if err != nil {
											return errors.Wrap(err, fmt.Sprintf("Error updating %q", p))
										}
										if verbose {
											fmt.Println(fmt.Sprintf("Updated %q", p))
										}
									}
								}
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
			if exitCodeOnChanges > 0 && changes > 0 {
				os.Exit(exitCodeOnChanges)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringP("input-dir", "i", "", "The input directory")
	flags.StringP("input-skip", "s", "", "Skip these directories and files.  A regular expression that is matched against the path.")
	flags.StringP("input-extension", "e", "", "Filter all files by this extesion")
	flags.IntP("fix-year", "y", -1, "The year to use.")
	flags.Int("exit-code-on-changes", 0, "The exit code to use when changes were made.")

	return cmd
}

func main() {
	rootCommand := createRootCommand()
	rootCommand.AddCommand(
		createVersionCommand(),
		createDumpCommand(),
		createFixCommand(),
	)

	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
