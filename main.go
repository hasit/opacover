package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	modules string
	input   string
	output  string
)

func init() {
	rootCmd.Flags().StringVarP(&modules, "modules", "m", "", "directory path of Rego modules")
	rootCmd.Flags().StringVarP(&input, "input", "i", "", "input file in JSON")
}

var rootCmd = &cobra.Command{
	Use:   "opacover",
	Short: "generate HTML representation of OPA test coverage",
	Run:   runOPACover,
}

// CoverageReport .
type CoverageReport struct {
	Coverage float32               `json:"coverage"`
	Files    map[string]FileReport `json:"files"`
}

// FileReport .
type FileReport struct {
	Covered    []LineReport `json:"covered,omitempty"`
	NotCovered []LineReport `json:"not_covered,omitempty"`
	Coverage   float32      `json:"coverage"`
	Index      int
	Body       string
}

// LineReport .
type LineReport struct {
	Start RowReport `json:"start"`
	End   RowReport `json:"end"`
}

// RowReport .
type RowReport struct {
	Row int `json:"row"`
}

func runOPACover(cmd *cobra.Command, args []string) {
	if input == "" {
		log.Fatal("input (--input) cannot be empty. Run 'opacover --help' for usage information.")
	}

	covBytes, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("error reading file %s: %v", input, err)
	}

	var covReport CoverageReport
	err = json.Unmarshal(covBytes, &covReport)
	if err != nil {
		log.Fatalf("error unmarshalling content of file %s: %v", input, err)
	}

	index := 0
	for file, fileReport := range covReport.Files {
		fileReport.Index = index

		fBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", modules, file))
		if err != nil {
			log.Fatalf("error opening modules file %s: %v", file, err)
		}

		fSlice := strings.Split(string(fBytes), "\n")

		// patch for Covered
		for _, covered := range fileReport.Covered {
			fSlice[covered.Start.Row-1] = fmt.Sprintf("<span class=\"covered\">%s", fSlice[covered.Start.Row-1])
			fSlice[covered.End.Row-1] = fmt.Sprintf("%s</span>", fSlice[covered.End.Row-1])
		}

		// patch fot Not Covered
		for _, notcovered := range fileReport.NotCovered {
			fSlice[notcovered.Start.Row-1] = fmt.Sprintf("<span class=\"not-covered\">%s", fSlice[notcovered.Start.Row-1])
			fSlice[notcovered.End.Row-1] = fmt.Sprintf("%s</span>", fSlice[notcovered.End.Row-1])
		}

		fileReport.Body = strings.Join(fSlice, "\n")
		index++

		covReport.Files[file] = fileReport
	}

	err = generateOutput(covReport)
	if err != nil {
		log.Fatalf("error generating output: %v", err)
	}
}

func generateOutput(report CoverageReport) error {
	tmpl, err := template.ParseFiles("index.gohtml")
	if err != nil {
		return err
	}

	err = tmpl.Execute(os.Stdout, report)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("error executing opacover: %v", err)
	}
}
