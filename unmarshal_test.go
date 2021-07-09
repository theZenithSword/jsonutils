// Copyright 2019 Yunion
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

package jsonutils

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"yunion.io/x/pkg/gotypes"
	"yunion.io/x/pkg/tristate"
)

type TestStruct struct {
	Name   string
	Age    int
	Grade  uint8
	Gender string
	Status string
	Json   JSONObject
	Json2  JSONObject
	Array  []string
	Tri    tristate.TriState
}

func TestJSONDictUnmarshal(t *testing.T) {
	var err error
	ts := TestStruct{Name: "test", Age: 23, Grade: 2, Gender: "Male", Status: "Enabled", Tri: tristate.True}
	t.Logf("%s", Marshal(ts))
	json := NewDict()
	json.Add(NewString("name1"), "name")
	json.Add(NewInt(19), "age")
	json.Add(NewInt(3), "grade")
	json.Add(NewStringArray([]string{"1", "2", "3"}), "array")
	json.Add(JSONFalse, "tri")
	subDict := NewDict()
	subDict.Add(NewString("value"), "key")
	subDict.Add(NewString("value2"), "key2")
	json.Add(subDict, "json")
	subArray := NewArray()
	subArray.Add(NewString("arr1"))
	subArray.Add(NewString("arr2"))
	subArray.Add(NewString("arr3"))
	subArray.Add(NewString("arr4"))
	subArray.Add(NewString("arr5"))
	json.Add(subArray, "json2")
	t.Logf("%s", json.String())
	err = json.Unmarshal(&ts)
	if err != nil {
		t.Errorf("unmarshal struct fail: %s", err)
	} else {
		t.Logf("%s", Marshal(ts))
	}

	val := make(map[string]string)
	err = json.Unmarshal(val)
	if err != nil {
		t.Errorf("unmarshal map fail: %s", err)
	} else {
		t.Logf("%s", Marshal(val))
	}
}

func TestJSONDict_Unmarshal(t *testing.T) {
	type TestStruct struct {
		Id   string
		Name string
		Dict JSONObject // *JSONDict
	}
	jsonDict := NewDict()
	jsonDict.Add(NewString("nameVal"), "name")
	jsonDict.Add(NewString("idVal"), "id")
	subDict, err := ParseString(`{"parent_task_id": "30247a37-0328-4c47-bf5e-796672118923", "__stages": [{"complete_at": "2018-05-24T03:00:43Z", "name": "on_init"}], "__request_context": {"request_id": "5c2bd"}}`)
	if err != nil {
		t.Errorf("Parse json error")
	}
	// subDict := NewDict()
	// subDict.Add(NewString("yes"), "answer")
	// subDict.Add(NewInt(24), "age")
	jsonDict.Add(subDict, "dict")
	t.Logf("%s", jsonDict.String())

	dest := TestStruct{}

	jsonDict.Unmarshal(&dest)

	t.Logf("%s", dest)
	t.Logf("%s", Marshal(dest).String())

}

func TestJSONValue_Unmarshal(t *testing.T) {
	jsonvalue := &JSONValue{}
	var (
		obj1 = "json"
		obj2 = 2
		obj3 = struct {
			Name string
		}{
			"json",
		}
		obj4 JSONObject
	)
	var (
		want1 = ""
		want2 = 0
		want3 = struct {
			Name string
		}{}
		want4 JSONObject
	)
	cases := []struct {
		name    string
		obj     interface{}
		keys    []string
		wantObj interface{}
		isErr   bool
	}{
		{
			name:    "string test",
			obj:     &obj1,
			wantObj: &want1,
		},
		{
			name:    "int test",
			obj:     &obj2,
			wantObj: &want2,
		},
		{
			name:    "struct test",
			obj:     &obj3,
			wantObj: &want3,
		},
		{
			name:    "interface test",
			obj:     obj4,
			wantObj: want4,
		},
		{
			name:  "haven key test",
			keys:  []string{"json"},
			isErr: true,
		},
	}
	for _, c := range cases {
		err := jsonvalue.Unmarshal(c.obj, c.keys...)
		if err != nil && !c.isErr {
			t.Fatalf("There shouldn't be an error for %s", c.name)
		}
		if err == nil && c.isErr {
			t.Fatalf("There should be an error for %s", c.name)
		}
		if !reflect.DeepEqual(c.obj, c.wantObj) {
			t.Fatalf("For %s, want: %#v, actual: %#v", c.name, c.wantObj, c.obj)
		}
	}
}

func TestUnmarshalTime(t *testing.T) {
	type TimeStruct struct {
		EndTime time.Time
	}
	jsonDict := NewDict()
	jsonDict.Add(NewString(""), "end_time")
	t.Logf("json: %s", jsonDict.String())
	ts := TimeStruct{}
	err := jsonDict.Unmarshal(&ts)
	if err != nil {
		t.Errorf("unmarshal timestruct error %s", err)
	} else if !ts.EndTime.IsZero() {
		t.Fatalf("unmarshal empty time should zero")
	} else {
		t.Logf("unmarshal result %s", ts)
	}
}

func TestMarshalPtr(t *testing.T) {
	type SPtrs struct {
		Bool   *bool
		Int    *int
		Float  *float64
		String *string
		Struct *struct{ Hmm int }
		Array  *[9]int
		Slice  *[]int
		Map    *map[string]int
	}
	// marshal nils
	ptrsNil := &SPtrs{}
	jsonNil := Marshal(ptrsNil)
	jsonStrNil := jsonNil.String()
	if jsonStrNil != "{}" {
		t.Errorf("Should omit nil values, got %s", jsonStrNil)
	}

	// parse null JSON values
	jsonStrNil2 := `
		{
			bool:    null,
			int:     null,
			float:   null,
			string:  null,
			struct:  null,
			array:   null,
			slice:   null,
			map:     null
		}
	`
	jsonObjNil, err := ParseString(jsonStrNil2)
	if err != nil {
		t.Errorf("parse json string error: %v", err)
	}
	jsonDictNil := jsonObjNil.(*JSONDict)
	if numFields := reflect.TypeOf(SPtrs{}).NumField(); jsonDictNil.Length() != numFields {
		t.Errorf("num fields want %d, got %d", numFields, jsonDictNil.Length())
	}

	// make nonNil
	vBool := true
	vInt := 99
	vFloat := 99.9
	vString := "9999"
	vStruct := struct{ Hmm int }{99999}
	vArray := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	vSlice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	vMap := map[string]int{"999999": 1234567}
	ptrsNonNil := &SPtrs{
		Bool:   &vBool,
		Int:    &vInt,
		Float:  &vFloat,
		String: &vString,
		Struct: &vStruct,
		Array:  &vArray,
		Slice:  &vSlice,
		Map:    &vMap,
	}
	jsonStrNonNil := Marshal(ptrsNonNil).String()

	// unmarshal nils to non nils should perform override, partial if the source is not FULL
	{
		nils := *ptrsNonNil
		pnils := &nils
		jsonObjNil.Unmarshal(pnils)
		jsonObj2 := Marshal(pnils)
		jsonObj2Str := jsonObj2.String()
		if jsonObj2Str != "{}" {
			t.Errorf("unmarshal result should be {}, got %s", jsonObj2Str)
		}
	}

	// unmarshal non nil str will restore correctly
	{
		jsonObj, err := ParseString(jsonStrNonNil)
		if err != nil {
			t.Errorf("parse error: %s", err)
		}
		ptrs := &SPtrs{}
		err = jsonObj.Unmarshal(ptrs)
		if err != nil {
			t.Errorf("unmarshal error: %s", err)
		}
		if !reflect.DeepEqual(ptrs, ptrsNonNil) {
			t.Errorf("reverse failed: want\n%#v\ngot\n%#v", ptrsNonNil, ptrs)
		}
	}
}

func TestUnmarshalNonNilPtr(t *testing.T) {
	t.Run("non-nil short-cap slice", func(t *testing.T) {
		s := `[43]`
		j, _ := ParseString(s)
		v := []int{}
		vp := &v
		err := j.Unmarshal(&vp)
		if err != nil {
			t.Errorf("expect no error, got %s", err)
			return
		}
		if len(v) != 1 {
			t.Errorf("expect length 0, got %d", len(v))
			return
		}
		if v[0] != 43 {
			t.Errorf("expect [43], got %#v", v)
			return
		}
	})
	t.Run("non-nil over-cap slice", func(t *testing.T) {
		s := `[43]`
		j, _ := ParseString(s)
		v := []int{1, 2}
		vp := &v
		err := j.Unmarshal(&vp)
		if err != nil {
			t.Errorf("expect no error, got %s", err)
			return
		}
		if len(v) != 1 {
			t.Errorf("expect length 0, got %d", len(v))
			return
		}
		if v[0] != 43 {
			t.Errorf("expect [43], got %#v", v)
			return
		}
	})
	t.Run("non-nil map", func(t *testing.T) {
		s := `{"a": "b"}`
		j, _ := ParseString(s)
		v := struct {
			A string
		}{}
		vp := &v
		err := j.Unmarshal(&vp)
		if err != nil {
			t.Errorf("expect no error, got %s", err)
			return
		}
		if v.A != "b" {
			t.Errorf("expect v.A == \"b\", got %#v", v)
			return
		}
	})
}

func TestJSONArrayUnmarshal(t *testing.T) {
	s := `[{"conf":{"cachedbadbbu":false,"conf":"none","count":0,"direct":false,"ra":false,"range":[],"size":[],"strip":0,"type":"hybrid","wt":false},"disks":[{"adapter":0,"driver":"Linux","enclousure":0,"index":0,"max_strip_size":0,"min_strip_size":0,"rotate":true,"size":100000,"slot":0}],"size":100000}]`
	jsonArr, err := ParseString(s)
	if err != nil {
		t.Errorf("parse json error")
	}

	dest := JSONArray{}
	jsonArr.Unmarshal(&dest)
	if jsonArr1 := Marshal(dest); !jsonArr1.Equals(jsonArr) {
		t.Errorf("json array unmarshal, want:\n%s\ngot:\n%s", jsonArr.String(), jsonArr1.String())
	}
}

func TestUnmarshalCurrency(t *testing.T) {
	type SAccountBalance struct {
		USBalance     float64
		GermanBalance float32
	}
	jsonStr := `{"us_balance":"3,118.54", "german_balance":"3.490.000,89"}`
	json, err := ParseString(jsonStr)
	if err != nil {
		t.Errorf("parse %s error %s", jsonStr, err)
		return
	}
	balance := &SAccountBalance{}
	err = json.Unmarshal(&balance)
	if err != nil {
		t.Errorf("unmarshal %s fail %s", jsonStr, err)
		return
	}
	if balance.USBalance != float64(3118.54) {
		t.Fatalf("unmarshal us balance fail")
	}
	if balance.GermanBalance != float32(3490000.89) {
		t.Fatalf("unmarshal german balance fail!")
	}
}

func TestUnmarshalJsonTags(t *testing.T) {
	type SJsonTagStruct struct {
		Name    string `json:"OS:Name,omitempty"`
		Keyword string `json:"key_word,omitempty"`
	}
	cases := []struct {
		in   string
		want SJsonTagStruct
	}{
		{`{"name":"John","keyword":"json"}`, SJsonTagStruct{Name: "John", Keyword: "json"}},
		{`{"OS:Name":"John1","key_word":"json2"}`, SJsonTagStruct{Name: "John1", Keyword: "json2"}},
		{`{"Name":"John1","Keyword":"json2"}`, SJsonTagStruct{Name: "John1", Keyword: "json2"}},
		{`{"name":"John1","Keyword":"json2"}`, SJsonTagStruct{Name: "John1", Keyword: "json2"}},
		{`{"NAME":"John1","Keyword":"json2"}`, SJsonTagStruct{Name: "John1", Keyword: "json2"}},
	}
	for _, c := range cases {
		json, _ := ParseString(c.in)
		got := SJsonTagStruct{}
		err := json.Unmarshal(&got)
		if err != nil {
			t.Fatalf("unmarshal %s fail: %s", json, err)
		}
		if c.want.Name != got.Name || c.want.Keyword != got.Keyword {
			t.Fatalf("want %#v got %#v", c.want, got)
		}
	}
}

func TestUnmarshalEmbbedPtr(t *testing.T) {
	type OneStruct struct {
		Name string `json:"levelone:name"`
	}
	type TwoStruct struct {
		*OneStruct
		Gender string `json:"leveltwo:gender"`
	}

	cases := []struct {
		in   string
		want string
	}{
		{`{"levelone:name":"jack", "leveltwo:gender":"male"}`, "jack"},
		{`{"leveltwo:gender":"male"}`, ""},
	}
	for _, c := range cases {
		json, err := ParseString(c.in)
		if err != nil {
			t.Fatalf("fail to parse json %s %s", c.in, err)
		}
		got := TwoStruct{}
		err = json.Unmarshal(&got)
		if err != nil {
			t.Fatalf("fail to unmarshal %s %s", json.String(), err)
		}
		if got.Name != c.want {
			t.Fatalf("want %s got %s", c.want, got.Name)
		}
	}
}

type TestUnmarshalInterfaceI interface {
	String() string
	IsZero() bool
}
type TestUnmarshalInterfaceSI struct {
	Si int
}
type TestUnmarshalInterfaceS struct {
	M TestUnmarshalInterfaceI
}

func (si *TestUnmarshalInterfaceSI) IsZero() bool {
	return si.Si == 0
}

func (si *TestUnmarshalInterfaceSI) String() string {
	return fmt.Sprintf("%d", si.Si)
}

func TestUnmarshalInterface(t *testing.T) {
	t.Run("as-map-val", func(t *testing.T) {
		metadata := NewDict()
		metadata.Add(NewString("john"), "name")
		metadata.Add(NewInt(12), "age")
		metadata.Add(JSONTrue, "is_student")
		metadata.Add(NewFloat64(1.2), "weight")

		type MapKeyType string

		meta := make(map[MapKeyType]interface{}, 0)
		err := metadata.Unmarshal(meta)
		if err != nil {
			t.Fatalf("Get VM Metadata error: %v", err)
		}
	})

	t.Run("as-member", func(t *testing.T) {
		s := &TestUnmarshalInterfaceS{}
		gotypes.RegisterSerializable(reflect.TypeOf((*TestUnmarshalInterfaceI)(nil)).Elem(), func() gotypes.ISerializable {
			return &TestUnmarshalInterfaceSI{}
		})
		wantNum := 0xdeadbeef
		jsonStr := fmt.Sprintf(`{"m": {"si": %d}}`, wantNum)
		jo, err := ParseString(jsonStr)
		if err != nil {
			t.Fatalf("parse %q failed: %v", jsonStr, err)
		}
		err = jo.Unmarshal(s)
		if err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		si, ok := s.M.(*TestUnmarshalInterfaceSI)
		if !ok {
			t.Fatalf("expecting type *TestUnmarshalInterfaceSI, got %#v", s.M)
		}
		if si.Si != wantNum {
			t.Fatalf("want %x, got %x", wantNum, si.Si)
		}
	})
}

func TestUnmarshalString2Array(t *testing.T) {
	type sStruct struct {
		Provider []string `json:"provider"`
	}
	json := NewDict()
	json.Add(NewString("Aliyun"), "provider")
	s := sStruct{}
	err := json.Unmarshal(&s)
	if err != nil {
		t.Errorf("TestUnmarshalString2Array fail %s", err)
	}
	t.Logf("%s", s)
}

type ObsoleteStruct struct {
	CloudEnv  string `json:"cloud_env"`
	IsPublic  *bool  `json:"is_public"`
	Project   string `json:"project"`
	ProjectId string `json:"project_id" yunion-deprecated-by:"project"`
	Tenant    string `json:"tenant" yunion-deprecated-by:"project_id"`
	TenantId  string `json:"tenant_id" yunion-deprecated-by:"tenant"`
	Loop1     string `json:"loop1" yunion-deprecated-by:"loop2"`
	Loop2     string `json:"loop2" yunion-deprecated-by:"loop1"`
}

func (s *ObsoleteStruct) AfterUnmarshal() {
	if s.CloudEnv == "" && s.IsPublic != nil {
		if *s.IsPublic {
			s.CloudEnv = "public"
		} else {
			s.CloudEnv = "private"
		}
	}
}

type ObsoleteStruct2 struct {
	Hypervisors []string `json:"hypervisors"`
	Baremetal   *bool    `json:"baremetal"`
}

func (s *ObsoleteStruct2) AfterUnmarshal() {
	if s.Baremetal != nil && *s.Baremetal {
		s.Hypervisors = append(s.Hypervisors, "baremetal")
	}
}

type EmbedObsoleteStruct struct {
	ObsoleteStruct
	ObsoleteStruct2

	Name string `json:"name"`
}

type EmbedObsoleteStruct2 struct {
	*ObsoleteStruct
	*ObsoleteStruct2

	Name string `json:"name"`
}

func TestObsoleteBy(t *testing.T) {
	jsonVal := NewDict()
	jsonVal.Add(JSONTrue, "is_public")
	jsonVal.Add(NewString("testproject"), "tenant_id")
	jsonVal.Add(NewString("loop"), "loop1")
	jsonVal.Add(JSONTrue, "baremetal")

	t.Logf("origin: %s", jsonVal)
	s := ObsoleteStruct{}
	err := jsonVal.Unmarshal(&s)
	if err != nil {
		t.Fatalf("fail to unmarshal %s", err)
	}
	t.Logf("s: %s", Marshal(s))
	if s.CloudEnv != "public" || s.Project != "testproject" {
		t.Errorf("obsoleteby not work!")
	}

	s1 := EmbedObsoleteStruct{}
	err = jsonVal.Unmarshal(&s1)
	if err != nil {
		t.Fatalf("fail to unmarshal %s", err)
	}
	s1.Name = "s1"
	t.Logf("s1: %s", Marshal(s1))
	if s1.CloudEnv != "public" || s1.Project != "testproject" || len(s1.Hypervisors) == 0 || s1.Hypervisors[0] != "baremetal" {
		t.Errorf("obsoleteby not work!")
	}

	s2 := EmbedObsoleteStruct2{}
	err = jsonVal.Unmarshal(&s2)
	if err != nil {
		t.Fatalf("fail to unmarshal %s", err)
	}
	s2.Name = "s1"
	t.Logf("s2: %s", Marshal(s1))
	if s2.CloudEnv != "public" || s2.Project != "testproject" || len(s2.Hypervisors) == 0 || s2.Hypervisors[0] != "baremetal" {
		t.Errorf("obsoleteby not work!")
	}
}

func TestUnmarshalDuplicateEmbed(t *testing.T) {
	type Embeded struct {
		Provider string `json:"provider"`
	}
	type Struct1 struct {
		Embeded
		Cloudregion string `json:"cloudregion"`
	}
	type Struct2 struct {
		Embeded
		Cloudaccount string `json:"cloudaccount"`
	}
	type TopStruct struct {
		Struct1
		Struct2
	}

	json := NewDict()
	json.Add(NewString("Aliyun"), "provider")
	json.Add(NewString("region1"), "cloudregion")
	json.Add(NewString("account1"), "cloudaccount")

	s := TopStruct{}
	err := json.Unmarshal(&s)
	if err != nil {
		t.Fatalf("unmarshal fail %s", err)
	}

	if s.Struct1.Provider != s.Struct2.Provider {
		t.Fatalf("s.Struct1.Provier(%s) != s.Struct2.Provier(%s)", s.Struct1.Provider, s.Struct2.Provider)
	}

	t.Logf("%s", Marshal(s))
}

func TestUnmarshalString2Int(t *testing.T) {
	type sStruct struct {
		VlanId int     `json:"vlan_id"`
		Amount float64 `json:"amount"`
		BillId uint64  `json:"bill_id"`
	}
	cases := []struct {
		in   string
		want int
	}{
		{
			in:   `{"vlan_id":"", "amount":"", "bill_id":"200"}`,
			want: 0,
		},
		{
			in:   `{"vlan_id":"10", "amount":"10.0"}`,
			want: 10,
		},
		{
			in:   `{"vlan_id":"10,000", "amount":"10,000.00"}`,
			want: 10000,
		},
		{
			in:   "{\"vlan_id\":3\r}",
			want: 3,
		},
	}
	for _, c := range cases {
		json, err := ParseString(c.in)
		if err != nil {
			t.Fatalf("ParseString error %s", err)
		}
		s := sStruct{}
		err = json.Unmarshal(&s)
		if err != nil {
			t.Errorf("unmarshal vlan_id fail %s", err)
		}
		if s.VlanId != c.want {
			t.Errorf("Unmarshal %s got %d want %d", c.in, s.VlanId, c.want)
		}
	}
}

func TestUnmarshalJSONDictPtr(t *testing.T) {
	cases := []string{
		`{}`,
		`{"name":"jack"}`,
		`{"__request_context":{"request_id":"cbde96","service_name":"region","trace":{"debug":true,"duration":0,"id":"0","kind":"SERVER","local_endpoint":{"port":0,"service_name":"region"},"name":"delete","remote_endpoint":{"addr":"10.168.222.188","port":57240,"service_name":"(unknown_service)"},"shared":false,"tags":{"resource":"cloudaccounts"},"timestamp":"2020-04-11T06:37:07.100833Z","trace_id":"1866608c"}},"__stages":[{"complete_at":"2020-04-11T14:44:57Z","name":"on_init"}],"parent_task_id":"65617a87-3ecd-40c0-8add-70e17fec8ab2"}`,
	}
	for _, str := range cases {
		json, err := ParseString(str)
		if err != nil {
			t.Fatalf("parsestring fail %s", err)
		}
		dest := NewDict()
		err = json.Unmarshal(dest)
		if err != nil {
			t.Fatalf("Unmarshal *JSONDict fail %s", err)
		}
		if !json.Equals(dest) {
			t.Fatalf("want %s got %s", json, dest)
		}
		t.Logf("want %s", dest)
	}
}
