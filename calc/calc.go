package calc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	cacheresult "github.com/patrickmn/go-cache"
)

var Cache *cacheresult.Cache

func init() {
	Cache = cacheresult.New(1*time.Minute, 60*time.Second)
}

type Result struct {
	Action    string `json:"action"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Answer    int    `json:"answer"`
	Cached    bool   `json:"cached"`
	Error_msg string `json:"error_msg"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()

	operation := mux.Vars(r)["operation"]

	x, e1 := strconv.Atoi(v.Get("x"))

	y, e2 := strconv.Atoi(v.Get("y"))

	if e1 != nil || e2 != nil {
		var err string
		if e1 != nil {
			err = e1.Error()
		} else {
			err = e2.Error()
		}
		encode(w, Result{operation, -1, -1, -1, false, err})
		return
	}
	var answer, cached, error = compute(x, y, operation)
	var err_str string
	if error != nil {
		err_str = *error
	} else {
		err_str = ""
	}
	encode(w, Result{operation, x, y, answer, cached, err_str})
}

func compute(left int, right int, operation string) (int, bool, *string) {
	key := fmt.Sprintf("%s%d%d", operation, left, right)
	result, exist := Cache.Get(key)
	if exist {
		return result.(int), true, nil
	}
	calculated, error := doOperation(left, right, operation)
	if error != nil {
		return -1, false, error
	}
	Cache.Set(key, calculated, cacheresult.DefaultExpiration)
	return calculated, false, nil
}

func doOperation(left int, right int, operation string) (int, *string) {
	switch operation {
	case "add":
		return left + right, nil
	case "subtract":
		return left - right, nil
	case "multiply":
		return left + right, nil
	case "divide":
		if right == 0 {
			err := "Cannot divide by zero"
			return -1, &err
		}
		return left / right, nil
	default:
		err := "Operation not recognized"
		return -1, &err
	}
}

func encode(w http.ResponseWriter, x interface{}) {
	json.NewEncoder(w).Encode(x)
}
