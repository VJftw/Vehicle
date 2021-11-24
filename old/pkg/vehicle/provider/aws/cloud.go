package aws

import (
	"fmt"
	"time"

	"github.com/VJftw/vehicle/pkg/vehicle/provider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"golang.org/x/crypto/ssh"
)

// AWS represents the AWS Cloud vehicle instance provider
type AWS struct {
	*provider.Base
	*Config

	ec2Client *ec2.EC2

	instance *ec2.Instance

	privateKey      []byte
	sshClientConfig *ssh.ClientConfig
	instanceID      *string
}

// New returns a new AWS Cloud vehicle instance
func New(uuid string) *AWS {
	a := &AWS{
		Base:   &provider.Base{Provider: "aws", UUID: uuid},
		Config: NewConfig(),
	}
	sess := session.Must(session.NewSession(&aws.Config{
		MaxRetries: aws.Int(3),
	}))
	a.ec2Client = ec2.New(sess)
	return a
}

// GetSSHInfo returns the ssh configuration for a provisioned AWS instance
func (a *AWS) GetSSHInfo() (string, uint16, time.Duration, *ssh.ClientConfig) {
	return aws.StringValue(a.instance.PublicIpAddress), a.SSH.Port, a.SSH.Timeout, a.sshClientConfig
}

// ResolveFuncs returns the functions involved in resolving dynamic resource IDs
func (a *AWS) ResolveFuncs() []func() (error, []string) {
	return []func() (error, []string){}
}

// StartFuncs returns the functions involved in starting an AWS instance
func (a *AWS) StartFuncs() []func() error {
	return []func() error{
		a.createPrivateKey,
		a.generateSSHConfig,
		a.launchInstance,
		a.waitForInstance,
	}
}

// Stop stops and all AWS resources created by vehicle
func (a *AWS) Stop() error {
	if a.instanceID != nil {
		fmt.Printf("terminating instance: %s\n", aws.StringValue(a.instanceID))
		_, err := a.ec2Client.TerminateInstances(&ec2.TerminateInstancesInput{
			InstanceIds: []*string{a.instanceID},
		})
		if err != nil {
			return err
		}
		err = a.ec2Client.WaitUntilInstanceTerminated(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{a.instanceID},
		})
		if err != nil {
			return err
		}
	}
	if len(a.privateKey) > 0 {
		fmt.Printf("deleting key pair: %s\n", a.UUID)
		_, err := a.ec2Client.DeleteKeyPair(&ec2.DeleteKeyPairInput{
			KeyName: aws.String(a.UUID),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AWS) createPrivateKey() error {
	fmt.Printf("creating key pair: %s\n", a.UUID)
	keyPairOut, err := a.ec2Client.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(a.UUID),
	})
	if err != nil {
		return err
	}

	a.privateKey = []byte(aws.StringValue(keyPairOut.KeyMaterial))
	return nil
}

func (a *AWS) generateSSHConfig() error {
	fmt.Println("generating ssh config")
	signer, err := ssh.ParsePrivateKey(a.privateKey)
	if err != nil {
		return err
	}

	a.sshClientConfig = &ssh.ClientConfig{
		User: a.SSH.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// HostKeyCallback: ssh.FixedHostKey(keyPairOut.)
	}
	return nil
}

func (a *AWS) launchInstance() error {
	fmt.Println("launching instance")
	reservations, err := a.ec2Client.RunInstances(&ec2.RunInstancesInput{
		// DryRun:       aws.Bool(true),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		ImageId:      aws.String(a.AMI.ID),
		InstanceType: aws.String(a.Type),
		KeyName:      aws.String(a.UUID),
		SecurityGroupIds: []*string{
			aws.String(a.SecurityGroups[0].IDs[0]),
		},
		SubnetId: aws.String(a.Subnet.ID),
		TagSpecifications: []*ec2.TagSpecification{
			&ec2.TagSpecification{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					&ec2.Tag{
						Key:   aws.String("Name"),
						Value: aws.String(a.UUID),
					},
				},
			},
		},
		// IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
		// 	Name: "",
		// 	Arn: ""
		// }
	})

	if err != nil {
		// log.Fatal(err)
		return err
	}

	// log.Printf("%+v", reservations)
	a.instanceID = reservations.Instances[0].InstanceId

	return nil
}

func (a *AWS) waitForInstance() error {
	fmt.Printf("waiting for instance to be OK: %s\n", aws.StringValue(a.instanceID))
	if err := a.ec2Client.WaitUntilInstanceStatusOk(&ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{a.instanceID},
	}); err != nil {
		return err
	}
	describeOut, err := a.ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{a.instanceID},
	})
	if err != nil {
		return err
	}
	a.instance = describeOut.Reservations[0].Instances[0]

	return nil
}
