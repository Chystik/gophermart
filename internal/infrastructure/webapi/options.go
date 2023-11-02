package webapi

type Options func(a *accrual)

func Address(adr string) Options {
	return func(a *accrual) {
		a.url = adr
	}
}
