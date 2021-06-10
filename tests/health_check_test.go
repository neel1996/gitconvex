package tests

import (
	"github.com/neel1996/gitconvex/api"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthCheckApi(t *testing.T) {
	tests := []struct {
		name string
		want *model.HealthCheckParams
	}{
		{name: "HealthCheck OS test case", want: &model.HealthCheckParams{
			Os: "linux||darwin||windows",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := api.HealthCheckApi()
			assert.Contains(t, tt.want.Os, got.Os)
			assert.NotEmpty(t, got.Gitconvex)
		})
	}
}
