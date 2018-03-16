package register

import (
	"github.com/golang/protobuf/proto"
	pbregister "github.com/zk101/nixie/proto/auth/register"
)

// Unpack returns a RegisterRequest pointer
func Unpack(data *[]byte) (*pbregister.RegisterRequest, error) {
	requestData := pbregister.RegisterRequest{}

	if err := proto.Unmarshal(*data, &requestData); err != nil {
		return nil, err
	}

	return &requestData, nil
}

// EOF
