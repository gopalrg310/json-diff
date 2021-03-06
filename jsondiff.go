package jsondiff

import (
	"reflect"
)

func GetDiffJSONValue(types string, Exact interface{}, comparitive interface{}, Changes bool) (bool, interface{}, interface{}, error) {
	var data1, data2 interface{}
	exactduplicate := make(map[string]interface{})
	var err error
	switch reflect.TypeOf(comparitive).Kind().String() {
	case "map":
		if reflect.TypeOf(comparitive) == reflect.TypeOf(Exact) {
			result_map := make(map[string]interface{})
			old_map := make(map[string]interface{})
			for key, value := range comparitive.(map[string]interface{}) {
				if avail_value, ok := Exact.(map[string]interface{})[key]; !ok {
					Changes = true
					instant_result1, instant_result2 := GetDiffOutput(types, avail_value, value)
					result_map[key] = instant_result1
					old_map[key] = instant_result2
				} else if reflect.DeepEqual(avail_value, value) {
					exactduplicate[key] = true
				} else {
					exactduplicate[key] = true
					value_type := reflect.ValueOf(value).Kind().String()
					if reflect.TypeOf(avail_value) != reflect.TypeOf(value) {
						Changes = true
						instant_result1, instant_result2 := GetDiffOutput(types, avail_value, value)
						result_map[key] = instant_result1
						old_map[key] = instant_result2
					} else if value_type == "map" || value_type == "slice" {
						if Changes, data1, data2, err = GetDiffJSONValue(types, avail_value, value, Changes); data1 != nil {
							result_map[key] = data1
							old_map[key] = data2
						} else if err != nil {
							return false, nil, nil, err
						}
					} else {
						Changes = true
						instant_result1, instant_result2 := GetDiffOutput(types, avail_value, value)
						result_map[key] = instant_result1
						old_map[key] = instant_result2
					}
				}
			}
			for key, value := range Exact.(map[string]interface{}) {
				if _, ok := exactduplicate[key]; !ok {
					Changes = true
					instant_result1, instant_result2 := GetDiffOutput(types, value, nil)
					result_map[key] = instant_result1
					old_map[key] = instant_result2
				}
			}
			if len(result_map) > 0 {
				return Changes, result_map, old_map, nil
			}
			return Changes, nil, nil, nil
		} else {
			Changes = true
			instant_result1, instant_result2 := GetDiffOutput(types, Exact, comparitive)
			return Changes, instant_result1, instant_result2, nil
		}

	case "slice":
		var resultant_val, newslice_val interface{}
		if reflect.TypeOf(Exact) == reflect.TypeOf(comparitive) {
			result_slice := make([]interface{}, 0)
			old_slice := make([]interface{}, 0)
			if _, ok := comparitive.([]string); ok {
				for i, j := range comparitive.([]string) {
					if len(Exact.([]string)) <= i {
						if Changes, resultant_val, newslice_val, err = GetDiffJSONValue(types, nil, j, Changes); resultant_val != nil {
							result_slice = append(result_slice, resultant_val)
							old_slice = append(old_slice, newslice_val)
						}
					} else if Changes, resultant_val, newslice_val, err = GetDiffJSONValue(types, Exact.([]string)[i], j, Changes); resultant_val != nil {
						result_slice = append(result_slice, resultant_val)
						old_slice = append(old_slice, newslice_val)
					}
				}
				if len(Exact.([]string)) > len(comparitive.([]string)) {
					for _, j := range Exact.([]string)[len(comparitive.([]string)):] {
						instant_result1, instant_result2 := GetDiffOutput(types, j, nil)
						result_slice = append(result_slice, instant_result1)
						old_slice = append(old_slice, instant_result2)
					}
				}
				if len(result_slice) > 0 {
					Changes = true
					return Changes, result_slice, old_slice, nil
				}
				return Changes, nil, nil, nil
			} else {
				for i, j := range comparitive.([]interface{}) {
					if len(Exact.([]interface{})) <= i {
						if Changes, resultant_val, newslice_val, err = GetDiffJSONValue(types, nil, j, Changes); resultant_val != nil {
							result_slice = append(result_slice, resultant_val)
							old_slice = append(old_slice, newslice_val)
						}
					} else if Changes, resultant_val, newslice_val, err = GetDiffJSONValue(types, Exact.([]interface{})[i], j, Changes); resultant_val != nil {
						result_slice = append(result_slice, resultant_val)
						old_slice = append(old_slice, newslice_val)
					}
				}
				if len(Exact.([]interface{})) > len(comparitive.([]interface{})) {
					for _, j := range Exact.([]interface{})[len(comparitive.([]interface{})):] {
						instant_result1, instant_result2 := GetDiffOutput(types, j, nil)
						result_slice = append(result_slice, instant_result1)
						old_slice = append(old_slice, instant_result2)
					}
				}
				if len(result_slice) > 0 {
					Changes = true
					return Changes, result_slice, old_slice, nil
				}
				return Changes, nil, nil, nil
			}
		} else {
			Changes = true
			return Changes, comparitive, Exact, nil
		}

	default:
		if !reflect.DeepEqual(comparitive, Exact) {
			Changes = true
			instant_result1, instant_result2 := GetDiffOutput(types, Exact, comparitive)
			return Changes, instant_result1, instant_result2, nil
		} else {
			return Changes, nil, nil, nil
		}
	}
}
func GetDiffJSON(types string, Exact interface{}, comparitive interface{}) (bool, interface{}, interface{}, error) {
	return GetDiffJSONValue(types, Exact, comparitive, false)
}

func GetDiffOutput(types string, old, new interface{}) (interface{}, interface{}) {
	switch types {
	case "bool":
		return true, nil
	default:
		return new, old
	}
}

func GetjsonDiffInBool(Exact, comparative interface{}) (bool, interface{}, interface{}, error) {
	return GetDiffJSON("bool", Exact, comparative)
}

func GetjsonDiffInValue(Exact, comparative interface{}) (bool, interface{}, interface{}, error) {
	return GetDiffJSON("value", Exact, comparative)
}
