package upgrade

import (
	"bytes"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/util"
)

func Run(log *util.Logger) error {
	return upgradeFromFs(log, &afero.Afero{Fs: afero.NewOsFs()})
}

func upgradeFromFs(log *util.Logger, afs *afero.Afero) error {
	conf, err := config.GetRawFromFS(afs)
	if err != nil {
		return err
	}

	if conf.ApiVersion == config.LatestVersion {
		log.Success("Configuration already at version", conf.ApiVersion)
		return nil
	}

	if conf.ApiVersion != 1 {
		return errors.Errorf("Unsupported version %d, unable to upgrade to %d", conf.ApiVersion, config.LatestVersion)
	}

	for _, app := range conf.Apps {
		basePaths := findBasePaths(log, afs, conf.Targets)
		for _, base := range basePaths {
			err := moveBaseFiles(log, afs, &app, base)
			if err != nil {
				return err
			}
		}

		for _, target := range conf.Targets {
			renameTargetManifest(log, afs, &app, &target)
			deleteNamespace(log, afs, &app, &target)
		}
	}

	err = updateVersion(log, afs)
	if err != nil {
		return err
	}

	log.Success("Successfully upgraded configuration to version", config.LatestVersion)
	return nil
}

func findBasePaths(log *util.Logger, afs *afero.Afero, targets []config.Target) []string {
	var paths []string
	m := map[string]bool{}

	for _, target := range targets {
		path := filepath.Dir(target.Path)
		if !m[path] {
			m[path] = true
			paths = append(paths, path)
		}
	}

	return paths
}

func moveBaseFiles(log *util.Logger, afs *afero.Afero, app *config.App, path string) error {
	baseOld := filepath.Join(app.Path, path)
	baseNew := filepath.Join(app.Path, path, "_base")

	entries, err := afs.ReadDir(baseOld)
	if err != nil {
		return err
	}

	for _, file := range entries {
		name := file.Name()
		pathOld := filepath.Join(baseOld, name)
		pathNew := filepath.Join(baseNew, name)

		if name == "kube-manifest.yaml" {
			pathNew = filepath.Join(baseNew, "kustomization.yaml")
		}

		afs.Mkdir(baseNew, 0755)

		if !file.IsDir() && filepath.Ext(name) == ".yaml" {
			err := afs.Rename(pathOld, pathNew)
			if err == nil {
				log.Note("Renamed:", pathOld, "->", pathNew)
			} else {
				log.Error("Could not rename:", err)
			}
		}
	}

	return nil
}

func deleteNamespace(log *util.Logger, afs *afero.Afero, app *config.App, target *config.Target) {
	namespacePath := filepath.Join(app.Path, target.Path, "namespace.yaml")

	if exists, err := afs.Exists(namespacePath); exists && err == nil {
		err := afs.Remove(namespacePath)
		if err == nil {
			log.Note("Removed:", namespacePath)
		} else {
			log.Error("Could not remove:", err)
		}
	}
}

func renameTargetManifest(log *util.Logger, afs *afero.Afero, app *config.App, target *config.Target) {
	manifestOld := filepath.Join(app.Path, target.Path, "kube-manifest.yaml")
	manifestNew := filepath.Join(app.Path, target.Path, "kustomization.yaml")

	if exists, err := afs.Exists(manifestNew); exists && err == nil {
		log.Warn("File exists:", manifestNew)
	} else if exists, err := afs.Exists(manifestOld); exists && err == nil {
		err := afs.Rename(manifestOld, manifestNew)
		if err != nil {
			log.Error("Could not rename:", err)
			return
		}

		log.Note("Renamed:", manifestOld, "->", manifestNew)
		yml, err := afs.ReadFile(manifestNew)
		if err != nil {
			log.Error("Could not modify:", err)
			return
		}

		yml = bytes.Replace(yml, []byte("- ..\n"), []byte("- ../_base\n"), -1)
		yml = bytes.Replace(yml, []byte("- namespace.yaml\n"), []byte{}, -1)

		err = afs.WriteFile(manifestNew, yml, 0644)
		if err != nil {
			log.Error("Could not modify:", err)
			return
		}
	} else {
		log.Warn("Could not find:", manifestOld)
	}
}

func updateVersion(log *util.Logger, afs *afero.Afero) error {
	yml, err := afs.ReadFile(config.ConfigName)
	if err != nil {
		return err
	}

	yml = bytes.Replace(yml, []byte("version: 1\n"), []byte("version: 2\n"), -1)

	err = afs.WriteFile(config.ConfigName, yml, 0644)
	if err != nil {
		return err
	}

	log.Note("Updated:", config.ConfigName)
	return nil
}
