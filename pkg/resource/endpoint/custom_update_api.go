// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Use this file to add custom implementation for any operation of intercept
// the autogenerated code that trigger an update on an endpoint

package endpoint

import (
	"context"
	"errors"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

// customUpdateEndpoint adds specialized logic to requeueAfter until endpoint is in
// InService state before an updateEndpoint can be called
func (rm *resourceManager) customUpdateEndpoint(
	ctx context.Context,
	desired *resource,
	latest *resource,
	diffReporter *ackcompare.Reporter,
) (*resource, error) {
	return nil, rm.endpointStatusAllowUpdates(ctx, latest)
}

// customDeleteEndpoint adds specialized logic to requeueAfter until endpoint is in
// InService or Failed state before a deleteEndpoint can be called
func (rm *resourceManager) customDeleteEndpoint(
	ctx context.Context,
	latest *resource,
) error {
	latestStatus := latest.ko.Status.EndpointStatus
	if latestStatus != nil && *latestStatus == svcsdk.EndpointStatusFailed {
		return nil
	}
	return rm.endpointStatusAllowUpdates(ctx, latest)
}

// endpointStatusAllowUpdates is a helper method to determine if endpoint allows modification
func (rm *resourceManager) endpointStatusAllowUpdates(
	ctx context.Context,
	r *resource,
) error {
	latestStatus := r.ko.Status.EndpointStatus
	if latestStatus != nil && *latestStatus != svcsdk.EndpointStatusInService {
		return requeue.NeededAfter(
			errors.New("endpoint status does not allow modification, it is not in 'InService' state"),
			requeue.DefaultRequeueAfterDuration)
	}

	return nil
}