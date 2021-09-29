package plugins

import (
	"github.com/mdlayher/wol"
	"net"
)

type WoL struct{}

type WoLCfg struct {
	Addr         string
	HardwareAddr net.HardwareAddr
	Password     *string
}

func (cfg *WoLCfg) Assert() error {
	return nil
}

func (w *WoL) Run(cfg Config) error {
	if err := cfg.Assert(); err != nil {
		return err
	}
	wolCfg, ok := cfg.(*WoLCfg)
	if !ok {
		return ErrConfigCastingFailed
	}

	client, err := wol.NewClient()
	if err != nil {
		return err
	}

	if wolCfg.Password != nil {
		if err := client.WakePassword(wolCfg.Addr, wolCfg.HardwareAddr, []byte(*(wolCfg.Password))); err != nil {
			return err
		}
	} else {
		if err := client.Wake(wolCfg.Addr, wolCfg.HardwareAddr); err != nil {
			return err
		}
	}

	if err := client.Close(); err != nil {
		return err
	}

	return nil
}
