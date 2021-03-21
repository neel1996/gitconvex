package tests

import (
	"github.com/neel1996/gitconvex-server/api"
	"github.com/neel1996/gitconvex-server/graph/model"
	"strings"
	"testing"
)

func TestHealthCheckApi(t *testing.T) {
	tests := []struct {
		name string
		want *model.HealthCheckParams
	}{
		{name: "HealthCheck OS test case", want: &model.HealthCheckParams{
			Os: "linux||darwin||win",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := api.HealthCheckApi(); strings.Contains(got.Os, tt.want.Os) {
				t.Errorf("HealthCheckApi() = %v, want %v", got, tt.want)
			}

			if got := api.HealthCheckApi(); got.Gitconvex == "" {
				t.Errorf("HealthCheckApi() = git version string empty")
			}
		})
	}
}
