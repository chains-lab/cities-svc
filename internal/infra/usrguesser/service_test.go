package usrguesser

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

const testURL = "http://localhost:8002/profiles-svc/v1/profiles"

var userIDs = []uuid.UUID{
	uuid.MustParse("3673b5e6-3b74-4c94-a735-8185971529fc"),
	uuid.MustParse("f9b60155-9ed3-4489-86a3-e93545a7f624"),
	uuid.MustParse("d392708d-a6d6-4014-b989-065a5e910fc5"),
}

func Test(t *testing.T) {
	svc := NewService(testURL, nil)

	userIDs := append(userIDs, uuid.New()) // add a non-existing user ID for testing

	res, err := svc.Guess(context.Background(), userIDs...)
	if err != nil {
		t.Fatalf("guess failed: %v", err)
	}
	t.Logf("got %d profiles", len(res))

	for id, profile := range res {
		t.Logf("id: %s, username: %s", id, profile.Username)
	}
}
