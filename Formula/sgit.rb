class Sgit < Formula
  desc "Solar LLM-powered git wrapper with AI-enhanced commit messages"
  homepage "https://github.com/hunkim/sgit"
  url "https://github.com/hunkim/sgit/archive/v0.1.0.tar.gz"
  sha256 "REPLACE_WITH_ACTUAL_SHA256"
  license "MIT"

  depends_on "go" => :build
  depends_on "git"

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"sgit"
  end

  test do
    # Test that the binary works
    assert_match "Solar LLM-powered git wrapper", shell_output("#{bin}/sgit --help")
    
    # Test git passthrough functionality
    system "git", "init"
    assert_match "fatal: your current branch 'main' does not have any commits yet", 
                 shell_output("#{bin}/sgit status", 1)
  end
end 