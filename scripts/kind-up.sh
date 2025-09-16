#!/usr/bin/env bash
set -euo pipefail

CLUSTER_FILE="kind/kind-cluster.yaml"

if ! command -v kind >/dev/null 2>&1; then
  echo "kind is required (https://kind.sigs.k8s.io)." >&2
  exit 1
fi

if [ -f "$CLUSTER_FILE" ]; then
  kind create cluster --config "$CLUSTER_FILE"
else
  echo "No kind cluster config found, creating default 'kind' cluster" >&2
  kind create cluster
fi

echo "Cluster is up. You can now run: tilt up (optional)"

