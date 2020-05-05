package k8s

import (
	"github.com/dbunion/com/scheduler"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	listOptionsKeyLabelSelector        = "LabelSelector"
	listOptionsKeyFieldSelector        = "FieldSelector"
	listOptionsKeyWatch                = "Watch"
	listOptionsKeyAllowWatchBookmarks  = "AllowWatchBookmarks"
	listOptionsKeyResourceVersion      = "ResourceVersion"
	listOptionsKeyTimeoutSeconds       = "TimeoutSeconds"
	listOptionsKeyLimit                = "Limit"
	listOptionsKeyContinue             = "Continue"
	deleteOptionsKeyGracePeriodSeconds = "GracePeriodSeconds"
)

// convertToListOptions - convert scheduler's Options to k8s's ListOptions
func convertToListOptions(options scheduler.Options) (op meta_v1.ListOptions) {
	// LabelSelector
	if v, found := options[listOptionsKeyLabelSelector]; found {
		if result, ok := v.(string); ok {
			op.LabelSelector = result
		}
	}

	// FieldSelector
	if v, found := options[listOptionsKeyFieldSelector]; found {
		if result, ok := v.(string); ok {
			op.FieldSelector = result
		}
	}

	// Watch
	if v, found := options[listOptionsKeyWatch]; found {
		if result, ok := v.(bool); ok {
			op.Watch = result
		}
	}

	// AllowWatchBookmarks
	if v, found := options[listOptionsKeyAllowWatchBookmarks]; found {
		if result, ok := v.(bool); ok {
			op.AllowWatchBookmarks = result
		}
	}

	// ResourceVersion
	if v, found := options[listOptionsKeyResourceVersion]; found {
		if result, ok := v.(string); ok {
			op.ResourceVersion = result
		}
	}

	// TimeoutSeconds
	if v, found := options[listOptionsKeyTimeoutSeconds]; found {
		if result, ok := v.(int64); ok {
			op.TimeoutSeconds = &result
		}
	}

	// Limit
	if v, found := options[listOptionsKeyLimit]; found {
		if result, ok := v.(int64); ok {
			op.Limit = result
		}
	}

	// Continue
	if v, found := options[listOptionsKeyContinue]; found {
		if result, ok := v.(string); ok {
			op.Continue = result
		}
	}
	return op
}

// convertToDeleteOptions - convert scheduler's Options to k8s's DeleteOptions
func convertToDeleteOptions(options scheduler.Options) (op meta_v1.DeleteOptions) {
	// GracePeriodSeconds
	if v, found := options[deleteOptionsKeyGracePeriodSeconds]; found {
		if result, ok := v.(int64); ok {
			op.GracePeriodSeconds = &result
		}
	}

	return op
}

// convertToCreateOptions - convert scheduler's Options to k8s's CreateOptions
func convertToCreateOptions(options scheduler.Options) (op meta_v1.CreateOptions) {
	return op
}

// convertToUpdateOptions - convert scheduler's Options to k8s's UpdateOptions
func convertToUpdateOptions(options scheduler.Options) (op meta_v1.UpdateOptions) {
	return op
}
