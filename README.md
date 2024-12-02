# GitHub Janitor

## Overview

This project is a command-line tool to archive GitHub repositories that are considered empty. It scans all repositories
in a specified organization, identifies repositories that contain only a `README.md`, `LICENSE`, or `.gitignore` file,
and have not been updated in the last month, and archives them.

## Features

- **Scan Repositories**: Scans all repositories in a specified GitHub organization.
- **Identify Empty Repositories**: Identifies repositories that are considered empty based on specific criteria.
- **Archive Repositories**: Archives the identified empty repositories.

## Installation

1. Install the application using `go install`:
    ```sh
    go install github.com/wernerstrydom/github-janitor@latest
    ```

## Usage


1. Run the `scan` command:
    ```sh
    github-janitor scan
    ```
   
2. Run the `archive` command:
   ```sh
   github-janitor archive
   ```
