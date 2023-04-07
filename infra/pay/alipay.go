package pay

import (
	"github.com/KAMIENDER/golang-scaffold/infra/config"
	"github.com/pkg/errors"
	"github.com/smartwalle/alipay/v3"
)

type AliPay struct {
	client *alipay.Client
}

func NewAliPay(conf *config.Config) (*AliPay, error) {
	client, err := alipay.New(conf.PayConf.AppID, conf.PayConf.AppPrivateKey, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	client.LoadAliPayPublicKey(conf.PayConf.AliPublicKey)
	return &AliPay{
		client: client,
	}, nil
}
