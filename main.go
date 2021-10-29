package main

import (
	"fmt"
)

func main() {
	f := FuzzyAssetAllocator{}
	f.Initialize()
	f.InitializeFuzzyInputs()
	f.ProcessRules()
	f.CalculateResultantAssetAllocation()
	f.FactorizeReasons()
}


const(
	MeanCAPERatio = 16.85
	stocks = "stocks"
	bonds = "bonds"
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

type Weight struct {
	Asset string
	Weight float64
	Reason string
	Contribution float64
}

type FuzzyAssetAllocator struct {
	CrispInputs map[string]Input
	FuzzyInputs map[string]Input
	Weights []Weight
	TotalWeight float64
}


func (f *FuzzyAssetAllocator) AddCrispInput(name string, value float64) {
	f.CrispInputs[name] = Input{
		Name:  name,
		Value: value,
	}
	fmt.Printf("%v initialized to: %v\n",name, value)
}

func (f *FuzzyAssetAllocator) AddFuzzyVariable(name string, value float64) {
	f.FuzzyInputs[name] = Input{
		Name:  name,
		Value: value,
	}
	fmt.Printf("%v initialized to: %v\n",name, value)

}

func (f *FuzzyAssetAllocator) Initialize() {
	f.CrispInputs = map[string]Input{}
	fmt.Println("Crisp Inputs")
	f.AddCrispInput(VariableNames.CAPERatio, 37.63)
	f.AddCrispInput(VariableNames.FEDFundsRate, 0.25)
	f.AddCrispInput(VariableNames.TimeHorizon, 20.0)
	f.AddCrispInput(VariableNames.Cashflow, 1.0)
	f.AddCrispInput(VariableNames.CPI, 4.8)
}

func normalize(val float64) float64 {
	if val < 0.0 {
		return 0.0
	} else if val > 1.0 {
		return 1.0
	}
	return val
}



func Line(v Input, min, max float64) float64 {
	return normalize((v.Value-min)/(max-min))
}

func Triangle(v Input, min, middle, max float64) float64 {
	if v.Value < middle {
		return normalize((v.Value - min) / (middle - min))
	} else {
		return normalize((v.Value - max) / (middle - max))
	}
}


// InitializeFuzzyInputs sets the membership functions for
func (f *FuzzyAssetAllocator) InitializeFuzzyInputs() {
	fmt.Println("\n\nFuzzy Inputs")
	f.AddFuzzyVariable(VariableNames.HighGrowth, Line(f.CrispInputs[VariableNames.CAPERatio], MeanCAPERatio, 12.0))
	f.AddFuzzyVariable(VariableNames.LowGrowth, Line(f.CrispInputs[VariableNames.CAPERatio], MeanCAPERatio, 30.0))
	f.AddFuzzyVariable(VariableNames.LowRate, Line(f.CrispInputs[VariableNames.FEDFundsRate], 2.0, -2.0))
	f.AddFuzzyVariable(VariableNames.HighRate, Line(f.CrispInputs[VariableNames.FEDFundsRate], 2.0, 10.0))
	f.AddFuzzyVariable(VariableNames.ShortTerm, Line(f.CrispInputs[VariableNames.TimeHorizon], 5.0, 1.0))
	f.AddFuzzyVariable(VariableNames.LongTerm, Line(f.CrispInputs[VariableNames.TimeHorizon], 5.0, 15.0))
	f.AddFuzzyVariable(VariableNames.Accumulation, Line(f.CrispInputs[VariableNames.Cashflow], 0, 20.0))
	f.AddFuzzyVariable(VariableNames.Drawdown, Triangle(f.CrispInputs[VariableNames.Cashflow], 0, -3.0, -10.0))
	f.AddFuzzyVariable(VariableNames.Deflation, Line(f.CrispInputs[VariableNames.TimeHorizon], 0.0, -3.0))
	f.AddFuzzyVariable(VariableNames.LowInflation, Triangle(f.CrispInputs[VariableNames.TimeHorizon], 0.0, 3.0, 7.0))
	f.AddFuzzyVariable(VariableNames.HighInflation, Line(f.CrispInputs[VariableNames.TimeHorizon], 5.0, 15.0))
}

func (f *FuzzyAssetAllocator) ProcessRules() {

}

func (f *FuzzyAssetAllocator) CalculateResultantAssetAllocation() float64 {
	s := 0.0
	b := 0.0
	for _, w := range f.Weights {
		if w.Asset == stocks {
			s += w.Weight
		} else {
			b += w.Weight
		}
	}
	f.TotalWeight = s+b
	return s/f.TotalWeight
}

func (f *FuzzyAssetAllocator) FactorizeReasons() {

}