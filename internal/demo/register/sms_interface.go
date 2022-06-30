package register

type SMS interface {
	SendSms(code int, phoneNumber string) (resp interface{}, err error)
}
