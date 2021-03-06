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
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/liquidweb/liquidweb-cli/instance"
	"github.com/liquidweb/liquidweb-cli/types/api"
	"github.com/liquidweb/liquidweb-cli/types/cmd"
	"github.com/liquidweb/liquidweb-cli/validate"
)

var networkLoadBalancerUpdateNodesCmd []string
var networkLoadBalancerUpdateServicesCmd []string
var healthChecksMapUpdate map[string]string
var networkLoadBalancerServicesHealthChecksHelp = `--services flag values are ',' delimited. Each value should be in format:

  'sourcePort:destinationPort',

such as '80:80,443:443'.

--health-check flag values represent custom health check paramaters for a service on a Load Balancer. Valid health check parameters:

  failure_threshold -> int // permissible failures before node is taken out of services pool (default 3)
  http_body_match -> string // when protocol is http, the string to look for in the http body to determine if health is ok (default unset)
  http_path -> string // when protocol is http, the http path to hit when performing a health check (default /)
  http_response_codes -> string // when protocol is http, http response codes to consider "success" when performing a health check (default 200-206:300-304)
  http_use_tls -> "bool" // when protocol is http, uses https when "true" for health check (default false)
  interval -> int // time duration between health checks (default 30)
  protocol -> string // *Required (valid values: tcp, http)
  timeout -> int // timeout value for the health check probe (default 5)

For example, to set these values for the service with source port 443, the flag could look like this:

  --health-check 443_failure_threshold=12,443_http_body_match=hello,443_http_path=/status,443_http_response_codes=200:201:202,443_http_use_tls=true,443_interval=10,443_protocol=http,443_timeout=99

Notice the leading '443_' before the parameter name. To create a health check for service 80 as well, follow the same pattern, but
replacing '443_' with '80_'.`
var networkLoadBalancerServicesHealthCheckFileHelp = `--health-check-file value should be the path to a yaml file containing the health check(s) to apply for each service(s). Here is
an example of how that file should look:

443:
  protocol: http
  timeout: 5
  interval: 10
  http_use_tls: true
  http_response_codes: 200-202,404
  http_path: /status-443
  http_body_match:
  failure_threshold: 3
80:
  protocol: http
  timeout: 10
  interval: 20
  http_use_tls: false
  http_response_codes: 200-202,404
  http_path: /status-80
  http_body_match:
  failure_threshold: 3

It is an error to provide both --health-check and --health-check-file flags.`

var networkLoadBalancerUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Change configuration of an existing Load Balancer",
	Long: fmt.Sprintf(`Change configuration of an existing Load Balancer

A Load Balancer allows you to distribute traffic to multiple endpoints.

%s

%s

To remove a health check from a service, simply call update for the service(s) omitting their --health-check entries. For example,
this would remove any set health checks for services 443:443,80:80 (as well as remove any other services entirely):

network load-balancer update --uniq-id ABC123 --services 443:443,80:80

Similarly to remove a health check when using --health-check-file, simply remove the health check from the file.
`, networkLoadBalancerServicesHealthChecksHelp, networkLoadBalancerServicesHealthCheckFileHelp),
	Run: func(cmd *cobra.Command, args []string) {
		uniqIdFlag, _ := cmd.Flags().GetString("uniq-id")
		nameFlag, _ := cmd.Flags().GetString("name")
		strategyFlag, _ := cmd.Flags().GetString("strategy")
		enableSslTerminationFlag, _ := cmd.Flags().GetBool("enable-ssl-termination")
		disableSslTerminationFlag, _ := cmd.Flags().GetBool("disable-ssl-termination")
		sslPrivateKeyFlag, _ := cmd.Flags().GetString("ssl-private-key")
		sslCertFlag, _ := cmd.Flags().GetString("ssl-certificate")
		sslIntermediateCertFlag, _ := cmd.Flags().GetString("ssl-intermediate-certificate")
		enableSslIncludesFlag, _ := cmd.Flags().GetBool("enable-ssl-includes")
		disableSslIncludesFlag, _ := cmd.Flags().GetBool("disable-ssl-includes")
		healthCheckFileFlag, _ := cmd.Flags().GetString("health-check-file")

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
		if len(healthChecksMapUpdate) > 0 && healthCheckFileFlag != "" {
			lwCliInst.Die(fmt.Errorf("cannot pass conflicting flags --health-check and --health-check-file"))
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
			contents, err := ioutil.ReadFile(filepath.Clean(sslCertFlag))
			if err != nil {
				lwCliInst.Die(err)
			}
			strContents := cast.ToString(contents)
			apiArgs["ssl_cert"] = strContents
			validateFields[strContents] = "NonEmptyString"
		}

		// read and set ssl private key
		if sslPrivateKeyFlag != "" {
			contents, err := ioutil.ReadFile(filepath.Clean(sslPrivateKeyFlag))
			if err != nil {
				lwCliInst.Die(err)
			}
			strContents := cast.ToString(contents)
			apiArgs["ssl_key"] = strContents
			validateFields[strContents] = "NonEmptyString"
		}

		// read and set intermediate cert
		if sslIntermediateCertFlag != "" {
			contents, err := ioutil.ReadFile(filepath.Clean(sslIntermediateCertFlag))
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

		// services
		if len(networkLoadBalancerUpdateServicesCmd) > 0 {
			var servicesToBalance []map[string]interface{}
			// a service is permitted to have one health check.

			var healthChecks map[string]map[string]interface{}

			// health check, command line flags.
			if len(healthChecksMapUpdate) > 0 {
				healthChecksFromCmdLine, err := cmdTypes.LoadBalancerHealthCheckCmdLine{HealthCheck: healthChecksMapUpdate}.Transform()
				if err != nil {
					lwCliInst.Die(err)
				}
				healthChecks = healthChecksFromCmdLine
			} else if healthCheckFileFlag != "" {
				// health check, yaml file
				contents, err := ioutil.ReadFile(filepath.Clean(healthCheckFileFlag))
				if err != nil {
					lwCliInst.Die(fmt.Errorf("error reading given --health-check-file [%s]: %s", healthCheckFileFlag, err))
				}
				if err = yaml.Unmarshal(contents, &healthChecks); err != nil {
					lwCliInst.Die(fmt.Errorf("error yaml decoding [%s] (see help for an example of the file); %s", healthCheckFileFlag, err))
				}
			}

			// validate
			for _, healthCheck := range healthChecks {
				var obj apiTypes.NetworkLoadBalancerDetailsServiceHealthCheck
				if err := instance.CastFieldTypes(healthCheck, &obj); err != nil {
					lwCliInst.Die(fmt.Errorf(
						"failed casting --health-check-file [%s] to expected structure (see help for an example of the file): %s",
						healthCheckFileFlag, err))
				}
				if err := obj.Validate(); err != nil {
					lwCliInst.Die(err)
				}
			}

			// build services api argument
			for _, pair := range networkLoadBalancerUpdateServicesCmd {
				err := validate.Validate(map[interface{}]interface{}{pair: "NetworkPortPair"})
				if err != nil {
					lwCliInst.Die(err)
				}

				splitPair := strings.Split(pair, ":")
				srcPort := cast.ToInt(splitPair[0])
				destPort := cast.ToInt(splitPair[1])

				serviceToBalance := map[string]interface{}{
					"src_port":  srcPort,
					"dest_port": destPort,
				}

				// if a health check exists for this service set it
				if _, exists := healthChecks[splitPair[0]]; exists {
					serviceToBalance["health_check"] = healthChecks[splitPair[0]]
				}

				servicesToBalance = append(servicesToBalance, serviceToBalance)
			}
			apiArgs["services"] = servicesToBalance
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
	networkLoadBalancerUpdateCmd.Flags().String("uniq-id", "", "uniq-id of Load Balancer")
	networkLoadBalancerUpdateCmd.Flags().String("strategy", "", "Load Balancer strategy (see 'network load-balancer get-strategies')")
	networkLoadBalancerUpdateCmd.Flags().String("name", "", "name of Load Balancer")
	networkLoadBalancerUpdateCmd.Flags().Bool("enable-ssl-termination", false, "enable ssl termination")
	networkLoadBalancerUpdateCmd.Flags().Bool("disable-ssl-termination", false, "disable ssl termination")
	networkLoadBalancerUpdateCmd.Flags().String("ssl-private-key", "", "path to ssl private key")
	networkLoadBalancerUpdateCmd.Flags().String("ssl-certificate", "", "path to ssl certificate")
	networkLoadBalancerUpdateCmd.Flags().String("ssl-intermediate-certificate", "", "path to ssl ssl intermediate certificate")
	networkLoadBalancerUpdateCmd.Flags().Bool("enable-ssl-includes", false, "enable ssl includes")
	networkLoadBalancerUpdateCmd.Flags().Bool("disable-ssl-includes", false, "disable ssl includes")

	networkLoadBalancerUpdateCmd.Flags().StringSliceVar(&networkLoadBalancerUpdateNodesCmd, "nodes",
		[]string{}, "nodes (ips) separated by ',' to balance via the Load Balancer (see 'network load-balancer get-possible-nodes')")

	networkLoadBalancerUpdateCmd.Flags().StringSliceVar(&networkLoadBalancerUpdateServicesCmd, "services",
		[]string{}, "source/destination port pairs (such as 80:80) separated by ',' to balance via the Load Balancer")

	networkLoadBalancerUpdateCmd.Flags().StringToStringVar(&healthChecksMapUpdate, "health-check", nil,
		"Health check defintions for the service matching source port. Should not be combined with --health-check.")

	networkLoadBalancerUpdateCmd.Flags().String("health-check-file", "",
		"A file containing valid yaml describing the LoadBalancer health checks to add for the service(s). Should not be combined with --health-check.")

	if err := networkLoadBalancerUpdateCmd.MarkFlagRequired("uniq-id"); err != nil {
		lwCliInst.Die(err)
	}
}
