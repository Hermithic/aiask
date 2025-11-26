# Homebrew formula for aiask
# To use: brew tap Hermithic/aiask && brew install aiask
# Or create your own tap repository: https://github.com/Hermithic/homebrew-aiask

class Aiask < Formula
  desc "AI-powered command line assistant - turn plain English into shell commands"
  homepage "https://github.com/Hermithic/aiask"
  version "2.0.1"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/Hermithic/aiask/releases/download/v2.0.1/aiask-2.0.1-darwin-arm64.tar.gz"
      sha256 "ee7eeeb1145eea6e990a46fc9de9c9e572943ee433dc16c4ac37f06664837ae2"

      def install
        bin.install "aiask-darwin-arm64" => "aiask"
        
        # Generate shell completions
        generate_completions_from_executable(bin/"aiask", "completion")
      end
    else
      url "https://github.com/Hermithic/aiask/releases/download/v2.0.1/aiask-2.0.1-darwin-amd64.tar.gz"
      sha256 "a8afe66ac957b191fe1d200ab3f7d834c02e6bb9d3677ebf10cf9bb1cef628f2"

      def install
        bin.install "aiask-darwin-amd64" => "aiask"
        
        # Generate shell completions
        generate_completions_from_executable(bin/"aiask", "completion")
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/Hermithic/aiask/releases/download/v2.0.1/aiask-2.0.1-linux-arm64.tar.gz"
      sha256 "92e4dd7b9f3dacdbc11e8a04f0a0993d63f421f5d6ba3a6b28ae7faab983f7d0"

      def install
        bin.install "aiask-linux-arm64" => "aiask"
        
        # Generate shell completions
        generate_completions_from_executable(bin/"aiask", "completion")
      end
    else
      url "https://github.com/Hermithic/aiask/releases/download/v2.0.1/aiask-2.0.1-linux-amd64.tar.gz"
      sha256 "d80eeb8c3b880b23bfce774a6652d88dfdcf45fcfb6f2780b92a98d4f733f053"

      def install
        bin.install "aiask-linux-amd64" => "aiask"
        
        # Generate shell completions
        generate_completions_from_executable(bin/"aiask", "completion")
      end
    end
  end

  test do
    assert_match "aiask version 2.0.1", shell_output("#{bin}/aiask version")
  end
end
