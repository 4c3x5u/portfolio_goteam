//go:build utest

package main

import (
	"os"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
)

func TestEnv(t *testing.T) {
	errPostfix := " environment variable was empty"
	for _, c := range []struct {
		name       string
		envVarName string
		setup      func() error
		teardown   func() error
	}{
		{
			name:       "PORTEmpty",
			envVarName: "PORT",
			setup:      func() error { return nil },
			teardown:   func() error { return nil },
		},
		{
			name:       "DBCONNSTREmpty",
			envVarName: "DBCONNSTR",
			setup:      func() error { return os.Setenv("PORT", "8000") },
			teardown:   func() error { return os.Unsetenv("PORT") },
		},
		{
			name:       "JWTKEYEmpty",
			envVarName: "JWTKEY",
			setup: func() error {
				if err := os.Setenv("PORT", "8000"); err != nil {
					return err
				}
				if err := os.Setenv("DBCONNSTR", "some//connection.string"); err != nil {
					return err
				}
				return nil
			},
			teardown: func() error {
				if err := os.Unsetenv("PORT"); err != nil {
					return err
				}
				if err := os.Unsetenv("DBCONNSTR"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			name:       "CLIENTORIGINEmpty",
			envVarName: "CLIENTORIGIN",
			setup: func() error {
				if err := os.Setenv("PORT", "8000"); err != nil {
					return err
				}
				if err := os.Setenv("DBCONNSTR", "connection.string"); err != nil {
					return err
				}
				if err := os.Setenv("JWTKEY", "secretjwtsigningkey"); err != nil {
					return err
				}
				return nil
			},
			teardown: func() error {
				if err := os.Unsetenv("PORT"); err != nil {
					return err
				}
				if err := os.Unsetenv("DBCONNSTR"); err != nil {
					return err
				}
				if err := os.Unsetenv("JWTKEY"); err != nil {
					return err
				}
				return nil
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			wantErrMsg := c.envVarName + errPostfix
			if err := c.setup(); err != nil {
				t.Fatal(err)
			}

			if err := newEnv().validate(); err == nil {
				t.Error("validate() returned nil")
			} else {
				assert.Equal(t.Error, err.Error(), wantErrMsg)
			}

			if err := c.teardown(); err != nil {
				t.Fatal(err)
			}
		})
	}

	t.Run("Success", func(t *testing.T) {
		if err := os.Setenv("PORT", "8000"); err != nil {
			t.Fatal(err)
		}
		if err := os.Setenv("DBCONNSTR", "connection.string"); err != nil {
			t.Fatal(err)
		}
		if err := os.Setenv("JWTKEY", "secretjwtsigningkey"); err != nil {
			t.Fatal(err)
		}
		if err := os.Setenv("CLIENTORIGIN", "client:origin"); err != nil {
			t.Fatal(err)
		}

		if err := assert.Nil(newEnv().validate()); err != nil {
			t.Error(err)
		}

		if err := os.Unsetenv("PORT"); err != nil {
			t.Fatal(err)
		}
		if err := os.Unsetenv("DBCONNSTR"); err != nil {
			t.Fatal(err)
		}
		if err := os.Unsetenv("JWTKEY"); err != nil {
			t.Fatal(err)
		}
		if err := os.Unsetenv("CLIENTORIGIN"); err != nil {
			t.Fatal(err)
		}
	})
}
