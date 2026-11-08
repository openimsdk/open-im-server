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

package mgo

import (
	"errors"
	"testing"

	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestIsAllDuplicateKeyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "generic non-duplicate error",
			err:      errors.New("connection refused"),
			expected: false,
		},
		{
			name: "single CommandError duplicate key 11000",
			err: mongo.CommandError{
				Code:    11000,
				Message: "E11000 duplicate key error collection: openim_v3.conversation",
			},
			expected: true,
		},
		{
			name: "single CommandError non-duplicate key 121",
			err: mongo.CommandError{
				Code:    121,
				Message: "Document failed validation",
			},
			expected: false,
		},
		{
			name: "BulkWriteException with single duplicate key write error",
			err: mongo.BulkWriteException{
				WriteErrors: []mongo.BulkWriteError{
					{
						WriteError: mongo.WriteError{
							Code:    11000,
							Message: "E11000 duplicate key error collection: openim_v3.conversation index: owner_user_id_1_conversation_id_1",
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "BulkWriteException with multiple duplicate key write errors",
			err: mongo.BulkWriteException{
				WriteErrors: []mongo.BulkWriteError{
					{
						WriteError: mongo.WriteError{
							Code:    11000,
							Message: "E11000 duplicate key error index 0",
						},
					},
					{
						WriteError: mongo.WriteError{
							Code:    11001,
							Message: "duplicate key error on update index 1",
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "BulkWriteException with mixed duplicate key and other write errors",
			err: mongo.BulkWriteException{
				WriteErrors: []mongo.BulkWriteError{
					{
						WriteError: mongo.WriteError{
							Code:    11000,
							Message: "E11000 duplicate key error",
						},
					},
					{
						WriteError: mongo.WriteError{
							Code:    121,
							Message: "Document failed validation",
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "BulkWriteException with duplicate key but has WriteConcernError",
			err: mongo.BulkWriteException{
				WriteConcernError: &mongo.WriteConcernError{
					Code:    64,
					Message: "waiting for replication timed out",
				},
				WriteErrors: []mongo.BulkWriteError{
					{
						WriteError: mongo.WriteError{
							Code:    11000,
							Message: "E11000 duplicate key error",
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "BulkWriteException with empty WriteErrors",
			err: mongo.BulkWriteException{
				WriteErrors: []mongo.BulkWriteError{},
			},
			expected: false,
		},
		{
			name: "WriteException with single duplicate key write error",
			err: mongo.WriteException{
				WriteErrors: []mongo.WriteError{
					{
						Code:    11000,
						Message: "E11000 duplicate key error",
					},
				},
			},
			expected: true,
		},
		{
			name: "WriteException with mixed errors",
			err: mongo.WriteException{
				WriteErrors: []mongo.WriteError{
					{
						Code:    11000,
						Message: "E11000 duplicate key error",
					},
					{
						Code:    14,
						Message: "TypeMismatch",
					},
				},
			},
			expected: false,
		},
		{
			name: "Wrapped BulkWriteException using errs.WrapMsg",
			err: errs.WrapMsg(mongo.BulkWriteException{
				WriteErrors: []mongo.BulkWriteError{
					{
						WriteError: mongo.WriteError{
							Code:    11000,
							Message: "E11000 duplicate key error",
						},
					},
				},
			}, "mongo bulk write"),
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsAllDuplicateKeyError(tc.err)
			if got != tc.expected {
				t.Errorf("IsAllDuplicateKeyError() = %v, want %v", got, tc.expected)
			}
		})
	}
}

