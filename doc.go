/*
One or two sentence synopsys of the module.  The Senzing XXXX module is...

# Overview

One or two paragraph overview of this module...

(This page describes the nature of the entire Go module or the top-level command, not an individual package.)

More information at https://github.com/senzing-garage/demo-quickstart

# Example Package

The example package gives an example layout for a package.  This includes how
documentation, tests, and examples should be done.  This paragraph (or two)
should provide a brief overview while linking the reader to the documentation
included in the package itself.

More information can be found in the [pkg/github.com/senzing-garage/demo-quickstart/examplepackage] documentation.

# Another Header: Package or other module features...

More details about the module...
Lorem ipsum dolor sit amet, consectetur adipiscing elit...

# Examples

The examples given here should show a holistic view of the module, if appropriate.

Examples of use can be seen in the main_test.go files.

	package main
	import (
		fmt

		"github.com/senzing-garage/demo-quickstart/examplepackage"
		"github.com/senzing-garage/demo-quickstart/anotherpackage"
	)

	func main() {
		ctx := context.TODO()
		testObject := &ExamplePackageImpl{
			Something: "I'm here",
		}
		err := testObject.SaySomething(ctx)
		if err != nil {
			fmt.Println("whoops")
		}
		anotherpackage.DoSomething(ctx)...
	}
*/
package main
