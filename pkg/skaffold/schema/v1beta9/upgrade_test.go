/*
Copyright 2019 The Skaffold Authors

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

package v1beta9

import (
	"testing"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/GoogleContainerTools/skaffold/testutil"
	yaml "gopkg.in/yaml.v2"
)

func TestUpgrade(t *testing.T) {
	yaml := `apiVersion: skaffold/v1beta9
kind: Config
build:
  artifacts:
  - image: gcr.io/k8s-skaffold/skaffold-example
    plugin:
      name: docker
      properties:
        dockerfile: path/to/Dockerfile
  - image: gcr.io/k8s-skaffold/bazel
    plugin:
      name: bazel
      properties:
        target: //mytarget
  executionEnvironment:
    name: local
test:
  - image: gcr.io/k8s-skaffold/skaffold-example
    structureTests:
     - ./test/*
deploy:
  kubectl:
    manifests:
    - k8s-*
profiles:
  - name: test profile
    build:
      artifacts:
      - image: gcr.io/k8s-skaffold/skaffold-example
        kaniko:
          buildContext:
            gcsBucket: skaffold-kaniko
          cache: {}
      cluster:
        pullSecretName: e2esecret
        namespace: default
    test:
     - image: gcr.io/k8s-skaffold/skaffold-example
       structureTests:
         - ./test/*
    deploy:
      kubectl:
        manifests:
        - k8s-*
`
	expected := `apiVersion: skaffold/v1beta10
kind: Config
build:
  artifacts:
  - image: gcr.io/k8s-skaffold/skaffold-example
    docker:
      dockerfile: path/to/Dockerfile
  - image: gcr.io/k8s-skaffold/bazel
    bazel:
      target: //mytarget
test:
  - image: gcr.io/k8s-skaffold/skaffold-example
    structureTests:
     - ./test/*
deploy:
  kubectl:
    manifests:
    - k8s-*
profiles:
  - name: test profile
    build:
      artifacts:
      - image: gcr.io/k8s-skaffold/skaffold-example
        kaniko:
          buildContext:
            gcsBucket: skaffold-kaniko
          cache: {}
      cluster:
        pullSecretName: e2esecret
        namespace: default
    test:
     - image: gcr.io/k8s-skaffold/skaffold-example
       structureTests:
         - ./test/*
    deploy:
      kubectl:
        manifests:
        - k8s-*
`
	verifyUpgrade(t, yaml, expected)
}

func verifyUpgrade(t *testing.T, input, output string) {
	pipeline := NewSkaffoldConfig()
	err := yaml.UnmarshalStrict([]byte(input), pipeline)
	testutil.CheckErrorAndDeepEqual(t, false, err, Version, pipeline.GetVersion())

	upgraded, err := pipeline.Upgrade()
	testutil.CheckError(t, false, err)

	expected := latest.NewSkaffoldConfig()
	err = yaml.UnmarshalStrict([]byte(output), expected)

	testutil.CheckErrorAndDeepEqual(t, false, err, expected, upgraded)
}
