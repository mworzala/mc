package model

type OutputFormat string

const (
	OutputDefault OutputFormat = "default"
	OutputJson    OutputFormat = "json"
	OutputYaml    OutputFormat = "yaml"
)

var OutputFormatValidationMap = map[string]bool{
	string(OutputDefault): true,
	string(OutputJson):    true,
	string(OutputYaml):    true,
}

type PrettyPrintable interface {
	Print()
}

type CliError struct {
}
