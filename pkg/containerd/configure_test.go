/*
   Copyright The KWasm Authors.

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

package containerd

import (
	"testing"

	"github.com/kwasm/kwasm-node-installer/pkg/config"
	"github.com/spf13/afero"
)

func newGlobalConfig(configFile string) *config.Config {
	return &config.Config{
		Runtime: struct {
			Name       string
			ConfigPath string
		}{Name: "containerd", ConfigPath: configFile},
	}
}

func TestConfig_AddRuntime(t *testing.T) {
	type args struct {
		shimPath string
	}
	tests := []struct {
		name                     string
		args                     args
		configFile               string
		initialConfigFileContent string
		createFile               bool
		wantErr                  bool
		wantFileContent          string
	}{
		{"foobar", args{"/assets/foobar"}, "/etc/containerd/config.toml", "Hello World\n", true, false,
			`Hello World

# KWASM runtime config for foobar
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.foobar]
runtime_type = "/assets/foobar"
`},
		{"foobar", args{"/assets/foobar"}, "/etc/config.toml", "", false, true, ``},
		{"foobar", args{"/assets/foobar"}, "/etc/containerd/config.toml", `Hello World

# KWASM runtime config for foobar
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.foobar]
runtime_type = "/assets/foobar"

Foobar
`, true, false,
			`Hello World

# KWASM runtime config for foobar
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.foobar]
runtime_type = "/assets/foobar"

Foobar
`},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			if tt.createFile {
				file, err := fs.Create(tt.configFile)
				if err != nil {
					t.Fatal(err)
				}

				_, err = file.WriteString(tt.initialConfigFileContent)
				if err != nil {
					t.Fatal(err)
				}
			}

			c := &Config{
				config: newGlobalConfig(tt.configFile),
				fs:     fs,
			}
			got, err := c.AddRuntime(tt.args.shimPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.AddRuntime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != c.config.Runtime.ConfigPath {
				t.Errorf("Config.AddRuntime() = %v, want %v", got, c.config.Runtime.ConfigPath)
			}

			if tt.wantErr {
				return
			}

			gotFileContent, err := afero.ReadFile(fs, tt.configFile)
			if err != nil {
				t.Fatal(err)
			}

			if string(gotFileContent) != tt.wantFileContent {
				t.Errorf("runtimeConfigFile content: %v, want %v", string(gotFileContent), tt.wantFileContent)
			}

		})
	}
}
