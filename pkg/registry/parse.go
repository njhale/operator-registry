package registry

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type bundleParser struct {
	log *logrus.Entry
}

func newBundleParser(log *logrus.Entry) *bundleParser {
	return &bundleParser{
		log: log,
	}
}

// Parse parses the given FS into a Bundle.
func (b *bundleParser) Parse(root fs.FS) (*Bundle, error) {
	if root == nil {
		return nil, fmt.Errorf("filesystem is nil")
	}

	bundle := &Bundle{}
	manifests, err := fs.Sub(root, "manifests")
	if err != nil {
		return nil, err
	}
	if err := b.addManifests(manifests, bundle); err != nil {
		return nil, err
	}

	metadata, err := fs.Sub(root, "metadata")
	if err != nil {
		return nil, err
	}
	if err := b.addMetadata(metadata, bundle); err != nil {
		return nil, err
	}

	derived, err := b.derivedProperties(bundle)
	if err != nil {
		return nil, fmt.Errorf("failed to derive properties: %s", err)
	}

	bundle.Properties = joinProperties(append(bundle.Properties, derived...))

	if err := expandProperties(bundle); err != nil {
		return nil, fmt.Errorf("property expansion failed, properties: %s, err: %s", bundle.Properties, err)
	}

	return bundle, nil
}

// addManifests adds the result of parsing the manifests directory to a bundle.
func (b *bundleParser) addManifests(manifests fs.FS, bundle *Bundle) error {
	files, err := fs.ReadDir(manifests, ".")
	if err != nil {
		return err
	}

	var csvFound bool
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		obj := &unstructured.Unstructured{}
		if err = decodeFileFS(manifests, name, obj); err != nil {
			b.log.Warnf("failed to decode: %s", err)
			continue
		}

		// Only include the first CSV we find in the
		if obj.GetKind() == "ClusterServiceVersion" {
			if csvFound {
				continue
			}
			csvFound = true
		}

		if obj.Object != nil {
			bundle.Add(obj)
		}
	}

	if bundle.Size() == 0 {
		return fmt.Errorf("no bundle objects found")
	}

	csv, err := bundle.ClusterServiceVersion()
	if err != nil {
		return err
	}
	if csv == nil {
		return fmt.Errorf("no csv in bundle")
	}

	bundle.Name = csv.GetName()
	if err := bundle.AllProvidedAPIsInBundle(); err != nil {
		return fmt.Errorf("error checking provided apis in bundle %s: %s", bundle.Name, err)
	}

	return err
}

// addManifests adds the result of parsing the metadata directory to a bundle.
func (b *bundleParser) addMetadata(metadata fs.FS, bundle *Bundle) error {
	files, err := fs.ReadDir(metadata, ".")
	if err != nil {
		return err
	}

	var (
		af *AnnotationsFile
		df *DependenciesFile
		pf *PropertiesFile
	)
	for _, f := range files {
		name := f.Name()
		if af == nil {
			decoded := AnnotationsFile{}
			if err = decodeFileFS(metadata, name, &decoded); err == nil {
				if decoded != (AnnotationsFile{}) {
					af = &decoded
				}
			}
		}
		if df == nil {
			decoded := DependenciesFile{}
			if err = decodeFileFS(metadata, name, &decoded); err == nil {
				if len(decoded.Dependencies) > 0 {
					df = &decoded
				}
			}
		}
		if pf == nil {
			decoded := PropertiesFile{}
			if err = decodeFileFS(metadata, name, &decoded); err == nil {
				if len(decoded.Properties) > 0 {
					pf = &decoded
				}
			}
		}
	}

	if af != nil {
		bundle.Annotations = &af.Annotations
		bundle.Channels = af.GetChannels()
	} else {
		return fmt.Errorf("Could not file annotations file")
	}

	if df != nil {
		bundle.Dependencies = append(bundle.Dependencies, df.GetDependencies()...)
	} else {
		b.log.Info("Could not find optional dependencies file")
	}

	if pf != nil {
		bundle.Properties = append(bundle.Properties, pf.Properties...)
	} else {
		b.log.Info("Could not find optional properties file")
	}

	return nil
}

func (b *bundleParser) derivedProperties(bundle *Bundle) ([]Property, error) {
	// Add properties from CSV annotations
	csv, err := bundle.ClusterServiceVersion()
	if err != nil {
		return nil, fmt.Errorf("error getting csv: %s", err)
	}
	if csv == nil {
		return nil, fmt.Errorf("bundle missing csv")
	}

	var derived []Property
	if len(csv.GetAnnotations()) > 0 {
		properties, ok := csv.GetAnnotations()[PropertyKey]
		if ok {
			if err := json.Unmarshal([]byte(properties), &derived); err != nil {
				b.log.Warnf("failed to unmarshal csv annotation properties: %s", err)
			}
		}
	}

	if bundle.Annotations != nil && bundle.Annotations.PackageName != "" {
		pkg := bundle.Annotations.PackageName
		version, err := bundle.Version()
		if err != nil {
			return nil, err
		}

		value, err := json.Marshal(PackageProperty{
			PackageName: pkg,
			Version:     version,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal package property: %s", err)
		}

		// Annotations file takes precedent over CSV annotations
		derived = append([]Property{{Type: PackageType, Value: value}}, derived...)
	}

	providedAPIs, err := bundle.ProvidedAPIs()
	if err != nil {
		return nil, fmt.Errorf("error getting provided apis: %s", err)
	}

	for api := range providedAPIs {
		value, err := json.Marshal(GVKProperty{
			Group:   api.Group,
			Kind:    api.Kind,
			Version: api.Version,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal gvk property: %s", err)
		}
		derived = append(derived, Property{Type: GVKType, Value: value})
	}

	return joinProperties(derived), nil
}

// joinProperties deduplicates and returns the given list of properties.
// It ensures the single package type invariant by retaining the first occurrence of that type in the returned list.
func joinProperties(properties []Property) []Property {
	var (
		foundPkg bool
		joined   []Property
		visited  = map[string]struct{}{}
	)
	for _, p := range properties {
		if _, ok := visited[p.String()]; ok {
			continue
		}
		visited[p.String()] = struct{}{}

		switch p.Type {
		case PackageType:
			// Only the first olm.package property should be retained
			if foundPkg {
				continue
			}
			foundPkg = true
		}

		joined = append(joined, p)
	}

	return joined
}

// expandProperties expands bundle properties into related bundle fields.
// e.g. olm.package -> bundle.PackageName
func expandProperties(bundle *Bundle) error {
	// Sort properties into type buckets and validate
	byType := map[string][]Property{}
	for _, p := range bundle.Properties {
		byType[p.Type] = append(byType[p.Type], p)
	}

	// Validate the package property
	ofPkg := byType[PackageType]
	if len(ofPkg) < 1 {
		return fmt.Errorf("missing package property")
	}
	if len(ofPkg) > 1 {
		return fmt.Errorf("too many package properties, must specify exactly one")
	}

	var pkg PackageProperty
	if err := json.Unmarshal([]byte(ofPkg[0].Value), &pkg); err != nil {
		return fmt.Errorf("failed to unmarshal package property: %s", err)
	}

	// The package property's version must match the bundle version (from the CSV)
	version, err := bundle.Version()
	if err != nil {
		return err
	}
	if pkg.Version != version {
		return fmt.Errorf("package property version does not match bundle csv version: (%s != %s)", pkg.Version, version)
	}

	bundle.Package = pkg.PackageName

	return nil
}
