package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
)

type StackProps struct {
	awscdk.StackProps
}

func S3BucketStack(scope constructs.Construct, id string, props StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	bucketName := "s3-demo-go-bucket"
	removalPolicy := awscdk.RemovalPolicy_DESTROY
	if &props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	awss3.NewBucket(stack, &bucketName, &awss3.BucketProps{BucketName: &bucketName, RemovalPolicy: removalPolicy})

	return stack
}

func EcsClusterStack(scope constructs.Construct, id string, props StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	clusterName := "go-demo-cluster"
	containerInsights := false
	if &props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	awsecs.NewCluster(stack, &clusterName, &awsecs.ClusterProps{ClusterName: &clusterName, ContainerInsights: &containerInsights})
	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	S3BucketStack(app, "S3BucketStack", StackProps{})

	EcsClusterStack(app, "EcsClusterStack", StackProps{})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
