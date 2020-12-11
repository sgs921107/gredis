/*************************************************************************
	> File Name: commands.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月11日 星期五 21时44分43秒
*************************************************************************/

package gredis

type ScansCmd struct {
	scanCmds []*ScanCmd
}

func (c *ScansCmd) Val() (keys []string) {
	for _, scanCmd := range c.scanCmds {
		page, _ := scanCmd.Val()
		keys = append(keys, page...)
	}
	return keys
}

func (c *ScansCmd) Result() (keys []string, err error) {
	for _, scanCmd := range c.scanCmds {
		page, _, e := scanCmd.Result()
		if e != nil {
			err = e
		}
		keys = append(keys, page...)
	}
	return keys, err
}

func (c *ScansCmd) addScanCmd(scanCmd *ScanCmd) {
	c.scanCmds = append(c.scanCmds, scanCmd)
}

func NewScansCmd() *ScansCmd {
	return &ScansCmd{}
}
