package webapi

import "fmt"

type Options func(a *accrual)

func Address(adr string) Options {
	return func(a *accrual) {
		a.url = fmt.Sprintf("http://%s", adr)
	}
}
