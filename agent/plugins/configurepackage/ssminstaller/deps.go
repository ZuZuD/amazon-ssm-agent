// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the
// License is located at
//
// http://aws.amazon.com/apache2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package ssminstaller

import (
	"github.com/aws/amazon-ssm-agent/agent/context"
	"github.com/aws/amazon-ssm-agent/agent/contracts"
	"github.com/aws/amazon-ssm-agent/agent/docmanager/model"
	"github.com/aws/amazon-ssm-agent/agent/docparser"
	"github.com/aws/amazon-ssm-agent/agent/fileutil"
	"github.com/aws/amazon-ssm-agent/agent/framework/runpluginutil"
	"github.com/aws/amazon-ssm-agent/agent/platform"

	"encoding/json"
	"io/ioutil"
)

// dependency on action execution
type execDep interface {
	ParseDocument(context context.T, documentRaw []byte, orchestrationDir string, s3Bucket string, s3KeyPrefix string, messageID string, documentID string, defaultWorkingDirectory string) (pluginsInfo []model.PluginState, err error)
	ExecuteDocument(runner runpluginutil.PluginRunner, context context.T, pluginInput []model.PluginState, documentID string, documentCreatedDate string) (pluginOutputs map[string]*contracts.PluginResult)
}

type execDepImp struct {
}

func (m *execDepImp) ParseDocument(context context.T, documentRaw []byte, orchestrationDir string, s3Bucket string, s3KeyPrefix string, messageID string, documentID string, defaultWorkingDirectory string) (pluginsInfo []model.PluginState, err error) {
	log := context.Log()
	parserInfo := docparser.DocumentParserInfo{
		OrchestrationDir:  orchestrationDir,
		S3Bucket:          s3Bucket,
		S3Prefix:          s3KeyPrefix,
		MessageId:         messageID,
		DocumentId:        documentID,
		DefaultWorkingDir: defaultWorkingDirectory,
	}

	var docContent contracts.DocumentContent
	err = json.Unmarshal(documentRaw, &docContent)
	if err != nil {
		return
	}
	// TODO Add parameters
	return docparser.ParseDocument(log, &docContent, parserInfo, nil)
}

func (m *execDepImp) ExecuteDocument(runner runpluginutil.PluginRunner, context context.T, pluginInput []model.PluginState, documentID string, documentCreatedDate string) (pluginOutputs map[string]*contracts.PluginResult) {
	log := context.Log()
	log.Debugf("Running subcommand")
	return runner.ExecuteDocument(context, pluginInput, documentID, documentCreatedDate)
}

// dependency on filesystem and os utility functions
type fileSysDep interface {
	Exists(filePath string) bool
	ReadFile(filename string) ([]byte, error)
}

type fileSysDepImp struct{}

func (fileSysDepImp) Exists(filePath string) bool {
	return fileutil.Exists(filePath)
}

func (fileSysDepImp) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

var instance instanceInfo = &instanceInfoImp{}

// system represents the dependency for platform
type instanceInfo interface {
	InstanceID() (string, error)
	Region() (string, error)
}

type instanceInfoImp struct{}

// InstanceID wraps platform InstanceID
func (instanceInfoImp) InstanceID() (string, error) { return platform.InstanceID() }

// Region wraps platform Region
func (instanceInfoImp) Region() (string, error) { return platform.Region() }
