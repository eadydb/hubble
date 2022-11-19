package cmd

import (
	"context"
	"fmt"
	"reflect"

	"github.com/eadydb/hubble/pkg/output/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Nillable is used to reset objects that implement pflag's `Value` and `SliceValue`.
// Some flags, like `--default-repo`, use nil to indicate that they are unset, which
// is different from the empty string.
type Nillable interface {
	SetNil() error
}

// Flag defines a Skaffold CLI flag which contains a list of
// subcommands the flag belongs to in `DefinedOn` field.
// See https://pkg.go.dev/github.com/spf13/pflag#Flag
type Flag struct {
	Name                 string
	Shorthand            string
	Usage                string
	Value                interface{}
	DefValue             interface{}
	DefValuePerCommand   map[string]interface{}
	DeprecatedPerCommand map[string]interface{}
	NoOptDefVal          string
	FlagAddMethod        string
	Deprecated           string
	DefinedOn            []string
	Hidden               bool
	IsEnum               bool
}

// flagRegistry is a list of all Skaffold CLI flags.
// When adding a new flag to the registry, please specify the
// command/commands to which the flag belongs in `DefinedOn` field.
// If the flag is a global flag, or belongs to all the subcommands,
// specify "all"
// FlagAddMethod is method which defines a flag value with specified
// name, default value, and usage string. e.g. `StringVar`, `BoolVar`
var flagRegistry = []Flag{}

func methodNameByType(v reflect.Value) string {
	t := v.Type().Kind()
	switch t {
	case reflect.Bool:
		return "BoolVar"
	case reflect.Int:
		return "IntVar"
	case reflect.String:
		return "StringVar"
	case reflect.Slice:
		return "StringSliceVar"
	case reflect.Struct:
		return "Var"
	case reflect.Ptr:
		return methodNameByType(reflect.Indirect(v))
	}
	return ""
}

func (fl *Flag) flag(cmdName string) *pflag.Flag {
	methodName := fl.FlagAddMethod
	if methodName == "" {
		methodName = methodNameByType(reflect.ValueOf(fl.Value))
	}
	isVar := methodName == "Var"
	// pflags' Var*() methods do not take a default value but instead
	// assume the value is already set to its default value.  So we
	// explicitly set the default value here to ensure help text is correct.
	if isVar {
		setDefaultValues(fl.Value, fl, cmdName)
	}

	inputs := []interface{}{fl.Value, fl.Name}
	if !isVar {
		if d, found := fl.DefValuePerCommand[cmdName]; found {
			inputs = append(inputs, d)
		} else {
			inputs = append(inputs, fl.DefValue)
		}
	}
	inputs = append(inputs, fl.Usage)

	fs := pflag.NewFlagSet(fl.Name, pflag.ContinueOnError)
	reflect.ValueOf(fs).MethodByName(methodName).Call(reflectValueOf(inputs))

	f := fs.Lookup(fl.Name)
	if len(fl.NoOptDefVal) > 0 {
		// f.NoOptDefVal may be set depending on value type
		f.NoOptDefVal = fl.NoOptDefVal
	}
	f.Shorthand = fl.Shorthand
	f.Hidden = fl.Hidden || (fl.Deprecated != "")
	f.Deprecated = fl.Deprecated

	// Deprecations can be applied per command
	if _, found := fl.DeprecatedPerCommand[cmdName]; found {
		f.Deprecated = fl.Deprecated
	}
	return f
}

func ResetFlagDefaults(cmd *cobra.Command, flags []*Flag) {
	// Update default values.
	for _, fl := range flags {
		flag := cmd.Flag(fl.Name)
		if !flag.Changed {
			setDefaultValues(flag.Value, fl, cmd.Name())
		}
		if fl.IsEnum {
			// instrumentation.AddFlag(flag)
		}
	}
}

// setDefaultValues sets the default value (or values) for the given flag definition.
// This function handles pflag's SliceValue and Value interfaces.
func setDefaultValues(v interface{}, fl *Flag, cmdName string) {
	d, found := fl.DefValuePerCommand[cmdName]
	if !found {
		d = fl.DefValue
	}
	if nv, ok := v.(Nillable); ok && d == nil {
		nv.SetNil()
	} else if sv, ok := v.(pflag.SliceValue); ok {
		sv.Replace(asStringSlice(d))
	} else if val, ok := v.(pflag.Value); ok {
		val.Set(fmt.Sprintf("%v", d))
	} else {
		log.Entry(context.TODO()).Fatalf("%s --%s: unhandled value type: %v (%T)", cmdName, fl.Name, v, v)
	}
}

// AddFlags adds to the command the common flags that are annotated with the command name.
func AddFlags(cmd *cobra.Command) {
	var flagsForCommand []*Flag

	for i := range flagRegistry {
		fl := &flagRegistry[i]
		if !hasCmdAnnotation(cmd.Use, fl.DefinedOn) {
			continue
		}

		cmd.Flags().AddFlag(fl.flag(cmd.Use))

		flagsForCommand = append(flagsForCommand, fl)
	}

	// Apply command-specific default values to flags.
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		ResetFlagDefaults(cmd, flagsForCommand)
		// Since PersistentPreRunE replaces the parent's PersistentPreRunE,
		// make sure we call it, if it is set.
		if parent := cmd.Parent(); parent != nil {
			if preRun := parent.PersistentPreRunE; preRun != nil {
				if err := preRun(cmd, args); err != nil {
					return err
				}
			} else if preRun := parent.PersistentPreRun; preRun != nil {
				preRun(cmd, args)
			}
		}

		return nil
	}
}

func hasCmdAnnotation(cmdName string, annotations []string) bool {
	for _, a := range annotations {
		if cmdName == a || a == "all" {
			return true
		}
	}
	return false
}

func reflectValueOf(values []interface{}) []reflect.Value {
	var results []reflect.Value
	for _, v := range values {
		results = append(results, reflect.ValueOf(v))
	}
	return results
}

func asStringSlice(v interface{}) []string {
	vt := reflect.TypeOf(v)
	if vt == reflect.TypeOf([]string{}) {
		return v.([]string)
	}
	switch vt.Kind() {
	case reflect.Array, reflect.Slice:
		value := reflect.ValueOf(v)
		var slice []string
		for i := 0; i < value.Len(); i++ {
			slice = append(slice, fmt.Sprintf("%v", value.Index(i)))
		}
		return slice
	default:
		return []string{fmt.Sprintf("%v", v)}
	}
}
