// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package database

import "context"

// SubscribeUser Operation interface of user mongodb.
type SubscribeUser interface {
	// AddSubscriptionList Subscriber's handling of thresholds.
	AddSubscriptionList(ctx context.Context, userID string, userIDList []string) error
	// UnsubscriptionList Handling of unsubscribe.
	UnsubscriptionList(ctx context.Context, userID string, userIDList []string) error
	// RemoveSubscribedListFromUser Among the unsubscribed users, delete the user from the subscribed list.
	RemoveSubscribedListFromUser(ctx context.Context, userID string, userIDList []string) error
	// GetAllSubscribeList Get all users subscribed by this user
	GetAllSubscribeList(ctx context.Context, id string) (userIDList []string, err error)
	// GetSubscribedList Get the user subscribed by those users
	GetSubscribedList(ctx context.Context, id string) (userIDList []string, err error)
}
