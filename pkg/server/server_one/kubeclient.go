package server_one

import (
	"errors"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func (k *KubeClient) Init(kubeconfig string) error {
	var config *rest.Config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		k.clientset = nil
		return err
	}
	k.clientset = clientset
	return nil
}

func (k *KubeClient) DeployJob(numberOfPods int32, testName, namespace, logtime, msgcount, TimeToWaitAfterLoggingSec, alasticAPI, master string) ([]byte, error) {
	if k.clientset == nil {
		return []byte{}, errors.New("missing or unvalid kubeconfig file")
	}
	jobClient := k.clientset.BatchV1().Jobs(namespace)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: namespace,
			Labels: map[string]string{
				"app":     "flood-and-analyse",
				"role":    "test",
				"section": namespace,
			},
		},
		Spec: batchv1.JobSpec{
			Parallelism: int32Ptr(numberOfPods),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     "flood-and-analyse",
						"role":    "test",
						"section": namespace,
					},
				},
				Spec: apiv1.PodSpec{
					RestartPolicy: apiv1.RestartPolicyNever,
					Containers: []apiv1.Container{
						{
							Name:            testName,
							Image:           "hisshadow85/flood-and-analyze:latest",
							ImagePullPolicy: apiv1.PullAlways,
							Env: []apiv1.EnvVar{
								{
									Name: "POD_NAME",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
							},
							Command: []string{
								"/app/flood_and_analyze",
								"--workername=$(POD_NAME)",
								"--logtime=" + logtime,
								"--msgcount=" + msgcount,
								"--time_to_wait_after_logging_sec=" + TimeToWaitAfterLoggingSec,
								"--elastic_end_point=" + alasticAPI,
								"--master=" + master,
								"--testname=" + testName,
							},
						},
					},
				},
			},
		},
	}
	result, err := jobClient.Create(job)
	if err != nil {
		return []byte{}, err
	}
	return []byte("Created job " + result.GetObjectMeta().GetName()), nil
}

func int32Ptr(i int32) *int32 { return &i }

func (k *KubeClient) DeleteJob(namespace, name string) ([]byte, error) {
	if k.clientset == nil {
		return []byte{}, errors.New("missing or unvalid kubeconfig file")
	}
	jobClient := k.clientset.BatchV1().Jobs(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := jobClient.Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return []byte{}, err
	}
	return []byte("Job " + name + " deleted."), nil
}
