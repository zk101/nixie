package login

import (
	"github.com/golang/protobuf/proto"
	pblogin "github.com/zk101/nixie/proto/auth/login"
)

// PackReply returns a pointer to a Protobuf Loginreply message ready to ship
func PackReply(key, sign, cipher string) ([]byte, error) {
	login := pblogin.LoginReply{
		Error:  pblogin.LoginReply_OKAY,
		Key:    key,
		Sign:   sign,
		Cipher: cipher,
	}

	data, err := proto.Marshal(&login)
	if err != nil {
		return []byte(""), err
	}

	return data, nil
}

// PackRequest returns a marshalled LoginRequest message ready to ship
func PackRequest(user, pass string) ([]byte, error) {
	login := pblogin.LoginRequest{
		Username: user,
		Password: pass,
	}

	data, err := proto.Marshal(&login)
	if err != nil {
		return []byte(""), err
	}

	return data, nil
}

// EOF
