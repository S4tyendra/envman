#!/bin/zsh

# List of target architectures
targets=(
  "linux/amd64"
  "linux/arm64"
  "linux/arm"
  "linux/386"
  "linux/ppc64le"
  "linux/mips64"
  "linux/mips64le"
)

# Project name
project_name="envman"

# Create output directory
output_dir="build"
mkdir -p $output_dir

# Build for each target
for target in "${targets[@]}"; do
  IFS="/" read -r GOOS GOARCH <<< "$target"
  output_name="${project_name}-${GOOS}-${GOARCH}"

  echo "Building for $GOOS/$GOARCH..."
  env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_dir/$output_name

  if [ $? -ne 0 ]; then
    echo "Failed to build for $GOOS/$GOARCH"
  else
    echo "Successfully built $output_name"
  fi
done

echo "Builds completed. Binaries are in the $output_dir directory."