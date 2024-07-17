package spdx

import (
	"fmt"
	"strings"
)

func DiffPackages(filename string, docA, docB *Document) (string, error) {
	builder := &strings.Builder{}
	fmt.Println(builder, "diffing file: "+filename)
	aNotInB := []*Package{}

	packagesToCheck := []packagesToCheckRelations{}
	for _, pkgA := range docA.Packages {
		found := false
		for _, pkgB := range docB.Packages {
			if pkgA.Name == pkgB.Name {
				found = true
				packagesToCheck = append(packagesToCheck, packagesToCheckRelations{PackageA: pkgA, PackageB: pkgB})
				break
			}
		}
		if !found {
			aNotInB = append(aNotInB, pkgA)
		}
	}

	if len(aNotInB) < 1 {
		fmt.Println(builder, "No missing packages")
	} else {
		fmt.Println(builder, "missing pacakges")
		fmt.Println(builder, packageNames(aNotInB))
	}

	for _, pkChk := range packagesToCheck {
		diffRelations := pkChk.diffRelations()
		if len(diffRelations) > 0 {
			fmt.Println(builder, "Missing relations")
			for _, rel := range diffRelations {
				fmt.Println(builder, rel)
			}
		}
	}

	return builder.String(), nil
}

func packageNames(packages []*Package) []string {
	names := []string{}
	for _, pkg := range packages {
		names = append(names, pkg.Name)
	}
	return names
}

type packagesToCheckRelations struct {
	PackageA *Package
	PackageB *Package
}

func (p *packagesToCheckRelations) Name() string {
	return p.PackageA.Name
}

func (p *packagesToCheckRelations) diffRelations() []string {
	missingRelations := []string{}

	for _, rel := range p.PackageA.Relationships {
		_, ok := rel.Peer.(*Package)
		if !ok {
			continue
		}
		found := false

		nameA := rel.Peer.(*Package).drawName(&DrawingOptions{Version: true})
		for _, relB := range p.PackageB.Relationships {
			_, ok := relB.Peer.(*Package)
			if !ok {
				continue
			}
			nameB := relB.Peer.(*Package).drawName(&DrawingOptions{Version: true})
			if nameA == nameB {
				found = true
				break
			}

		}
		if !found {
			missingRelations = append(missingRelations, rel.Peer.(*Package).drawName(&DrawingOptions{Version: true, Purls: true}))
		}

	}
	return missingRelations
}
