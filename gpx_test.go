package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_makeGpx(t *testing.T) {
	type args struct {
		status *ChatStatus
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
		wantErr    bool
	}{
		{
			name: "empty",
			args: args{
				status: &ChatStatus{
					chatId: 123456,
					positions: []*Position{
						{
							when:      time.Now(),
							longitude: 42.344,
							latitude:  1.00,
						},
					},
				},
			},
			wantResult: nil,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := makeGpx(tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeGpx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("makeGpx() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
