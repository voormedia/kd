package kubectl

import (
	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/internal/kubectl"
	"github.com/voormedia/kd/pkg/util"
)

func Run(log *util.Logger, args ...string) error {
	conf, err := config.Load()
	if err != nil {
		return err
	}

	tgt, err := conf.ResolveTarget(args[0])
	if err != nil {
		return err
	}

	return kubectl.RunForTarget([]byte{}, tgt, args[1:]...)
}
