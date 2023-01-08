package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func main() {
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		fmt.Println("Error creating session ", err)
	}
	svc := ecs.New(sess)
	fmt.Println(getEcsCluster(svc))
	taskDefArn := getServiceTaskDefinitionArn(svc)
	taskDefName := strings.Split(strings.Split(taskDefArn, "/")[1], ":")[0]
	imageVersion := getContainerImage(svc, taskDefName)
	fmt.Println(strings.Split(imageVersion, ":")[1])
}

func getEcsCluster(svc *ecs.ECS) string {
	result, err := svc.DescribeClusters(&ecs.DescribeClustersInput{Clusters: []*string{aws.String("default")}})
	if err != nil {
		panic(err)
	}
	return *result.Clusters[0].ClusterArn
}

func getServiceTaskDefinitionArn(svc *ecs.ECS) string {
	result, err := svc.DescribeServices(&ecs.DescribeServicesInput{Services: []*string{aws.String("nginx-service")}})
	if err != nil {
		panic(err)
	}
	return *result.Services[0].Deployments[0].TaskDefinition
}

func getContainerImage(svc *ecs.ECS, taskDefName string) string {
	result, err := svc.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{TaskDefinition: aws.String(taskDefName)})
	if err != nil {
		panic(err)
	}
	return *result.TaskDefinition.ContainerDefinitions[0].Image
}
