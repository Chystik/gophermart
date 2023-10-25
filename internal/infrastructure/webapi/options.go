package webapi

import "fmt"

type Options func(a *accrual)

func Address(adr string) Options {
	return func(a *accrual) {
		a.url = fmt.Sprintf("%s%s", a.url, adr)
	}
}

func Scheme(sch string) Options {
	return func(a *accrual) {
		a.url = fmt.Sprintf("%s://%s", sch, a.url)
	}
}
