package main

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/eks"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// BEGIN Input Parameters
		iamUserName, configExists := ctx.GetConfig("iamUserName")
		if !configExists {
			return fmt.Errorf("iamUserName config is required")
		}
		subnetList, configExists := ctx.GetConfig("subnetList")
		if !configExists {
			return fmt.Errorf("subnetList config is required")
		}

		subnetIds := pulumi.StringArray{}
		for _, subnet := range subnetList {
			subnetIds = append(subnetIds, pulumi.String(subnet))
		}
		// END Input Parameters
		tmpJSON0, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				map[string]interface{}{
					"Action": []string{
						"sts:AssumeRole",
					},
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": "ec2.amazonaws.com",
					},
				},
			},
		})
		if err != nil {
			return err
		}
		json0 := string(tmpJSON0)
		node, err := iam.NewRole(ctx, "node", &iam.RoleArgs{
			Name:             pulumi.String("blog-eks-auto-node"),
			AssumeRolePolicy: pulumi.String(json0),
		})
		if err != nil {
			return err
		}
		tmpJSON1, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				map[string]interface{}{
					"Action": []string{
						"sts:AssumeRole",
						"sts:TagSession",
					},
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": "eks.amazonaws.com",
					},
				},
			},
		})
		if err != nil {
			return err
		}
		json1 := string(tmpJSON1)
		clusterRole, err := iam.NewRole(ctx, "blog-eks-auto-mode-cluster", &iam.RoleArgs{
			Name:             pulumi.String("blog-eks-auto-mode-cluster"),
			AssumeRolePolicy: pulumi.String(json1),
		})
		if err != nil {
			return err
		}
		clusterAmazonEKSClusterPolicy, err := iam.NewRolePolicyAttachment(ctx, "cluster_AmazonEKSClusterPolicy", &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
			Role:      clusterRole.Name,
		})
		if err != nil {
			return err
		}
		clusterAmazonEKSComputePolicy, err := iam.NewRolePolicyAttachment(ctx, "cluster_AmazonEKSComputePolicy", &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSComputePolicy"),
			Role:      clusterRole.Name,
		})
		if err != nil {
			return err
		}
		clusterAmazonEKSBlockStoragePolicy, err := iam.NewRolePolicyAttachment(ctx, "cluster_AmazonEKSBlockStoragePolicy", &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSBlockStoragePolicy"),
			Role:      clusterRole.Name,
		})
		if err != nil {
			return err
		}
		clusterAmazonEKSLoadBalancingPolicy, err := iam.NewRolePolicyAttachment(ctx, "cluster_AmazonEKSLoadBalancingPolicy", &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSLoadBalancingPolicy"),
			Role:      clusterRole.Name,
		})
		if err != nil {
			return err
		}
		clusterAmazonEKSNetworkingPolicy, err := iam.NewRolePolicyAttachment(ctx, "cluster_AmazonEKSNetworkingPolicy", &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSNetworkingPolicy"),
			Role:      clusterRole.Name,
		})
		if err != nil {
			return err
		}
		eksCluster, err := eks.NewCluster(ctx, "blog-cluster", &eks.ClusterArgs{
			Name: pulumi.String("blog-cluster"),
			AccessConfig: &eks.ClusterAccessConfigArgs{
				AuthenticationMode: pulumi.String("API"),
			},
			RoleArn: clusterRole.Arn,
			Version: pulumi.String("1.30"),
			ComputeConfig: &eks.ClusterComputeConfigArgs{
				Enabled: pulumi.Bool(true),
				NodePools: pulumi.StringArray{
					pulumi.String("general-purpose"),
				},
				NodeRoleArn: node.Arn,
			},
			KubernetesNetworkConfig: &eks.ClusterKubernetesNetworkConfigArgs{
				ElasticLoadBalancing: &eks.ClusterKubernetesNetworkConfigElasticLoadBalancingArgs{
					Enabled: pulumi.Bool(true),
				},
			},
			StorageConfig: &eks.ClusterStorageConfigArgs{
				BlockStorage: &eks.ClusterStorageConfigBlockStorageArgs{
					Enabled: pulumi.Bool(true),
				},
			},
			VpcConfig: &eks.ClusterVpcConfigArgs{
				EndpointPrivateAccess: pulumi.Bool(true),
				EndpointPublicAccess:  pulumi.Bool(true),
				SubnetIds:             pulumi.StringArray(subnetIds),
			},
			BootstrapSelfManagedAddons: pulumi.Bool(false),
		}, pulumi.DependsOn([]pulumi.Resource{
			clusterAmazonEKSClusterPolicy,
			clusterAmazonEKSComputePolicy,
			clusterAmazonEKSBlockStoragePolicy,
			clusterAmazonEKSLoadBalancingPolicy,
			clusterAmazonEKSNetworkingPolicy,
		}))
		if err != nil {
			return err
		}
		_, err = iam.NewRolePolicyAttachment(ctx, "node_AmazonEKSWorkerNodeMinimalPolicy", &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSWorkerNodeMinimalPolicy"),
			Role:      node.Name,
		})
		if err != nil {
			return err
		}
		_, err = iam.NewRolePolicyAttachment(ctx, "node_AmazonEC2ContainerRegistryPullOnly", &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPullOnly"),
			Role:      node.Name,
		})
		if err != nil {
			return err
		}

		// Reference an existing IAM user
		existingUser, err := iam.LookupUser(ctx, &iam.LookupUserArgs{
			UserName: iamUserName,
		}, nil)
		if err != nil {
			return err
		}

		// Create an AccessEntry to grant the IAM user access to the EKS cluster
		_, err = eks.NewAccessEntry(ctx, "eksAccessEntry", &eks.AccessEntryArgs{
			ClusterName:  eksCluster.Name,
			PrincipalArn: pulumi.String(existingUser.Arn),
		})
		if err != nil {
			return err
		}

		// Create an AccessPolicyAssociation to associate the policy with the EKS cluster
		_, err = eks.NewAccessPolicyAssociation(ctx, "eksAccessPolicyAssociation", &eks.AccessPolicyAssociationArgs{
			ClusterName:  eksCluster.Name,
			PrincipalArn: pulumi.String(existingUser.Arn),
			PolicyArn:    pulumi.String("arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"),
			AccessScope: eks.AccessPolicyAssociationAccessScopeArgs{
				Type: pulumi.String("cluster"),
			},
		})
		if err != nil {
			return err
		}

		// Export the cluster name as an output
		ctx.Export("clusterName", eksCluster.Name)
		return nil
	})
}
