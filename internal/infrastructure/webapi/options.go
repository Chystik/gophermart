package webapi

import "fmt"

type Options func(a *accrual)

func Address(adr string, scheme string) Options {
	return func(a *accrual) {
		a.url = fmt.Sprintf("%s://%s", scheme, adr)
	}
}
