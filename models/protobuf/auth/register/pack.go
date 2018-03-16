package register

import (
	"github.com/golang/protobuf/proto"
	pbregister "github.com/zk101/nixie/proto/auth/register"
)

// Pack returns a pointer to a Protobuf RegisterRequest message
func Pack(user, pass, name string) ([]byte, error) {
	register := pbregister.RegisterRequest{
		Username: user,
		Password: pass,
		Name:     name,
	}

	data, err := proto.Marshal(&register)
	if err != nil {
		return []byte(""), err
	}

	return data, nil
}

// EOF
