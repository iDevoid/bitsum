package routing

import (
	"reflect"
	"testing"

	"github.com/iDevoid/bitsum/internal/handler/rest"
	coins_mock "github.com/iDevoid/bitsum/mocks/coins"
	"github.com/iDevoid/bitsum/platform/routers"
)

func TestCoinsRouting(t *testing.T) {
	mocked := coins_mock.MockCoinsHandler{}
	type args struct {
		handler rest.CoinsHandler
	}
	tests := []struct {
		name string
		args args
		want []routers.Router
	}{
		{
			name: "success",
			args: args{
				handler: &mocked,
			},
			want: []routers.Router{
				{
					Method:  "GET",
					Path:    "/wallet",
					Handler: mocked.Wallet,
				},
				{
					Method:  "POST",
					Path:    "/pay",
					Handler: mocked.Payment,
				},
				{
					Method:  "POST",
					Path:    "/receive",
					Handler: mocked.Receive,
				},
				{
					Method:  "POST",
					Path:    "/history",
					Handler: mocked.History,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CoinsRouting(tt.args.handler)
			for i, val := range got {
				if !reflect.DeepEqual(val.Method, tt.want[i].Method) {
					t.Errorf("CoinsRouting() method = %v, want %v", val.Method, tt.want[i].Method)
				}
				if !reflect.DeepEqual(val.Path, tt.want[i].Path) {
					t.Errorf("CoinsRouting() path = %v, want %v", val.Path, tt.want[i].Path)
				}
				if val.Handler == nil && tt.want[i].Handler != nil {
					t.Error("got nih handler")
				}
			}
		})
	}
}
