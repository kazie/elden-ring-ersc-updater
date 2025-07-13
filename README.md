## Elden Ring Seamless Co-op Updater

I have created a simple program to help with updating the [Seamless Co-op](https://www.nexusmods.com/eldenring/mods/510) 
mod for Elden Ring, by using the latest release from the [github mirror](https://github.com/LukeYui/EldenRingSeamlessCoopRelease/releases)

One can typically just update it by hand, but I made this simple program just to help my own lazy ass.

![](./elden.gif)

### Config

The config is done in a toml file `config.toml`, that shall be in the same folder as the running program.

```toml
# The last downloaded version v1.7.8 for example
current_version = "v1.7.8"
# The full path to your elden ring game folder
elden_ring_game_path = "/home/username/.local/share/Steam/steamapps/common/ELDEN RING/Game/"
# Github API key with reading public repositories, make one under developer settings -> personal access token
# This is optional. If left empty, the application will use GitHub's API without authentication (rate limited by GitHub)
github_read_token = ""
# If it shall ignore writing the ini file (writing ini file may reset password)
ignore_ini_file = true
```

> NOTE: When using Microsoft Windows, file paths need to be escaped like so:  
> `C:\\some\\path\\to\\somewhere` as they use `\ ` natively, but needs to be escaped, thus `\\` 

### Development

#### Testing

The project includes unit tests for key functions. To run the tests and generate a coverage report, use the following commands:

```bash
# Run tests with race detection and generate coverage profile
go test -race -coverprofile=coverage.txt -covermode=atomic ./...

# Display coverage report in terminal
go tool cover -func=coverage.txt

# Generate HTML coverage report
go tool cover -html=coverage.txt -o coverage.html
```

This will:
1. Run all tests with race detection
2. Generate a coverage profile (coverage.txt)
3. Display a summary of the coverage in the terminal
4. Create an HTML coverage report (coverage.html) that you can view in your browser

#### Continuous Integration

The project uses GitHub Actions for continuous integration. The workflow:
1. Builds the application for Linux, Windows, and MacOS
2. Runs tests with coverage
3. Uploads the coverage report to Codecov

You can view the coverage report on the Codecov dashboard after pushing your changes.
