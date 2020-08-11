package k8s

import (
	k8sapiv1 "github.com/lyft/clutch/backend/api/k8s/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8slabels "k8s.io/apimachinery/pkg/labels"
)

func ApplyListOptions(listOpts *k8sapiv1.ListOptions) metav1.ListOptions {
	opts := metav1.ListOptions{}
	if len(listOpts.Labels) > 0 {
		opts.LabelSelector = k8slabels.FormatLabels(listOpts.Labels)
	}

	if len(listOpts.FieldSelectors) > 0 {
		opts.FieldSelector = listOpts.FieldSelectors
	}

	return opts
}
