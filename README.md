# go-reader-cli

A command-line tool that reads files from a directory and outputs their contents to a markdown file.

## Description

This CLI tool recursively walks through a specified directory, reads the content of each file (excluding images and configuration files), and appends them to a markdown file with proper formatting. Each file's content is preceded by its relative path as a header and wrapped in code blocks.

## Installation

```bash
go install github.com/oumaoumag/go-reader-cli@latest
```

## Usage

```bash
go-cli-file-reader <directory-path> <output-md-file>
```

### Arguments:
- `<directory-path>`: Path to the directory containing files to process
- `<output-md-file>`: Path to the markdown file where content will be appended

### Examples:

```bash
# Process all files in the current project and output to docs.md
go-cli-file-reader . docs.md

# Process files from a specific directory
go-cli-file-reader ./src output.md

# Process files from an absolute path
go-cli-file-reader /home/user/projects/myproject documentation.md
```

### Output Format:

The generated markdown file will have the following structure:

```markdown
# relative/path/to/file1.ext

```js
// File content goes here
```

# relative/path/to/file2.ext

```js
// File content goes here
```


## File Filtering

The tool automatically skips:
- Hidden files and directories (starting with `.`)
- Image files (`.jpg`, `.jpeg`, `.png`, `.gif`, `.bmp`, `.tiff`, `.svg`, `.ico`)
- Configuration files (`.config`, `.ini`, `.yaml`, `.yml`, `.toml`)

## License

MIT License - See [LICENSE](LICENSE) file for details
