package user

import (
	"testing"

	"github.com/pkg/errors"
)

type TestDataHash struct {
	Input    string
	Output   string
	Expected string
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
}

func TestHashPass(t *testing.T) {
	data := []TestDataHash{
		{
			Input:    "mango",
			Expected: "aa00faf97d042c13a59da4d27eb32358",
		},
		{
			Input:    "Beer322",
			Expected: "14d45bbe24703ac76c092f2d57408d3e",
		},
		{
			Input:    "Bear",
			Expected: "372137ebb0d053fecd7a594ec5cb5971",
		},
	}
	for _, val := range data {
		var err error
		val.Output, err = HashPass(val.Input)
		checkError(err, t)
		if val.Expected != val.Output {
			t.Errorf("Hash passwords are different. Expected %s .\n Got %s instead", val.Expected, val.Output)
		}
	}
}

type TestDataValid struct {
	Input    User
	Output   error
	Expected error
}

func TestCheckValidUser(t *testing.T) {
	users := []TestDataValid{
		{
			Input:    User{Login: "", Password: "Mongo"},
			Expected: errors.New("enter login"),
		},
		{
			Input:    User{Login: "Mango", Password: ""},
			Expected: errors.New("enter password"),
		},
	}

	for _, val := range users {
		val.Output = CheckValidUser(&val.Input)

		if val.Expected.Error() != val.Output.Error() {
			t.Errorf("Errors are different. Expected %s.\n Got %s instead", val.Expected.Error(), val.Output.Error())
		}
	}
}
