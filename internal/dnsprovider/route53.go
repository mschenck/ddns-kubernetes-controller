package dnsprovider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	r53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

var (
	AWS_REGION = "us-east-1"
)

type Aws struct{}

func (a *Aws) UpdateRecord(record, zone, ip string, TTL int64) (err error) {
	var cfg aws.Config
	var zoneId *string
	ctx := context.TODO()

	// Configure AWS client
	cfg, err = config.LoadDefaultConfig(ctx,
		config.WithRegion(AWS_REGION),
	)
	if err != nil {
		return err
	}
	r53Client := route53.NewFromConfig(cfg)

	// Lookup Zone ID
	var zoneOutput *route53.ListHostedZonesOutput
	zoneOutput, err = r53Client.ListHostedZones(ctx, nil)
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

	_, err = r53Client.ChangeResourceRecordSets(ctx, &params)
	if err != nil {
		return err
	}

	return nil
}
