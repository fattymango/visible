package visible

import (
	"reflect"
	"testing"

	"github.com/fatih/structs"
)

type Data struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	AliceData       string `json:"alice_data" visible:"alice"`
	BobData         string `json:"bob_data" visible:"bob"`
	AliceAndBobData string `json:"alice_and_bob_data" visible:"alice,bob"`
	PrivateData     string `json:"-"`
	PrivateData2    string `json:"-" visible:"bob"`
}

var data = Data{
	ID:              1,
	Name:            "name",
	AliceData:       "alice_data",
	BobData:         "bob_data",
	AliceAndBobData: "alice_and_bob_data",
	PrivateData:     "private_data",
	PrivateData2:    "private_data2",
}

func Test_isJsonIgnored(t *testing.T) {
	type args struct {
		field *structs.Field
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no visible tag, not ignored",
			args: args{
				field: structs.Fields(data)[0],
			},
			want: false,
		},
		{
			name: "visible tag, not ignored",
			args: args{
				field: structs.Fields(data)[2],
			},
			want: false,
		},

		{
			name: "no visible tag, ignored",
			args: args{
				field: structs.Fields(data)[5],
			},
			want: true,
		},
		{
			name: "visible tag, ignored",
			args: args{
				field: structs.Fields(data)[6],
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJsonIgnored(tt.args.field); got != tt.want {
				t.Errorf("isJsonIgnored() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractVisibleTag(t *testing.T) {
	type args struct {
		field *structs.Field
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no visible tag",
			args: args{
				field: structs.Fields(data)[0],
			},
			want: "",
		},
		{
			name: "visible tag",
			args: args{
				field: structs.Fields(data)[2],
			},
			want: "alice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractVisibleTag(tt.args.field); got != tt.want {
				t.Errorf("extractVisibleTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isFieldVisible(t *testing.T) {
	type args struct {
		tag     string
		visible string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{

		{
			name: "visible tag, visible",
			args: args{
				tag:     "alice",
				visible: "alice",
			},
			want: true,
		},
		{
			name: "visible tag, not visible",
			args: args{
				tag:     "alice",
				visible: "bob",
			},
			want: false,
		},
		{
			name: "multiple visible tag, visible",
			args: args{
				tag:     "alice,bob",
				visible: "alice",
			},
			want: true,
		},
		{
			name: "multiple visible tag, not visible",
			args: args{
				tag:     "alice,bob",
				visible: "charlie",
			},
			want: false,
		},
		{
			name: "no visible tag",
			args: args{
				tag:     "",
				visible: "alice",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isFieldVisible(tt.args.tag, tt.args.visible); got != tt.want {
				t.Errorf("isFieldVisible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanStruct(t *testing.T) {
	type args struct {
		data    interface{}
		visible string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "no data",
			args: args{
				data:    nil,
				visible: "alice",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not a struct",
			args: args{
				data:    1,
				visible: "alice",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no visible field",
			args: args{
				data:    struct{ ID int }{ID: 1},
				visible: "alice",
			},
			want:    map[string]interface{}{"ID": 1},
			wantErr: false,
		},
		{
			name: "visible field",
			args: args{
				data:    data,
				visible: "alice",
			},
			want:    map[string]interface{}{"ID": 1, "Name": "name", "AliceData": "alice_data", "AliceAndBobData": "alice_and_bob_data"},
			wantErr: false,
		},
		{
			name: "visible field",
			args: args{
				data:    data,
				visible: "bob",
			},
			want:    map[string]interface{}{"ID": 1, "Name": "name", "BobData": "bob_data", "AliceAndBobData": "alice_and_bob_data"},
			wantErr: false,
		},
		{
			name: "no visible field",
			args: args{
				data:    data,
				visible: "charlie",
			},
			want:    map[string]interface{}{"ID": 1, "Name": "name"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CleanStruct(tt.args.data, tt.args.visible)
			if (err != nil) != tt.wantErr {
				t.Errorf("CleanStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CleanStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}
