#!/usr/bin/env bash
# Print the fetchurl SRI hashes for a prebuilt-binary package in packages.nix
# (ax, hunk, agent-browser). `nix store prefetch-file` hashes the downloaded file
# flat, exactly as fetchurl pins it, so its output is drop-in.
#
# To bump: edit `version` for the tool in packages.nix, run this, paste both hashes.
#   ./fetch-prebuilt-hashes.sh hunk 0.18.0
set -euo pipefail

tool=${1:-}
version=${2:-}
[[ -n $tool && -n $version ]] || { echo "usage: $0 <ax|hunk|agent-browser|agtlog> <version>" >&2; exit 1; }

# %s is the release asset's platform suffix; each tool also names its own system->suffix
# pairs, since the suffix spelling is not shared (e.g. agtlog uses linux-amd64, not linux-x64).
case $tool in
  ax)            template="https://github.com/yusukebe/ax/releases/download/v$version/ax-%s";                          suffixes=("x86_64-linux linux-x64" "aarch64-darwin darwin-arm64") ;;
  hunk)          template="https://github.com/modem-dev/hunk/releases/download/v$version/hunkdiff-%s.tar.gz";          suffixes=("x86_64-linux linux-x64" "aarch64-darwin darwin-arm64") ;;
  agent-browser) template="https://github.com/vercel-labs/agent-browser/releases/download/v$version/agent-browser-%s"; suffixes=("x86_64-linux linux-x64" "aarch64-darwin darwin-arm64") ;;
  agtlog)        template="https://github.com/motoki317/agtlog/releases/download/v$version/agtlog-$version-%s.tar.gz";  suffixes=("x86_64-linux linux-amd64" "aarch64-darwin darwin-arm64") ;;
  *)             echo "unknown tool: $tool (want ax, hunk, agent-browser, or agtlog)" >&2; exit 1 ;;
esac

echo "$tool $version"
for pair in "${suffixes[@]}"; do
  read -r system suffix <<<"$pair"
  hash=$(nix store prefetch-file --json "${template//%s/$suffix}" | jq -r .hash)
  printf '  %-15s %s\n' "$system" "$hash"
done
