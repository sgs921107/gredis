/*************************************************************************
	> File Name: commands.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月11日 星期五 21时44分43秒
*************************************************************************/

package gredis

/*
ScansCmd 定义ScansCmd结构类型, 用于scaniter方法的返回值类型
*/
type ScansCmd struct {
	scanCmds []*ScanCmd
}

/*
Val 获取命令执行的结果
*/
func (c *ScansCmd) Val() (keys []string) {
	for _, scanCmd := range c.scanCmds {
		page, _ := scanCmd.Val()
		keys = append(keys, page...)
	}
	return keys
}

/*
Result 获取命令执行的结果
*/
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

/*
String 以字符串形式展示每个命令的执行结果
*/
func (c *ScansCmd) String() (strSlice []string) {
	for _, scanCmd := range c.scanCmds {
		strSlice = append(strSlice, scanCmd.String())
	}
	return strSlice
}

func (c *ScansCmd) addScanCmd(scanCmd *ScanCmd) {
	c.scanCmds = append(c.scanCmds, scanCmd)
}

/*
NewScansCmd 实例化
*/
func NewScansCmd() *ScansCmd {
	return &ScansCmd{}
}
