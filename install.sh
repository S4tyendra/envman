#!/bin/zsh

# base
base_url="https://envman.devh.in/install"


detect_arch() {
  case "$(uname -m)" in
    x86_64)
      echo "linux/amd64"
      ;;
    armv7l)
      echo "linux/arm"
      ;;
    aarch64)
      echo "linux/arm64"
      ;;
    i386)
      echo "linux/386"
      ;;
    ppc64le)
      echo "linux/ppc64le"
      ;;
    mips64)
      echo "linux/mips64"
      ;;
    mips64le)
      echo "linux/mips64le"
      ;;
    *)
      echo "unsupported"
      ;;
  esac
}

#  download and install
install_binary() {
  local arch=$1
  local url="${base_url}/${arch}"
  local output_file="envman"

  echo "Downloading binary for architecture: $arch from $url..."
  curl -L -o $output_file $url

  if [ $? -ne 0 ]; then
    echo "Failed to download binary for architecture: $arch"
    return 1
  fi

  chmod +x $output_file
  sudo mv $output_file /usr/local/bin/envman

  if [ $? -ne 0 ]; then
    echo "Failed to install binary for architecture: $arch"
    return 1
  fi

  echo "Successfully installed envman for architecture: $arch"
  return 0
}

# detect
arch=$(detect_arch)
if [ "$arch" = "unsupported" ]; then
  echo "Unsupported architecture: $(uname -m)"
  exit 1
fi

install_binary $arch