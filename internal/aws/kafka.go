package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
)

type kafkaClient interface {
	ListClusters(ctx context.Context, params *kafka.ListClustersInput, optFns ...func(*kafka.Options)) (*kafka.ListClustersOutput, error)
	ListNodes(ctx context.Context, params *kafka.ListNodesInput, optFns ...func(*kafka.Options)) (*kafka.ListNodesOutput, error)
}

type ClusterDetails struct {
	ClusterName string
	ClusterArn  string
	Brokers     kafka.GetBootstrapBrokersOutput
}

func GetClusterConfig(svc *kafka.Client) (ClusterDetails, error) {
	clusters, err := getClusters(svc)
	if err != nil {
		return ClusterDetails{}, err
	}
	input := &kafka.GetBootstrapBrokersInput{ClusterArn: clusters.ClusterInfoList[0].ClusterArn}
	brokers, err := svc.GetBootstrapBrokers(context.TODO(), input)
	if err != nil {
		return ClusterDetails{}, err
	}
	cluster := ClusterDetails{
		ClusterName: *clusters.ClusterInfoList[0].ClusterName,
		ClusterArn:  *clusters.ClusterInfoList[0].ClusterArn,
		Brokers:     *brokers,
	}
	return cluster, nil
}

func getClusters(svc kafkaClient) (*kafka.ListClustersOutput, error) {
	input := &kafka.ListClustersInput{}
	output := &kafka.ListClustersOutput{}

	p := kafka.NewListClustersPaginator(svc, input)
	for p.HasMorePages() {
		page, err := p.NextPage(context.TODO())
		if err != nil {
			return &kafka.ListClustersOutput{}, err
		}
		output.ClusterInfoList = append(output.ClusterInfoList, page.ClusterInfoList...)
	}
	return output, nil
}
