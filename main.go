package main

import (
	"fmt"
	"math"
	"sort"
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
	CAPERatio      string
	FEDFundsRate   string
	TimeHorizon    string
	Cashflow       string
	CPI            string
	LowRate        string
	HighRate       string
	ShortTerm      string
	LongTerm       string
	Accumulation   string
	Drawdown       string
	Deflation      string
	LowInflation   string
	HighInflation  string
	HighGrowth     string
	LowGrowth      string
	Contractionary string
	Expansionary   string
	RisingRate     string
	FallingRate    string
	Stagflationary string
}

var VariableNames = VariableNameEnum{
	//crisp variables
	CAPERatio:     "Cyclically Adjusted Price to Earnings (CAPE) Ratio",
	FEDFundsRate:  "US Federal Reserve Funds Rate",
	TimeHorizon:   "Investing Time Horizon",
	Cashflow:      "Annual Cashflow into or out of the investment portfolio as a percentage",
	CPI:           "Consumer Price Index as a measure of inflation",

	//fuzzy variables
	HighGrowth: "low valuation for stocks",
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

	//compound variables
	Contractionary: "risk of low growth and high rates",
	Expansionary: "accomodative policy and expected high growth",
	RisingRate: "rates expected to rise due to high inflation",
	FallingRate: "rates expected to fall due to low inflation",
	Stagflationary: "risk of low growth and high inflation",

}

type Input struct {
	Name string
	Value float64
}

type Weight struct {
	Asset string
	Weight float64
	Reason string
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
	f.FuzzyInputs = map[string]Input{}
	f.Weights = []Weight{}
	fmt.Println("Crisp Inputs")
	f.AddCrispInput(VariableNames.CAPERatio, 37.63)
	f.AddCrispInput(VariableNames.FEDFundsRate, 0.25)
	f.AddCrispInput(VariableNames.TimeHorizon, 20.0)
	f.AddCrispInput(VariableNames.Cashflow, 1.0)
	f.AddCrispInput(VariableNames.CPI, 4.8)
}

func (f *FuzzyAssetAllocator) Process(m map[string]float64) float64 {
	f.Initialize()
	f.SetCrispInputs(m)
	f.InitializeFuzzyInputs()
	f.ProcessRules()
	allocation := f.CalculateResultantAssetAllocation()
	f.FactorizeReasons()
	return allocation
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

func (f *FuzzyAssetAllocator) SetCrispInputs(m map[string]float64) {
	f.AddCrispInput(VariableNames.CAPERatio, m[VariableNames.CAPERatio])
	f.AddCrispInput(VariableNames.FEDFundsRate, m[VariableNames.FEDFundsRate])
	f.AddCrispInput(VariableNames.TimeHorizon, m[VariableNames.TimeHorizon])
	f.AddCrispInput(VariableNames.Cashflow, m[VariableNames.Cashflow])
	f.AddCrispInput(VariableNames.CPI, m[VariableNames.CPI])
}


// InitializeFuzzyInputs sets the membership functions for
func (f *FuzzyAssetAllocator) InitializeFuzzyInputs() {
	fmt.Println("\n\nFuzzy Inputs")
	f.AddFuzzyVariable(VariableNames.HighGrowth, Line(f.CrispInputs[VariableNames.CAPERatio], MeanCAPERatio, 10.0))
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

func (f *FuzzyAssetAllocator) WeightVariable(varName, asset string, weight float64) {
	v := f.FuzzyInputs[varName]
	w := Weight{
		Asset:        asset,
		Weight:       v.Value*weight,
		Reason:       varName,
	}
	f.Weights = append(f.Weights, w)
}

func (f *FuzzyAssetAllocator) And(var1, var2, varName, asset string, weight float64) {
	min := math.Min(f.FuzzyInputs[var1].Value, f.FuzzyInputs[var2].Value)
	f.FuzzyInputs[varName] = Input{Name: varName, Value: min}
	f.WeightVariable(varName, asset, weight)
}

func (f *FuzzyAssetAllocator) ProcessRules() {
	//High Growth environment good for stocks
	f.WeightVariable(VariableNames.HighGrowth, stocks, 10)

	//Low Growth environment not good for stocks
	f.WeightVariable(VariableNames.LowGrowth, bonds, 10)

	//Low Rate environment good for stocks
	f.WeightVariable(VariableNames.LowRate, bonds, 20)

	//High Rate environment not good for stocks
	f.WeightVariable(VariableNames.HighRate, bonds, 20)

	//bonds mitigate risk for the short term
	f.WeightVariable(VariableNames.ShortTerm, bonds, 10)

	//stocks perform better in the long term
	f.WeightVariable(VariableNames.LongTerm, stocks, 10)

	//High Rate environment not good for stocks
	f.WeightVariable(VariableNames.Deflation, bonds, 50)

	//low inflation environment generally good for both bonds and stocks
	f.WeightVariable(VariableNames.LowInflation, stocks, 10)
	f.WeightVariable(VariableNames.LowInflation, bonds, 10)

	//High Inflation environment is generally bad, but much better for stocks relative to bonds
	f.WeightVariable(VariableNames.HighInflation, stocks, 30)

	//stocks likely to fall due to high rates and low expected growth
	f.And(VariableNames.LowGrowth, VariableNames.HighRate, VariableNames.Contractionary, bonds, 50)

	//stocks likely to increase greatly due to low rates and high growth expectations
	f.And(VariableNames.HighGrowth, VariableNames.LowRate, VariableNames.Expansionary, stocks, 100)

	//bonds likely to fall due to rising rates
	f.And(VariableNames.HighInflation, VariableNames.LowRate, VariableNames.RisingRate, stocks, 50)

	//bonds likely to rise due to falling rates
	f.And(VariableNames.LowInflation, VariableNames.HighRate, VariableNames.Contractionary, bonds, 50)

	//bonds likely to rise due to falling rates
	f.And(VariableNames.LowInflation, VariableNames.HighRate, VariableNames.Contractionary, bonds, 50)

	//bonds likely to rise due to falling rates
	f.And(VariableNames.LowGrowth, VariableNames.HighInflation, VariableNames.Stagflationary, stocks, 20)



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
	sort.SliceStable(f.Weights, func(i, j int) bool {
		return f.Weights[i].Weight > f.Weights[j].Weight
	})

	for i:= 0; i<3; i++ {
		w := f.Weights[i]
		fmt.Printf("%.2f indicating %v due to %v.\n", w.Weight/f.TotalWeight, w.Asset, w.Reason)
	}
}