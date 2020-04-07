#!/usr/bin/env bash 

function vendor_knative_serving {
  local version=$1
  logger.info 'Vendoring knative-serving'
  if [ -d knative-serving ]; then
    logger.info "Removing extant knative-serving"
    rm -rf "knative-serving"
  fi
  git clone -b "release-${version}" --depth 1 https://github.com/openshift/knative-serving.git "knative-serving"
  rm -rf knative-serving/.git
  if [ ! -f knative-serving/go.mod ]; then
    logger.info "Priming knative-serving go.mod"
    echo "module knative.dev/serving" >> knative-serving/go.mod
  fi

}

