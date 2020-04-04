// Copyright 2020 PerceptiSys Ltd, (Stephen Gaito)
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

package webserver

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

// Test finding and describing routes in a webserver.
//
func TestFindRoute(t *testing.T) {

  ws := WS{}
  ws.BaseRoute = &Route{}

  aRoute, err := ws.FindRoute("/this/is/a/test")
  assert.NotNil(t, aRoute)
  assert.NotNil(t, err)
  assert.Equal(t, err.NumPartsFound, 0)
  assert.Equal(t, err.NumParts, 4)

  stdErr := ws.DescribeRoute("/this", "this description", true)
  assert.Nil(t, stdErr)

  stdErr = ws.DescribeRoute("/this/is", "this is description", true)
  assert.Nil(t, stdErr)

  aRoute, err = ws.FindRoute("/this/is")
  assert.NotNil(t, aRoute)
  assert.Equal(t, aRoute.Path, "/this/is")
  assert.Equal(t, aRoute.Prefix, "is")
  assert.Equal(t, aRoute.Desc,  "this is description")
  assert.Nil(t, err)
//  assert.Equal(t, err.NumPartsFound, 2)
//  assert.Equal(t, err.NumParts, 4)
//  assert.Equal(t, err.CurPrefix, "a")
//  assert.Contains(t, err.Message, "/this/is/a")

  aRoute, err = ws.FindRoute("/this/not")
  assert.NotNil(t, aRoute)
  assert.Equal(t, aRoute.Path, "/this")
  assert.Equal(t, aRoute.Prefix, "this")
  assert.Equal(t, aRoute.Desc,  "this description")
  assert.NotNil(t, err)
  assert.Equal(t, err.NumPartsFound, 1)
  assert.Equal(t, err.NumParts, 2)
  assert.Equal(t, err.CurPrefix, "not")
  assert.Contains(t, err.Message, "/this/not")

  aRoute, err = ws.FindRoute("/this/is/a/test")
  assert.NotNil(t, aRoute)
  assert.Equal(t, aRoute.Path, "/this/is")
  assert.Equal(t, aRoute.Prefix, "is")
  assert.Equal(t, aRoute.Desc,  "this is description")
  assert.NotNil(t, err)
  assert.Equal(t, err.NumPartsFound, 2)
  assert.Equal(t, err.NumParts, 4)
  assert.Equal(t, err.CurPrefix, "a")
  assert.Contains(t, err.Message, "/this/is/a")
}
