package myhttp

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"encoding/json"
	"sync"
	"time"
	"math"
	"strconv"
)

var CacheMapWithTime sync.Map
var CacheMap sync.Map

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()  //解析参数，默认是不会解析的
	fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("method", r.Method)
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	Map :=make( map[int] string)
	Map[1]="ok"
	jMap,err:=json.Marshal(Map)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Fprintf(w, string(jMap)) //这个写入到w的是输出到客户端的
}

func getValue(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	key :=  strings.Join(r.Form["key"], "")
	fmt.Println("key"+key)
	if key==""{
		fmt.Fprintf(w, "key为空ss")
		return
	}
	now := time.Now().Unix()
	v,ok:=CacheMapWithTime.Load(key)
	fmt.Println(v)
	fmt.Println(ok)
	if ok{
		var v_time int64=v.(int64)
		if v_time-now>0 {
			v2,_:=CacheMap.Load(key)
			var v2_int int64 =v2.(int64)
			fmt.Fprintf(w,"value:"+ strconv.FormatInt(v2_int,10))
		}
		}



}

func getkeyValue(key string) int64{
	now := time.Now().Unix()
	v,ok:=CacheMapWithTime.Load(key)
	fmt.Println(v)
	fmt.Println(ok)
	if ok{
		var v_time int64=v.(int64)
		if v_time-now>0 {
			v2,_:=CacheMap.Load(key)
			var v2_int int64 =v2.(int64)
			return v2_int
		}
	}
	return math.MinInt64
}


var isStart bool=false
func setKeyValue(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	fmt.Println("key:"+key+",value:"+value)
	if key==""||value==""{
		fmt.Fprintf(w, "key或value为空")
		return
	}
	expire := r.Form.Get("expire")
	fmt.Println("expire:"+expire)
	if expire!=""{
		_expire, _ := strconv.ParseInt(expire, 10, 64)
		_expire+=time.Now().Unix()
		CacheMapWithTime.Store(key,_expire)
	}else {
		max:=strconv.FormatInt(math.MaxInt64,10)
		_max,_:=strconv.ParseInt(max, 10, 64)
		CacheMapWithTime.Store(key,_max)
	}
	_value, _ := strconv.ParseInt(value, 10, 64)
	CacheMap.Store(key,_value)
	fmt.Println("set success!")
	if isStart {
		
	}
	fmt.Fprintf(w,"set success!")
}

func clearUpThread()  {
	for  {
		time.Sleep(60*time.Second)
		fmt.Println("sleep 60 second")
		CacheMapWithTime.Range( func(key, value interface{}) bool{
			v:=value.(int64)
			now:=time.Now().Unix()
			if v<now {
				CacheMapWithTime.Delete(key)
				CacheMap.Delete(key)
				fmt.Println("delete overtime..."+key.(string))
			}
			return true
		})
	}

}

func init() {
	go clearUpThread()
}

func statistics(w http.ResponseWriter, r *http.Request)  {
     data:=make(map[string]interface{})
	expireData :=make(map[string]int64)
	validData :=make(map[string]int64)
	CacheMapWithTime.Range(func(k, v interface{}) bool {
		expireTime := v.(int64)
		if  expireTime!=math.MaxInt64 && time.Now().Unix()<expireTime{
			expireData[k.(string)]= expireTime
			validData[k.(string)]= expireTime
		}
		if expireTime==math.MaxInt64{
			validData[k.(string)]= expireTime
		}
		return true
	})

	storeData :=make(map[string]int64)

	for k,_:=range validData{
		value, ok := CacheMap.Load(k)
		if ok {
			storeData[k]=value.(int64)
		}
	}
	/*CacheMap.Range(func(k, v interface{}) bool {
		storeData[k.(string)]=v.(int64)
		return true
	})*/
	data["expireCt"]=len(expireData)
	data["expireMap"]=expireData

	data["storeCt"]=len(storeData)
	data["storeMap"]=storeData
	jMap,err:=json.Marshal(data)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Fprintf(w,string(jMap))

}


func getAll(w http.ResponseWriter, r *http.Request)  {
	data:=make(map[string]interface{})
	expireData :=make(map[string]int64)
	validData :=make(map[string]int64)
	CacheMapWithTime.Range(func(k, v interface{}) bool {
		expireTime := v.(int64)
        validData[k.(string)]= expireTime
		expireData[k.(string)]= expireTime
		return true
	})

	storeData :=make(map[string]int64)

	for k,_:=range validData{
		value, ok := CacheMap.Load(k)
		if ok {
			storeData[k]=value.(int64)
		}
	}
	data["expireCt"]=len(expireData)
	data["expireMap"]=expireData

	data["storeCt"]=len(storeData)
	data["storeMap"]=storeData
	jMap,err:=json.Marshal(data)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Fprintf(w,string(jMap))

}

func delete(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	key := r.Form.Get("key")
	CacheMapWithTime.Delete(key)
	CacheMap.Delete(key)
	fmt.Fprintf(w,"delete success!")

}
func increase(w http.ResponseWriter, r *http.Request)  {

	r.ParseForm()
	key :=  strings.Join(r.Form["key"], "")
	fmt.Println("key"+key)
	if key==""{
		fmt.Fprintf(w, "key为空ss")
		return
	}

	value := getkeyValue(key)
	if value!=math.MinInt64 {
        value++
        CacheMap.Store(key,value)
	}
	fmt.Fprintf(w,"increase success!")
}


func StartServer() {
//	http.HandleFunc("/", sayhelloName) //设置访问的路由
	http.HandleFunc("/getValue", getValue) //设置访问的路由
	http.HandleFunc("/setKeyValue", setKeyValue) //设置访问的路由
	http.HandleFunc("/statistics", statistics) //设置访问的路由
	http.HandleFunc("/delete", delete) //设置访问的路由
	http.HandleFunc("/increase", increase) //设置访问的路由
	http.HandleFunc("/getAll", getAll) //设置访问的路由

	log.Println("start server at port 8080")
	err := http.ListenAndServe(":8080", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
