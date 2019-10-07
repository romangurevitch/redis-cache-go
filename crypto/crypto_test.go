package crypto

import (
	"reflect"
	"testing"
)

func TestHash(t *testing.T) {
	type args struct {
		password string
		salt     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "simple", args: args{password: "pass", salt: ""}, want: "a3de2301be7727891f635192670704d041319b589aac61452f9ec615d998f9b8"},
		{name: "test with salt", args: args{password: "pass", salt: "salt"}, want: "02fcc288e7bc681cb111817f981bc8ff7824fa38fc61c7817f2fba7f5b5b4b0d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(tt.args.password, tt.args.salt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashDifferentSalt(t *testing.T) {
	pass := "password"
	hash1 := Hash(pass, "salt1")
	hash2 := Hash(pass, "salt2")

	if hash1 == hash2 {
		t.Errorf("expecting different hash values")
	}
}
