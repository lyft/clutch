package gateway

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	type File struct {
		name  string
		value string
	}
	var testCases = []struct {
		files         []File
		envVar        string
		expectedValue string
	}{
		{
			files: []File{
				{
					name:  ".env.dev",
					value: "FOOBAR1=true",
				},
				{
					name:  ".env",
					value: "FOOBAR1=false",
				},
			},
			envVar:        "FOOBAR1",
			expectedValue: "true",
		},
		{
			files: []File{
				{
					name:  ".env.dev",
					value: "",
				},
				{
					name:  ".env",
					value: "FOOBAR2=true",
				},
			},
			envVar:        "FOOBAR2",
			expectedValue: "true",
		},
	}
	for idx, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			t.Parallel()
			tmpVar := os.Getenv(tc.envVar)
			fileNames := []string{}
			for _, f := range tc.files {
				envFile, err := ioutil.TempFile(".", f.name)
				if err != nil {
					log.Fatal(err)
				}
				defer os.Remove(envFile.Name())
				ioutil.WriteFile(envFile.Name(), []byte(f.value), 0644)
				fileNames = append(fileNames, envFile.Name())
			}

			flags := &Flags{EnvFiles: fileNames}

			loadEnv(flags)
			v := os.Getenv(tc.envVar)
			assert.Equal(t, tc.expectedValue, v)
			os.Setenv(tc.envVar, tmpVar)
		})
	}

}
