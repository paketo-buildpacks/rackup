package fakes

import "sync"

type GemParser struct {
	ParseCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Path string
		}
		Returns struct {
			RackFound bool
			Err       error
		}
		Stub func(string) (bool, error)
	}
}

func (f *GemParser) Parse(param1 string) (bool, error) {
	f.ParseCall.Lock()
	defer f.ParseCall.Unlock()
	f.ParseCall.CallCount++
	f.ParseCall.Receives.Path = param1
	if f.ParseCall.Stub != nil {
		return f.ParseCall.Stub(param1)
	}
	return f.ParseCall.Returns.RackFound, f.ParseCall.Returns.Err
}
