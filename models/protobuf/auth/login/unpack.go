package login

import (
	"github.com/golang/protobuf/proto"
	pblogin "github.com/zk101/nixie/proto/auth/login"
)

// UnpackReply returns a Login Request pointer
func UnpackReply(data *[]byte) (*pblogin.LoginReply, error) {
	replyData := pblogin.LoginReply{}

	if err := proto.Unmarshal(*data, &replyData); err != nil {
		return nil, err
	}

	return &replyData, nil
}

// UnpackRequest returns a Login Request pointer
func UnpackRequest(data *[]byte) (*pblogin.LoginRequest, error) {
	requestData := pblogin.LoginRequest{}

	if err := proto.Unmarshal(*data, &requestData); err != nil {
		return nil, err
	}

	return &requestData, nil
}

// EOF
