package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/service/ecs"
	"github.com/gocircuit/runtime/prov"
)

// New returns a new worker provisioner based on Amazon spot instances.
func New() prov.Provisioner {
	x
}

type provisioner struct {
	ecs *ecs.ECS
}

func init() {
	aws.DefaultConfig.Region = "us-east-1a"
}
