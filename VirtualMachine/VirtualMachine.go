package VirtualMachine

type VirtualMachine struct {
	memory   map[int16]uint16
	register [8]uint16
	stack    []uint16
}

func Initialize() VirtualMachine {
	return VirtualMachine{
		memory:   make(map[int16]uint16),
		register: [8]uint16{},
		stack:    make([]uint16, 32),
	}
}

func (vm *VirtualMachine) Execute(bytes *[]byte) {

}
