/*
Copyright 2020 The Knative Authors
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

package v1beta1

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	eventingduck "knative.dev/eventing/pkg/apis/duck/v1"
	duck "knative.dev/pkg/apis/duck/v1"

	"knative.dev/eventing-kafka/pkg/common/constants"
)

const (
	testNumPartitions     = 10
	testReplicationFactor = 5
	testRetentionDuration = "P1D"
)

func TestKafkaChannelDefaults(t *testing.T) {
	testCases := map[string]struct {
		initial  KafkaChannel
		expected KafkaChannel
	}{
		"nil spec": {
			initial: KafkaChannel{},
			expected: KafkaChannel{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"messaging.knative.dev/subscribable": "v1"},
				},
				Spec: KafkaChannelSpec{
					NumPartitions:     constants.DefaultNumPartitions,
					ReplicationFactor: constants.DefaultReplicationFactor,
					RetentionDuration: constants.DefaultRetentionISO8601Duration,
				},
			},
		},
		"numPartitions not set": {
			initial: KafkaChannel{
				Spec: KafkaChannelSpec{
					ReplicationFactor: testReplicationFactor,
					RetentionDuration: testRetentionDuration,
				},
			},
			expected: KafkaChannel{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"messaging.knative.dev/subscribable": "v1"},
				},
				Spec: KafkaChannelSpec{
					NumPartitions:     constants.DefaultNumPartitions,
					ReplicationFactor: testReplicationFactor,
					RetentionDuration: testRetentionDuration,
				},
			},
		},
		"replicationFactor not set": {
			initial: KafkaChannel{
				Spec: KafkaChannelSpec{
					NumPartitions:     testNumPartitions,
					RetentionDuration: testRetentionDuration,
				},
			},
			expected: KafkaChannel{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"messaging.knative.dev/subscribable": "v1"},
				},
				Spec: KafkaChannelSpec{
					NumPartitions:     testNumPartitions,
					ReplicationFactor: constants.DefaultReplicationFactor,
					RetentionDuration: testRetentionDuration,
				},
			},
		},
		"retentionDuration not set": {
			initial: KafkaChannel{
				Spec: KafkaChannelSpec{
					NumPartitions:     testNumPartitions,
					ReplicationFactor: testReplicationFactor,
				},
			},
			expected: KafkaChannel{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"messaging.knative.dev/subscribable": "v1"},
				},
				Spec: KafkaChannelSpec{
					NumPartitions:     testNumPartitions,
					ReplicationFactor: testReplicationFactor,
					RetentionDuration: constants.DefaultRetentionISO8601Duration,
				},
			},
		},
		"delivery.deadLetterSink.ref.namespace not set": {
			initial: KafkaChannel{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:   "custom",
					Name:        "ch",
					Annotations: map[string]string{"messaging.knative.dev/subscribable": "v1"},
				},
				Spec: KafkaChannelSpec{
					NumPartitions:     testNumPartitions,
					ReplicationFactor: testReplicationFactor,
					RetentionDuration: constants.DefaultRetentionISO8601Duration,
					ChannelableSpec: eventingduck.ChannelableSpec{
						Delivery: &eventingduck.DeliverySpec{
							DeadLetterSink: &duck.Destination{
								Ref: &duck.KReference{
									APIVersion: "v1",
									Name:       "svc",
									Kind:       "Service",
								},
							},
						},
					},
				},
			},
			expected: KafkaChannel{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:   "custom",
					Name:        "ch",
					Annotations: map[string]string{"messaging.knative.dev/subscribable": "v1"},
				},
				Spec: KafkaChannelSpec{
					NumPartitions:     testNumPartitions,
					ReplicationFactor: testReplicationFactor,
					RetentionDuration: constants.DefaultRetentionISO8601Duration,
					ChannelableSpec: eventingduck.ChannelableSpec{
						Delivery: &eventingduck.DeliverySpec{
							DeadLetterSink: &duck.Destination{
								Ref: &duck.KReference{
									APIVersion: "v1",
									Name:       "svc",
									Namespace:  "custom",
									Kind:       "Service",
								},
							},
						},
					},
				},
			},
		},
	}
	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			tc.initial.SetDefaults(context.TODO())
			if diff := cmp.Diff(tc.expected, tc.initial); diff != "" {
				t.Fatalf("Unexpected defaults (-want, +got): %s", diff)
			}
		})
	}
}
