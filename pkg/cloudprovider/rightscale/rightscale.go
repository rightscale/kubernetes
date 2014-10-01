/*
Copyright 2014 RightScale Inc. All rights reserved.

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

package rightscale_cloud

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"code.google.com/p/gcfg"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/cloudprovider"
)

// RightScaleCloud is an implementation of Interface, TCPLoadBalancer and Instances for RightScale instances.
type RightScaleCloud struct {
	Rll struct {
		Secret string // Auth secret for local API proxy requests
		Port   string // Port to local API proxy
	}
}

func init() {
	cloudprovider.RegisterCloudProvider("rightscale",
		func(config io.Reader) (cloudprovider.Interface, error) {
			return newRightScaleCloud(config)
		})
}

// newRightScaleCloud creates a new instance of RightScaleCloud configured to talk to the RightScale API.
func newRightScaleCloud(config io.Reader) (*RightScaleCloud, error) {
	if config == nil {
		return nil, fmt.Errorf("missing configuration file for RightScale cloud provider")
	}
	var cloud RightScaleCloud
	if err := gcfg.ReadInto(&cloud, config); err != nil {
		return nil, err
	}

	return &cloud, nil
}

// Instances returns an implementation of Instances for RightScale
func (r *RightScaleCloud) Instances() (cloudprovider.Instances, bool) {
	return r, true
}

// Zones not supported yet
func (r *RightScaleCloud) Zones() (cloudprovider.Zones, bool) {
	return nil, false
}

// External load balancer not supported yet
func (r *RightScaleCloud) TCPLoadBalancer() (cloudprovider.TCPLoadBalancer, bool) {
	return nil, false
}

// IPAddress returns the address of a particular machine instance.
func (r *RightScaleCloud) IPAddress(instance string) (net.IP, error) {
	resp, err := r.sendRequest("/api/servers?filter[]=name%3d%3d" + instance + "?view=instance_detail")
	if err != nil {
		return net.IP{}, nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return net.IP{}, err
	}

	var servers []map[string]interface{}
	err = json.Unmarshal(body, &servers)
	if err != nil {
		return net.IP{}, err
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("Could not retrieve server with name %s", instance)
	}
	if len(servers) > 1 {
		return nil, fmt.Errorf("More than one server with name %s", instance)
	}

	current := servers[0]["current_instance"]
	if current == nil {
		return nil, fmt.Errorf("Failed to retrieve current instance of %s", instance)
	}

	cur := current.(map[string]interface{})
	ipAddresses := cur["private_ip_addresses"].([]string)
	if current == nil {
		return nil, fmt.Errorf("Failed to retrieve private IPs of %s", instance)
	}

	ip := net.ParseIP(ipAddresses[0])
	if ip == nil {
		return nil, fmt.Errorf("Invalid network IP: %s", ipAddresses[0])
	}

	return ip, nil
}

// List enumerates the set of minions instances known by the cloud provider.
func (r *RightScaleCloud) List(filter string) ([]string, error) {
	resp, err := r.sendRequest("/api/servers?filter[]=name%3d%3d" + url.QueryEscape(filter))
	if err != nil {
		return []string{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	var servers []map[string]interface{}
	err = json.Unmarshal(body, &servers)
	if err != nil {
		return []string{}, err
	}

	names := make([]string, len(servers))
	for i, s := range servers {
		names[i] = s["name"].(string)
	}

	return names, nil
}

// Send GET request to given RightScale API uri (e.g. "/api/instances")
// uri must start with "/"
func (r *RightScaleCloud) sendRequest(uri string) (resp *http.Response, err error) {
	url := "http://localhost:" + r.Rll.Port + uri
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("RLL_SECRET", r.Rll.Secret)
	cl := http.Client{}
	resp, err = cl.Do(req)

	return
}
