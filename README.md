# envman

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

`envman` is a command-line tool that simplifies the management of environment variables using profiles. It allows you to create, edit, list, view, and delete different sets of environment variables, streamlining your workflow when working on multiple projects or environments.

> **Note:** This tool is currently only available on Linux systems. Support for macOS and Windows is coming soon.


## ✨ Features

-   **Profile Creation:** Easily create new environment variable profiles.
-   **Profile Editing:** Edit existing profiles using a user-friendly text-based editor.
-   **Profile Listing:** View a summary of available profiles, including the number of entries and last modification time.
-   **Profile Viewing:** Inspect a profile's content with syntax highlighting in a read-only viewer.
-   **Profile Deletion:** Remove profiles that are no longer needed.
-   **Interactive Interface:** Provides a smooth text-based interface for managing profiles.

## 📦 TODO

- [x] Load profiles and environment variables to session.
- [ ] Add Comments to code.
- [ ] Add support for exporting and importing/exporting profiles.
- [ ] Implement Encrypted profiles.

## 🚀 Installation

```shell
curl -L http://envman.devh.in/install | bash
```

## 📚 Usage

### **Help**

```shell
❯ envman -h
envman v0.1.0
Usage:
  envman [command] [flags]

Available Commands:
  init        Initialize envman in your shell
  profile     Manage environment profiles
  load        Load a profile into the current shell

Profile Subcommands:
  create      Create a new environment profile
  edit        Edit an existing environment profile
  show        Display profile contents
  delete      Delete an environment profile
  list        List all available profiles

Examples:
  # Initialize envman
  $ envman init

  # Profile Management
  $ envman profile create server-test    # Create new profile
  $ envman profile list                  # List all profiles
  $ envman profile show server-test      # Show profile contents
  $ envman profile edit server-test      # Edit existing profile
  $ envman profile delete server-test    # Delete profile

  # Load Profile
  $ envman load server-test             # Load profile into current shell

Flags:
  -h, --help    Display help information
  -v, --version Display version information

```


**Quick Start Example:**

To create a new profile named `dev`:

```bash
envman profile create dev
```

To edit the newly created profile:

```bash
envman profile edit dev
```

To view all profiles:

```bash
envman profile list
```

To see the contents of a profile:

```bash
envman profile view dev
```

## 🔨 Building from Source

If you prefer to build `envman` yourself, here are the instructions.

**📋 Prerequisites:**

-   Go (version 1.20 or later)

**Build Steps:**

1.  Clone the repository:
    ```bash
    git clone https://github.com/s4tyendra/envman.git
    cd envman
    ```
2.  Build the binary:
    ```bash
    go build -o envman main.go
    ```
3. (Optional) Install the binary to `/usr/local/bin`:
    ```bash
    sudo cp envman /usr/local/bin
    ```

**System Requirements:**

- Any system with Go support

## 💻 Development

For those interested in contributing to `envman`, here's how to set up the development environment.

**Development Environment Setup:**

1.  Ensure that Go is installed and properly configured.
2.  Clone the repository as described in the "Building from Source" section.

**Build Commands:**

```bash
go build -o envman main.go
```

**Test Commands:**

Currently, no automated testing is implemented, but we encourage manual testing of new features and bug fixes.

**Directory Structure Overview:**

```
envman/
├── main.go          # Main entry point of the application
├── helpers.go       # Helper functions and initial setup
├── models.go        # Definitions of app's data structures
├── outputs.go       # Constants and formatted output strings
├── profile_create.go  # Profile creation logic
├── profile_delete.go  # Profile deletion logic
├── profile_edit.go    # Profile editing logic
├── profile_list.go    # Profile listing logic
├── profile_view.go    # Profile viewing logic
```

## 🤝 Contributing

We welcome contributions to `envman`! Here's how you can help:

1.  **Fork the repository** on GitHub.
2.  **Make your changes** on your fork.
3.  **Submit a pull request** with a clear description of your changes.

**Code Style Guidelines:**
- Adhere to standard Go coding conventions.
- Provide descriptive commit messages.

**Issue Reporting:**

- Report bugs, propose feature requests, or ask questions in the issues.
- Include steps to reproduce any bugs.

**Pull Request Process:**

- Keep pull requests focused on single issues or features.
- Provide detailed description about the changes you made.
- Respond to any review comments in a timely manner.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/licenses/MIT) file for details.
