#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo "Usage: $0 [URL] [TIMEOUT_SECONDS]"
  echo "Example: $0 http://localhost:8080/health 5"
}

check_health() {
  local target_url="$1"
  local max_time="$2"
  local http_code

  # Ambil HTTP status code; jika curl gagal (misal connection refused), set ke 000
  http_code="$(curl -s -o /dev/null -w "%{http_code}" --max-time "$max_time" "$target_url")" || http_code="000"

  if [[ "$http_code" -ge 200 && "$http_code" -lt 400 ]]; then
    echo "[OK] ERP sehat: $target_url (HTTP $http_code)"
    return 0
  else
    echo "[FAIL] ERP tidak merespons normal: $target_url (HTTP $http_code)" >&2
    return 1
  fi
}

main() {
  if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
    usage
    exit 0
  fi

  local url="${1:-http://localhost:8080/health}"
  local timeout="${2:-5}"

  check_health "$url" "$timeout"
}

main "$@"
