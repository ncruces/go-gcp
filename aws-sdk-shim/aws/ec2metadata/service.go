package ec2metadata

type EC2Metadata struct{}

func New(p *struct{}) *EC2Metadata {
	return nil
}

func (c *EC2Metadata) Available() bool {
	return false
}

func (c *EC2Metadata) GetInstanceIdentityDocument() (d EC2InstanceIdentityDocument, err error) {
	return
}

type EC2InstanceIdentityDocument struct {
	Region     string
	InstanceID string
	AccountID  string
}
