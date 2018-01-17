package utils

import (
	"os"
	"github.com/ewalker544/libsvm-go"
	"path/filepath"
	"runtime"
	"path"
	"image/gif"
	"image/png"
	"image"
	"fmt"
	"bufio"
	"io"
	"strings"
	"net/http"
	"io/ioutil"
	"bytes"
)

func getSourceCodePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func getExecPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func gifToPng(imagePath string) {
	fp, _ := os.Open(imagePath)
	__gif, _ := gif.Decode(fp)
	out, _ := os.Create(imagePath)
	png.Encode(out, __gif)
	fp.Close()
	out.Close()

}

func trainPreprocess(__path string) {
	filepath.Walk(path.Join(getSourceCodePath(), __path), func(subPath string, f os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if f.IsDir() {
			return nil
		}
		gifToPng(subPath)
		return nil
	})
}

func crop(src image.Image, name string) map[string][]float64 {
	vec := make(map[string][]float64, 0)
	rgbImg := src
	index := 0
	for i := 2; i < 50; i += 12 {
		var tmp []float64
		for y := 1; y < 22; y++ {
			for x := 0; x <= 16; x++ {
				pixel := rgbImg.At(x+i, y)
				r, g, b, _ := pixel.RGBA()
				y := float64(0.3*float64(r)+0.59*float64(g)+0.11*float64(b)) / 257.0
				tmp = append(tmp, y/255.0)
			}
		}
		vec[fmt.Sprintf("%s-%d", name, index)] = tmp
		index++
	}
	return vec
}

func loadSamples(__path string) map[string][]float64 {
	vector := make(map[string][]float64, 0)
	filepath.Walk(path.Join(getSourceCodePath(), __path), func(subPath string, f os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if f.IsDir() {
			return nil
		}
		fp, _ := os.Open(subPath)
		imOfGif, _ := gif.Decode(fp)
		name := string([]byte(subPath)[strings.LastIndex(subPath, "\\")+1:strings.LastIndex(subPath, ".")])
		subImgVector := crop(imOfGif, name)
		for key, value := range subImgVector {
			vector[key] = value
		}

		defer fp.Close()
		return nil
	})
	return vector
}

func loadLabels(__path string) map[string]byte {
	labels := make(map[string]byte, 0)
	f, _ := os.Open(path.Join(getSourceCodePath(), __path))
	inputReader := bufio.NewReader(f)
	index := 0
	for {
		line, err := inputReader.ReadSlice('\n')
		if err == io.EOF {
			break
		}

		for i, c := range line {
			labels[fmt.Sprintf("%d-%d", index, i)] = c
			//fmt.Println(c)
		}
		index++
	}
	return labels
}

func localTrain() {
	//loadSamples("datasets/samples/0.png")
	//trainPreprocess("datasets/samples/")
	samples := loadSamples("datasets/samples/")
	labels := loadLabels("datasets/answer.txt")

	f, _ := os.Create("ty.txt")

	for key, value := range samples {
		f.WriteString(fmt.Sprintf("%d ", int(labels[key])))
		for j, val := range value {
			f.WriteString(fmt.Sprintf("%d:%.3f ", j+1, val))
		}
		f.WriteString("\n")
	}

	defer f.Close()

	param := libSvm.NewParameter()
	param.KernelType = libSvm.C_SVC
	param.C = 128.0
	param.Gamma = 0.00048828125

	model := libSvm.NewModel(param)

	problem, _ := libSvm.NewProblem(path.Join(getSourceCodePath(), "ty.txt"), param)

	model.Train(problem)

	model.Dump(path.Join(getSourceCodePath(), "zf.model"))
}

func Predict(save bool) string {
	rep, _ := http.Get("http://zfxk.zjtcm.net/CheckCode.aspx")
	pix, _ := ioutil.ReadAll(rep.Body)
	defer rep.Body.Close()
	if save {
		fp, _ := os.Create(path.Join(getSourceCodePath(), "tmp.png"))
		io.Copy(fp, bytes.NewReader(pix))
		defer fp.Close()
	}
	model := libSvm.NewModelFromFile(path.Join(getSourceCodePath(), "zf.model"))
	im, _, _ := image.Decode(bytes.NewReader(pix))
	vec := crop(im, "loc")
	ret := make([]byte, 0)
	x := make(map[int]float64)
	for ind := 0; ind < 4; ind++ {
		for index, value := range vec[fmt.Sprintf("loc-%d", ind)] {
			x[index+1] = value
		}
		predictLabel := model.Predict(x)
		ans := byte(predictLabel)
		ret = append(ret, ans)
	}
	return string(ret)
}

func main() {
	//localTrain()
	fmt.Println(Predict(true))
}
