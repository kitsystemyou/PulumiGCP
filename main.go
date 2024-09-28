package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const kubeport = "8443"

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Compute Engineインスタンスの作成
		instance, err := compute.NewInstance(ctx, "my-instance", &compute.InstanceArgs{
			MachineType: pulumi.String("n1-standard-2"),
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
			// インスタンスの停止を許可したいときはtrueに設定
			// AllowStoppingForUpdate: pulumi.Bool(true),

			// インスタンスのタグを設定: http, https のトラフィックを許可
			Tags: pulumi.StringArray{
				pulumi.String("http-server"),
				// pulumi.String("https-server"),
			},
		})
		if err != nil {
			return err
		}


		// HTTP用ファイアウォールルールの作成
		_, err = compute.NewFirewall(ctx, "allow-http", &compute.FirewallArgs{
			Network: pulumi.String("default"),
			Allows: compute.FirewallAllowArray{
				&compute.FirewallAllowArgs{
					Protocol: pulumi.String("tcp"),
					Ports:    pulumi.StringArray{pulumi.String(kubeport)},
				},
			},
			TargetTags: pulumi.StringArray{pulumi.String("http-server")},
			
			// 実際には適切なソース範囲を指定する
			SourceRanges: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
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
