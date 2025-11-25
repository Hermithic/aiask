# Homebrew formula for aiask
# To use: brew tap Hermithic/aiask && brew install aiask
# Or create your own tap repository: https://github.com/Hermithic/homebrew-aiask

class Aiask < Formula
  desc "AI-powered command line assistant"
  homepage "https://github.com/Hermithic/aiask"
  version "1.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/Hermithic/aiask/releases/download/v1.0.0/aiask-1.0.0-darwin-arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_DARWIN_ARM64"

      def install
        bin.install "aiask-darwin-arm64" => "aiask"
      end
    else
      url "https://github.com/Hermithic/aiask/releases/download/v1.0.0/aiask-1.0.0-darwin-amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_DARWIN_AMD64"

      def install
        bin.install "aiask-darwin-amd64" => "aiask"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/Hermithic/aiask/releases/download/v1.0.0/aiask-1.0.0-linux-arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_ARM64"

      def install
        bin.install "aiask-linux-arm64" => "aiask"
      end
    else
      url "https://github.com/Hermithic/aiask/releases/download/v1.0.0/aiask-1.0.0-linux-amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_AMD64"

      def install
        bin.install "aiask-linux-amd64" => "aiask"
      end
    end
  end

  test do
    assert_match "aiask version", shell_output("#{bin}/aiask version")
  end
end

