/*
Copyright © LiquidWeb

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/validate"
)

var networkLoadBalancerUpdateNodesCmd []string
var networkLoadBalancerUpdateServicesCmd []string

var networkLoadBalancerUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Change configuration of an existing Load Balancer",
	Long:  `Change configuration of an existing Load Balancer`,
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq_id")
		nameFlag, _ := cmd.Flags().GetString("name")
		strategyFlag, _ := cmd.Flags().GetString("strategy")
		enableSslTerminationFlag, _ := cmd.Flags().GetBool("enable-ssl-termination")
		disableSslTerminationFlag, _ := cmd.Flags().GetBool("disable-ssl-termination")
		sslPrivateKeyFlag, _ := cmd.Flags().GetString("ssl-private-key")
		sslCertFlag, _ := cmd.Flags().GetString("ssl-certificate")
		sslIntermediateCertFlag, _ := cmd.Flags().GetString("ssl-intermediate-certificate")
		enableSslIncludesFlag, _ := cmd.Flags().GetBool("enable-ssl-includes")
		disableSslIncludesFlag, _ := cmd.Flags().GetBool("disable-ssl-includes")

		if enableSslTerminationFlag && disableSslTerminationFlag {
			lwCliInst.Die(fmt.Errorf("can't both enable and disable ssl termination"))
		}
		if enableSslIncludesFlag && disableSslIncludesFlag {
			lwCliInst.Die(fmt.Errorf("can't both enable and disable ssl includes"))
		}
		if sslIntermediateCertFlag != "" {
			if !enableSslIncludesFlag {
				lwCliInst.Die(fmt.Errorf("when using --ssl-intermediate-certificate --enable-ssl-includes must be passed"))
			}
		}
		if sslCertFlag != "" || sslPrivateKeyFlag != "" {
			if !enableSslTerminationFlag {
				lwCliInst.Die(fmt.Errorf("when using --ssl-certificate or --ssl-private-key --enable-ssl-termination must be passed"))
			}
		}

		validateFields := map[interface{}]interface{}{
			uniqIdFlag: "UniqId",
		}

		apiArgs := map[string]interface{}{
			"uniq_id": uniqIdFlag,
		}

		// ssl termination
		if enableSslTerminationFlag {
			apiArgs["ssl_termination"] = true
		}
		if disableSslTerminationFlag {
			apiArgs["ssl_termination"] = false
		}

		// name
		if nameFlag != "" {
			apiArgs["name"] = nameFlag
		}

		// strategy
		if strategyFlag != "" {
			validateFields[strategyFlag] = "LoadBalancerStrategy"
			apiArgs["strategy"] = strategyFlag
		}

		// ssl includes
		if enableSslIncludesFlag {
			apiArgs["ssl_includes"] = true
			validateFields[sslIntermediateCertFlag] = "NonEmptyString"
		}
		if disableSslIncludesFlag {
			apiArgs["ssl_includes"] = false
		}

		// read and set ssl cert
		if sslCertFlag != "" {
			contents, err := ioutil.ReadFile(sslCertFlag)
			if err != nil {
				lwCliInst.Die(err)
			}
			strContents := cast.ToString(contents)
			apiArgs["ssl_cert"] = strContents
			validateFields[strContents] = "NonEmptyString"
		}

		// read and set ssl private key
		if sslPrivateKeyFlag != "" {
			contents, err := ioutil.ReadFile(sslPrivateKeyFlag)
			if err != nil {
				lwCliInst.Die(err)
			}
			strContents := cast.ToString(contents)
			apiArgs["ssl_key"] = strContents
			validateFields[strContents] = "NonEmptyString"
		}

		// read and set intermediate cert
		if sslIntermediateCertFlag != "" {
			contents, err := ioutil.ReadFile(sslIntermediateCertFlag)
			if err != nil {
				lwCliInst.Die(err)
			}
			strContents := cast.ToString(contents)
			apiArgs["ssl_int"] = strContents
			validateFields[strContents] = "NonEmptyString"
		}

		// nodes
		if len(networkLoadBalancerUpdateNodesCmd) > 0 {
			apiArgs["nodes"] = networkLoadBalancerUpdateNodesCmd
			for _, ip := range networkLoadBalancerUpdateNodesCmd {
				validateFields[ip] = "IP"
			}
		}

		// services TODO:  update its src:dest,src:dst form
		if len(networkLoadBalancerUpdateServicesCmd) > 0 {
			apiArgs["nodes"] = networkLoadBalancerUpdateServicesCmd
			for _, ip := range networkLoadBalancerUpdateServicesCmd {
				validateFields[ip] = "IP"
			}
		}

		if len(apiArgs) == 1 {
			lwCliInst.Die(fmt.Errorf("Must pass something to update. See 'help network load-balancer update'"))
		}

		// validate built input
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		// call the method, display results
		var update apiTypes.NetworkLoadBalancerDetails
		if err := lwCliInst.CallLwApiInto("bleed/network/loadbalancer/update", apiArgs,
			&update); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(update)
	},
}

func init() {
	networkLoadBalancerCmd.AddCommand(networkLoadBalancerUpdateCmd)
	networkLoadBalancerUpdateCmd.Flags().String("uniq_id", "", "uniq_id of Load Balancer")
	networkLoadBalancerUpdateCmd.Flags().String("strategy", "", "Load Balancer strategy")
	networkLoadBalancerUpdateCmd.Flags().String("name", "", "name of Load Balancer")
	networkLoadBalancerUpdateCmd.Flags().Bool("enable-ssl-termination", false, "enable ssl termination")
	networkLoadBalancerUpdateCmd.Flags().Bool("disable-ssl-termination", false, "disable ssl termination")
	networkLoadBalancerUpdateCmd.Flags().String("ssl-private-key", "", "path to ssl private key")
	networkLoadBalancerUpdateCmd.Flags().String("ssl-certificate", "", "path to ssl certificate")
	networkLoadBalancerUpdateCmd.Flags().String("ssl-intermediate-certificate", "", "path to ssl ssl intermediate certificate")
	networkLoadBalancerUpdateCmd.Flags().Bool("enable-ssl-includes", false, "enable ssl includes")
	networkLoadBalancerUpdateCmd.Flags().Bool("disable-ssl-includes", false, "disable ssl includes")

	networkLoadBalancerUpdateCmd.Flags().StringSliceVar(&networkLoadBalancerUpdateNodesCmd, "nodes",
		[]string{}, "nodes (ips) separated by ',' to balance via the Load Balancer")

	networkLoadBalancerUpdateCmd.Flags().StringSliceVar(&networkLoadBalancerUpdateServicesCmd, "services",
		[]string{}, "source/destination port pairs (such as 80:80) separated by ',' to balance via the Load Balancer")

	networkLoadBalancerUpdateCmd.MarkFlagRequired("uniq_id")
}
