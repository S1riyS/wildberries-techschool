package main

import "fmt"

// ModernPrinter is an interface that the client expects
type ModernPrinter interface {
	Print(data []byte) error
}

// LegacyPrinter is the existing legacy printer
type LegacyPrinter struct{}

func (lp *LegacyPrinter) PrintLegacy(s string) {
	fmt.Println("Legacy printer: " + s)
}

// The adapter that implements ModernPrinter but uses the legacy printer
type PrinterAdapter struct {
	legacyPrinter *LegacyPrinter
}

func (pa *PrinterAdapter) Print(data []byte) error {
	str := string(data)
	pa.legacyPrinter.PrintLegacy(str)
	return nil
}

// The code that works only with ModernPrinter
func modernCode(printer ModernPrinter) {
	printer.Print([]byte("Hello from modern client!"))
}

func main() {
	// Create the legacy printer
	legacyPrinter := &LegacyPrinter{}

	// Create an adapter that wraps it
	adapter := &PrinterAdapter{
		legacyPrinter: legacyPrinter,
	}

	modernCode(adapter)
}
