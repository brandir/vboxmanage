/* Time-stamp: <2020-04-27 19:20:28 (mellon@macbook) /Users/mellon/projects/vboxmanage/vboxmanage.go>
 * vboxmanage.go - Program to query VirtualBox
 *
 * vboxmanage project, created 04/23/2020
 *
 * https://github.com/brandir/vboxmanage
 *
 * vboxmanage provides the following information 
 *
 *   - VM start date if running
 *   - VM last running time
 *   - VM resource consumption (disk)
 *   - VM configuration (os, #cpus, network, disk size, ...)
 *   - VirtualBox version information
 *   - VirtualBox last version update(s)
 *   - VirtualBox configuration
 *
 * vboxmanage should run on windows and linux equally.
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	author = "dmon"                                  // The usual suspect ;-)
	version = "1.0"                                  // Initial version
	cdate = "04/23/2020"                             // Development start date
	carch = "macosx"                                 // Development main platform
	program = "vboxmgr"                              // Program Name
	github = "https://github.com/brandir/vboxmanage" // github repo
)

var vmguest = []string{"alpine", "minikube", "nenya", "oneq84", "pluto", "yavanna"}
var vmos = []string{"Alpine Linux", "minikube", "OpenBSD", "Solaris", "Kali Linux", "FreeBSD"}
var vmhost = "MacBook macOS Catalina 10.15.4"

// Get the state of a VM. Output looks like the following, returned from
// 'vboxmanage showvminfo nenya' and 'vboxmanage showvminfo yavanna'.
// Multiple spaces in a row have been deleted.
//   State: powered off (since 2020-04-22T08:02:42.000000000)
//   State: running (since 2020-04-24T16:48:10.825000000)
func GetVMState(vm string) string {
	var res string = ""
	
	f, err := os.Create("./out.tmp")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// cmd := fmt.Sprintf("vboxmanage showvminfo %s | grep State", vm)
	cmd := exec.Command("vboxmanage", "showvminfo", vm)
	cmd.Stdout = f
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
	
	// scan the file for 'State: ...'
	r, _ := regexp.Compile("^State: ")
	g, _ := os.Open("./out.tmp")
	defer g.Close()
	scanner := bufio.NewScanner(g)
	for scanner.Scan() {
		if r.MatchString(scanner.Text()) {
			res := fmt.Sprintf("%s", scanner.Text())
			return res
		}
		if err := scanner.Err(); err != nil {
			res := ""
			return res
		}
	}
	return res
}

// Display the VMs and their Operating System.
func showVMandOS() {
	var i int = 0
	fmt.Printf("%-8s - %s\n", "Name", "Operating System")
	dashes := strings.Repeat("-", 27)
	fmt.Println(dashes)
	for _, vm := range vmguest {
		fmt.Printf("%-8s - %s\n", vm, vmos[i])
		i += 1
	}
}
		
func main() {
	// Initialize flags, the ugly '_ = ..' avoids the Go compiler conplaining
	
	v := flag.Bool("v", false, "Makes vboxmgr verbose during operation"); _ = v
        l := flag.Bool("l", false, "Log vboxmgr output to given logfile"); _ = l
        g := flag.String("g", "yavanna", "Display information for given guest"); _ = g            
        V := flag.String("V", version, "Display vboxmgr version infformation"); _ = V
	a := flag.Bool("a", false, "Display information for all guest vms"); _ = a

	fmt.Printf("--- %s V%s [(c) %s %s)] ---\n", program, version, author, cdate)
	flag.Parse()
	fmt.Println("v:", *v)
	fmt.Println("l:", *l)
	fmt.Println("g:", *g)
	fmt.Println("V:", *V)
	fmt.Println("a:", *a)

	showVMandOS()
	fmt.Printf("%s\n", GetVMState("nenya"))
	fmt.Printf("%s\n", GetVMState("yavanna"))
}

