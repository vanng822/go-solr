package solr

import "testing"
import "fmt"

func TestConnection(t *testing.T) {
	fmt.Println("Running TestConnection")
	/*body,_ := HTTPGet("http://igeonote.com/api/geoip/country/66.249.66.20")

	res,_ := bytes2Json(&body)
	fmt.Println(fmt.Sprintf("%s", *res))
	*/
}

func TestBytes2Json(t *testing.T) {
	data := []byte(`{"t":"s","two":2,"obj":{"c":"b","j":"F"},"a":[1,2,3]}`)
	d, _ := bytes2Json(&data)
	PrintMapInterface(d)
}

func PrintMapInterface(d map[string]interface{}) {
	for k, v := range d {
		switch vv := v.(type) {
		case string:
			fmt.Println(fmt.Sprintf("%s:%s", k, v))
		case int:
			fmt.Println(k, "is int", vv)
		case float64:
			fmt.Println(k, "is float", vv)
		case map[string]interface{}:
			fmt.Println(k, "type is map[string]interface{}")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		case []interface{}:
			fmt.Println(k, "type is []interface{}")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle", vv)
		}
	}
}

func TestJson2Bytes(t *testing.T) {

	test_json := map[string]interface{}{
		"t":   "s",
		"two": 2,
		"obj": map[string]interface{}{"c": "b", "j": "F"},
		"a":   []interface{}{1, 2, 3},
	}

	b, err := json2Bytes(test_json)
	if err != nil {
		fmt.Println(err)
	}
	d, _ := bytes2Json(b)

	PrintMapInterface(d)
}
