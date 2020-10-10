/*
 @Desc

 @Date 2020-10-10 14:16
 @Author inori
*/
package util

import (
	"fmt"
	"log"
	"os"
	"path"
)

type Processor struct {
	percent int64  //百分比
	cur     int64  //当前进度位置
	total   int64  //总进度
	rate    string //进度条
	graph   string //显示符号
}

func NewProcessor(start, total int64) *Processor {
	return NewProcessorWithGraph(start, total, "█")
}

func NewProcessorWithGraph(start, total int64, graph string) *Processor {
	p := &Processor{}
	p.cur = start
	p.total = total
	p.graph = graph
	p.percent = p.getPercent()
	for i := 0; i < int(p.percent); i += 2 {
		p.rate += p.graph //初始化进度条位置
	}
	return p
}

func (p *Processor) getPercent() int64 {
	return int64(float32(p.cur) / float32(p.total) * 100)
}

func (p *Processor) Play(cur int64) {
	p.cur = cur
	last := p.percent
	p.percent = p.getPercent()
	if p.percent != last && p.percent%2 == 0 {
		p.rate += p.graph
	}
	fmt.Printf("\r[%-50s]%3d%%  %8s/%s", p.rate, p.percent, FormatSize(p.cur), FormatSize(p.total))
}

func (p *Processor) Finish() {
	fmt.Println()
}

type WriteProcessor struct {
	dir     string   //文件目录
	name    string   //文件名
	file    *os.File //文件
	percent int64    //百分比
	cur     int64    //当前进度位置
	total   int64    //总进度
	rate    string   //进度条
	graph   string   //显示符号
}

func NewWriteProcessor(start, total int64, fileName, dir string) *WriteProcessor {
	wc := &WriteProcessor{}
	wc.dir = dir
	wc.name = fileName
	wc.cur = start
	wc.total = total
	wc.graph = "█"
	wc.percent = wc.getPercent()
	for i := 0; i < int(wc.percent); i += 2 {
		wc.rate += wc.graph //初始化进度条位置
	}
	return wc
}
func (wc *WriteProcessor) Start() (err error) {
	file, err := os.Create(path.Join(wc.dir, "/."+wc.name+".process"))
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	wc.file = file
	return nil
}

func (wc *WriteProcessor) getPercent() int64 {
	return int64(float32(wc.cur) / float32(wc.total) * 100)
}

func (wc *WriteProcessor) Write(p []byte) (int, error) {
	n := len(p)
	wc.cur += int64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *WriteProcessor) PrintProgress() {
	last := wc.percent
	wc.percent = wc.getPercent()
	if wc.percent != last && wc.percent%2 == 0 {
		wc.rate += wc.graph
	}
	_, _ = wc.file.Seek(0, 0)
	_, err := wc.file.WriteString(fmt.Sprintf("\r[%-50s]%3d%%  %8s/%s", wc.rate, wc.percent, FormatSize(wc.cur), FormatSize(wc.total)))
	if err != nil {
		log.Printf(err.Error())
	}
}

func (wc *WriteProcessor) Finish() (err error) {
	_ = wc.file.Close()
	if err = os.RemoveAll(path.Join(wc.dir, "/."+wc.name+".process")); err != nil {
		log.Printf(err.Error())
		return err
	}
	return nil
}
