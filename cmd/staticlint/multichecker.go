package main

import (
	"github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/shekshuev/shortener/cmd/staticlint/analyzers"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var mychecks []*analysis.Analyzer
	stdAnalyzers := []*analysis.Analyzer{
		shadow.Analyzer,
		printf.Analyzer,
		structtag.Analyzer,
	}
	mychecks = append(mychecks, stdAnalyzers...)
	for _, v := range staticcheck.Analyzers {
		if v.Analyzer.Name[:2] == "SA" || v.Analyzer.Name == "S1000" {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	mychecks = append(mychecks, analyzer.Analyzer)
	mychecks = append(mychecks, analyzers.Analyzer)
	multichecker.Main(mychecks...)
}
