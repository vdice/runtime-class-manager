// /*
// Copyright 2024.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	rcmv1 "github.com/spinkube/runtime-class-manager/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NodeReconciler reconciles a Node object
type NodeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const ControllerPrefix = "kwasm.sh/"

// TODO: remove create;delete?
//+kubebuilder:rbac:groups=runtime.kwasm.sh,resources=nodes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=runtime.kwasm.sh,resources=nodes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=runtime.kwasm.sh,resources=nodes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Node object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *NodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.With().Str("node", req.Name).Logger()
	node := &corev1.Node{}

	log.Info().Msgf("%s> Reconciliation called", node.Name)

	if err := r.Get(ctx, req.NamespacedName, node); err != nil {
		if apierrors.IsNotFound(err) {
			// we'll ignore not-found errors, since they can't be fixed by an immediate
			// requeue (we'll need to wait for a new notification), and we can get them
			// on deleted requests.
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get node: %w", err)
	}

	log.Info().Msgf("%s> Getting Shim Labels...", node.Name)

	// TODO: Not sure if this flow is still what we want.. or we at least need/should compare with Shim resources that exist on the cluster right?

	// Get all Labels that contain a Shim name
	// These are labels like
	// * "<shimname>=true"
	// * example: "spin=true", defined under the nodeSelector section of a Shim spec
	var shimList []string
	for i := range node.Labels {
		if strings.HasPrefix(i, ControllerPrefix) {
			b := strings.SplitN(i, ControllerPrefix, 2)
			shimList = append(shimList, b[1])
		}
		log.Debug().Msgf("Shim Labels> %v", shimList)
	}

	if len(shimList) == 0 {
		log.Info().Msgf("%s> No Shim Labels found", node.Name)
	}

	// Check if Shim Resources with name exists
	for _, shimName := range shimList {
		_, err := r.findShim(ctx, shimName)
		if err != nil {
			log.Err(err).Msgf("No Shim resource with name %s existing.", shimName)
		} else {
			log.Info().Msgf("Shim resource with name %s existing.", shimName)
		}
	}

	// Check if node has provisioned labels for every shim
	// * "<shimname>" = "provisioned"
	// * example: "shim-spin-v2" = "provisioned" means node has "shim-spin-v2" installed
	log.Info().Msgf("%s> Check if Shim provisioned...", node.Name)

	for _, shimName := range shimList {
		provisioned := false
		for labelName := range node.Labels {
			if strings.HasPrefix(labelName, ControllerPrefix) {
				log.Debug().Msgf("Label: %s", labelName)

				// Node has label, so it is provisioned
				if labelName == (ControllerPrefix + shimName) {
					provisioned = true
					log.Info().Msgf("Shim %s provisioned.", shimName)
				}
			}
		}

		// Node does not have label, so we need to start provisioning job
		if !provisioned {
			log.Info().Msgf("Shim %s not provisioned, start provisioning...", shimName)
			err := r.reconcileDeployment(ctx, node, req, shimName)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to reconcile deployment for shim %s", shimName)
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// findShim finds a shim resource by name.
func (r *NodeReconciler) findShim(ctx context.Context, shimName string) (*rcmv1.Shim, error) {
	var shimResource rcmv1.Shim
	err := r.Client.Get(ctx, types.NamespacedName{Name: shimName, Namespace: "default"}, &shimResource)
	if err != nil {
		return nil, err
	}
	return &shimResource, nil
}

// reconcileDeployment creates a Job if one does not exist and reconciles it if it does.
func (r *NodeReconciler) reconcileDeployment(ctx context.Context, node *corev1.Node, req ctrl.Request, shimName string) error {
	log := log.With().Str("node", node.Name).Str("shim", shimName).Logger()
	log.Info().Msgf("Reconciling Provisioning...")

	// TODO: We should be reusing the job creation, etc from shim_controller.go here
	// desiredJob, err := r.constructDeployJob(node, req, shimName)
	// if err != nil {
	// 	log.Error().Err(err).Msg("Failed to construct Job description")
	// 	return err
	// }

	log.Debug().Msgf("Reconciling Deployment Job")

	if err = r.Create(ctx, desiredJob); err != nil {
		log.Err(err).Msg("Failed to create new Job " + req.Namespace + " Job.Name " + req.Name)
		return fmt.Errorf("failed to create new job: %w", err)
	}

	return nil
}

// func (r *NodeReconciler) constructDeployJob(node *corev1.Node, req ctrl.Request, shimName string) (*batchv1.Job, error) {
// 	priv := true
// 	name := node.Name + "-" + shimName
// 	nameMax := int(math.Min(float64(len(name)), 63))

// 	installerImage := "ghcr.io/kwasm/kwasm-node-installer:main"

// 	dep := &batchv1.Job{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      name[:nameMax],
// 			Namespace: os.Getenv("CONTROLLER_NAMESPACE"),
// 			Labels: map[string]string{
// 				"kwasm.sh/job":         "true",
// 				"kwasm.sh/" + shimName: "true",
// 				"node":                 node.Name,
// 			},
// 		},
// 		Spec: batchv1.JobSpec{
// 			Template: corev1.PodTemplateSpec{
// 				Spec: corev1.PodSpec{
// 					NodeName: req.Name,
// 					HostPID:  true,
// 					Volumes: []corev1.Volume{{
// 						Name: "root-mount",
// 						VolumeSource: corev1.VolumeSource{
// 							HostPath: &corev1.HostPathVolumeSource{
// 								Path: "/",
// 							},
// 						},
// 					}},
// 					Containers: []corev1.Container{{
// 						Image: installerImage,
// 						Name:  "kwasm-provision",
// 						SecurityContext: &corev1.SecurityContext{
// 							Privileged: &priv,
// 						},
// 						Env: []corev1.EnvVar{
// 							{
// 								Name:  "NODE_ROOT",
// 								Value: "/mnt/node-root",
// 							},
// 						},
// 						VolumeMounts: []corev1.VolumeMount{
// 							{
// 								Name:      "root-mount",
// 								MountPath: "/mnt/node-root",
// 							},
// 						},
// 					}},
// 					RestartPolicy: corev1.RestartPolicyNever,
// 				},
// 			},
// 		},
// 	}
// 	if err := ctrl.SetControllerReference(node, dep, r.Scheme); err != nil {
// 		return nil, fmt.Errorf("failed to set controller reference: %w", err)
// 	}

// 	return dep, nil
// }

// SetupWithManager sets up the controller with the Manager.
func (r *NodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&corev1.Node{}).Complete(r)
}
