package ioc

import "github.com/bwmarrin/snowflake"

func InitNode() *snowflake.Node {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return node
}
