package semver_test

import (
	"fmt"
	"sort"

	"pkg.package-operator.run/semver"
)

func ExampleMustNewVersion() {
	v := semver.MustNewVersion("1.2.4-alpha.0+meta")

	fmt.Println(v.Major, v.Minor, v.Patch, v.PreRelease, v.BuildMetadata)
	// Output: 1 2 4 alpha.0 [meta]
}

func ExampleNewVersion() {
	v, err := semver.NewVersion("1.2.4-alpha.0+meta")
	if err != nil {
		panic(err)
	}

	fmt.Println(v.Major, v.Minor, v.Patch, v.PreRelease, v.BuildMetadata)
	// Output: 1 2 4 alpha.0 [meta]
}

func ExampleNewConstraint_version() {
	constraint := "1.0.0 - 2.0.0"
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		panic(err)
	}

	notContained := semver.MustNewVersion("3.0.0")
	contained := semver.MustNewVersion("1.2.0")

	fmt.Printf("%s is contained in range %q: %v\n", contained.String(), constraint, c.Check(contained))
	fmt.Printf("%s is contained in range %q: %v\n", notContained.String(), constraint, c.Check(notContained))
	// Output:
	// 1.2.0 is contained in range "1.0.0 - 2.0.0": true
	// 3.0.0 is contained in range "1.0.0 - 2.0.0": false
}

func ExampleNewConstraint_range() {
	constraint := "1.0.0 - 2.0.0"
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		panic(err)
	}

	notContained := semver.MustNewConstraint("2.0.0 - 3.0.0")
	contained := semver.MustNewConstraint("1.0.0 - 1.4.0")

	fmt.Printf("1.0.0 - 1.4.0 is contained in range %q: %v\n", constraint, c.Contains(contained))
	fmt.Printf("2.0.0 - 3.0.0 is contained in range %q: %v\n", constraint, c.Contains(notContained))
	// Output:
	// 1.0.0 - 1.4.0 is contained in range "1.0.0 - 2.0.0": true
	// 2.0.0 - 3.0.0 is contained in range "1.0.0 - 2.0.0": false
}

func ExampleAscending() {
	versions := []semver.Version{
		semver.MustNewVersion("1.2.4"),
		semver.MustNewVersion("1.2.3"),
		semver.MustNewVersion("1.0.0"),
		semver.MustNewVersion("1.3.0"),
		semver.MustNewVersion("2.0.0"),
		semver.MustNewVersion("0.4.2"),
	}

	sort.Sort(semver.Ascending(versions))

	fmt.Println(semver.VersionList(versions).String())
	// Output: 0.4.2, 1.0.0, 1.2.3, 1.2.4, 1.3.0, 2.0.0
}

func ExampleDescending() {
	versions := []semver.Version{
		semver.MustNewVersion("1.2.4"),
		semver.MustNewVersion("1.2.3"),
		semver.MustNewVersion("1.0.0"),
		semver.MustNewVersion("1.3.0"),
		semver.MustNewVersion("2.0.0"),
		semver.MustNewVersion("0.4.2"),
	}

	sort.Sort(semver.Descending(versions))

	fmt.Println(semver.VersionList(versions).String())
	// Output: 2.0.0, 1.3.0, 1.2.4, 1.2.3, 1.0.0, 0.4.2
}
