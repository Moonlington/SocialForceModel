package main

import "fmt"

type Spacial interface {
	XY() (float64, float64)
}

type EmptyBin[T Spacial] struct {
	data [][][]T
	xmin float64
	xmax float64
	ymin float64
	ymax float64
}

func newEmptyBin[T Spacial](xbins, ybins int, xmin, xmax, ymin, ymax float64) *EmptyBin[T] {
	data := make([][][]T, ybins)
	for i := range data {
		data[i] = make([][]T, xbins)
		for j := range data[i] {
			data[i][j] = []T{}
		}
	}
	return &EmptyBin[T]{
		data: data,
		xmin: xmin,
		xmax: xmax,
		ymin: ymin,
		ymax: ymax,
	}
}

func (b *EmptyBin[T]) Get(x, y int) []T {
	return b.data[y][x]
}

func (b *EmptyBin[T]) Add(key T) {
	ibinX, ibinY := b.GetBinXY(key)
	b.data[ibinY][ibinX] = append(b.data[ibinY][ibinX], key)
}

func (b *EmptyBin[T]) GetBinXY(key T) (int, int) {
	x, y := key.XY()

	deltaY := (b.ymax - b.ymin) / float64(len(b.data))
	deltaX := (b.xmax - b.xmin) / float64(len(b.data[0]))

	ibinY := int((y - b.ymin) / deltaY)
	ibinX := int((x - b.xmin) / deltaX)

	return ibinX, ibinY
}

func (b *EmptyBin[T]) Remove(x, y, i int) {
	bin := b.data[y][x]
	bin = append(bin[:i], bin[i+1:]...)
}

func (b *EmptyBin[T]) GetAll() []T {
	var out []T
	for y := range b.data {
		for x := range b.data[y] {
			for _, v := range b.data[y][x] {
				out = append(out, v)
			}
		}
	}
	return out
}

func (b *EmptyBin[T]) Update() {
	for y := range b.data {
		for x := range b.data[y] {
			for i, v := range b.data[y][x] {
				ibinX, ibinY := b.GetBinXY(v)
				if ibinX == x && ibinY == y {
					continue
				}
				b.Remove(x, y, i)
				b.Add(v)
			}
		}
	}
	fmt.Println(b.data)
}
