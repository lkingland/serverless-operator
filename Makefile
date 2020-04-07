# Useful for local development
dev:
	./hack/dev.sh

# General purpose targets
images:
	./hack/images.sh $(DOCKER_REPO_OVERRIDE)

install:
	./hack/install.sh

teardown:
	./hack/teardown.sh

# vendors dependencies which are not able to be automatically
# pulled in via dep (for instance: Knative Serving for running
# E2E and upgrade tests, which has no compile-time references).
vendor:
	./hack/vendor.sh

# Test targets for CI operator
test-unit:
	go test ./knative-operator/...
	go test ./serving/ingress/...

test-e2e:
	./test/e2e-tests.sh

test-upgrade:
	./test/upgrade-tests.sh

# Generates a ci-operator configuration for a specific branch.
generate-ci-config:
	./openshift/ci-operator/generate-ci-config.sh $(BRANCH) > ci-operator-config.yaml
