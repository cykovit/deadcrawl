# deadcrawl

## Overview

This Go-based web crawler checks for dead links on a specified website. The crawler goes through all `<a>` tags in the HTML document, verifies the validity of each link, and reports any dead links it finds. 

## Features

- Recursively checks all links on the provided webpage
- Detects and reports dead links (HTTP status code 400 or higher)
- Prints "No dead links found" if all links are valid
- Provides an interactive prompt to enter the website URL

## Prerequisites

- Go 1.22.5 or later

## Usage

1. **Run the program**:
   - Compile and run the program using:
     ```bash
     go run main.go
     ```

2. **Enter the URL**:
   - When prompted, enter the URL of the website you want to check for dead links.

3. **Interpret the output**:
   - The program will output any dead links it finds along with their HTTP status codes.
   - If no dead links are found, it will print:
     ```
     No dead links found.
     ```

## Dependencies

- **`golang.org/x/net`**
- **`net/http`**
- **`net/url`**

## License

This project is licensed under the MIT License.

Feel free to adjust the content based on your specific needs or any additional information you might want to include.

---
