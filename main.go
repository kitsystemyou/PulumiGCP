package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Compute Engineインスタンスの作成
		instance, err := compute.NewInstance(ctx, "my-instance", &compute.InstanceArgs{
			MachineType: pulumi.String("f1-micro"),
			Zone:        pulumi.String("us-central1-a"),
			BootDisk: &compute.InstanceBootDiskArgs{
				InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
					Image: pulumi.String("debian-cloud/debian-11"),
				},
			},
			NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
				&compute.InstanceNetworkInterfaceArgs{
					Network: pulumi.String("default"),
					AccessConfigs: compute.InstanceNetworkInterfaceAccessConfigArray{
						&compute.InstanceNetworkInterfaceAccessConfigArgs{},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// インスタンスの外部IPアドレスをエクスポート
		ctx.Export(
			"instanceIP",
			instance.NetworkInterfaces.Index(pulumi.Int(0)).
				AccessConfigs().
				Index(pulumi.Int(0)).
				NatIp(),
		)

		return nil
	})
}
