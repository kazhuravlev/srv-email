package smtp

import (
	"bytes"
	"net/mail"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_encodeAddrs(t *testing.T) {
	type args struct {
		addrs []mail.Address
		bufTo *bytes.Buffer
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "check nil addrs",
			args: args{
				addrs: nil,
				bufTo: new(bytes.Buffer),
			},
			want: "",
		},
		{
			name: "check one addr",
			args: args{
				addrs: []mail.Address{
					{
						Name:    "Kirill Zhuravlev",
						Address: "kazhuravlev@fastmail.com",
					},
				},
				bufTo: new(bytes.Buffer),
			},
			want: `"Kirill Zhuravlev" <kazhuravlev@fastmail.com>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encodeAddrs(tt.args.addrs, tt.args.bufTo)
			assert.Equal(t, tt.want, tt.args.bufTo.String())
		})
	}
}
