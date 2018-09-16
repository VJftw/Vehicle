package vehicle

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"golang.org/x/crypto/ssh"
)

type AWS struct {
	config    *AWSConfig
	ec2Client *ec2.EC2

	uuid     string
	instance *ec2.Instance

	privateKey      []byte
	sshClientConfig *ssh.ClientConfig
	instanceID      *string
}

func NewAWS(config *AWSConfig, uuid string) *AWS {
	a := &AWS{config: config, uuid: uuid}
	sess := session.Must(session.NewSession(&aws.Config{
		MaxRetries: aws.Int(3),
	}))
	a.ec2Client = ec2.New(sess)
	return a
}

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
		fmt.Printf("deleting key pair: %s\n", a.uuid)
		_, err := a.ec2Client.DeleteKeyPair(&ec2.DeleteKeyPairInput{
			KeyName: aws.String(a.uuid),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AWS) createPrivateKey() error {
	fmt.Printf("creating key pair: %s\n", a.uuid)
	keyPairOut, err := a.ec2Client.CreateKeyPair(&ec2.CreateKeyPairInput{
		KeyName: aws.String(a.uuid),
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
		User: a.config.SSH.User,
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
		ImageId:      aws.String(a.config.AMI.ID),
		InstanceType: aws.String(a.config.Type),
		KeyName:      aws.String(a.uuid),
		SecurityGroupIds: []*string{
			aws.String(a.config.SecurityGroups[0].IDs[0]),
		},
		SubnetId: aws.String(a.config.Subnet.ID),
		TagSpecifications: []*ec2.TagSpecification{
			&ec2.TagSpecification{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					&ec2.Tag{
						Key:   aws.String("Name"),
						Value: aws.String(a.uuid),
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
		return err
	}

	a.instanceID = reservations.Instances[0].InstanceId

	return nil
}

func (a *AWS) waitForInstance() error {
	fmt.Printf("waiting for instance to be OK: %s\n", aws.StringValue(a.instanceID))
	err := a.ec2Client.WaitUntilInstanceStatusOk(&ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{a.instanceID},
	})
	if err != nil {
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

func (a *AWS) GetSSHInfo() (string, uint16, time.Duration, *ssh.ClientConfig) {
	return aws.StringValue(a.instance.PublicIpAddress), a.config.SSH.Port, a.config.SSH.Timeout, a.sshClientConfig
}

func (a *AWS) ResolveFuncs() []func() (error, []string) {
	return []func() (error, []string){}
}

func (a *AWS) StartFuncs() []func() error {
	return []func() error{
		a.createPrivateKey,
		a.generateSSHConfig,
		a.launchInstance,
		a.waitForInstance,
	}
}
