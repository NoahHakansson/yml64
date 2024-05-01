/*
Copyright © 2024 Noah Hakansson noah.hakansson@protonmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/NoahHakansson/yml64/model"
	"github.com/NoahHakansson/yml64/process"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var f model.Flags
var input model.Input

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yml64",
	Short: "yml64 is a CLI tool for encoding and decoding Kubernetes secret data properties in YAML/YML files using Base64.",
	Long:  `yml64 is a CLI tool for encoding/decoding Kubernetes secrets in YAML/YML files using Base64. It supports input from files or stdin and output to files or stdout. Written in Go, it offers fast performance and an easy-to-use command-line interface.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// check if the number of arguments is exactly 1
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return err
		}

		// if no input file is provided, the -i (inplace) flag is not allowed
		if len(args) == 0 && f.Inplace {
			return fmt.Errorf("inplace flag is not allowed when reading from stdin")
		}

		// assign the input file path to a global variable
		if len(args) > 0 && args[0] != "-" && args[0] != "" {
			input.Exists = true
			input.Path = args[0]
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("Running yml64 in debug mode\n")
		// get stdin
		inputReader := cmd.InOrStdin()

		// if an input file is provided, read from it
		if input.Exists {
			log.Debugf("Reading from file %s", args[0])
			file, err := os.Open(args[0])
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("file %s does not exist", args[0])
				}
				return fmt.Errorf("could not read file %s: %v", args[0], err)
			}
			defer file.Close()
			inputReader = file
		}

		log.Debugf("Decode: %v", f.Decode)
		log.Debugf("Metadata: %v", f.Metadata)
		log.Debugf("Output: %v", f.Output)
		log.Debugf("Inplace: %v", f.Inplace)
		log.Debugf("Args: %s", args)
		log.Debugf("Nr of args: %d", len(args))
		log.Debugf("File:\n%v\n\n", inputReader)

		result, err := processData(inputReader, f.Decode)
		if err != nil {
			return err
		}
		return outputResult(result)
	},
}

func processData(input io.Reader, decode bool) ([]byte, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	if decode {
		return process.DecodeDataProps(data, f)
	}
	return process.EncodeDataProps(data, f)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func getDecodeOrEncodeMessage(decode bool) string {
	if decode {
		return "Decoded"
	}
	return "Encoded"
}

func writeResultToFile(result []byte, filePath string) error {
	// write to a new file
	err := os.WriteFile(filePath, result, 0644)
	if err != nil {
		return err
	}
	log.Debugf("Wrote %d bytes to file %s", len(result), filePath)
	fmt.Println(getDecodeOrEncodeMessage(f.Decode), "content has been saved to file:", filePath)
	return nil
}

func outputResult(result []byte) error {
	// output result to a file, stdout or overwrite the input file
	if f.Inplace && input.Exists {
		// overwrite the input file
		log.Debugf("Overwriting file %s", input.Path)
		return writeResultToFile(result, input.Path)
	} else if f.Output != "" {
		// write to a new file
		log.Debugf("Writing to file %s", f.Output)
		return writeResultToFile(result, f.Output)
	}
	// print to stdout
	fmt.Printf("%s", string(result))
	return nil
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.yml64.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&f.Output, "output", "o", "", "Output file (if not specified, the output will be printed to the console)")
	rootCmd.Flags().BoolVarP(&f.Decode, "decode", "d", false, "Decode the base64 in the input")
	rootCmd.Flags().BoolVarP(&f.Inplace, "inplace", "i", false, "Replace the input file contents with the output")
	rootCmd.Flags().BoolVarP(&f.Metadata, "metadata", "m", false, "Include ALL metadata in the output and not just the name and namespace fields")
	rootCmd.MarkFlagsMutuallyExclusive("output", "inplace")
}
