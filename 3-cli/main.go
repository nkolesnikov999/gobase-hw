package main

import (
	"cli/bins"
)

func main() {
	// Example usage
	bin := bins.NewBin("bin-001", "My First Bin", false)

	binList := bins.NewBinList()
	binList.AddBin(bin)
}
