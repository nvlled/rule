
package rule

import (
	"testing"
	"fmt"
)

type account struct  {
	username string
	password string
}

type Fn func(x int) int

func cast(t interface{}, fn func(account) Error) Error {
	switch t := t.(type) {
		case account: return fn(t)
	}
	return AnError("message" , "invalid type, not an account")
}

func usernameIsDumb(t interface{}) Error {
	return cast(t, func(acc account) Error {
		if acc.username == "dumb" {
			return AnError("username", "username is dumb")
		}
		return nil
	})
}

func passwordIsDumb(t interface{}) Error {
	return cast(t, func(acc account) Error {
		if acc.password == "dumb" {
			return AnError("password", "password is dumb")
		}
		return nil
	})
}

func passwordIsShort(t interface{}) Error {
	return cast(t, func(acc account) Error {
		if len(acc.password) < 5 {
			return AnError("password", "password is too short")
		}
		return nil
	})
}

func TestAll(t *testing.T) {
	acc1 := account{"dumb", "dumb"}
	acc2 := account{"nope", "dumb"}
	acc3 := account{"nope", "password"}

	rule := All(T(usernameIsDumb), T(passwordIsDumb), T(passwordIsShort))

	err := rule(acc1)
	if err["username"] == nil {
		t.Error("username should fail")
	}
	if err["password"] == nil {
		t.Error("password should fail")
	}
	if !contains(err["password"], "password is dumb") ||
	!contains(err["password"], "password is too short"){
		t.Error("password is short and dumb")
	}

	err = rule(acc2)
	if err["username"] != nil {
		t.Error("username should pass")
	}
	if err["password"] == nil {
		t.Error("password should fail")
	}

	err = rule(acc3)
	if Failed(err) {
		t.Error("account should pass")
	}
}

func TestOne(t *testing.T) {
	acc1 := account{"dumb", "dumb"}
	acc2 := account{"nope", "dumb"}
	acc3 := account{"nope", "password"}

	rule := One(T(usernameIsDumb), T(passwordIsDumb))

	err := rule(acc1)
	if !contains(err["username"], "username is dumb") {
		t.Error("username should be dumb")
	}
	if err["password"] != nil {
		t.Error("One rule should only be tested")
	}

	err = rule(acc2)
	if err["username"] != nil {
		t.Error("One rule should only be tested")
	}
	if !contains(err["password"], "password is dumb") {
		t.Error("password should be dumb")
	}

	err = rule(acc3)
	if Failed(err) {
		t.Error("Rule should pass")
	}

	fmt.Printf("%v %v %v\n", acc1, acc2, rule)
}

func contains(msgs Msgs, s string) bool {
	for _,m := range msgs {
		if m == s { return true }
	}
	return false
}




