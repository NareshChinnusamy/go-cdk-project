package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsautoscaling"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type StackProps struct {
	awscdk.StackProps
}

func EcsClusterStack(scope constructs.Construct, id string, props StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	clusterName := "go-demo-cluster"
	vpcId := "vpc-0228fb56e3512f031"
	vpcName := "lookupvpc"
	autoscalineGroupName := "autoscalingGroup"
	containerInsights := false
	if &props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	vpc := awsec2.Vpc_FromLookup(stack, &vpcName, &awsec2.VpcLookupOptions{VpcId: &vpcId})

	cluster := awsecs.NewCluster(stack, &clusterName, &awsecs.ClusterProps{ClusterName: &clusterName, ContainerInsights: &containerInsights, Vpc: vpc})

	autoscaling := awsautoscaling.NewAutoScalingGroup(stack, &autoscalineGroupName, &awsautoscaling.AutoScalingGroupProps{InstanceType: awsec2.NewInstanceType(jsii.String("t2.small")),
		MachineImage:         awsecs.EcsOptimizedImage_AmazonLinux2(awsecs.AmiHardwareType_STANDARD, &awsecs.EcsOptimizedImageOptions{}),
		DesiredCapacity:      jsii.Number(1),
		Vpc:                  vpc,
		AutoScalingGroupName: jsii.String("GoDemoAutoscalingGroup")})

	capacity_provider := awsecs.NewAsgCapacityProvider(stack, jsii.String("AsgCapacityProvider"), &awsecs.AsgCapacityProviderProps{
		AutoScalingGroup: autoscaling, CapacityProviderName: jsii.String("GoDemoCapacityProvider"),
	})

	cluster.AddAsgCapacityProvider(capacity_provider, &awsecs.AddAutoScalingGroupCapacityOptions{})
	return stack
}

func EcsService(scope constructs.Construct, id string, props StackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if &props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	taskDef := awsecs.NewEc2TaskDefinition(stack, jsii.String("EcsDemoTaskDefinition"), &awsecs.Ec2TaskDefinitionProps{NetworkMode: awsecs.NetworkMode_BRIDGE, Family: jsii.String("EcsDemoTaskDefinition")})

	container := taskDef.AddContainer(jsii.String("FargoContainer"), &awsecs.ContainerDefinitionOptions{
		Image: awsecs.AssetImage_FromRegistry(jsii.String("docker pull nginx"), &awsecs.RepositoryImageProps{}), ContainerName: jsii.String("NginxDemo"), MemoryLimitMiB: jsii.Number(950), Cpu: jsii.Number(1024), Essential: jsii.Bool(true),
	})
	container.AddPortMappings(&awsecs.PortMapping{
		ContainerPort: jsii.Number(80),
		HostPort:      jsii.Number(80),
		Protocol:      awsecs.Protocol_TCP,
	})

	cluster := awsecs.Cluster_FromClusterArn(stack, jsii.String("lookUpCluster"), jsii.String("arn:aws:ecs:us-east-1:305251478828:cluster/go-demo-cluster"))

	awsecs.NewEc2Service(stack, jsii.String("EcsDemoServie"), &awsecs.Ec2ServiceProps{Cluster: cluster, DesiredCount: jsii.Number(1), ServiceName: jsii.String("EcsDemoServie"), TaskDefinition: taskDef})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	EcsClusterStack(app, "EcsClusterStack", StackProps{awscdk.StackProps{
		Env: env(),
	}})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("305251478828"),
		Region:  jsii.String("us-east-1"),
	}
}
