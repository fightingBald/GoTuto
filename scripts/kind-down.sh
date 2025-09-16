#!/usr/bin/env bash
set -euo pipefail

NAME=${1:-kind}
if ! command -v kind >/dev/null 2>&1; then
  echo "kind is required (https://kind.sigs.k8s.io)." >&2
  exit 1
fi

kind delete cluster --name "$NAME"
