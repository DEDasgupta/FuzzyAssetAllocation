package main

import "fmt"

const(
	MeanCAPERatio = 16.85
)

type VariableNameEnum struct {
	CAPERatio     string
	FEDFundsRate  string
	TimeHorizon   string
	Cashflow      string
	CPI           string
	LowValue      string
	HighValue     string
	LowRate       string
	HighRate      string
	ShortTerm     string
	LongTerm      string
	Accumulation  string
	Drawdown      string
	Deflation     string
	LowInflation  string
	HighInflation string

	HighGrowth    string
	LowGrowth     string
}

var VariableNames = VariableNameEnum{
	//crisp variables
	CAPERatio:     "Cyclically Adjusted Price to Earnings (CAPE) Ratio",
	FEDFundsRate:  "US Federal Reserve Funds Rate",
	TimeHorizon:   "Investing Time Horizon",
	Cashflow:      "Annual Cashflow into or out of the investment portfolio as a percentage",
	CPI:           "Consumer Price Index as a measure of inflation",

	//fuzzy variables
	LowValue: "low valuation for stocks",
	LowGrowth: "high valuation for stocks",
	LowRate: "low interest rate environment",
	HighRate: "high interest rate environment",
	ShortTerm: "short investment period",
	LongTerm:  "long investment period",
	Accumulation: "accumulation period",
	Drawdown: "drawdown period",
	Deflation: "a deflationary environment",
	LowInflation: "an optimal inflationary environment",
	HighInflation: "a highly inflationary environment",


}

 type Input struct {
 	Name string
 	Value float64
 }

var CrispInputs = map[string]Input{}
var FuzzyInputs = map[string]Input{}


func AddCrispInput(name string, value float64) {
	CrispInputs[name] = Input{
		Name:  name,
		Value: value,
	}
	fmt.Printf("%v initialized to: %v\n",name, value)
}

func AddFuzzyVariable(name string, value float64) {
	FuzzyInputs[name] = Input{
		Name:  name,
		Value: value,
	}
	fmt.Printf("%v initialized to: %v\n",name, value)

}

func InitializeCrispInputs() {
	fmt.Println("Crisp Inputs")
	AddCrispInput(VariableNames.CAPERatio, 37.63)
	AddCrispInput(VariableNames.FEDFundsRate, 0.25)
	AddCrispInput(VariableNames.TimeHorizon, 20.0)
	AddCrispInput(VariableNames.Cashflow, 1.0)
	AddCrispInput(VariableNames.CPI, 4.8)
}

func normalize(val float64) float64 {
	if val < 0.0 {
		return 0.0
	} else if val > 1.0 {
		return 1.0
	}
	return val
}



func Line(varName string, min, max float64) float64 {
		v := CrispInputs[varName]
		return normalize((v.Value-min)/(max-min))
}

func Triangle(varName string, min, middle, max float64) float64 {
	v := CrispInputs[varName]
	if v.Value < middle {
		return normalize((v.Value - min) / (middle - min))
	} else {
		return normalize((v.Value - max) / (middle - max))
	}
}



func InitializeFuzzyInputs() {
	fmt.Println("\n\nFuzzy Inputs")
	AddFuzzyVariable(VariableNames.HighGrowth, Line(VariableNames.CAPERatio, MeanCAPERatio, 12.0))
	AddFuzzyVariable(VariableNames.LowGrowth, Line(VariableNames.CAPERatio, MeanCAPERatio, 30.0))
	AddFuzzyVariable(VariableNames.LowRate, Line(VariableNames.FEDFundsRate, 2.0, -2.0))
	AddFuzzyVariable(VariableNames.HighRate, Line(VariableNames.FEDFundsRate, 2.0, -6.0))
	AddFuzzyVariable(VariableNames.ShortTerm, Line(VariableNames.TimeHorizon, 5.0, 1.0))
	AddFuzzyVariable(VariableNames.LongTerm, Line(VariableNames.TimeHorizon, 5.0, 15.0))
	AddFuzzyVariable(VariableNames.Accumulation, Line(VariableNames.Cashflow, 0, 20.0))
	AddFuzzyVariable(VariableNames.Drawdown, Triangle(VariableNames.Cashflow, 0, -3.0, -10.0))
	AddFuzzyVariable(VariableNames.Deflation, Line(VariableNames.TimeHorizon, 0.0, -3.0))
	AddFuzzyVariable(VariableNames.LowInflation, Triangle(VariableNames.TimeHorizon, 0.0, 3.0, 7.0))
	AddFuzzyVariable(VariableNames.HighInflation, Line(VariableNames.TimeHorizon, 5.0, 15.0))
}

func ProcessRules() {

}