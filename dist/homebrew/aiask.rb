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
      sha256 "bf15905a91cc3832450ad24d0fc618f263d220724694e9092d85a5019fed6fd0"

      def install
        bin.install "aiask-darwin-arm64" => "aiask"
      end
    else
      url "https://github.com/Hermithic/aiask/releases/download/v1.0.0/aiask-1.0.0-darwin-amd64.tar.gz"
      sha256 "a54a44f8084eb13ee7b904fef330017a9a515664515f56990b1c82e5b44a77f6"

      def install
        bin.install "aiask-darwin-amd64" => "aiask"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/Hermithic/aiask/releases/download/v1.0.0/aiask-1.0.0-linux-arm64.tar.gz"
      sha256 "31b2a2f18b87091736dce128aab92a710bb229bc38d6ad3a58877735d6f43127"

      def install
        bin.install "aiask-linux-arm64" => "aiask"
      end
    else
      url "https://github.com/Hermithic/aiask/releases/download/v1.0.0/aiask-1.0.0-linux-amd64.tar.gz"
      sha256 "4e317cf7e6e41d7a246fe72f7c811dbe0eca7ce731a98d2e542585c72692eeab"

      def install
        bin.install "aiask-linux-amd64" => "aiask"
      end
    end
  end

  test do
    assert_match "aiask version", shell_output("#{bin}/aiask version")
  end
end

