package idgenutil

import "github.com/sony/sonyflake"

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{
	MachineID: func() (u uint16, err error) {
		return 1024, nil
	},
	CheckMachineID: func(u uint16) bool {
		return true
	},
})

func NextId() (uint64, error) {
	return idGenerator.NextID()
}
