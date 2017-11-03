package mapping

import (
	containers "github.com/ernoaapa/can/pkg/api/services/containers/v1"
	pods "github.com/ernoaapa/can/pkg/api/services/pods/v1"
	"github.com/ernoaapa/can/pkg/model"
)

// MapPodsToInternalModel maps API Pod model to internal model
func MapPodsToInternalModel(pods []*pods.Pod) (result []model.Pod) {
	for _, pod := range pods {
		result = append(result, MapPodToInternalModel(pod))
	}
	return result
}

// MapPodToInternalModel maps API Pod model to internal model
func MapPodToInternalModel(pod *pods.Pod) model.Pod {
	return model.Pod{
		Metadata: model.Metadata{
			Name:      pod.Metadata.Name,
			Namespace: pod.Metadata.Namespace,
		},
		Spec: model.PodSpec{
			Containers:  MapContainerToInternalModel(pod.Spec.Containers),
			HostNetwork: pod.Spec.HostNetwork,
			HostPID:     pod.Spec.HostPID,
		},
	}
}

// MapContainerToInternalModel maps API Container model to internal model
func MapContainerToInternalModel(containers []*containers.Container) (result []model.Container) {
	for _, container := range containers {
		result = append(result, model.Container{
			Name:       container.Name,
			Image:      container.Image,
			Tty:        container.Tty,
			Args:       container.Args,
			Env:        container.Env,
			WorkingDir: container.WorkingDir,
			Mounts:     mapMountsToInternalModel(container.Mounts),
			Io:         mapIOSetToInternalModel(container.Io),
		})
	}
	return result
}

func mapIOSetToInternalModel(io *containers.IOSet) model.IOSet {
	return model.IOSet{
		In:  io.In,
		Out: io.Out,
		Err: io.Err,
	}
}

func mapMountsToInternalModel(mounts []*containers.Mount) (result []model.Mount) {
	for _, mount := range mounts {
		result = append(result, model.Mount{
			Type:        mount.Type,
			Source:      mount.Source,
			Destination: mount.Destination,
			Options:     mount.Options,
		})
	}
	return result
}
