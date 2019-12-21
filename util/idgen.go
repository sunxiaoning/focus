package util

import "github.com/sony/sonyflake"

var IdGenerator = sonyflake.NewSonyflake(sonyflake.Settings{
	MachineID: func() (u uint16, err error) {
		return 1024, nil
	},
	CheckMachineID: func(u uint16) bool {
		return true
	},
})
