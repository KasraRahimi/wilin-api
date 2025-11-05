package services_test

import (
	"os"
	"slices"
	"testing"

	"wilin.info/api/server/services"
)

type OriginsValue struct {
	env      string
	expected []string
}

func TestNoEnv(t *testing.T) {
	origins := services.GetOrigins()
	if !slices.Equal(origins, services.DEFAULT_ORIGINS) {
		failTest(t, origins, services.DEFAULT_ORIGINS)
	}
}

var originsValues = []OriginsValue{
	{"", services.DEFAULT_ORIGINS},
	{"foo,bar", []string{"foo", "bar"}},
	{"foo", []string{"foo"}},
	{"foo, bar", []string{"foo", "bar"}},
	{"foo bar", []string{"foo bar"}},
	{"three,origin,elements", []string{"three", "origin", "elements"}},
	{"this, should, work", []string{"this", "should", "work"}},
	{",", []string{"", ""}},
	{"what,", []string{"what", ""}},
	{"http://localhost:3000", []string{"http://localhost:3000"}},
	{",,,", []string{"", "", "", ""}},
}

func TestDifferentEnvValues(t *testing.T) {
	for _, test := range originsValues {
		os.Setenv("ORIGINS", test.env)

		services.SetOrigins()
		origins := services.GetOrigins()
		if !slices.Equal(origins, test.expected) {
			failTest(t, origins, test.expected)
		}
	}
}
