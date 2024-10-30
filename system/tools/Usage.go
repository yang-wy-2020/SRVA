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
    memory
        - memory view
    disk
        - memory view
    disk
        - disk view
    all 
        - all resource view
    edit
        - edit config
e.g.
./main $(args1)
		`)
}
