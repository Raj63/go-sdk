package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func Test_GetHeader(t *testing.T) {
	type args struct {
		ctx    context.Context
		header string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "when header is not found",
			args: args{
				ctx:    metadata.NewIncomingContext(context.Background(), metadata.MD{}),
				header: "abc",
			},
			wantErr: true,
		},
		{
			name: "when header is found",
			args: args{
				ctx:    metadata.NewIncomingContext(context.Background(), metadata.MD{"abc": []string{"123"}}),
				header: "abc",
			},
			want: []string{"123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHeader(tt.args.ctx, tt.args.header)

			if tt.wantErr {
				assert.NotNilf(t, err, "GetHeader() error = %v, wantErr %v", err, tt.wantErr)

			} else {
				assert.Nil(t, err)
				assert.Equalf(t, tt.want, got, "GetHeader() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetHeaders(t *testing.T) {
	type args struct {
		ctx     context.Context
		headers []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "when one of the headers is not found",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
					"abc": []string{"123"},
				}),
				headers: []string{"abc", "defg"},
			},
			wantErr: true,
		},
		{
			name: "when two or more of the headers are not found",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
					"abc": []string{"123"},
				}),
				headers: []string{"abc", "defg", "hijk"},
			},
			wantErr: true,
		},
		{
			name: "when all headers are found",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), metadata.MD{
					"abc":  []string{"123"},
					"defg": []string{"456"},
					"hijk": []string{"6789"},
				}),
				headers: []string{"abc", "defg", "hijk"},
			},
			want: map[string][]string{
				"abc":  {"123"},
				"defg": {"456"},
				"hijk": {"6789"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHeaders(tt.args.ctx, tt.args.headers...)

			if tt.wantErr {
				assert.NotNilf(t, err, "GetHeaders() error = %v, wantErr %v", err, tt.wantErr)

			} else {
				assert.Nil(t, err)
				assert.Equalf(t, tt.want, got, "GetHeaders() got = %v, want %v", got, tt.want)
			}
		})
	}
}
