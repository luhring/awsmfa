package authenticator

import (
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/luhring/awsmfa/file_coordinator"
	"testing"
)

func TestNew(t *testing.T) {
	stsClient := &sts.STS{}
	fileCoordinator := &file_coordinator.Coordinator{}

	auth, err := New(stsClient, fileCoordinator)

	if err != nil {
		t.Error("New should never return a non-nil error")
	}

	if auth == nil {
		t.Error("New authenticator object should not be nil")
	}

	if auth.stsClient != stsClient {
		t.Error("new authenticator object had incorrect reference for stsClient")
	}

	if auth.fileCoordinator != fileCoordinator {
		t.Error("new authenticator object had incorrect reference for fileCoordinator")
	}
}
