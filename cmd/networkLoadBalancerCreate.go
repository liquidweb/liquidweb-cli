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
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/validate"
)

var networkLoadBalancerCreateNodesCmd []string
var networkLoadBalancerCreateServicesCmd []string

var networkLoadBalancerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Load Balancer",
	Long:  `Create a Load Balancer`,
	Run: func(cmd *cobra.Command, args []string) {
		nameFlag, _ := cmd.Flags().GetString("name")
		strategyFlag, _ := cmd.Flags().GetString("strategy")
		enableSslTerminationFlag, _ := cmd.Flags().GetBool("enable-ssl-termination")
		disableSslTerminationFlag, _ := cmd.Flags().GetBool("disable-ssl-termination")
		sslPrivateKeyFlag, _ := cmd.Flags().GetString("ssl-private-key")
		sslCertFlag, _ := cmd.Flags().GetString("ssl-certificate")
		sslIntermediateCertFlag, _ := cmd.Flags().GetString("ssl-intermediate-certificate")
		enableSslIncludesFlag, _ := cmd.Flags().GetBool("enable-ssl-includes")
		disableSslIncludesFlag, _ := cmd.Flags().GetBool("disable-ssl-includes")
		regionFlag, _ := cmd.Flags().GetInt("region")

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
			strategyFlag: "LoadBalancerStrategy",
			nameFlag:     "NonEmptyString",
		}

		apiArgs := map[string]interface{}{
			"name":     nameFlag,
			"strategy": strategyFlag,
		}

		if regionFlag != 0 {
			validateFields[regionFlag] = "PositiveInt"
			apiArgs["region"] = regionFlag
		}

		// ssl termination
		if enableSslTerminationFlag {
			apiArgs["ssl_termination"] = true
		}
		if disableSslTerminationFlag {
			apiArgs["ssl_termination"] = false
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
		if len(networkLoadBalancerCreateNodesCmd) > 0 {
			apiArgs["nodes"] = networkLoadBalancerCreateNodesCmd
			for _, ip := range networkLoadBalancerCreateNodesCmd {
				validateFields[ip] = "IP"
			}
		}

		// services
		if len(networkLoadBalancerCreateServicesCmd) == 0 {
			lwCliInst.Die(fmt.Errorf("--services must have source/destination port pairs (see 'help network load-balancer create')"))
		}
		// slice of maps with keys src_port, dest_port, with a value of its network port number.
		var servicesToBalance []map[string]int

		for _, pair := range networkLoadBalancerCreateServicesCmd {
			err := validate.Validate(map[interface{}]interface{}{pair: "NetworkPortPair"})
			if err != nil {
				lwCliInst.Die(err)
			}

			splitPair := strings.Split(pair, ":")
			srcPort := cast.ToInt(splitPair[0])
			destPort := cast.ToInt(splitPair[1])

			servicesToBalance = append(servicesToBalance, map[string]int{
				"src_port":  srcPort,
				"dest_port": destPort,
			})
		}

		apiArgs["services"] = servicesToBalance

		// validate built input
		if err := validate.Validate(validateFields); err != nil {
			lwCliInst.Die(err)
		}

		// call the method, display results
		var create apiTypes.NetworkLoadBalancerDetails
		if err := lwCliInst.CallLwApiInto("bleed/network/loadbalancer/create", apiArgs,
			&create); err != nil {
			lwCliInst.Die(err)
		}

		fmt.Print(create)
	},
}

func init() {
	networkLoadBalancerCmd.AddCommand(networkLoadBalancerCreateCmd)
	networkLoadBalancerCreateCmd.Flags().String("strategy", "", "Load Balancer strategy (see 'network load-balancer get-strategies')")
	networkLoadBalancerCreateCmd.Flags().String("name", "", "name of Load Balancer")
	networkLoadBalancerCreateCmd.Flags().Bool("enable-ssl-termination", false, "enable ssl termination")
	networkLoadBalancerCreateCmd.Flags().Bool("disable-ssl-termination", false, "disable ssl termination")
	networkLoadBalancerCreateCmd.Flags().String("ssl-private-key", "", "path to ssl private key")
	networkLoadBalancerCreateCmd.Flags().String("ssl-certificate", "", "path to ssl certificate")
	networkLoadBalancerCreateCmd.Flags().String("ssl-intermediate-certificate", "", "path to ssl ssl intermediate certificate")
	networkLoadBalancerCreateCmd.Flags().Bool("enable-ssl-includes", false, "enable ssl includes")
	networkLoadBalancerCreateCmd.Flags().Bool("disable-ssl-includes", false, "disable ssl includes")
	networkLoadBalancerCreateCmd.Flags().Int("region", 0, "region id to create a Load Balancer in (see 'cloud server options --zones')")

	networkLoadBalancerCreateCmd.Flags().StringSliceVar(&networkLoadBalancerCreateNodesCmd, "nodes",
		[]string{}, "nodes (ips) separated by ',' to balance via the Load Balancer (see 'network load-balancer get-possible-nodes')")

	networkLoadBalancerCreateCmd.Flags().StringSliceVar(&networkLoadBalancerCreateServicesCmd, "services",
		[]string{}, "source/destination port pairs (such as 80:80) separated by ',' to balance via the Load Balancer")

	networkLoadBalancerCreateCmd.MarkFlagRequired("name")
	networkLoadBalancerCreateCmd.MarkFlagRequired("services")
	networkLoadBalancerCreateCmd.MarkFlagRequired("strategy")
}
