package preset

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/spf13/afero"
	"github.com/spinkube/runtime-class-manager/internal/containerd"
)

type ContainerdConfig struct {
	DefaultPath string
	CustomPath  string
}

var (
	ContainerdConfigMicroK8s = ContainerdConfig{
		"/var/snap/microk8s/current/args/containerd-template.toml",
		"/var/snap/microk8s/current/args/containerd-template.toml",
	}
	ContainerdConfigRKE2 = ContainerdConfig{
		"/var/lib/rancher/rke2/agent/etc/containerd/config.toml",
		"/var/lib/rancher/rke2/agent/etc/containerd/config.toml.tmpl",
	}
	ContainerdConfigK3S = ContainerdConfig{
		"/var/lib/rancher/k3s/agent/etc/containerd/config.toml",
		"/var/lib/rancher/k3s/agent/etc/containerd/config.toml.tmpl",
	}
	ContainerdConfigK0S = ContainerdConfig{
		"/etc/k0s/containerd.toml",
		"/etc/k0s/containerd.d/config.toml",
	}
	ContainerdConfigDefault = ContainerdConfig{
		"/etc/containerd/config.toml",
		"/etc/containerd/config.toml",
	}
)

var SettingsMap = map[string]Settings{
	ContainerdConfigMicroK8s.DefaultPath: MicroK8s,
	ContainerdConfigRKE2.DefaultPath:     RKE2,
	ContainerdConfigK3S.DefaultPath:      K3s,
	ContainerdConfigK0S.DefaultPath:      K0s,
	ContainerdConfigDefault.DefaultPath:  Default,
}

type Settings struct {
	ConfigPath string
	Setup      func(Env) error
	Restarter  containerd.Restarter
}

type Env struct {
	HostFs     afero.Fs
	ConfigPath string
}

var Default = Settings{
	ConfigPath: ContainerdConfigDefault.CustomPath,
	Setup:      func(_ Env) error { return nil },
	Restarter:  containerd.NewRestarter(),
}

func (s Settings) WithConfigPath(path string) Settings {
	s.ConfigPath = path
	return s
}

func (s Settings) WithSetup(setup func(env Env) error) Settings {
	s.Setup = setup
	return s
}

var MicroK8s = Default.WithConfigPath(ContainerdConfigMicroK8s.CustomPath)

var RKE2 = Default.WithConfigPath(ContainerdConfigRKE2.CustomPath).
	WithSetup(func(env Env) error {
		_, err := env.HostFs.Stat(env.ConfigPath)
		if err == nil {
			return nil
		}

		if errors.Is(err, os.ErrNotExist) {
			// Copy base config into .tmpl version
			src, _ := strings.CutSuffix(env.ConfigPath, ".tmpl")
			in, err := env.HostFs.Open(src)
			if err != nil {
				return err
			}
			defer in.Close()
			out, err := env.HostFs.Create(env.ConfigPath)
			if err != nil {
				return err
			}
			defer func() {
				cerr := out.Close()
				if err == nil {
					err = cerr
				}
			}()
			if _, err = io.Copy(out, in); err != nil {
				return err
			}
			err = out.Sync()

			return nil
		}

		return err
	})

var K3s = RKE2.WithConfigPath(ContainerdConfigK3S.CustomPath)

var K0s = Default.WithConfigPath(ContainerdConfigK0S.CustomPath).
	WithSetup(func(env Env) error {
		_, err := env.HostFs.Stat(env.ConfigPath)
		if err == nil {
			return nil
		}

		if errors.Is(err, os.ErrNotExist) {
			_, err := env.HostFs.Create(env.ConfigPath)
			if err != nil {
				return err
			}
			return nil
		}

		return err
	})
