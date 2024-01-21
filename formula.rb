class Prepare < Formula
    desc "Go Env Prepare is a tool to help you prepare your development environment."
    homepage "https://github.com/felipewom/go-env-prepare"
    url "https://github.com/felipewom/go-env-prepare/archive/refs/tags/0.0.1.tar.gz"
    sha256 "d65f49004a9cfcd07eac01e6f9236cbf90d70e51d6a1ddf935e820978f9bdd2f"
  
    # depends_on "go" => :build
  
    def install
      # unzip the binary
      system "tar", "xvf", "prepare.tar.gz"
      #  copy the binary to the bin folder
      system "cp", "prepare", "#{bin}/prepare"
  
      # If not using Makefile, you may copy the binary manually
      # bin.install "prepare"
    end
  
    test do
      # Add test cases if applicable
      system "#{bin}/prepare", "version"
    end
  end
  