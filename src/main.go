package main

import (
	"bulk_pdf_to_images/degos_helpers"
	"bulk_pdf_to_images/execution_mode"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	pdf_utils "github.com/ledongthuc/pdf"
)

var ErrSkipDirectory = errors.New("Skip directory") // Should only skip the current directory, not stop the program

func DetectImageMagick() bool {
	command := exec.Command("magick", "--version")

	err := command.Run()

	return err == nil
}

// Returns the command arguments, based on the execution_mode configuration, it skips the [input-pdf] [output-images-name].
func CompileCommandArgs() []string {
	dpi_string := strconv.Itoa(execution_mode.ImagesDPI())
	memory_limit := execution_mode.MemoryLimit()
	virtual_memory_limit := execution_mode.VirtualMemoryLimit()
	disk_limit := execution_mode.DiskLimit()

	var command_args []string = []string{
		"magick",
		"convert",
		"-density", dpi_string,
		"-limit", "memory", memory_limit,
		"-limit", "map", virtual_memory_limit,
		"-limit", "disk", disk_limit,
	}

	if execution_mode.IsVerbose() {
		command_args = append(command_args, "-verbose")
	}

	if execution_mode.AntiAliasing() {
		command_args = append(command_args, "-antialias")
	}

	return command_args
}

func RobustConvertPDFToImages(pdf_path, storage_path string) error {
	var command *exec.Cmd
	var err error

	var output_path string = filepath.Join(storage_path, execution_mode.ImagesName())

	if execution_mode.ImagesName() == "" {
		var pdf_name string = filepath.Base(pdf_path)

		pdf_name = strings.TrimSuffix(pdf_name, filepath.Ext(pdf_name))

		output_path = filepath.Join(storage_path, pdf_name)
	}

	output_path = fmt.Sprintf("%s-%%d%s", output_path, execution_mode.ImageExt())

	execution_mode.PrintIfVerbose(fmt.Sprintf("Converting PDF<%s> to images\n", pdf_path))
	execution_mode.PrintIfVerbose(fmt.Sprintf("Output path: %s\n", output_path))

	command_args := CompileCommandArgs()

	command_args = append(command_args, pdf_path, output_path)

	command = exec.Command(command_args[0], command_args[1:]...)

	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	execution_mode.PrintIfVerbose(fmt.Sprintf("Command: %s", command.String()))

	if !execution_mode.IsDryRun() {
		err = command.Run()
		if err != nil {
			return err
		}
	}

	execution_mode.PrintIfVerbose(fmt.Sprintf("Conversion finished\n"))

	return err
}

func LightConvertPDFToImages(pages_count int, pdf_path, storage_path string) error {

	var output_path string = filepath.Join(storage_path, execution_mode.ImagesName())

	if execution_mode.ImagesName() == "" {
		var pdf_name string = filepath.Base(pdf_path)

		pdf_name = strings.TrimSuffix(pdf_name, filepath.Ext(pdf_name))

		output_path = filepath.Join(storage_path, pdf_name)
	}

	execution_mode.PrintIfVerbose(fmt.Sprintf("Converting %d pages", pages_count))
	var zero_indexed_pages_count int = pages_count - 1

	var static_command_args []string = CompileCommandArgs()

	for k := 0; k < zero_indexed_pages_count; k++ {
		fmt.Printf("Page %d of %d", k+1, pages_count)
		var iter_file_path string = fmt.Sprintf("%s[%d]", pdf_path, k)
		var iter_output_path string = fmt.Sprintf("%s-%d%s", output_path, k, execution_mode.ImageExt())

		if file_exists := degos_helpers.PathExists(iter_output_path); file_exists {
			execution_mode.PrintIfVerbose(fmt.Sprintf("Skipping page %d", k))
			continue
		}

		command_args := append(static_command_args, iter_file_path, iter_output_path)

		command := exec.Command(command_args[0], command_args[1:]...)

		command.Stderr = os.Stderr
		command.Stdout = os.Stdout

		execution_mode.PrintIfVerbose(fmt.Sprintf("Command: %s", command.String()))

		if !execution_mode.IsDryRun() {
			err := command.Run()
			if err != nil {
				return fmt.Errorf("Error while converting page %d: %s", k, err)
			}
		}

		execution_mode.PrintIfVerbose(fmt.Sprintf("Page %d converted", k))
	}

	return nil
}

func GetPDFPages(pdf_path string) (int, error) {
	f, pdf_file, err := pdf_utils.Open(pdf_path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	pages := pdf_file.NumPage()

	return pages, nil
}

func GetAllPDFsPaths(pdf_source string) ([]string, error) {
	var pdf_paths []string

	directory_entries, err := os.ReadDir(pdf_source)
	if err != nil {
		return pdf_paths, err
	}

	for _, entry := range directory_entries {
		if entry.IsDir() {
			continue
		}

		path := fmt.Sprintf("%s/%s", pdf_source, entry.Name())
		if is_pdf, _ := degos_helpers.IsPDF(path); is_pdf {
			pdf_paths = append(pdf_paths, path)
		}
	}

	return pdf_paths, nil
}

func VerifyExecutionMode() {
	err := execution_mode.PopulateState()
	if err != nil {
		fmt.Printf("Error on boot: %s\n", err)
		degos_helpers.PrintUsage()
		os.Exit(1)
	}

	if execution_mode.IsHelpMode() {
		degos_helpers.PrintHelp()
		os.Exit(0)
	}

	pdf_source := execution_mode.PDFSource()

	if !degos_helpers.IsDirectory(pdf_source) {
		fmt.Println("PDF Source must be an existing directory")
		os.Exit(1)
	}
}

func RunBulkConversion(pdf_paths []string, pdf_storage_root string) error {
	var processing_limit int = execution_mode.Limit()
	if processing_limit < 0 || processing_limit > len(pdf_paths) {
		processing_limit = len(pdf_paths)
	}
	for h := 0; h < processing_limit; h++ {
		pdf_path := pdf_paths[h]

		fmt.Printf("PDF %d of %d\n", h+1, processing_limit)

		if h >= processing_limit {
			fmt.Printf("Stopping at limit: %d\n", processing_limit)
			break
		}

		execution_mode.PrintIfVerbose(fmt.Sprintf("%s\nProcessing PDF<%s>\n", execution_mode.OUTPUT_THICK_DIVIDER, pdf_path))

		// Get PDF pages count
		pdf_pages_count, err := GetPDFPages(pdf_path)
		if err != nil {
			fmt.Printf("Error while getting PDF pages: %s\n", err)
			return err
		}
		execution_mode.PrintIfVerbose(fmt.Sprintf("PDF<%s> has %d pages", pdf_path, pdf_pages_count))

		// Move PDF to storage root
		new_pdf_path, err := MovePDF(pdf_path, pdf_storage_root)
		if err != nil {
			if errors.Is(err, ErrSkipDirectory) {
				fmt.Printf("Skipping directory\n")
				continue
			} else {
				return err
			}
		}

		execution_mode.PrintIfVerbose(fmt.Sprintf("PDF<%s> moved", new_pdf_path))

		pages_storage_path := degos_helpers.GetParentDirectory(new_pdf_path)

		execution_mode.PrintIfVerbose(fmt.Sprintf("Pages Storage Path: %s", pages_storage_path))

		if execution_mode.Optimize() || pdf_pages_count < 15 {
			err = LightConvertPDFToImages(pdf_pages_count, new_pdf_path, pages_storage_path)
			if err != nil {
				return err
			}
		} else {
			err = RobustConvertPDFToImages(new_pdf_path, pages_storage_path)
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("%s\nDone\n", execution_mode.OUTPUT_THIN_DIVIDER)

	return nil
}

func MovePDF(pdf_path string, pdf_storage_root string) (string, error) {
	var pdf_directory string

	pdf_directory = degos_helpers.GetPDFDirectoryName(pdf_path)
	pdf_directory = filepath.Join(pdf_storage_root, pdf_directory)
	execution_mode.PrintIfVerbose(fmt.Sprintf("PDF Directory: %s", pdf_directory))

	if pdf_directory == "" {
		return "", fmt.Errorf("Error while getting PDF directory name")
	}

	execution_mode.PrintIfVerbose(fmt.Sprintf("PDF Root Directory: %s", pdf_directory))
	new_pdf_path := filepath.Join(pdf_directory, filepath.Base(pdf_path))
	execution_mode.PrintIfVerbose(fmt.Sprintf("New PDF Path: %s", new_pdf_path))

	if execution_mode.IsDryRun() {
		fmt.Printf("Would move PDF<%s> to <%s>\n", pdf_path, new_pdf_path)
		return new_pdf_path, nil
	}

	err := degos_helpers.CreateDirectoryIfNotExists(pdf_directory)
	if err != nil {
		return "", err
	}

	var images_extension string = execution_mode.ImageExt()

	directory_has_pages, err := degos_helpers.DirectoryHasPages(pdf_directory, images_extension)
	if err != nil {
		execution_mode.PrintIfVerbose(fmt.Sprintf("Error while checking for PNGs in directory: %s", err))
		return "", err
	}

	if directory_has_pages {
		if execution_mode.OverwriteExistingDirectories() {
			err = degos_helpers.CleanDirectoryPreviousPages(pdf_directory, images_extension)
			if err != nil {
				return "", err
			}
		} else {
			return "", ErrSkipDirectory
		}
	}

	err = os.Rename(pdf_path, new_pdf_path)

	return new_pdf_path, err
}

func RunSetup() {
	VerifyExecutionMode()

	image_magick_installed := DetectImageMagick()

	if !image_magick_installed {
		fmt.Println("ImageMagick is not installed, convert is required to run this program")
		os.Exit(1)
	}

	execution_mode.PrintIfVerbose(fmt.Sprintln("ImageMagick is detected"))
}

func main() {
	RunSetup()

	var pdf_paths []string
	var err error

	pdf_paths, err = GetAllPDFsPaths(execution_mode.PDFSource())
	if err != nil {
		fmt.Printf("Error while getting PDFs paths: %s", err)
		os.Exit(1)
	}

	var pdf_storage_root string

	pdf_storage_root = degos_helpers.GetParentDirectory(execution_mode.PDFSource())

	if pdf_storage_root == "" {
		fmt.Printf("Parsing directory A from B<%s> resulted in: ''", execution_mode.PDFSource())
		os.Exit(1)
	}

	if execution_mode.IsDryRun() {
		fmt.Println("Dry run, no operations will be executed")
	}

	fmt.Sprintf("PDF Source: %s", execution_mode.PDFSource())
	fmt.Printf("PDF Storage Root: %s\n", pdf_storage_root)
	fmt.Printf("PDFs found: %d\n", len(pdf_paths))

	err = RunBulkConversion(pdf_paths, pdf_storage_root)
	if err != nil {
		fmt.Printf("Error while running bulk conversion:\n%s", err)
		os.Exit(1)
	}
}
