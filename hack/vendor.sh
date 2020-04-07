#!/usr/bin/env bash

# Copies into this repository (vendors) packages which are required to exist 
# locally, but have no compile-time references and are thus out of scope of
# the dependency manager.
#
# shellcheck disable=SC1091,SC1090
source "$(dirname "${BASH_SOURCE[0]}")/lib/__sources__.bash"

set -Eeuo pipefail

exitcode=0

(( !exitcode )) && vendor_knative_serving "$KNATIVE_VERSION" || exitcode=2

exit $exitcode
