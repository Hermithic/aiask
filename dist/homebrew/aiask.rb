# Homebrew formula for aiask
# To use: brew tap Hermithic/aiask && brew install aiask
# Or create your own tap repository: https://github.com/Hermithic/homebrew-aiask

class Aiask < Formula
  desc "AI-powered command line assistant - turn plain English into shell commands"
  homepage "https://github.com/Hermithic/aiask"
  version "2.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/Hermithic/aiask/releases/download/v2.0.0/aiask-2.0.0-darwin-arm64.tar.gz"
      sha256 "af2486d283cd66a3db466b2f387994ece3dd91c83b1b0362aeea488af76f495b"

      def install
        bin.install "aiask-darwin-arm64" => "aiask"
        
        # Generate shell completions
        generate_completions_from_executable(bin/"aiask", "completion")
      end
    else
      url "https://github.com/Hermithic/aiask/releases/download/v2.0.0/aiask-2.0.0-darwin-amd64.tar.gz"
      sha256 "0c3458536cf5c78f4a70d4c3b7f9dd608255dcceaab38cdc2e15cc09ce728ed4"

      def install
        bin.install "aiask-darwin-amd64" => "aiask"
        
        # Generate shell completions
        generate_completions_from_executable(bin/"aiask", "completion")
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/Hermithic/aiask/releases/download/v2.0.0/aiask-2.0.0-linux-arm64.tar.gz"
      sha256 "14025c11e6b4217390a5bad0c1149b901b6973897f006e42ed2fbd894b79c1c8"

      def install
        bin.install "aiask-linux-arm64" => "aiask"
        
        # Generate shell completions
        generate_completions_from_executable(bin/"aiask", "completion")
      end
    else
      url "https://github.com/Hermithic/aiask/releases/download/v2.0.0/aiask-2.0.0-linux-amd64.tar.gz"
      sha256 "f18f4c2029d355bc737e1a562bf5247c6897c61dece43441ce8ed33fc2442562"

      def install
        bin.install "aiask-linux-amd64" => "aiask"
        
        # Generate shell completions
        generate_completions_from_executable(bin/"aiask", "completion")
      end
    end
  end

  test do
    assert_match "aiask version 2.0.0", shell_output("#{bin}/aiask version")
  end
end
