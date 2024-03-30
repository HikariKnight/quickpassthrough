package params

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

/*
	The whole purpose of this module is to make a struct
	to just carry all our parsed arguments around between functions

	Create a Params struct with all the argparse arguments
	pArg := params.NewParams()
*/

type Params struct {
	Flag        map[string]bool
	FlagCounter map[string]int
	IntList     map[string][]int
	StringList  map[string][]string
	String      map[string]string
}

func (p *Params) addFlag(name string, flag bool) {
	p.Flag[name] = flag
}

func (p *Params) addFlagCounter(name string, flag int) {
	p.FlagCounter[name] = flag
}

func (p *Params) addIntList(name string, flag []int) {
	p.IntList[name] = flag
}

func (p *Params) addStringList(name string, flag []string) {
	p.StringList[name] = flag
}

func (p *Params) addString(name string, flag string) {
	p.String[name] = flag
}

func NewParams() *Params {
	// Setup the parser for arguments
	parser := argparse.NewParser("quickpassthrough", "A utility to help you configure your host for GPU Passthrough")

	// Configure arguments
	/*gui := parser.Flag("g", "gui", &argparse.Options{
		Required: false,
		Help:     "Launch GUI (placeholder for now)",
	})*/

	// Add version flag
	version := parser.Flag("v", "version", &argparse.Options{
		Required: false,
		Help:     "Display version",
	})

	// Parse arguments
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(4)
	}

	// Make our struct
	pArg := &Params{
		Flag:        make(map[string]bool),
		FlagCounter: make(map[string]int),
		IntList:     make(map[string][]int),
		StringList:  make(map[string][]string),
		String:      make(map[string]string),
	}

	// Add all parsed arguments to a struct for portability since we will use them all over the program
	pArg.addFlag("version", *version)
	/*pArg.addFlag("gui", *gui)
	pArg.addFlag("gpu", *gpu)
	pArg.addFlag("usb", *usb)
	pArg.addFlag("nic", *nic)
	pArg.addFlag("sata", *sata)
	pArg.addFlagCounter("related", *related)
	pArg.addStringList("ignore", *ignore)
	pArg.addIntList("iommu_group", *iommu_group)
	pArg.addFlag("kernelmodules", *kernelmodules)
	pArg.addFlag("legacyoutput", *legacyoutput)
	pArg.addFlag("id", *id)
	pArg.addFlag("pciaddr", *pciaddr)
	pArg.addFlag("rom", *rom)
	pArg.addString("format", *format)*/

	return pArg
}
