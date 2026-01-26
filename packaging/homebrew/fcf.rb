# typed: false
# frozen_string_literal: true

# Homebrew formula for fcf - Find File or Folder
# To use: brew install ReggieAlbiosA/tap/fcf
#
# To set up your own tap:
# 1. Create a repo named "homebrew-tap" on GitHub
# 2. Copy this file to Formula/fcf.rb in that repo
# 3. Users can then: brew tap ReggieAlbiosA/tap && brew install fcf

class Fcf < Formula
  desc "Fast, interactive file and folder finder with parallel search"
  homepage "https://github.com/ReggieAlbiosA/fcf"
  version "3.2.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/ReggieAlbiosA/fcf/releases/download/v#{version}/fcf-darwin-arm64"
      sha256 "PLACEHOLDER_SHA256_DARWIN_ARM64"
    else
      url "https://github.com/ReggieAlbiosA/fcf/releases/download/v#{version}/fcf-darwin-amd64"
      sha256 "PLACEHOLDER_SHA256_DARWIN_AMD64"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/ReggieAlbiosA/fcf/releases/download/v#{version}/fcf-linux-arm64"
      sha256 "PLACEHOLDER_SHA256_LINUX_ARM64"
    else
      url "https://github.com/ReggieAlbiosA/fcf/releases/download/v#{version}/fcf-linux-amd64"
      sha256 "PLACEHOLDER_SHA256_LINUX_AMD64"
    end
  end

  def install
    binary_name = "fcf"

    # Rename downloaded binary to fcf
    downloaded_file = Dir["fcf-*"].first || "fcf"
    mv downloaded_file, binary_name if File.exist?(downloaded_file)

    bin.install binary_name
  end

  def caveats
    <<~EOS
      To enable directory navigation, add the shell integration to your config:

      For Bash (~/.bashrc) or Zsh (~/.zshrc):
        fcf() {
          local nav_file="/tmp/fcf_nav_path_$(id -u)"
          rm -f "$nav_file"
          command fcf "$@"
          if [[ -f "$nav_file" ]]; then
            local target
            target=$(cat "$nav_file")
            rm -f "$nav_file"
            if [[ -d "$target" ]]; then
              cd "$target" || return
            fi
          fi
        }

      For Fish (~/.config/fish/config.fish):
        function fcf
          set nav_file /tmp/fcf_nav_path_(id -u)
          rm -f $nav_file
          command fcf $argv
          if test -f $nav_file
            set target (cat $nav_file)
            rm -f $nav_file
            if test -d $target
              cd $target
            end
          end
        end

      Or run: fcf install --shell-only
      Then reload your shell config.
    EOS
  end

  test do
    assert_match "fcf", shell_output("#{bin}/fcf --help")
  end
end
