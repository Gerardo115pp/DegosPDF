package main

import (
	"bulk_pdf_to_images/execution_mode"
	bpti_helpers "bulk_pdf_to_images/helpers"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var ErrSkipDirectory = errors.New("Skip directory") // Should only skip the current directory, not stop the program

func DetectImageMagick() bool {
	command := exec.Command("magick", "--version")

	err := command.Run()

	return err == nil
}

func ConvertPDFToImages(pdf_path, storage_path string) error {
	var command *exec.Cmd
	var err error

	var output_path string = filepath.Join(storage_path, execution_mode.ImagesName())

	if execution_mode.ImagesName() == "" {
		var pdf_name string = filepath.Base(pdf_path)

		pdf_name = strings.TrimSuffix(pdf_name, filepath.Ext(pdf_name))

		output_path = filepath.Join(storage_path, pdf_name)
	}

	output_path = fmt.Sprintf("%s-%%d%s", output_path, execution_mode.ImageExt())

	fmt.Printf("Converting PDF<%s> to images\n", pdf_path)
	fmt.Printf("Output path: %s\n", output_path)

	dpi_string := strconv.Itoa(execution_mode.ImagesDPI())

	if execution_mode.IsVerbose() {
		command = exec.Command("magick", "convert", "-verbose", "-density", dpi_string, pdf_path, output_path)
	} else {
		command = exec.Command("magick", "convert", "-density", dpi_string, pdf_path, output_path)
	}

	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	execution_mode.PrintIfVerbose(fmt.Sprintf("Command: %s", command.String()))

	if !execution_mode.IsDryRun() {
		err = command.Run()
		if err != nil {
			return err
		}
	}

	fmt.Printf("Conversion finished\n")

	return err
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
		if is_pdf, _ := bpti_helpers.IsPDF(path); is_pdf {
			pdf_paths = append(pdf_paths, path)
		}
	}

	return pdf_paths, nil
}

func VerifyExecutionMode() {
	err := execution_mode.PopulateState()
	if err != nil {
		fmt.Printf("Error on boot: %s\n", err)
		bpti_helpers.PrintUsage()
		os.Exit(1)
	}

	if execution_mode.IsHelpMode() {
		bpti_helpers.PrintHelp()
		os.Exit(0)
	}

	pdf_source := execution_mode.PDFSource()

	if !bpti_helpers.IsDirectory(pdf_source) {
		fmt.Println("PDF Source must be an existing directory")
		os.Exit(1)
	}
}

func RunBulkConversion(pdf_paths []string, pdf_storage_root string) error {
	var processing_limit int = execution_mode.Limit()
	for h, pdf_path := range pdf_paths {
		if processing_limit >= 1 && h >= processing_limit {
			fmt.Printf("Stopping at limit: %d\n", processing_limit)
			break
		}

		fmt.Printf("%s\nProcessing PDF<%s>\n", execution_mode.OUTPUT_THICK_DIVIDER, pdf_path)

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

		pages_storage_path := bpti_helpers.GetParentDirectory(new_pdf_path)

		execution_mode.PrintIfVerbose(fmt.Sprintf("Pages Storage Path: %s", pages_storage_path))

		err = ConvertPDFToImages(new_pdf_path, pages_storage_path)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s\nDone\n", execution_mode.OUTPUT_THIN_DIVIDER)

	return nil
}

func MovePDF(pdf_path string, pdf_storage_root string) (string, error) {
	var pdf_directory string

	pdf_directory = bpti_helpers.GetPDFDirectoryName(pdf_path)
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

	err := bpti_helpers.CreateDirectoryIfNotExists(pdf_directory)
	if err != nil {
		return "", err
	}

	var images_extension string = execution_mode.ImageExt()

	directory_has_pages, err := bpti_helpers.DirectoryHasPages(pdf_directory, images_extension)
	if err != nil {
		execution_mode.PrintIfVerbose(fmt.Sprintf("Error while checking for PNGs in directory: %s", err))
		return "", err
	}

	if directory_has_pages {
		if execution_mode.OverwriteExistingDirectories() {
			err = bpti_helpers.CleanDirectoryPreviousPages(pdf_directory, images_extension)
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

	fmt.Println("ImageMagick is detected")
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

	pdf_storage_root = bpti_helpers.GetParentDirectory(execution_mode.PDFSource())

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
