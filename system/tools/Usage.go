package tools

import "fmt"

func Usage() {
	fmt.Println(`args1:
    cpu
        - cpu view
    io 
        - io view
    system
        - cpu total view 
    network
        - network card view
    disk
        - disk monitor view
    memory
        - memory view
    edit
        - edit config
e.g.
sudo qomolo-sys-analysis $(args1)
		`)
}
