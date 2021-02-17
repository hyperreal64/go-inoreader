package main

import "flag"

func main() {

	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	subCmd := flag.NewFlagSet("sub", flag.ExitOnError)

}
