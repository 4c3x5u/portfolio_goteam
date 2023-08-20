//go:build utest

package main

import (
	"os"
	"server/assert"
	"testing"
)

func TestEnv(t *testing.T) {
	errPostfix := " environment variable was empty"
	for _, c := range []struct {
		Name       string
		EnvVarName string
		Setup      func() error
		Cleanup    func() error
	}{
		{
			Name:       "PORTEmpty",
			EnvVarName: "PORT",
			Setup:      func() error { return nil },
			Cleanup:    func() error { return nil },
		},
		{
			Name:       "DBCONNSTREmpty",
			EnvVarName: "DBCONNSTR",
			Setup:      func() error { return os.Setenv("PORT", "8000") },
			Cleanup:    func() error { return os.Unsetenv("PORT") },
		},
		{
			Name:       "JWTKEYEmpty",
			EnvVarName: "JWTKEY",
			Setup: func() error {
				if err := os.Setenv("PORT", "8000"); err != nil {
					return err
				}
				if err := os.Setenv("DBCONNSTR", "some//connection.string"); err != nil {
					return err
				}
				return nil
			},
			Cleanup: func() error {
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
			Name:       "CLIENTORIGINEmpty",
			EnvVarName: "CLIENTORIGIN",
			Setup: func() error {
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
			Cleanup: func() error {
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
		t.Run(c.Name, func(t *testing.T) {
			wantErrMsg := c.EnvVarName + errPostfix
			if err := c.Setup(); err != nil {
				t.Fatal(err)
			}

			if err := newEnv().validate(); err == nil {
				t.Error("validate() returned nil")
			} else if err = assert.Equal(wantErrMsg, err.Error()); err != nil {
				t.Error(err)
			}

			if err := c.Cleanup(); err != nil {
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
