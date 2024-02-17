# Nexus Repo Migration (Golang Version)

## Introduction
The Nexus Repo Migration tool, developed in Go, is designed to facilitate the migration of artifacts between two Nexus repositories. This utility aims to streamline the transfer process, ensuring efficient data movement with minimal downtime and effort.

## Features
- Concurrent artifact upload and download to expedite the migration process.
- Automatic retry mechanism for handling temporary network failures or server errors.
- Configurable HTTP client timeout to circumvent load balancer limitations.
- User-friendly command-line interface for ease of use.

## Requirements
Before using this tool, ensure you have the following:
- Go version 1.15 or later installed on your system.
- Access to both source and target Nexus repositories with appropriate read/write permissions.

## Installation
Clone the repository from GitHub:
```
git clone https://github.com/yourusername/Nexus-Repo-Migration.git
cd Nexus-Repo-Migration
git checkout Golang-version
```
Build the tool using Go:
```
go build -o nexus-migrator
```

## Usage
To start the migration, run the following command after configuring the necessary variables in the code (`sourceNexus`, `sourceRepo`, `sourceUser`, `sourcePassword`, `targetNexus`, `targetRepo`, `targetUser`, `targetPassword`):
```
./nexus-migrator
```

## Configuration
Adjust the configuration variables in the code to match your source and target Nexus repository details, including access credentials.

## Example
With the variables configured and the program compiled, execute it to begin the migration:
```
./nexus-migrator
```
The tool will provide real-time feedback during the migration process, indicating both successfully migrated artifacts and any errors encountered.
