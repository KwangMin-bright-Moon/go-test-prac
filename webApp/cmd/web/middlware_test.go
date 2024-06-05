package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var app application

func Test_application_addIPToContext(t *testing.T){
	tests := []struct{
		headerName string
		headerValue string
		addr string
		emptyAddr bool
	}{
		{"", "", "", false},
		{"", "", "", true},
		{"X-Forwarded-For", "192.3.2.1", "", false},
		{"", "", "hello:world", false},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		val := r.Context().Value(CONTEXT_USER_KEY)
		if val == nil {
			t.Error(CONTEXT_USER_KEY, "not present")
		}

		ip, ok := val.(string)
		if !ok {
			t.Error("not string")
		}
		t.Log("IP:", ip)
	})

	for _, e := range tests {
		handlerToTest := app.addIPToContext(nextHandler)

		req := httptest.NewRequest("GET", "http://testing", nil)

		if e.emptyAddr{
			req.RemoteAddr = ""	
		}

		if len(e.headerName) > 0 {
			req.Header.Add(e.headerName, e.headerValue)
		}

		if len(e.addr) > 0 {
			req.RemoteAddr = e.addr
		}

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}
}

func Test_application_ipFromContext(t *testing.T){
	var app application

	 ctx := context.Background()

	 ctx = context.WithValue(ctx, CONTEXT_USER_KEY, "whatever")

	 ip := app.ipFromContext(ctx)

	 if !strings.EqualFold("whatever", ip){
	 	t.Error("wrone value returned from context")
	 }
}