package domain

import (
	"testing"
)

func TestTransStatus_String(t *testing.T) {
	tests := []struct {
		status   TransStatus
		expected string
	}{
		{StatusPending, "pending"},
		{StatusApproved, "approved"},
		{StatusRejected, "rejected"},
		{StatusFraud, "fraud"},
		{TransStatus(999), "pending"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
