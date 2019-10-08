package e2e

import (
	"strings"
	"testing"

	"github.com/openshift-knative/serverless-operator/test"
)

const (
	knativeServing    = "knative-serving"
	testNamespace     = "serverless-tests"
	image             = "gcr.io/knative-samples/helloworld-go"
	helloworldService = "helloworld-go"
)

func TestKnativeServing(t *testing.T) {
	contexts := test.Setup(t)
	adminCtx := contexts[0]

	defer test.CleanupAll(contexts)
	test.CleanupOnInterrupt(t, func() { test.CleanupAll(contexts) })

	t.Run("create subscription and wait for CSV to succeed", func(t *testing.T) {
		_, err := test.WithOperatorReady(adminCtx, "serverless-operator-subscription")
		if err != nil {
			t.Fatal("Failed", err)
		}
	})

	t.Run("deploy knativeserving cr and wait for it to be ready", func(t *testing.T) {
		_, err := test.WithKnativeServingReady(adminCtx, knativeServing, knativeServing)
		if err != nil {
			t.Fatal("Failed to deploy KnativeServing", err)
		}
	})

	t.Run("deploy knative service using kubeadmin", func(t *testing.T) {
		_, err := test.WithServiceReady(adminCtx, helloworldService, testNamespace, image)
		if err != nil {
			t.Fatal("Knative Service not ready", err)
		}
	})

	t.Run("user permissions", func(t *testing.T) {
		testUserPermissions(t, contexts)
	})

	t.Run("undeploy serverless operator and check dependent operators removed", func(t *testing.T) {
		adminCtx.Cleanup()
		err := test.WaitForOperatorDepsDeleted(adminCtx)
		if err != nil {
			t.Fatalf("Operators still running: %v", err)
		}
	})
}

func testUserPermissions(t *testing.T, contexts []*test.Context) {
	editCtx := contexts[1]
	viewCtx := contexts[2]
	tests := []struct {
		name        string
		userRole    string
		userContext *test.Context
		operation   func(context *test.Context) error
		wantErrStr  string
	}{{
		name:     "user with view role can get",
		userRole: "view",
		operation: func(c *test.Context) error {
			_, err := test.GetService(c, helloworldService, testNamespace)
			return err
		},
		userContext: viewCtx,
	}, {
		name:     "user with view role can list",
		userRole: "view",
		operation: func(c *test.Context) error {
			_, err := test.ListServices(c, testNamespace)
			return err
		},
		userContext: viewCtx,
	}, {
		name:     "user with view role cannot create",
		userRole: "view",
		operation: func(c *test.Context) error {
			_, err := test.CreateService(c, "userview-service", testNamespace, image)
			return err
		},
		userContext: viewCtx,
		wantErrStr:  "is forbidden",
	}, {
		name:     "user with view role cannot delete",
		userRole: "view",
		operation: func(c *test.Context) error {
			return test.DeleteService(c, helloworldService, testNamespace)
		},
		userContext: viewCtx,
		wantErrStr:  "is forbidden",
	}, {
		name:     "user with edit role can get",
		userRole: "edit",
		operation: func(c *test.Context) error {
			_, err := test.GetService(c, helloworldService, testNamespace)
			return err
		},
		userContext: editCtx,
	}, {
		name:     "user with edit role can list",
		userRole: "edit",
		operation: func(c *test.Context) error {
			_, err := test.ListServices(c, testNamespace)
			return err
		},
		userContext: editCtx,
	}, {
		name:     "user with edit role can create",
		userRole: "edit",
		operation: func(c *test.Context) error {
			_, err := test.CreateService(c, "useredit-service", testNamespace, image)
			return err
		},
		userContext: editCtx,
	}, {
		name:     "user with edit role can delete",
		userRole: "edit",
		operation: func(c *test.Context) error {
			return test.DeleteService(c, "useredit-service", testNamespace)
		},
		userContext: editCtx,
	},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.operation(test.userContext)
			if (err != nil) != (test.wantErrStr != "") {
				t.Errorf("User with role %s has unexpected behavior on knative services. Error thrown: %v, error expected: %t", test.userRole, err, (test.wantErrStr != ""))
			}
			if err != nil && !strings.Contains(err.Error(), test.wantErrStr) {
				t.Errorf("Unexpected error for user with role %s: %v", test.userRole, err)
			}
		})
	}
}
