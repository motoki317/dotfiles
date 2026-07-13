#!/usr/bin/env bash
# Print the fetchurl SRI hashes for a prebuilt-binary package in packages.nix
# (ax, hunk). `nix store prefetch-file` hashes the downloaded file flat, exactly
# as fetchurl pins it, so its output is drop-in.
#
# To bump: edit `version` for the tool in packages.nix, run this, paste both hashes.
#   ./fetch-prebuilt-hashes.sh hunk 0.18.0
set -euo pipefail

tool=${1:-}
version=${2:-}
[[ -n $tool && -n $version ]] || { echo "usage: $0 <ax|hunk> <version>" >&2; exit 1; }

# %s is the release asset's platform suffix, filled in per system below.
case $tool in
  ax)   template="https://github.com/yusukebe/ax/releases/download/v$version/ax-%s" ;;
  hunk) template="https://github.com/modem-dev/hunk/releases/download/v$version/hunkdiff-%s.tar.gz" ;;
  *)    echo "unknown tool: $tool (want ax or hunk)" >&2; exit 1 ;;
esac

echo "$tool $version"
for pair in "x86_64-linux linux-x64" "aarch64-darwin darwin-arm64"; do
  read -r system suffix <<<"$pair"
  hash=$(nix store prefetch-file --json "${template//%s/$suffix}" | jq -r .hash)
  printf '  %-15s %s\n' "$system" "$hash"
done
