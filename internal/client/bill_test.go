package client

import (
	"context"
	"mawinter-server/internal/model"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

var (
	testEndpoint      = "http://example.com:8080/bill/"
	testEndpointError = "http://error.example.com:8080/bill/"
)

func TestMain(m *testing.M) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	time.Local = jst

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", testEndpoint+"202201",
		func(req *http.Request) (*http.Response, error) {
			res, err := httpmock.NewJsonResponse(200,
				[]model.BillAPIResponse{
					{
						BillName: "gas",
						Price:    123,
					},
					{
						BillName: "elect",
						Price:    1234,
					},
					{
						BillName: "water",
						Price:    12345,
					},
				},
			)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return res, nil
		},
	)
	httpmock.RegisterResponder("GET", testEndpointError+"202201",
		func(req *http.Request) (*http.Response, error) {
			res := httpmock.NewStringResponse(500, "error")
			return res, nil
		},
	)

	m.Run()
}

func TestBillFetcher_FetchBills(t *testing.T) {
	type fields struct {
		BillEndpoint string
	}
	type args struct {
		ctx    context.Context
		yyyymm string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantRess []model.BillAPIResponse
		wantErr  bool
	}{
		{
			name: "ok",
			fields: fields{
				BillEndpoint: testEndpoint,
			},
			args: args{ctx: context.Background(), yyyymm: "202201"},
			wantRess: []model.BillAPIResponse{
				{
					BillName: "gas",
					Price:    123,
				},
				{
					BillName: "elect",
					Price:    1234,
				},
				{
					BillName: "water",
					Price:    12345,
				},
			},
			wantErr: false,
		},
		{
			name: "500 error",
			fields: fields{
				BillEndpoint: testEndpointError,
			},
			args:     args{ctx: context.Background(), yyyymm: "202201"},
			wantRess: []model.BillAPIResponse{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BillFetcher{
				BillEndpoint: tt.fields.BillEndpoint,
			}
			gotRess, err := b.FetchBills(tt.args.ctx, tt.args.yyyymm)
			if (err != nil) != tt.wantErr {
				t.Errorf("BillFetcher.FetchBills() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRess, tt.wantRess) {
				t.Errorf("BillFetcher.FetchBills() = %v, want %v", gotRess, tt.wantRess)
			}
		})
	}
}
