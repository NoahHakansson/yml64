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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yml64",
	Short: "`yml64` is a CLI tool for encoding and decoding Kubernetes secret data properties in YAML/YML files using Base64.",
	Long:  `yml64 is a CLI tool for encoding/decoding Kubernetes secrets in YAML/YML files using Base64. It supports input from files or stdin and output to files or stdout. Written in Go, it offers fast performance and an easy-to-use command-line interface.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// check if the number of arguments is exactly 1
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("Running yml64 in debug mode\n")
		// get stdin
		inputReader := cmd.InOrStdin()

		// Read from stdin if no file is provided
		if len(args) > 0 && args[0] != "-" {
			log.Debugf("Reading from file %s", args[0])
			file, err := os.Open(args[0])
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("file %s does not exist", args[0])
				}
				return fmt.Errorf("could not read file %s: %v", args[0], err)
			}
			defer file.Close()
			// assign the file to the inputReader
			inputReader = file
		}

		log.Debugf("Decode: %v", f.Decode)
		log.Debugf("Output: %v", f.Output)
		log.Debugf("Inplace: %v", f.Inplace)
		log.Debugf("Args: %s", args)
		log.Debugf("File:\n%v", inputReader)
		log.Debug("\n\n")

		result, err := processData(inputReader, f.Decode)
		if err != nil {
			return err
		}
		return outputResult(result, args)

	},
}

func processData(input io.Reader, decode bool) ([]byte, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	if decode {
		// return process.DecodeDataProps(input)
		fmt.Println("Decode: Not implemented yet")
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

func outputResult(result []byte, args []string) error {
	// output result to a file, stdout or overwrite the input file
	if f.Output != "" {
		if f.Inplace && len(args) > 0 {
			// overwrite the input file
			log.Debugf("Overwriting file %s", args[0])
			file, err := os.Create(args[0])
			if err != nil {
				return err
			}
			log.Debugf("Writing to file %s", args[0])
			n, err := file.Write(result)
			if err != nil {
				return err
			}
			log.Debugf("Wrote %d bytes to file %s", n, args[0])
			fmt.Println("File ", args[0], " has been encoded/decoded and saved.")
		} else {
			// write to a new file
			err := os.WriteFile(f.Output, result, 0644)
			if err != nil {
				return err
			}
			fmt.Println("File ", f.Output, " has been created with the encoded/decoded content.")
		}
	} else {
		// print to stdout
		fmt.Print(string(result))
	}
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
