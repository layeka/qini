// Copyright 2014 layeka Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package qini

import (
	"testing"
)

var ini *QIni

func init() {
	ini = Load("example.conf")
}
func TestQini(t *testing.T) {
	s, err := ini.GetValue("Wine", "Grape")
	if err != nil || s != "Cabernet Sauvignon" {
		t.Error("Gets the value of the error or the value returned is incorrect")
	}
	if !ini.DefaultBool("Pizza", "Ham", false) {
		t.Error("the value returned is incorrect")
	}
	if ini.DefaultInt("Wine", "Year", 0) != 1989 {
		t.Error("the value returned is incorrect")
	}
	if ini.DefaultFloat("Wine", "Alcohol", 0) != 12.5 {
		t.Error("the value returned is incorrect")
	}
}
