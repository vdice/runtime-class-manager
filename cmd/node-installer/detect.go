/*
   Copyright The SpinKube Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/spf13/afero"
	"github.com/spinkube/runtime-class-manager/internal/preset"
)

// Ordered set of config path locations, leaving the default for last.
// There may be cases where the default config *and* another distro-specific
// config exist; in these cases we want to detect the distro-specific config first.
var containerdConfigLocations = [5]preset.ContainerdConfig{
	preset.ContainerdConfigMicroK8s,
	preset.ContainerdConfigRKE2,
	preset.ContainerdConfigK3S,
	preset.ContainerdConfigK0S,
	preset.ContainerdConfigDefault,
}

func DetectDistro(config Config, hostFs afero.Fs) (preset.Settings, error) {
	if config.Runtime.ConfigPath != "" {
		// containerd config path has been set explicitly
		if distro, ok := preset.SettingsMap[config.Runtime.ConfigPath]; ok {
			return distro, nil
		}
		slog.Warn("could not determine distro from containerd config, falling back to defaults", "config", config.Runtime.ConfigPath)
		return preset.Default.WithConfigPath(config.Runtime.ConfigPath), nil
	}

	var errs []error

	for _, containerdConfig := range containerdConfigLocations {
		_, err := hostFs.Stat(containerdConfig.DefaultPath)
		if err == nil {
			// config file found, return corresponding distro settings
			return preset.SettingsMap[containerdConfig.DefaultPath], nil
		}
		errs = append(errs, err)
	}

	return preset.Settings{}, fmt.Errorf("failed to detect containerd config path: %w", errors.Join(errs...))
}
