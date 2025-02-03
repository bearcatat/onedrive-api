package http

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
)

type JsonResponse struct {
	response *http.Response
	body     []byte
	header   http.Header
}

func NewJsonResponse(response *http.Response) *JsonResponse {
	return &JsonResponse{
		response: response,
	}
}

func (r *JsonResponse) UnMarshal(target interface{}) error {
	r.header = r.response.Header
	err := r.readAndSetBody()
	if err != nil {
		return err
	}
	return r.unMarshal(target)
}

func (r *JsonResponse) readAndSetBody() error {
	defer r.response.Body.Close()
	var err error
	r.body, err = io.ReadAll(r.response.Body)
	if err != nil {
		return err
	}
	return nil
}

func (r *JsonResponse) unMarshal(target interface{}) error {
	if err := r.unmarshalError(); err != nil {
		return err
	}
	if err := r.unmarshalBodyToTarget(target); err != nil {
		return err
	}
	return r.unmarshalHeaderToTarget(target)

}

func (r *JsonResponse) unmarshalError() error {
	if len(r.body) == 0 {
		return nil
	}
	var errorResponse *ErrorResponse
	if err := json.Unmarshal(r.body, &errorResponse); err != nil {
		return err
	}
	return errorResponse.GetError()
}

func (r *JsonResponse) unmarshalBodyToTarget(target interface{}) error {
	if r.response.StatusCode == http.StatusNoContent {
		return nil
	}
	if len(r.body) == 0 {
		return nil
	}
	return json.Unmarshal(r.body, target)
}

func (r *JsonResponse) unmarshalHeaderToTarget(target interface{}) error {
	targetV, targetT, isNil := indirect(reflect.ValueOf(target))
	if isNil {
		return nil
	}
	for i := 0; i < targetT.NumField(); i++ {
		r.setTargetFieldByHeaderTag(targetT.Field(i), targetV.Field(i))
	}

	return nil
}

func (r *JsonResponse) setTargetFieldByHeaderTag(fieldT reflect.StructField, fieldV reflect.Value) {
	tag := fieldT.Tag.Get("header")
	headerValue := r.header.Get(tag)
	if headerValue == "" {
		return
	}
	if fieldT.Type.Kind() != reflect.String {
		return
	}
	fieldV.SetString(headerValue)
}

func indirect(v reflect.Value) (reflect.Value, reflect.Type, bool) {
	if v.Kind() == reflect.Invalid {
		return v, nil, true
	}
	if v.Kind() == reflect.Interface && !v.IsNil() {
		return indirect(v.Elem())
	}
	if v.Kind() != reflect.Pointer {
		return v, v.Type(), false
	}
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return indirect(v.Elem())
}
