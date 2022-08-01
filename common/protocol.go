package common

import "encoding/json"

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type Response struct {

	Code string `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func BuildResp(code,msg string,data interface{})([]byte,error){
	resp:=&Response{
		Code: code,
		Msg: msg,
		Data: data,
	}
	return json.Marshal(resp)
}

// 反序列化Job
func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job *Job
	)

	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	ret = job
	return
}
