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

	for k, v := range d {
		switch vv := v.(type) {
		case string:
			fmt.Println(fmt.Sprintf("%s:%s", k, v))
		case int:
			fmt.Println(k, "is int", vv)
		case float64:
			fmt.Println(k, "is float", vv)
		case map[string]interface{}:
		case []interface{}:
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle", vv)
		}
	}

}
