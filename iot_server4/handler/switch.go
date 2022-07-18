package handler

import (
	"iot_server4/bean"
	"iot_server4/command"
	"iot_server4/config"
	"iot_server4/model"
	"sort"
	"strconv"
	"strings"
	"time"

	g "github.com/GramYang/gylog"
	"github.com/davyxu/cellnet"

	sc "iot_server4/sqlx_client"
)

func handler9(msg *bean.SetTimer, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "set timer operation invalid",
		})
		ses.Close()
		return
	}
	//设置定时器
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		switch msg.Cycle {
		case 0:
			command.Switch_SetTimer_Once(msg.DeviceId, msg.Product, uint8(msg.Addr), uint8(msg.Group), uint8(msg.Task), uint8(msg.State), msg.Timestamp)
		case 1:
			command.Switch_SetTimer_Daily(msg.DeviceId, msg.Product, uint8(msg.Addr), uint8(msg.Group), uint8(msg.Task), uint8(msg.State), uint8(msg.Hour), uint8(msg.Minute))
		case 2:
			command.Switch_SetTimer_Weekly(msg.DeviceId, msg.Product, uint8(msg.Addr), uint8(msg.Group), uint8(msg.Task), uint8(msg.State), uint8(msg.WeekDay), uint8(msg.Hour), uint8(msg.Minute))
		case 3:
			command.Switch_SetTimer_Monthly(msg.DeviceId, msg.Product, uint8(msg.Addr), uint8(msg.Group), uint8(msg.Task), uint8(msg.State), uint8(msg.Day), uint8(msg.Hour), uint8(msg.Minute))
		}
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler10(msg *bean.SwitchElectricLeakageTest, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch electric leakage test operation invalid",
		})
		ses.Close()
		return
	}
	//漏电测试
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.SwitchElectricLeakageTest(msg.DeviceId, msg.Product, uint8(msg.Addr))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler11(msg *bean.SwitchAlarmEnable, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch alarm enable operation invalid",
		})
		ses.Close()
		return
	}
	//预警开启
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		var list []uint8
		for _, v := range msg.Enables {
			list = append(list, uint8(v))
		}
		command.SwitchAlarmEnableTotal(msg.DeviceId, msg.Product, uint8(msg.Addr), list)
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler12(msg *bean.SwitchErrorEnable, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch error enable operation invalid",
		})
		ses.Close()
		return
	}
	//保护开启
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		var list []uint8
		for _, v := range msg.Enables {
			list = append(list, uint8(v))
		}
		command.SwitchErrorEnableTotal(msg.DeviceId, msg.Product, uint8(msg.Addr), list)
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler13(msg *bean.DeviceElectricQuantity, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Type < 1 || msg.Type > 3 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "device electric quantity operation invalid",
		})
		ses.Close()
		return
	}
	switch msg.Type {

	//日，直接返回
	case 1:
		result, err := sc.GetDeviceDataByDay(msg.DeviceId, msg.Year, msg.Month, msg.Day)
		if err != nil {
			ses.Send(&bean.ResultMessage{
				IsSuccess: false,
				Message:   err.Error(),
			})
			return
		} else {
			ses.Send(&bean.DeviceElectricQuantityResult{
				DeviceId:   result.DeviceId,
				Type:       1,
				PowerQ:     stringToIntSlice(result.PowerQ),
				EnergyPt:   stringToFloatSlice(result.EnergyPt),
				RmsI:       stringToFloatSlice(result.RmsI),
				RecordTime: result.RecordTime.UnixMilli(),
			})
		}

	//月，获取指定月的数据然后计算差值得到月用电量
	case 2:
		results, err := sc.GetDeviceDataByMonth(msg.DeviceId, msg.Year, msg.Month)
		if err != nil {
			ses.Send(&bean.ResultMessage{
				IsSuccess: false,
				Message:   err.Error(),
			})
			return
		} else {
			deqr := bean.DeviceElectricQuantityResult{
				DeviceId: results[0].DeviceId,
				Type:     2,
				Year:     msg.Year,
				Month:    msg.Month,
			}
			//处理统计数据
			handlerDeviceData(results, &deqr)
			ses.Send(&deqr)
		}

	//年
	case 3:
		results, err := sc.GetDeviceDataByYear(msg.DeviceId, msg.Year)
		if err != nil {
			ses.Send(&bean.ResultMessage{
				IsSuccess: false,
				Message:   err.Error(),
			})
			return
		} else {
			deqr := bean.DeviceElectricQuantityResult{
				DeviceId: results[0].DeviceId,
				Type:     2,
				Year:     msg.Year,
			}
			//处理统计数据
			handlerDeviceData(results, &deqr)
			ses.Send(&deqr)
		}
	}
}

func stringToIntSlice(s string) []int {
	ss := strings.Split(s, ",")
	var res []int
	if len(ss) == 0 {
		return res
	}
	for _, v := range ss {
		a, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return res
		}
		res = append(res, int(a))
	}
	return res
}

func stringToFloatSlice(s string) []float32 {
	ss := strings.Split(s, ",")
	var res []float32
	if len(ss) == 0 {
		return res
	}
	for _, v := range ss {
		a, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return res
		}
		res = append(res, float32(a))
	}
	return res
}

func handlerDeviceData(datas []sc.DeviceDataByDay, deqr *bean.DeviceElectricQuantityResult) {
	var matrix1 [][]int
	var matrix2, matrix3 [][]float32
	var res1 []int
	var res2, res3 []float32
	for _, v := range datas {
		matrix1 = append(matrix1, stringToIntSlice(v.PowerQ))
		matrix2 = append(matrix2, stringToFloatSlice(v.EnergyPt))
		matrix3 = append(matrix3, stringToFloatSlice(v.RmsI))
	}
	var maxLength = 0
	for i := 0; i < len(matrix1); i++ {
		if maxLength < len(matrix1[i]) {
			maxLength = len(matrix1[i])
		}
	}
	var reverseMatrix1 = make([][]int, maxLength)
	var reverseMatrix2 = make([][]float32, maxLength)
	var reverseMatrix3 = make([][]float32, maxLength)
	for i := 0; i < len(matrix1); i++ {
		for j := 0; j < len(matrix1[i]); j++ {
			reverseMatrix1[j] = append(reverseMatrix1[j], matrix1[i][j])
			reverseMatrix2[j] = append(reverseMatrix2[j], matrix2[i][j])
			reverseMatrix3[j] = append(reverseMatrix3[j], matrix3[i][j])
		}
	}
	for i := 0; i < len(reverseMatrix1); i++ {
		res1 = append(res1, intSortAndMinus(reverseMatrix1[i]))
		res2 = append(res2, floatSortAndMinus(reverseMatrix2[i]))
		res3 = append(res3, floatSortAndMinus(reverseMatrix3[i]))
	}
	deqr.PowerQ = res1
	deqr.EnergyPt = res2
	deqr.RmsI = res3
}

func intSortAndMinus(v []int) int {
	if len(v) == 0 {
		return 0
	}
	if len(v) == 1 {
		return v[0]
	}
	sort.Slice(v, func(i, j int) bool {
		return v[i] < v[j]
	})
	return v[len(v)-1] - v[0]
}

func floatSortAndMinus(v []float32) float32 {
	if len(v) == 0 {
		return 0
	}
	if len(v) == 1 {
		return v[0]
	}
	sort.Slice(v, func(i, j int) bool {
		return v[i] < v[j]
	})
	return v[len(v)-1] - v[0]
}

func handler15(msg *bean.GetSwitchSetting, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "get switch setting invalid",
		})
		ses.Close()
		return
	}
	//请求开关设置数据
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.SwitchSetting(msg.DeviceId, msg.Product, uint8(msg.Addr))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler16(msg *bean.SwitchLoopOn, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch loop on invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("开启设备 %s 开关 %d 轮询\n", msg.DeviceId, msg.Addr)
	model.UpdateDeviceCacheSwitchLoopAndSend(msg.DeviceId, msg.Addr, true, func(product string) {
		command.SwitchRuntime(msg.DeviceId, product, uint8(msg.Addr), uint8(config.Conf.LoopInterval), 255)
	}, config.Conf.CmdInterval)
}

func handler17(msg *bean.SwitchLoopOff, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch loop off invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("关闭设备 %s 开关 %d 轮询\n", msg.DeviceId, msg.Addr)
	model.UpdateDeviceCacheSwitchLoopAndSend(msg.DeviceId, msg.Addr, false, func(product string) {
		command.SwitchRuntime(msg.DeviceId, product, uint8(msg.Addr), uint8(config.Conf.LoopInterval), 0)
	}, config.Conf.CmdInterval)
}

func handler18(msg *bean.SwitchClearFault, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch clear fault invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("故障清除 设备 %s 开关 %d\n", msg.DeviceId, msg.Addr)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.SwitchClearCurrentError(msg.DeviceId, msg.Product, uint8(msg.Addr))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler19(msg *bean.VoltageLimitRstEnable, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch voltage limit reset enable invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("过欠压恢复设置 设备 %s 开关 %d 启动 %t\n", msg.DeviceId, msg.Addr, msg.Enable)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		if msg.Enable {
			command.SwitchVolLimitRst0(msg.DeviceId, msg.Product, uint8(msg.Addr))
		} else {
			command.SwitchVolLimitRst1(msg.DeviceId, msg.Product, uint8(msg.Addr))
		}
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler20(msg *bean.SetIHP, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set IH_P invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置IH_P 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_IH_P(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler21(msg *bean.SetIH, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set IH invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置IH 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_IH(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler22(msg *bean.SetUHP, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set UH_P invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置UH_P 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_UH_P(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler23(msg *bean.SetUH, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set UH invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置UH 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_UH(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler24(msg *bean.SetULP, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set UL_P invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置UL_P 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_UL_P(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler25(msg *bean.SetUL, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set UL invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置UL 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_UL(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler26(msg *bean.SetPHP, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set PH_P invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置PH_P 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_PH_P(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler27(msg *bean.SetPH, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set PH invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置PH 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_PH(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler28(msg *bean.SetEHP, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set EH_P invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置EH_P 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_EH_P(msg.DeviceId, msg.Product, uint8(msg.Addr), uint32(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler29(msg *bean.SetEH, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 || msg.Index == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set EH invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置EH 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_EH(msg.DeviceId, msg.Product, uint8(msg.Addr), uint8(msg.Index), msg.Value)
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler30(msg *bean.SetILP, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set IL_P invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置IL_P 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_IL_P(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler31(msg *bean.SetIL, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set IL invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置IL 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_IL(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler32(msg *bean.SetTHP, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set TH_P invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置TH_P 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_TH_P(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler33(msg *bean.SetTH, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set TH invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置TH 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_TH(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler34(msg *bean.SetUHLCT, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set UHL_CT invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置UHL_CT 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_UHL_CT(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler35(msg *bean.SetUHLRT, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set UHL_RT invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置UHL_RT 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_UHL_RT(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}

func handler36(msg *bean.SetIHPHCT, ses cellnet.Session) {
	ok := handleHeartbeat(&msg.Heartbeat, ses)
	if !ok {
		return
	}
	if msg.DeviceId == "" || msg.Product == "" || msg.Addr == 0 {
		ses.Send(&bean.Shutdown{
			Code:    4,
			Message: "switch set IH_PH_CT invalid",
		})
		ses.Close()
		return
	}
	g.Debugf("设置IH_PH_CT 设备 %s 开关 %d 值 %d\n", msg.DeviceId, msg.Addr, msg.Value)
	ok, msg1 := model.CheckDeviceCacheAndSessionAndSend(msg.DeviceId, time.Now().UnixMilli(), func() {
		command.Switch_IH_PH_CT(msg.DeviceId, msg.Product, uint8(msg.Addr), uint16(msg.Value))
	}, config.Conf.CmdInterval)
	if !ok {
		ses.Send(&bean.DeviceResultMessage{
			IsSuccess: ok,
			DeviceId:  msg.DeviceId,
			Message:   msg1,
		})
	}
}
