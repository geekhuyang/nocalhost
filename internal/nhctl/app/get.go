/*
Copyright 2020 The Nocalhost Authors.
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

package app

import (
	"nocalhost/internal/nhctl/nocalhost"
	"nocalhost/internal/nhctl/profile"
	"nocalhost/pkg/nhctl/log"
	"os"
	"path/filepath"
)

func (a *Application) GetDependencies() []*SvcDependency {
	result := make([]*SvcDependency, 0)

	if a.configV2 == nil {
		return nil
	}

	svcConfigs := a.configV2.ApplicationConfig.ServiceConfigs
	if len(svcConfigs) == 0 {
		return nil
	}

	for _, svcConfig := range svcConfigs {
		if svcConfig.DependLabelSelector == nil {
			continue
		}
		if svcConfig.DependLabelSelector.Pods == nil && svcConfig.DependLabelSelector.Jobs == nil {
			continue
		}
		svcDep := &SvcDependency{
			Name: svcConfig.Name,
			Type: string(svcConfig.Type),
			Jobs: svcConfig.DependLabelSelector.Jobs,
			Pods: svcConfig.DependLabelSelector.Pods,
		}
		result = append(result, svcDep)
	}
	return result
}

// Get local path of resource dirs
// If resource path undefined, use git url
func (a *Application) GetResourceDir() []string {
	var resourcePath []string
	if len(a.AppProfileV2.ResourcePath) != 0 {
		for _, path := range a.AppProfileV2.ResourcePath {
			fullPath := filepath.Join(a.getGitDir(), path)
			resourcePath = append(resourcePath, fullPath)
		}
		return resourcePath
	}
	return []string{a.getGitDir()}
}

//func (a *Application) getUpgradeResourceDir() []string {
//	var resourcePath []string
//	if len(a.AppProfileV2.ResourcePath) != 0 {
//		for _, path := range a.AppProfileV2.ResourcePath {
//			fullPath := filepath.Join(a.getUpgradeGitDir(), path)
//			resourcePath = append(resourcePath, fullPath)
//		}
//		return resourcePath
//	}
//	return []string{a.getUpgradeGitDir()}
//}
func (a *Application) getUpgradeResourceDir(upgradeResourcePath []string) []string {
	var resourcePath []string
	if len(upgradeResourcePath) != 0 {
		for _, path := range upgradeResourcePath {
			fullPath := filepath.Join(a.getUpgradeGitDir(), path)
			resourcePath = append(resourcePath, fullPath)
		}
		return resourcePath
	}
	return []string{a.getUpgradeGitDir()}
}

func (a *Application) getIgnoredPath() []string {
	results := make([]string, 0)
	for _, path := range a.AppProfileV2.IgnoredPath {
		results = append(results, filepath.Join(a.getGitDir(), path))
	}
	return results
}

func (a *Application) getPreInstallFiles() []string {
	return a.sortedPreInstallManifest
}

func (a *Application) getUpgradePreInstallFiles() []string {
	return a.upgradeSortedPreInstallManifest
}

func (a *Application) getUpgradeIgnoredPath() []string {
	results := make([]string, 0)
	for _, path := range a.AppProfileV2.IgnoredPath {
		results = append(results, filepath.Join(a.getUpgradeGitDir(), path))
	}
	return results
}

func (a *Application) GetDefaultWorkDir(svcName, container string) string {
	svcProfile := a.GetSvcProfileV2(svcName)
	if svcProfile != nil && svcProfile.GetContainerDevConfigOrDefault(container).WorkDir != "" {
		return svcProfile.GetContainerDevConfigOrDefault(container).WorkDir
	}
	return DefaultWorkDir
}

func (a *Application) GetPersistentVolumeDirs(svcName, container string) []*profile.PersistentVolumeDir {
	svcProfile := a.GetSvcProfileV2(svcName)
	if svcProfile != nil {
		return svcProfile.GetContainerDevConfigOrDefault(container).PersistentVolumeDirs
	}
	return nil
}

func (a *Application) GetDefaultSideCarImage(svcName string) string {
	return DefaultSideCarImage
}

func (a *Application) GetDefaultDevImage(svcName string, container string) string {
	svcProfile := a.GetSvcProfileV2(svcName)
	if svcProfile != nil && svcProfile.GetContainerDevConfigOrDefault(container).Image != "" {
		return svcProfile.GetContainerDevConfigOrDefault(container).Image
	}
	return DefaultDevImage
}

func (a *Application) GetDefaultDevPort(svcName string, container string) []string {
	config := a.GetSvcProfileV2(svcName)
	if config != nil && len(config.GetContainerDevConfigOrDefault(container).PortForward) > 0 {
		return config.GetContainerDevConfigOrDefault(container).PortForward
	}
	return []string{}
}

func (a *Application) GetApplicationSyncDir(deployment string) string {
	dirPath := filepath.Join(a.GetHomeDir(), nocalhost.DefaultBinSyncThingDirName, deployment)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0700)
		if err != nil {
			log.Fatalf("fail to create syncthing directory: %s", dirPath)
		}
	}
	return dirPath
}
