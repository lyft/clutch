package k8s

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"
)

func (s *svc) GetLogs(ctx context.Context, clientset, cluster, namespace, name string) error {
	cs, err := s.manager.GetK8sClientset(ctx, clientset, cluster, namespace)
	if err != nil {
		return err
	}

	req := cs.CoreV1().Pods(cs.Namespace()).GetLogs(name, &v1.PodLogOptions{
		TailLines: pointer.Int64Ptr(10),
	})

	res := req.Do(ctx)
	if err := res.Error(); err != nil {
		return err
	}

	rawLog, err := res.Raw()
	if err != nil {
		return err
	}

	fmt.Println(string(rawLog))

	return nil
}
