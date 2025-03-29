[![Open in Dev Containers](https://img.shields.io/static/v1?label=Dev%20Containers&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/philipf/gt)
![Build workflow](https://github.com/philipf/gt/actions/workflows/go.yml/badge.svg)

# GT - Go Time

GT (Go Time) is a CLI tool designed to improve productivity. It provides various utilities for task management, time tracking, and calendar management.

## Features

- **GTD (Getting Things Done)**: Create and manage tasks in an Obsidian Kanban board
- **Toggl Integration**: Track time, manage projects and clients, and generate reports
- **Calendar Management**: View free time slots in your calendar

## Installation

### Prerequisites

- Go 1.23 or higher
- Git

### Install from Source

```bash
go install github.com/philipf/gt@latest
```

### Configuration

Running `gt` for the first time will create a config file in the user's home directory at `~/.gt/config.yaml`. This file contains settings for:

- GTD Kanban board location
- Toggl API credentials
- Personal information
- AI settings (if using AI features)

You can edit this file manually to customize your settings.

## Usage

### General Commands

- `gt version`: Display the current version of GT
- `gt settings`: Display the current configuration settings

### GTD Commands

- `gt gtd add`: Add a new action to your GTD system
  - Flags:
    - `--ai`: Use AI to assist with creating a new action

- `gt gtd purge`: Remove archived actions from your Kanban board
  - Flags:
    - `-d, --dry-run`: Simulate the purge without actually deleting files
    - `-a, --clean-all`: Clean all notes, not only the Kanban board

### Toggl Commands

- `gt toggl time`: List time entries
  - Flags:
    - `--today`: Show entries for today
    - `--yesterday`: Show entries for yesterday
    - `--eyesterday`: Show entries for the day before yesterday
    - `--thisweek`: Show entries for this week
    - `--lastweek`: Show entries for last week
    - `-s, --startDate`: Specify a start date (format: YYYY/MM/DD)
    - `-e, --endDate`: Specify an end date (format: YYYY/MM/DD)
    - `--csv`: Export results to CSV format

- `gt toggl report`: Generate time tracking reports
  - Flags:
    - `-t, --text`: Generate a text report
    - `-j, --json`: Generate a JSON report
    - `--ot`: Write text report to a file (default: /tmp/time.txt)
    - `--oj`: Write JSON report to a file (default: /tmp/time.json)
    - Time period flags (same as `toggl time` command)

- `gt toggl stop`: Stop the current time entry

- `gt toggl resume`: Resume the last time entry

- `gt toggl edit`: Edit the description of the running time entry
  - Flags:
    - `-d, --description`: New description for the time entry

- `gt toggl project list`: List Toggl projects
  - Flags:
    - `--validate`: Validate that projects match the naming convention
    - `--includeArchived`: Include archived projects
    - `-c, --clientId`: Filter by client ID
    - `-n, --name`: Filter by project name

- `gt toggl client list`: List Toggl clients
  - Flags:
    - `-f, --filter`: Filter clients by name

### Calendar Commands

- `gt free`: View free time slots in your calendar (currently in development)

## Development

### Prerequisites for Development

- Go 1.23 or higher
- Git
- VSCode with Dev Containers extension (optional, for containerized development)

### Setting Up Development Environment

1. Clone the repository:
   ```bash
   git clone https://github.com/philipf/gt.git
   cd gt
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   go build
   ```

4. Run tests:
   ```bash
   ./run_tests.sh
   ```

## Contributing

Contributions are welcome! Here's how you can contribute:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/my-new-feature`
5. Submit a pull request

Please make sure your code passes all tests and follows the project's coding style.

## License

This project is licensed under the terms of the LICENSE file included in the repository.
