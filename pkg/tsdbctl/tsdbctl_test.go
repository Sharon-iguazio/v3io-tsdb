// +build unit

/*
Copyright 2018 Iguazio Systems Ltd.

Licensed under the Apache License, Version 2.0 (the "License") with
an addition restriction as set forth herein. You may not use this
file except in compliance with the License. You may obtain a copy of
the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
implied. See the License for the specific language governing
permissions and limitations under the License.

In addition, you may not use the software for any purposes that are
illegal under applicable law, and the grant of the foregoing license
under the Apache 2.0 license is conditioned upon your compliance with
such restriction.
*/

package tsdbctl

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"github.com/v3io/v3io-tsdb/internal/pkg/performance"
	"github.com/v3io/v3io-tsdb/pkg/config"
)

type testTsdbctlSuite struct {
	suite.Suite
}

func (suite *testTsdbctlSuite) TestPopulateConfigWithTenant() {
	rc := RootCommandeer{v3ioURL: "localhost:80123"}
	cfg := &config.V3ioConfig{Username: "Vel@Odar", Password: "p455w0rd", Container: "123", TablePath: "/x/y/z"}

	err := rc.populateConfig(cfg)
	suite.Require().Nil(err)

	metricReporter, err := performance.DefaultReporterInstance()
	if err != nil {
		err = errors.Wrap(err, "unable to initialize performance metrics reporter")
		return
	}
	suite.Require().Nil(err)

	expectedRc := RootCommandeer{
		v3iocfg:  cfg,
		v3ioURL:  "localhost:80123",
		Reporter: metricReporter,
	}
	expectedCfg := &config.V3ioConfig{
		WebAPIEndpoint: "localhost:80123",
		Container:      "123",
		TablePath:      "/x/y/z",
		Username:       "Vel@Odar",
		Password:       "p455w0rd",
		LogLevel:       "info",
	}

	suite.Require().Equal(expectedCfg, rc.v3iocfg)
	suite.Require().Equal(expectedRc, rc)
}

func (suite *testTsdbctlSuite) TestContainerConfig() {
	oldV3ioApi := os.Getenv("V3IO_API")
	err := os.Setenv("V3IO_API", "host-from-env:123")
	suite.Require().NoError(err)
	defer os.Setenv("V3IO_API", oldV3ioApi)

	oldAccessKey := os.Getenv("V3IO_ACCESS_KEY")
	err = os.Setenv("V3IO_ACCESS_KEY", "key-from-env")
	suite.Require().NoError(err)
	defer os.Setenv("V3IO_ACCESS_KEY", oldAccessKey)

	rc := RootCommandeer{v3ioURL: "localhost:80123", container: "test", accessKey: "acce55-key"}
	cfg := &config.V3ioConfig{Username: "Vel@Odar", Password: "p455w0rd", TablePath: "/x/y/z"}

	err = rc.populateConfig(cfg)
	expectedCfg := &config.V3ioConfig{
		WebAPIEndpoint: "localhost:80123",
		Container:      "test",
		TablePath:      "/x/y/z",
		Username:       "Vel@Odar",
		Password:       "p455w0rd",
		AccessKey:      "acce55-key",
		LogLevel:       "info",
	}

	suite.Require().Nil(err)
	suite.Require().Equal(expectedCfg, rc.v3iocfg)
}

func (suite *testTsdbctlSuite) TestConfigFromEnvVarsAndPassword() {
	oldV3ioApi := os.Getenv("V3IO_API")
	err := os.Setenv("V3IO_API", "host-from-env:123")
	suite.Require().NoError(err)
	defer os.Setenv("V3IO_API", oldV3ioApi)

	oldAccessKey := os.Getenv("V3IO_ACCESS_KEY")
	err = os.Setenv("V3IO_ACCESS_KEY", "key-from-env")
	suite.Require().NoError(err)
	defer os.Setenv("V3IO_ACCESS_KEY", oldAccessKey)

	rc := RootCommandeer{container: "test", username: "Vel@Odar", password: "p455w0rd"}
	cfg, err := config.GetOrLoadFromStruct(&config.V3ioConfig{TablePath: "/x/y/z"})
	suite.Require().Nil(err)

	expectedCfg := *cfg
	err = rc.populateConfig(cfg)
	expectedCfg.WebAPIEndpoint = "host-from-env:123"
	expectedCfg.Container = "test"
	expectedCfg.TablePath = "/x/y/z"
	expectedCfg.Username = "Vel@Odar"
	expectedCfg.Password = "p455w0rd"
	expectedCfg.LogLevel = "info"

	suite.Require().Nil(err)
	suite.Require().Equal(&expectedCfg, rc.v3iocfg)
}

func (suite *testTsdbctlSuite) TestConfigFromEnvVars() {
	oldV3ioApi := os.Getenv("V3IO_API")
	err := os.Setenv("V3IO_API", "host-from-env:123")
	suite.Require().NoError(err)
	defer os.Setenv("V3IO_API", oldV3ioApi)

	oldAccessKey := os.Getenv("V3IO_ACCESS_KEY")
	err = os.Setenv("V3IO_ACCESS_KEY", "key-from-env")
	suite.Require().NoError(err)
	defer os.Setenv("V3IO_ACCESS_KEY", oldAccessKey)

	rc := RootCommandeer{container: "test"}
	cfg, err := config.GetOrLoadFromStruct(&config.V3ioConfig{TablePath: "/x/y/z"})
	suite.NoError(err)

	expectedCfg := *cfg
	err = rc.populateConfig(cfg)
	expectedCfg.WebAPIEndpoint = "host-from-env:123"
	expectedCfg.AccessKey = "key-from-env"
	expectedCfg.Container = "test"
	expectedCfg.LogLevel = "info"

	suite.Require().Nil(err)
	suite.Require().Equal(&expectedCfg, rc.v3iocfg)
}

func TestTsdbctlSuite(t *testing.T) {
	suite.Run(t, new(testTsdbctlSuite))
}
