/*
One or two sentence synopsis of the package...

# Overview

One or two paragraph overview of the package...
(This page describes the nature of the individual package.)

More information at https://github.com/senzing-garage/demo-quickstart

# Another Header

Details of the package...
Lorem ipsum dolor sit amet, consectetur adipiscing elit...

# Examples

The examples given here should be specific to the package.

Examples of use can be seen in the examplepackage_test.go files.

	package main
	import (
		fmt

		"github.com/senzing-garage/demo-quickstart/examplepackage"
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
	}
*/
package httpserver
