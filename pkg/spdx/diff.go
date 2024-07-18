package spdx

import (
	"fmt"
	"sort"
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
			if stripPrefix("/", pkgA.Name) == stripPrefix("/", pkgB.Name) {
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
		missingRelations := pkChk.diffRelations()
		if len(missingRelations) > 0 {
			fmt.Println(builder, "the other sbom is Missing relations:")
			for _, rel := range missingRelations {
				if rel.hasPartialMatchs {
					fmt.Println(builder, "partialMatch: "+rel.purl+"->"+strings.Join(rel.partialPurls, ","))
				} else {
					fmt.Println(builder, "no match: "+rel.purl)
				}
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

type missingRelation struct {
	purl             string
	hasPartialMatchs bool
	partialPurls     []string
}

func (p *packagesToCheckRelations) diffRelations() []missingRelation {
	missingRelations := []missingRelation{}

	for _, rel := range p.PackageA.Relationships {
		partialMatches := []string{}
		relPkg, ok := rel.Peer.(*Package)
		if !ok {
			continue
		}
		found := false
		partial := false
		for _, relB := range p.PackageB.Relationships {
			pkg, ok := relB.Peer.(*Package)
			if !ok {
				continue
			}
			if relPkg.Purl().ToString() == pkg.Purl().ToString() {
				found = true
				break
			}
			// fmt.Println(found)
			if strings.Contains(relPkg.Name, pkg.Name) {
				partial = true
				partialMatches = append(partialMatches, pkg.Purl().ToString())
			}

		}
		if !found {
			missingRelations = append(
				missingRelations,
				missingRelation{purl: rel.Peer.(*Package).Purl().ToString(), hasPartialMatchs: partial, partialPurls: partialMatches})
		}

	}
	sort.Slice(missingRelations, func(i, j int) bool {
		return missingRelations[i].hasPartialMatchs
	})

	return missingRelations
}

func extractName(name string) string {
	parts := strings.Split(name, "/")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}

	return name
}

func stripPrefix(seperator, name string) string {
	parts := strings.Split(name, seperator)
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return name
}

func stripPurlV(purl string) string {
	parts := strings.Split(purl, "@")
	if len(parts) > 1 {
		s := parts[0] + strings.TrimPrefix(parts[1], "v")
		return trimStringFromHash(s)
	}
	return purl
}

func trimStringFromHash(s string) string {
	if idx := strings.Index(s, "#"); idx != -1 {
		return s[:idx]
	}
	return s
}
