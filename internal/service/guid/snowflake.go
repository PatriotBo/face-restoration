package guid

import "github.com/bwmarrin/snowflake"

type Service interface {
	GenerateString() string
}

type serviceImpl struct {
	n *snowflake.Node
}

func MustNew(node int64) Service {
	n, err := snowflake.NewNode(node)
	if err != nil {
		panic(err)
	}

	return &serviceImpl{
		n: n,
	}
}

func (s *serviceImpl) GenerateString() string {
	return s.n.Generate().String()
}
