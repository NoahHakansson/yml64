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
	"os"

	"github.com/NoahHakansson/yml64/process"
	"github.com/spf13/cobra"
)

var decode bool
var inplace bool
var output string
var input []byte

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ysf",
	Short: "`yml64` is a CLI tool for encoding and decoding Kubernetes secret data properties in YAML/YML files using Base64.",
	Long:  `yml64 is a CLI tool for encoding/decoding Kubernetes secrets in YAML/YML files using Base64. It supports input from files or stdin and output to files or stdout. Written in Go, it offers fast performance and an easy-to-use command-line interface.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// check if the number of arguments is exactly 1
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		// Read the file
		ymlFile, err := os.ReadFile(args[0])
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("file %s does not exist", args[0])
			}
			return fmt.Errorf("could not read file %s: %v", args[0], err)
		}
		// Set the input file and bytes
		input = ymlFile
		return nil

	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Decode: ", decode)
		fmt.Println("Output: ", output)
		fmt.Println("Inplace: ", inplace)
		fmt.Println("Args: ", args)
		fmt.Println("File: ", string(input))
		fmt.Print("\n\n")

		result, err := codeDataProps(input, decode)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		// output result to file, stdout or overwrite input file
		if output != "" {
			if inplace {
				// overwrite the input file
				err := os.WriteFile(args[0], result, 0644)
				if err != nil {
					fmt.Println("Error: ", err)
					os.Exit(1)
				}
				fmt.Println("File ", args[0], " has been encoded/decoded and saved.")
			} else {
				// write to a new file
				err := os.WriteFile(output, result, 0644)
				if err != nil {
					fmt.Println("Error: ", err)
					os.Exit(1)
				}
				fmt.Println("File ", output, " has been created with the encoded/decoded content.")
			}
		} else {
			// print to stdout
			fmt.Print(string(result))
		}
	},
}

func codeDataProps(input []byte, decode bool) ([]byte, error) {
	if decode {
		// return process.DecodeDataProps(input)
		fmt.Println("Decode: Not implemented yet")
	}
	return process.EncodeDataProps(input)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.yml64.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output file (if not specified, the output will be printed to the console)")
	rootCmd.Flags().BoolVarP(&decode, "decode", "d", false, "Decode the base64 in the input")
	rootCmd.Flags().BoolVarP(&inplace, "inplace", "i", false, "Replace the input file contents with the output")
	rootCmd.MarkFlagsMutuallyExclusive("output", "inplace")
}
