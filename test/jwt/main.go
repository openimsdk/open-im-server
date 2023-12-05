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

package main

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	rawJWT := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI4MjkzODEzMTgzIiwiUGxhdGZvcm1JRCI6NSwiZXhwIjoxNzA2NTk0MTU0LCJuYmYiOjE2OTg4MTc4NTQsImlhdCI6MTY5ODgxODE1NH0.QCJHzU07SC6iYBoFO6Zsm61TNDor2D89I4E3zg8HHHU`

	// Verify the token
	claims := &jwt.MapClaims{}
	parsedT, err := jwt.ParseWithClaims(rawJWT, claims, func(token *jwt.Token) (any, error) {
		// Validate the alg is HMAC signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if kid, ok := token.Header["kid"].(string); ok {
			fmt.Println("kid", kid)
		}

		return []byte("key1"), nil
	})

	if err != nil || !parsedT.Valid {
		fmt.Println("token valid failed", err)

		return
	}

	fmt.Println("ok")
}
