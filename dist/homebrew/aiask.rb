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
      sha256 "7139e012d5409c6bc4c3a5ee108b81ca81b34f67be4a8c824cbe6181fd397949"

      def install
        bin.install "aiask-darwin-arm64" => "aiask"
      end
    else
      url "https://github.com/Hermithic/aiask/releases/download/v1.0.0/aiask-1.0.0-darwin-amd64.tar.gz"
      sha256 "ddae5cffc58a57bf299ebd914c1b066247a169d16d5e85c848d9b8c1da7f6d7a"

      def install
        bin.install "aiask-darwin-amd64" => "aiask"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/Hermithic/aiask/releases/download/v1.0.0/aiask-1.0.0-linux-arm64.tar.gz"
      sha256 "eeb767b9362ff28a5cb3831cfc744703ecf084dfab3f490b51368f02ed682848"

      def install
        bin.install "aiask-linux-arm64" => "aiask"
      end
    else
      url "https://github.com/Hermithic/aiask/releases/download/v1.0.0/aiask-1.0.0-linux-amd64.tar.gz"
      sha256 "3382cd4c11214847c87b9fd3d35b58ef6b388026976a3e57907aac2995b228fb"

      def install
        bin.install "aiask-linux-amd64" => "aiask"
      end
    end
  end

  test do
    assert_match "aiask version", shell_output("#{bin}/aiask version")
  end
end

