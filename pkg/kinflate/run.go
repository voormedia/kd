package kinflate

import (
	"fmt"
	"io"
	"sort"

	// "k8s.io/apimachinery/pkg/api/meta"
	// "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/kubectl/pkg/scheme"
)

type groupVersionKindName struct {
	gvk  schema.GroupVersionKind
	name string
}

func Run(dir string, out io.Writer) error {
	baseResources, overlayResource, overlayPkg, err := loadBaseAndOverlayPkg(dir)
	if err != nil {
		return err
	}

	gvknToNewNameObject := map[groupVersionKindName]newNameObject{}

	// map from GroupVersionKind to marshaled json bytes
	overlayResouceMap := map[groupVersionKindName][]byte{}
	err = populateResourceMap(overlayResource.resources, overlayResouceMap)
	if err != nil {
		return err
	}

	err = populateMapOfConfigMapAndSecret(overlayResource, gvknToNewNameObject)
	if err != nil {
		return err
	}

	// map from GroupVersionKind to marshaled json bytes
	baseResouceMap := map[groupVersionKindName][]byte{}
	for _, baseResource := range baseResources {
		err = populateResourceMap(baseResource.resources, baseResouceMap)
		if err != nil {
			return err
		}
		err = populateMapOfConfigMapAndSecret(baseResource, gvknToNewNameObject)
		if err != nil {
			return err
		}
	}

	// Strategic merge the resources exist in both base and overlay.
	for gvkn, base := range baseResouceMap {
		// Merge overlay with base resource.
		if overlay, found := overlayResouceMap[gvkn]; found {
			versionedObj, err := scheme.Scheme.New(gvkn.gvk)
			if err != nil {
				switch {
				case runtime.IsNotRegisteredError(err):
					return fmt.Errorf("CRD and TPR are not supported now: %v", err)
				default:
					return err
				}
			}
			merged, err := strategicpatch.StrategicMergePatch(base, overlay, versionedObj)
			if err != nil {
				return err
			}
			baseResouceMap[gvkn] = merged
			delete(overlayResouceMap, gvkn)
		}
	}

	// If there are resources in overlay that are not defined in base, just add it to base.
	if len(overlayResouceMap) > 0 {
		for gvkn, jsonObj := range overlayResouceMap {
			baseResouceMap[gvkn] = jsonObj
		}
	}

	cmAndSecretGVKN := []groupVersionKindName{}
	for gvkn := range gvknToNewNameObject {
		cmAndSecretGVKN = append(cmAndSecretGVKN, gvkn)
	}
	sort.Sort(ByGVKN(cmAndSecretGVKN))
	for _, gvkn := range cmAndSecretGVKN {
		nameAndobj := gvknToNewNameObject[gvkn]
		yamlObj, err := updateObjectMetadata(nameAndobj.obj, overlayPkg)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "---\n%s", yamlObj)
	}

	// Inject the labels, annotations and name prefix.
	// Then print the object.
	resourceGVKN := []groupVersionKindName{}
	for gvkn := range baseResouceMap {
		resourceGVKN = append(resourceGVKN, gvkn)
	}
	sort.Sort(ByGVKN(resourceGVKN))
	for _, gvkn := range resourceGVKN {
		yamlObj, err := updateMetadata(baseResouceMap[gvkn], overlayPkg, gvknToNewNameObject)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "---\n%s", yamlObj)
	}
	return nil
}
