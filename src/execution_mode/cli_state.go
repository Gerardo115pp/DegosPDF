package execution_mode

import (
	"flag"
	"fmt"
)

const (
	OUTPUT_THICK_DIVIDER string = "========================================"
	OUTPUT_THIN_DIVIDER  string = "----------------------------------------"
)

type ToolState struct {
	IsDryRun                     bool
	OverwriteExistingDirectories bool
	PDFSource                    string
	ImagesDPI                    int
	AntiAliasing                 bool
	ImagesName                   string
	ImageExt                     string
	RunHelpMode                  bool
	Limit                        int
	Verbose                      bool
	MemoryLimit                  string
	VirtualMemoryLimit           string
	DiskLimit                    string
	NotOptimize                  bool
	StateReliable                bool
}

var tool_state *ToolState

func PopulateState() error {
	tool_state = &ToolState{}

	flag.BoolVar(&tool_state.IsDryRun, "d", false, "Runs the program logging operations but not executing them")
	flag.BoolVar(&tool_state.OverwriteExistingDirectories, "overwrite", false, "If a directory with the name of the PDF exists and contains images, delete them and reprocess the PDF")
	flag.IntVar(&tool_state.ImagesDPI, "dpi", 300, "DPI of the images to be generated. For reference: 150 is low quality, 300 is default, 600 is high quality")
	flag.BoolVar(&tool_state.AntiAliasing, "aa", false, "Enable anti-aliasing for the images to be generated")
	flag.StringVar(&tool_state.ImagesName, "custom-name", "", "Custom name for the images to be generated. If not provided, the PDF name will be used")
	flag.StringVar(&tool_state.ImageExt, "ext", "webp", "Image extension to be used")
	flag.BoolVar(&tool_state.RunHelpMode, "help", false, "Help message")
	flag.BoolVar(&tool_state.Verbose, "v", false, "Verbose output")
	flag.IntVar(&tool_state.Limit, "limit", -1, "Max number of PDFs to process")
	flag.StringVar(&tool_state.MemoryLimit, "memory-limit", "2G", "Memory limit for ImageMagick operations")
	flag.StringVar(&tool_state.VirtualMemoryLimit, "vm-limit", "1G", "Virtual memory limit for ImageMagick operations")
	flag.StringVar(&tool_state.DiskLimit, "disk-limit", "5G", "Disk limit for ImageMagick operations")
	flag.BoolVar(&tool_state.NotOptimize, "not-optimize", false, "Delegates all the work to ImageMagick, which can be faster and be able to handle more PDFs, but the memory usage will enormous")
	flag.Parse()
	tool_state.StateReliable = true

	if tool_state.RunHelpMode {
		return nil
	}

	if flag.NArg() == 0 {
		return fmt.Errorf("No PDF source provided")
	}

	tool_state.PDFSource = flag.Arg(0)

	return nil
}

func stateReliableOrPanic() {
	if tool_state == nil || !tool_state.StateReliable {
		panic("Invalid attempt to access tool state before calling execution_mode.PopulateState()")
	}
}

func IsDryRun() bool {
	stateReliableOrPanic()

	return tool_state.IsDryRun
}

func OverwriteExistingDirectories() bool {
	stateReliableOrPanic()

	return tool_state.OverwriteExistingDirectories
}

func PDFSource() string {
	stateReliableOrPanic()

	return tool_state.PDFSource
}

func IsHelpMode() bool {
	stateReliableOrPanic()

	return tool_state.RunHelpMode
}

func PrintIfVerbose(message string) {
	stateReliableOrPanic()

	if tool_state.Verbose {
		fmt.Println(message)
	}
}

func IsVerbose() bool {
	stateReliableOrPanic()

	return tool_state.Verbose
}

func ImagesDPI() int {
	stateReliableOrPanic()

	return tool_state.ImagesDPI
}

func AntiAliasing() bool {
	stateReliableOrPanic()

	return tool_state.AntiAliasing
}

func ImagesName() string {
	stateReliableOrPanic()

	return tool_state.ImagesName
}

// ImageExt returns the image extension to be used. the image extension is always returned with a leading dot, e.g. ".webp"
func ImageExt() string {
	stateReliableOrPanic()

	var image_ext string = tool_state.ImageExt

	if image_ext[0] != '.' {
		image_ext = "." + image_ext
	}

	return image_ext
}

func Limit() int {
	stateReliableOrPanic()

	return tool_state.Limit
}

func MemoryLimit() string {
	stateReliableOrPanic()

	return tool_state.MemoryLimit
}

func VirtualMemoryLimit() string {
	stateReliableOrPanic()

	return tool_state.VirtualMemoryLimit
}

func DiskLimit() string {
	stateReliableOrPanic()

	return tool_state.DiskLimit
}
