package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/llbec/robotech/squadrons/ethsquadron"
	"github.com/llbec/robotech/utils"
	"github.com/xuri/excelize/v2"
)

func init() {
	initCmd()
}

var (
	gEth *ethsquadron.EthSquadron
)

func main() {
	parseCmd()

	if err := utils.LocalEnv(); err != nil {
		fmt.Printf("env file error: %v\n", err)
		os.Exit(-1)
	}

	/*workDir, e := os.Getwd()
	if e != nil {
		fmt.Printf("os.Getwd failed: %v\n", e)
		os.Exit(-2)
	}

	d := daemon.NewDaemon(6960, thread)
	d.Run(filepath.Join(workDir, "run.fmt"))*/
	thread()
}

func thread( /*chSig, chExit chan int*/ ) {
	gEth, _ = ethsquadron.NewEthSquadron(os.Getenv("RPC"))
	fmt.Printf("new rpc at: %v\n", os.Getenv("RPC"))

	if iHeight > 0 {
		os.Exit(avrDiffs(iHeight))
	}
	if iTest > 0 {
		os.Exit(presageDiff(iTest))
	}

	height, err := gEth.BlockNumber()
	if err != nil {
		fmt.Printf("BlockNumber: %v\n", err)
		os.Exit(-3)
	}

	if fPower {
		os.Exit(avrDiffs(height))
	}
	if fOracle {
		os.Exit(presageDiff(height))
	}
	if iLine > 0 {
		os.Exit(outTime(height, iLine))
	}
	if fShedule {
		os.Exit(schedule(height))
	}
	if fBomb {
		os.Exit(bombDiff())
	}
	if fDownload {
		os.Exit(downloadHeader(height))
	}

	help()
}

//2
func presageDiff(height uint64) int {
	block, err := gEth.HeaderByNumber(height)
	if err != nil {
		fmt.Printf("Block(%v): %v\n", height, err)
		return -21
	}
	speed, ok := new(big.Int).SetString(os.Getenv("speed"), 10)
	if !ok {
		fmt.Printf("read speed error\n")
		return -22
	}
	//target, err := strconv.ParseUint(os.Getenv("target"), 10, 64)
	target, ok := new(big.Int).SetString(os.Getenv("target"), 10)
	if !ok {
		fmt.Printf("read target error\n")
		return -23
	}
	cur := gEth.MakeDiffHeader(block.Time, height, block.Difficulty)
	blockPeriod := new(big.Int).Div(block.Difficulty, speed)
	var (
		now  uint64
		diff *big.Int
	)
	max := blockPeriod
	for i := uint64(1); ; i++ {
		now = cur.Time + blockPeriod.Uint64()
		if blockPeriod.Cmp(target) > 0 {
			break
		}
		diff = gEth.CalcDifficulty(now, cur)
		cur = gEth.MakeDiffHeader(now, height+i, diff)
		blockPeriod = new(big.Int).Div(diff, speed)
		if blockPeriod.Cmp(max) > 0 {
			max = blockPeriod
			fmt.Printf("block %v at %v,next period %v\n", cur.Number, formatTime(int64(now)), blockPeriod)
		}
	}
	fmt.Printf("block %v at %v, next period is %v\n", cur.Number, formatTime(int64(now)), blockPeriod)
	return 0
}

func formatTime(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

//1
func avrDiffs(height uint64) int {
	//chain := make(map[uint64]*types.Header)
	vDiff := big.NewInt(0)
	vCount := big.NewInt(0)
	vDelta := big.NewInt(0)
	period, err := strconv.ParseUint(os.Getenv("Days"), 10, 64)
	if err != nil {
		fmt.Printf("read Days error: %v", err)
		return -11
	}
	blockPeriod, err := strconv.ParseUint(os.Getenv("blockPeriod"), 10, 64)
	if err != nil {
		fmt.Printf("read blockPeriod error: %v", err)
		return -12
	}
	period = period * 24 * 60 * 60
	if height < (period / blockPeriod) {
		fmt.Printf("height %v is too small\n", height)
		return -13
	}
	start := height - (period / blockPeriod)
	preBlk, err := gEth.HeaderByNumber(start - 1)
	if err != nil {
		fmt.Printf("Block(%v): %v\n", start-1, err)
		return -14
	}
	preTime := preBlk.Time
	for i := start; i <= height; i++ {
		block, err := gEth.HeaderByNumber(i)
		if err != nil {
			fmt.Printf("Block(%v): %v\n", i, err)
			continue
		}
		vDiff.Add(vDiff, block.Difficulty)
		vCount.Add(vCount, big.NewInt(1))
		vDelta.Add(vDelta, new(big.Int).SetUint64(block.Time-preTime))
		preTime = block.Time
		//chain[i] = block
	}
	vDiff.Div(vDiff, vCount)
	vDelta.Div(vDelta, vCount)
	fmt.Printf("from %v total %v blocks, average difficulty is %v, average period is %v\n", height, vCount, vDiff, vDelta)
	bp, ok := new(big.Int).SetString(os.Getenv("blockPeriod"), 10)
	if ok {
		fmt.Printf("%v diff per second at %v\n", new(big.Int).Div(vDiff, bp), bp)
	}
	fmt.Printf("average %v diff per second\n", new(big.Int).Div(vDiff, vDelta))
	return 0
}

//3
func outTime(height, line uint64) int {
	period, err := strconv.ParseUint(os.Getenv("Days"), 10, 64)
	if err != nil {
		fmt.Printf("read Days error: %v", err)
		return -31
	}
	blockPeriod, err := strconv.ParseUint(os.Getenv("blockPeriod"), 10, 64)
	if err != nil {
		fmt.Printf("read blockPeriod error: %v", err)
		return -32
	}
	period = period * 24 * 60 * 60
	if height < (period / blockPeriod) {
		fmt.Printf("height %v is too small\n", height)
		return -33
	}
	start := height - (period / blockPeriod)
	preBlk, err := gEth.HeaderByNumber(start - 1)
	if err != nil {
		fmt.Printf("Block(%v): %v\n", start-1, err)
		return -34
	}
	for i := start; i <= height; i++ {
		block, err := gEth.HeaderByNumber(i)
		if err != nil {
			fmt.Printf("Block(%v): %v\n", i, err)
			return -35
		}
		if delta := block.Time - preBlk.Time; delta >= line {
			fmt.Printf("block %v %v, use %v s\n", i, formatTime(int64(block.Time)), delta)
		}
		preBlk = block
	}
	return 0
}

//4
//从 14-900-000 开始，进行先统计，再推算
func schedule(height uint64) int {
	curH := uint64(14900000)
	step := uint64(10000)
	vDiff := big.NewInt(0)
	vCount := big.NewInt(0)
	vTime := big.NewFloat(0)
	vBomb := big.NewInt(0)

	preBlk, err := gEth.HeaderByNumber(curH - 1)
	if err != nil {
		fmt.Printf("Block(%v): %v\n", curH-1, err)
		return -41
	}
	preTime := preBlk.Time
	for curH <= height {
		for i := uint64(0); i < step; i++ {
			block, err := gEth.HeaderByNumber(curH)
			if err != nil {
				fmt.Printf("Block(%v): %v\n", curH, err)
				return -42
			}
			vBomb.Add(vBomb, gEth.BombDiff(new(big.Int).SetUint64(curH)))
			vDiff.Add(vDiff, block.Difficulty)
			vCount.Add(vCount, big.NewInt(1))
			vTime.Add(vTime, new(big.Float).Mul(new(big.Float).SetUint64(block.Time-preTime), big.NewFloat(1000)))
			preTime = block.Time
			curH += 1
			if curH > height {
				break
			}
		}
		vDiff.Div(vDiff, vCount)
		vTime.Quo(vTime, new(big.Float).SetInt(vCount))
		vBomb.Div(vBomb, vCount)
		fmt.Printf(
			"from %v total %v blocks, average difficulty is %v, average period is %v s\n",
			curH-step,
			vCount,
			vDiff,
			new(big.Float).Quo(vTime, big.NewFloat(1000)),
		)
		fmt.Printf("average %v diff per mille second\n", new(big.Float).Quo(new(big.Float).SetInt(vDiff), vTime))
		fmt.Printf("average bomb is %v\n", vBomb)
	}
	return 0
}

//5
func bombDiff() int {
	start := big.NewInt(13_773_000)
	end, ok := new(big.Int).SetString(os.Getenv("end"), 10)
	if !ok {
		fmt.Printf("read end error\n")
		return -51
	}

	maxDiff := big.NewInt(0)
	f := excelize.NewFile()

	// 创建一个工作表
	//index := f.NewSheet("Sheet2")
	// 设定单元格的值
	f.SetCellValue("Sheet1", "B1", "Difficulty")
	index := 2

	for i := start; i.Cmp(end) < 0; i.Add(i, big.NewInt(1)) {
		diff := gEth.BombDiff(i)
		if diff.Cmp(maxDiff) > 0 {
			maxDiff = diff
			fmt.Printf("block %v, bomb difficulty is %v\n", i, diff)
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", index), fmt.Sprintf("%v", i))
			f.SetCellValue("Sheet1", fmt.Sprintf("B%d", index), fmt.Sprintf("%v", diff))
			index++
		}
	}

	if err := f.SaveAs("BombDifficulty.xlsx"); err != nil {
		fmt.Println(err)
		return -52
	}

	return 0
}

//6
func downloadHeader(heights ...uint64) int {
	var (
		start uint64
	)

	end := heights[0]
	if len(heights) > 1 {
		start = heights[1]
	} else {
		start = 13_773_000
	}

	cur := start
	step := uint64(100_000)
	preBlk, err := gEth.HeaderByNumber(cur - 1)
	if err != nil {
		fmt.Printf("Block(%v): %v\n", cur-1, err)
		return -61
	}

	f := excelize.NewFile()
	sheetId := 1

	for cur <= end {
		if sheetId != 1 {
			f.NewSheet(fmt.Sprintf("Sheet%v", sheetId))
		}
		f.SetCellValue(fmt.Sprintf("Sheet%d", sheetId), "A1", "块高")
		f.SetCellValue(fmt.Sprintf("Sheet%d", sheetId), "B1", "块时间")
		f.SetCellValue(fmt.Sprintf("Sheet%d", sheetId), "C1", "间隔时间")
		f.SetCellValue(fmt.Sprintf("Sheet%d", sheetId), "D1", "难度")
		f.SetCellValue(fmt.Sprintf("Sheet%d", sheetId), "E1", "调整难度")
		f.SetCellValue(fmt.Sprintf("Sheet%d", sheetId), "F1", "炸弹难度")
		f.SetCellValue(fmt.Sprintf("Sheet%d", sheetId), "G1", "是否有叔块")
		round := cur/step*step + step
		for i := 2; cur < round; cur += 1 {
			if header, err := gEth.HeaderByNumber(cur); err != nil {
				fmt.Printf("Block(%v): %v\n", cur, err)
			} else {
				//block number
				f.SetCellValue(
					fmt.Sprintf("Sheet%d", sheetId),
					fmt.Sprintf("A%d", i),
					fmt.Sprintf("%v", header.Number),
				)
				//block timestamp
				f.SetCellValue(
					fmt.Sprintf("Sheet%d", sheetId),
					fmt.Sprintf("B%d", i),
					fmt.Sprintf("%v", formatTime(int64(header.Time))),
				)
				//time - pretime
				f.SetCellValue(
					fmt.Sprintf("Sheet%d", sheetId),
					fmt.Sprintf("C%d", i),
					fmt.Sprintf("%v", header.Time-preBlk.Time),
				)
				//difficulty
				f.SetCellValue(
					fmt.Sprintf("Sheet%d", sheetId),
					fmt.Sprintf("D%d", i),
					fmt.Sprintf("%v", header.Difficulty),
				)
				//difficulty index
				f.SetCellValue(
					fmt.Sprintf("Sheet%d", sheetId),
					fmt.Sprintf("E%d", i),
					fmt.Sprintf("%v", gEth.DifficultyIndex(header.Time, preBlk)),
				)
				//bomb difficulty
				f.SetCellValue(
					fmt.Sprintf("Sheet%d", sheetId),
					fmt.Sprintf("F%d", i),
					fmt.Sprintf("%v", gEth.BombDiff(header.Number)),
				)
				//uncle block
				f.SetCellValue(
					fmt.Sprintf("Sheet%d", sheetId),
					fmt.Sprintf("G%d", i),
					fmt.Sprintf("%v", gEth.HasUncle(header.UncleHash)),
				)
			}
			i += 1
		}
		sheetId += 1
	}

	return 0
}
