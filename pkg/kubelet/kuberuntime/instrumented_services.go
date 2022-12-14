/*
Copyright 2016 The Kubernetes Authors.

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

package kuberuntime

import (
	"time"

	internalapi "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/kubernetes/pkg/kubelet/metrics"
)

// instrumentedRuntimeService wraps the RuntimeService and records the operations
// and errors metrics.
type instrumentedRuntimeService struct {
	service internalapi.RuntimeService
}

// Creates an instrumented RuntimeInterface from an existing RuntimeService.
func newInstrumentedRuntimeService(service internalapi.RuntimeService) internalapi.RuntimeService {
	return &instrumentedRuntimeService{service: service}
}

// instrumentedImageManagerService wraps the ImageManagerService and records the operations
// and errors metrics.
type instrumentedImageManagerService struct {
	service internalapi.ImageManagerService
}

// Creates an instrumented ImageManagerService from an existing ImageManagerService.
func newInstrumentedImageManagerService(service internalapi.ImageManagerService) internalapi.ImageManagerService {
	return &instrumentedImageManagerService{service: service}
}

// recordOperation records the duration of the operation.
func recordOperation(operation string, start time.Time) {
	metrics.RuntimeOperations.WithLabelValues(operation).Inc()
	metrics.RuntimeOperationsDuration.WithLabelValues(operation).Observe(metrics.SinceInSeconds(start))
}

// recordError records error for metric if an error occurred.
func recordError(operation string, err error) {
	if err != nil {
		metrics.RuntimeOperationsErrors.WithLabelValues(operation).Inc()
	}
}

func (in instrumentedRuntimeService) Version(apiVersion string) (*runtimeapi.VersionResponse, error) {
	const operation = "version"
	defer recordOperation(operation, time.Now())

	out, err := in.service.Version(apiVersion)
	recordError(operation, err)
	return out, err
}

func (in instrumentedRuntimeService) Status(verbose bool) (*runtimeapi.StatusResponse, error) {
	const operation = "status"
	defer recordOperation(operation, time.Now())

	out, err := in.service.Status(verbose)
	recordError(operation, err)
	return out, err
}

func (in instrumentedRuntimeService) CreateContainer(podSandboxID string, config *runtimeapi.ContainerConfig, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	const operation = "create_container"
	defer recordOperation(operation, time.Now())

	out, err := in.service.CreateContainer(podSandboxID, config, sandboxConfig)
	recordError(operation, err)
	return out, err
}

func (in instrumentedRuntimeService) StartContainer(containerID string) error {
	const operation = "start_container"
	defer recordOperation(operation, time.Now())

	err := in.service.StartContainer(containerID)
	recordError(operation, err)
	return err
}

func (in instrumentedRuntimeService) StopContainer(containerID string, timeout int64) error {
	const operation = "stop_container"
	defer recordOperation(operation, time.Now())

	err := in.service.StopContainer(containerID, timeout)
	recordError(operation, err)
	return err
}

func (in instrumentedRuntimeService) RemoveContainer(containerID string) error {
	const operation = "remove_container"
	defer recordOperation(operation, time.Now())

	err := in.service.RemoveContainer(containerID)
	recordError(operation, err)
	return err
}

func (in instrumentedRuntimeService) ListContainers(filter *runtimeapi.ContainerFilter) ([]*runtimeapi.Container, error) {
	const operation = "list_containers"
	defer recordOperation(operation, time.Now())

	out, err := in.service.ListContainers(filter)
	recordError(operation, err)
	return out, err
}

func (in instrumentedRuntimeService) ContainerStatus(containerID string, verbose bool) (*runtimeapi.ContainerStatusResponse, error) {
	const operation = "container_status"
	defer recordOperation(operation, time.Now())

	out, err := in.service.ContainerStatus(containerID, verbose)
	recordError(operation, err)
	return out, err
}

func (in instrumentedRuntimeService) UpdateContainerResources(containerID string, resources *runtimeapi.ContainerResources) error {
	const operation = "update_container"
	defer recordOperation(operation, time.Now())

	err := in.service.UpdateContainerResources(containerID, resources)
	recordError(operation, err)
	return err
}

func (in instrumentedRuntimeService) ReopenContainerLog(containerID string) error {
	const operation = "reopen_container_log"
	defer recordOperation(operation, time.Now())

	err := in.service.ReopenContainerLog(containerID)
	recordError(operation, err)
	return err
}

func (in instrumentedRuntimeService) ExecSync(containerID string, cmd []string, timeout time.Duration) ([]byte, []byte, error) {
	const operation = "exec_sync"
	defer recordOperation(operation, time.Now())

	stdout, stderr, err := in.service.ExecSync(containerID, cmd, timeout)
	recordError(operation, err)
	return stdout, stderr, err
}

func (in instrumentedRuntimeService) Exec(req *runtimeapi.ExecRequest) (*runtimeapi.ExecResponse, error) {
	const operation = "exec"
	defer recordOperation(operation, time.Now())

	resp, err := in.service.Exec(req)
	recordError(operation, err)
	return resp, err
}

func (in instrumentedRuntimeService) Attach(req *runtimeapi.AttachRequest) (*runtimeapi.AttachResponse, error) {
	const operation = "attach"
	defer recordOperation(operation, time.Now())

	resp, err := in.service.Attach(req)
	recordError(operation, err)
	return resp, err
}

func (in instrumentedRuntimeService) RunPodSandbox(config *runtimeapi.PodSandboxConfig, runtimeHandler string) (string, error) {
	const operation = "run_podsandbox"
	startTime := time.Now()
	defer recordOperation(operation, startTime)
	defer metrics.RunPodSandboxDuration.WithLabelValues(runtimeHandler).Observe(metrics.SinceInSeconds(startTime))

	out, err := in.service.RunPodSandbox(config, runtimeHandler)
	recordError(operation, err)
	if err != nil {
		metrics.RunPodSandboxErrors.WithLabelValues(runtimeHandler).Inc()
	}
	return out, err
}

func (in instrumentedRuntimeService) StopPodSandbox(podSandboxID string) error {
	const operation = "stop_podsandbox"
	defer recordOperation(operation, time.Now())

	err := in.service.StopPodSandbox(podSandboxID)
	recordError(operation, err)
	return err
}

func (in instrumentedRuntimeService) RemovePodSandbox(podSandboxID string) error {
	const operation = "remove_podsandbox"
	defer recordOperation(operation, time.Now())

	err := in.service.RemovePodSandbox(podSandboxID)
	recordError(operation, err)
	return err
}

func (in instrumentedRuntimeService) PodSandboxStatus(podSandboxID string, verbose bool) (*runtimeapi.PodSandboxStatusResponse, error) {
	const operation = "podsandbox_status"
	defer recordOperation(operation, time.Now())

	out, err := in.service.PodSandboxStatus(podSandboxID, verbose)
	recordError(operation, err)
	return out, err
}

func (in instrumentedRuntimeService) ListPodSandbox(filter *runtimeapi.PodSandboxFilter) ([]*runtimeapi.PodSandbox, error) {
	const operation = "list_podsandbox"
	defer recordOperation(operation, time.Now())

	out, err := in.service.ListPodSandbox(filter)
	recordError(operation, err)
	return out, err
}

func (in instrumentedRuntimeService) ContainerStats(containerID string) (*runtimeapi.ContainerStats, error) {
	const operation = "container_stats"
	defer recordOperation(operation, time.Now())

	out, err := in.service.ContainerStats(containerID)
	recordError(operation, err)
	return out, err
}

func (in instrumentedRuntimeService) ListContainerStats(filter *runtimeapi.ContainerStatsFilter) ([]*runtimeapi.ContainerStats, error) {
	const operation = "list_container_stats"
	defer recordOperation(operation, time.Now())

	out, err := in.service.ListContainerStats(filter)
	recordError(operation, err)
	return out, err
}

func (in instrumentedRuntimeService) PodSandboxStats(podSandboxID string) (*runtimeapi.PodSandboxStats, error) {
	const operation = "podsandbox_stats"
	defer recordOperation(operation, time.Now())

	out, err := in.service.PodSandboxStats(podSandboxID)
	recordError(operation, err)
	return out, err
}

func (in instrumentedRuntimeService) ListPodSandboxStats(filter *runtimeapi.PodSandboxStatsFilter) ([]*runtimeapi.PodSandboxStats, error) {
	const operation = "list_podsandbox_stats"
	defer recordOperation(operation, time.Now())

	out, err := in.service.ListPodSandboxStats(filter)
	recordError(operation, err)
	return out, err
}

func (in instrumentedRuntimeService) PortForward(req *runtimeapi.PortForwardRequest) (*runtimeapi.PortForwardResponse, error) {
	const operation = "port_forward"
	defer recordOperation(operation, time.Now())

	resp, err := in.service.PortForward(req)
	recordError(operation, err)
	return resp, err
}

func (in instrumentedRuntimeService) UpdateRuntimeConfig(runtimeConfig *runtimeapi.RuntimeConfig) error {
	const operation = "update_runtime_config"
	defer recordOperation(operation, time.Now())

	err := in.service.UpdateRuntimeConfig(runtimeConfig)
	recordError(operation, err)
	return err
}

func (in instrumentedImageManagerService) ListImages(filter *runtimeapi.ImageFilter) ([]*runtimeapi.Image, error) {
	const operation = "list_images"
	defer recordOperation(operation, time.Now())

	out, err := in.service.ListImages(filter)
	recordError(operation, err)
	return out, err
}

func (in instrumentedImageManagerService) ImageStatus(image *runtimeapi.ImageSpec, verbose bool) (*runtimeapi.ImageStatusResponse, error) {
	const operation = "image_status"
	defer recordOperation(operation, time.Now())

	out, err := in.service.ImageStatus(image, verbose)
	recordError(operation, err)
	return out, err
}

func (in instrumentedImageManagerService) PullImage(image *runtimeapi.ImageSpec, auth *runtimeapi.AuthConfig, podSandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	const operation = "pull_image"
	defer recordOperation(operation, time.Now())

	imageRef, err := in.service.PullImage(image, auth, podSandboxConfig)
	recordError(operation, err)
	return imageRef, err
}

func (in instrumentedImageManagerService) RemoveImage(image *runtimeapi.ImageSpec) error {
	const operation = "remove_image"
	defer recordOperation(operation, time.Now())

	err := in.service.RemoveImage(image)
	recordError(operation, err)
	return err
}

func (in instrumentedImageManagerService) ImageFsInfo() ([]*runtimeapi.FilesystemUsage, error) {
	const operation = "image_fs_info"
	defer recordOperation(operation, time.Now())

	fsInfo, err := in.service.ImageFsInfo()
	recordError(operation, err)
	return fsInfo, nil
}

func (in instrumentedRuntimeService) CheckpointContainer(options *runtimeapi.CheckpointContainerRequest) error {
	const operation = "checkpoint_container"
	defer recordOperation(operation, time.Now())

	err := in.service.CheckpointContainer(options)
	recordError(operation, err)
	return err
}

func (in instrumentedRuntimeService) GetContainerEvents(containerEventsCh chan *runtimeapi.ContainerEventResponse) error {
	const operation = "get_container_events"
	defer recordOperation(operation, time.Now())

	err := in.service.GetContainerEvents(containerEventsCh)
	recordError(operation, err)
	return err
}
