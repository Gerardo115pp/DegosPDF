<h1 align="center" id="title">DegosPDF</h1>

<p align="center">
    <img src="https://socialify.git.ci/Gerardo115pp/DegosPDF/image?description=1&amp;descriptionEditable=Efficient%20PDF%20to%20image%20conversion&amp;language=1&amp;name=1&amp;owner=1&amp;pattern=Circuit%20Board&amp;theme=Light" alt="project-image">
</p>

<p align="center">
    <img src="https://img.shields.io/badge/License-GPL%20v2-gree.svg" alt="shields">
    <img src="https://img.shields.io/badge/Arch_Linux-1793D1?style=for-the-badge&amp;logo=arch-linux&amp;logoColor=white" alt="shields">
    <img src="https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&amp;logo=linux&amp;logoColor=black" alt="shields">
    <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&amp;logo=go&amp;logoColor=white" alt="shields">
    <img src="https://img.shields.io/badge/ImageMagick-000000?style=for-the-badge&amp;logo=imagemagick&amp;logoColor=white" alt="shields">
    <img src="https://img.shields.io/badge/Ghostscript-5642f4?style=for-the-badge&amp;logo=ghostscript&amp;logoColor=white" alt="shields">
</p>

## Description
Basically what the project does is:

Scans directory `A/B` for PDF files, moves each one to new directories `A/<pdf-name>/<pdf-name>.pdf` and and parses all pages in the PDF to images which are stored in `A/<pdf-name>/`

    
## Features
Here're some of the project's best features:

*   PDF bulk conversion into images of multiple formats
*   Conversion process is completely resumable
*   Much less memory usage than if you use ImageMagick directly


## Installation Steps:
- ### For Arch Users download the latest release:
Get the latest release by running:
```bash
mkdir -p /tmp/degos
latest_tarball_url="$(wget -qO- https://api.github.com/repos/Gerardo115pp/DegosPDF/releases/latest | grep tarball_url | cut -d '"' -f 4)"
wget -O /tmp/degos/degos.tar.gz "$latest_tarball_url"
tar -xvzf /tmp/degos/degos.tar.gz -C /tmp/degos --strip-components=1
cd /tmp/degos
makepkg -sic
```

- ### For other distros:
Make sure you have `go` and `imagemagick` and `ghostscript` installed in your system. Then run:
```bash
go install github.com/Gerardo115pp/DegosPDF@latest
```





## Usage/Examples  

Basic usage
```bash
    degos "<path/to/pdf/file>"
```
<br/>
Specify an Mime type by extension

```bash
    degos -ext png "<path/to/pdf/file>"
```
<br/>
Specify a MIME type and a custom name for the images

```bash
    degos -ext png -custom-name "my-custom-name" "<path/to/pdf/file>"
```
The previous example will generate a sequence of images(one for each page) with the names `my-custom-name-1.png`, `my-custom-name-2.png`, ..., `my-custom-name-N.png`

## Available flags/options

Yo can alternate the default behavior of the program by using the following flags:
|flag name|type|description|
|---|---|---|
|`-aa`|Boolean|Enable anti-aliasing for the images to be generated|
|`-custom-name`|string|Custom name for the images to be generated. If not provided, the PDF name will be used|
|`-d`|Boolean|Runs the program logging operations but not executing them|
|`-disk-limit`|string| Disk limit for ImageMagick operations (default "5G")|
|`-dpi`|int|DPI of the images to be generated. For reference: 150 is low quality, 300 is default, 600 is high quality (default 300)|
|`-ext`|string|Image extension to be used (default "webp")|
|`-help`|Boolean|Help message|
|`-limit`|int|Max number of PDFs to process (default -1)|
|`-memory-limit`|string|Memory limit for ImageMagick operations (default "2G")|
|`-not-optimize`|Boolean|Delegates all the work to ImageMagick, which can be faster and be able to handle more PDFs, but the memory usage will enormous|
|`-overwrite`|Boolean|If a directory with the name of the PDF exists and contains images, delete them and reprocess the PDF|
|`-v`|Boolean|Verbose output|
|`-vm-limit`|string|Virtual memory limit for ImageMagick operations (default "1G")|



## License
[GPLv2-only](https://sources.debian.org/data/main/libc/libcamera/0.3.1-2~bpo12%2B1/LICENSES/GPL-2.0-only.txt)