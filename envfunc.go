// Package envfunc defines the environment function.
package envfunc

import (
	"fmt"
	"os"
	"strconv"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type boolEnv struct {
	name, desc string
	d, v       bool
}

type intEnv struct {
	name, desc string
	d, v       int
}

type stringEnv struct {
	name, desc, d, v string
}

var (
	boolEnvs   = make(map[string]boolEnv, 32)
	stringEnvs = make(map[string]stringEnv, 32)
	intEnvs    = make(map[string]intEnv, 32)
)

// EnvBoolFunc returns a function that read bool type from registered env.
var EnvBoolFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "name",
			Type:             cty.String,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(cty.Bool),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		e := args[0].AsString()

		s, ok := boolEnvs[e]
		if !ok {
			return cty.Value{}, fmt.Errorf("env: %s does not register", e)
		}

		v := os.Getenv(e)
		if v == "" {
			s.v = s.d
			return cty.BoolVal(s.d), nil
		}

		var err error
		s.v, err = strconv.ParseBool(v)
		if err != nil {
			s.v = s.d
		}

		return cty.BoolVal(s.v), nil
	},
})

// EnvIntFunc returns a function that read int type from registered env.
var EnvIntFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "name",
			Type:             cty.String,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(cty.Number),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		e := args[0].AsString()

		s, ok := intEnvs[e]
		if !ok {
			return cty.Value{}, fmt.Errorf("env: %s does not register", e)
		}

		v := os.Getenv(e)
		if v == "" {
			s.v = s.d
			return cty.NumberIntVal(int64(s.d)), nil
		}

		i, err := strconv.ParseInt(v, 10, 64)
		s.v = int(i)
		return cty.NumberIntVal(i), err
	},
})

// EnvStringFunc returns a function that read string type from registered env.
var EnvStringFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "name",
			Type:             cty.String,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		e := args[0].AsString()

		s, ok := stringEnvs[e]
		if !ok {
			return cty.Value{}, fmt.Errorf("env: %s does not register", e)
		}

		v := os.Getenv(e)
		if v == "" {
			v = s.d
		}
		s.v = v

		return cty.StringVal(v), nil
	},
})

// RegisterIntEnv registers a int to env.
func RegisterIntEnv(name, desc string, d int) {
	intEnvs[name] = intEnv{
		name: name,
		desc: desc,
		d:    d,
	}
}

// RangeIntEnv iterates over the int environment.
func RangeIntEnv(fn func(name, desc string, d int)) {
	for _, v := range intEnvs {
		fn(v.name, v.desc, v.d)
	}
}

// RegisterStringEnv registers a string to env.
func RegisterStringEnv(name, desc, d string) {
	stringEnvs[name] = stringEnv{
		name: name,
		desc: desc,
		d:    d,
	}
}

// RangeStringEnv iterates over the string environment.
func RangeStringEnv(fn func(name, desc, d string)) {
	for _, v := range stringEnvs {
		fn(v.name, v.desc, v.d)
	}
}

// RegisterBoolEnv registers a bool to env.
func RegisterBoolEnv(name, desc string, d bool) {
	boolEnvs[name] = boolEnv{
		name: name,
		desc: desc,
		d:    d,
	}
}

// RangeBoolEnv iterates over the bool environment.
func RangeBoolEnv(fn func(name, desc string, d bool)) {
	for _, v := range boolEnvs {
		fn(v.name, v.desc, v.d)
	}
}
