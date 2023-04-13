package dnsprovider

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	r53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

const (
	AWS_ACCESS_KEY = "AWS_ACCESS_KEY_ID"
	AWS_SECRET_KEY = "AWS_SECRET_ACCESS_KEY"
)

var (
	AWS_REGION = "us-east-1"
)

type Aws struct {
	Client *route53.Client
}

func NewAws(ctx context.Context, accessKeyId, secretAccessKey string) (Aws, error) {
	var cfg aws.Config
	var err error

	a := Aws{}

	// Configure AWS client
	cfg, err = config.LoadDefaultConfig(ctx,
		config.WithRegion(AWS_REGION),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     accessKeyId,
				SecretAccessKey: secretAccessKey,
				SessionToken:    "",
				Source:          "ddns-kubernetes-controller",
			},
		}),
	)
	if err != nil {
		return a, err
	}
	a.Client = route53.NewFromConfig(cfg)

	return a, nil
}

func (a Aws) UpdateRecord(ctx context.Context, record, zone, ip string, TTL int64) (err error) {
	var zoneId *string

	// Lookup Zone ID
	if !strings.HasSuffix(zone, ".") {
		zone = fmt.Sprintf("%s.", zone)
	}

	var zoneOutput *route53.ListHostedZonesOutput
	zoneOutput, err = a.Client.ListHostedZones(ctx, nil)
	if err != nil {
		return err
	}

	for _, hostedZone := range zoneOutput.HostedZones {
		if *hostedZone.Name == zone {
			zoneId = hostedZone.Id
		}
	}

	// Upsert Dynamic DNS record
	fqdn := fmt.Sprintf("%s.%s", record, zone)
	comment := fmt.Sprintf("Updating DDNS Record %q for zone %q", record, zone)
	params := route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &r53types.ChangeBatch{ // Required
			Changes: []r53types.Change{
				{
					Action: r53types.ChangeActionUpsert,
					ResourceRecordSet: &r53types.ResourceRecordSet{
						Name: &fqdn,
						Type: "A",
						TTL:  &TTL,
						ResourceRecords: []r53types.ResourceRecord{
							{
								Value: &ip,
							},
						},
					},
				},
			},
			Comment: &comment,
		},
		HostedZoneId: zoneId,
	}

	_, err = a.Client.ChangeResourceRecordSets(ctx, &params)
	if err != nil {
		return err
	}

	return nil
}
