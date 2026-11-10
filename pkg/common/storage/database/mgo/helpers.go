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
	"strings"

	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsNotFound(err error) bool {
	return errs.Unwrap(err) == mongo.ErrNoDocuments
}

// IsAllDuplicateKeyError returns true if err is a duplicate key error or
// if err is a bulk write error where all write errors are duplicate key errors.
func IsAllDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	var bwe mongo.BulkWriteException
	if errors.As(err, &bwe) {
		if bwe.WriteConcernError != nil {
			return false
		}
		if len(bwe.WriteErrors) == 0 {
			return false
		}
		for _, we := range bwe.WriteErrors {
			if !isDuplicateKeyWriteError(we.WriteError) {
				return false
			}
		}
		return true
	}
	var we mongo.WriteException
	if errors.As(err, &we) {
		if we.WriteConcernError != nil {
			return false
		}
		if len(we.WriteErrors) == 0 {
			return false
		}
		for _, w := range we.WriteErrors {
			if !isDuplicateKeyWriteError(w) {
				return false
			}
		}
		return true
	}
	return mongo.IsDuplicateKeyError(err)
}

func isDuplicateKeyWriteError(we mongo.WriteError) bool {
	return we.Code == 11000 || we.Code == 11001 || we.Code == 12582 ||
		strings.Contains(we.Message, " E11000 ") ||
		strings.Contains(we.Message, "duplicate key") ||
		mongo.IsDuplicateKeyError(we)
}
