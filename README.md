# brewfile-desc

`brewfile-desc` add descriptions of formulae to Brewfile.

## Usage

``` console
$ cat path/to/Brewfile
tap "golangci/tap"
tap "k1low/tap"
tap "mas-cli/tap"

brew "act"
brew "glib"
brew "autoconf"
brew "automake"
brew "docker"
brew "docker-compose"
brew "git"
brew "make"
brew "golangci/tap/golangci-lint"
brew "k1low/tap/octocov"
brew "mas-cli/tap/mas"
cask "1password"
cask "aws-vault"
cask "bitbar"
cask "deepl"
cask "font-ipafont"
cask "iterm2"
$ brewfile-desc path/to/Brewfile
tap "golangci/tap"
tap "k1low/tap"
tap "mas-cli/tap"

brew "act"                        # Run your GitHub Actions locally 🚀
brew "glib"                       # Core application library for C
brew "autoconf"                   # Automatic configure script builder
brew "automake"                   # Tool for generating GNU Standards-compliant Makefiles
brew "docker"                     # Pack, ship and run any application as a lightweight container
brew "docker-compose"             # Isolated development environments using Docker
brew "git"                        # Distributed revision control system
brew "make"                       # Utility for directing compilation
brew "golangci/tap/golangci-lint" # Fast linters runner for Go.
brew "k1low/tap/octocov"          # octocov is a toolkit for collecting code metrics.
brew "mas-cli/tap/mas"            # Mac App Store command-line interface
cask "1password"                  # Password manager that keeps all passwords secure behind one password
cask "aws-vault"                  # Securely stores and accesses AWS credentials in a development environment
cask "bitbar"                     # Utility to display the output from any script or program in the menu bar
cask "deepl"                      # Trains AIs to understand and translate texts
cask "font-ipafont"               # IPA Fonts
cask "iterm2"                     # Terminal emulator as alternative to Apple's Terminal app
```

## Install

``` console
$ brew install k1LoW/tap/brewfile-desc
```

