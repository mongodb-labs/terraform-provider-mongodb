#!/usr/bin/env bash
set -ueo pipefail

PKG_NAME=mongodb
platforms=("linux/amd64" "windows/amd64" "darwin/amd64")

# Compute a version string based on GIT data
function git_version_string() {
  gitsha="$(git log -n1 --pretty='%h')"
  tag=$(git describe --exact-match --tags "${gitsha}" 2>/dev/null || echo "")
  if [ -n "$tag" ]; then
    # The current commit is tagged
    version="${tag}"
  else
    # Otherwise use the short git sha
    version="${gitsha}"
  fi

  if ! git diff-index --quiet HEAD --; then
    # If we have changes in the working directory, augment the version string
    version="${version}-modified"
  fi

  echo "$version"
}

# Output directory for builds
version="$(git_version_string)"
output="out/${version}"
mkdir -p "$output"

if [ -d "${output}" ]; then
  # Clear out any previously existing builds for the current version
  rm -rf "${output:?}"/*
fi

# Augment output with package name and version
output_package="${output}/terraform-provider-${PKG_NAME}_v${version}"

# Based on https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04
for platform in "${platforms[@]}"; do
  #shellcheck disable=SC2206
  platform_split=(${platform//\// })
  GOOS="${platform_split[0]}"
  GOARCH="${platform_split[1]}"

  output_name="${output_package}"
  if [ "${GOOS}" = "windows" ]; then
    output_name+=".exe"
  fi

  echo "Building and compressing ${output_name} ..."
  env GOOS="${GOOS}" GOARCH="${GOARCH}" CGO_ENABLED=0 go build -a -o "${output_name}"
  zip -m "${output_name}_${GOOS}_${GOARCH}" "${output_name}"
  echo
done

echo "Generating SHA256SUMS..."
sha256sum -b "${output}/"* >"${output_name}_SHA256SUMS"

echo "Done."
