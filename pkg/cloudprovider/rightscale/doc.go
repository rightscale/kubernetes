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

// Package rightscale_cloud is an implementation of Interface, TCPLoadBalancer
// and Instances for RightScale.
//
// The RightScale provider relies on the config to retrieve the value of the
// RLL_SECRET HTTP header and the port number used to make API calls to
// RightScale through the local proxy operated by the RLL agent running on the
// instance.
package rightscale_cloud
