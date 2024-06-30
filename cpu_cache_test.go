package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
	"unsafe"

	cpuid "github.com/klauspost/cpuid/v2"
)

func MatrixCombination(data *BenchmarkData) {
	for i := 0; i < data.MatrixSize; i++ {
		for j := 0; j < data.MatrixSize; j++ {
			data.MatrixA[i][j] = data.MatrixA[i][j] + data.MatrixB[i][j]
		}
	}
}

func MatrixReversedCombination(data *BenchmarkData) {
	for i := 0; i < data.MatrixSize; i++ {
		for j := 0; j < data.MatrixSize; j++ {
			data.MatrixA[i][j] = data.MatrixA[i][j] + data.MatrixB[j][i]
		}
	}
}

func MatrixCombinationPerBlock(data *BenchmarkData) {
	for i := 0; i < data.MatrixSize; i += data.BlockSize {
		for j := 0; j < data.MatrixSize; j += data.BlockSize {
			for ii := i; ii < i+data.BlockSize; ii++ {
				for jj := j; jj < j+data.BlockSize; jj++ {
					data.MatrixA[ii][jj] = data.MatrixA[ii][jj] + data.MatrixB[ii][jj]
				}
			}
		}
	}
}

func MatrixReversedCombinationPerBlock(data *BenchmarkData) {
	for i := 0; i < data.MatrixSize; i += data.BlockSize {
		for j := 0; j < data.MatrixSize; j += data.BlockSize {
			for ii := i; ii < i+data.BlockSize; ii++ {
				for jj := j; jj < j+data.BlockSize; jj++ {
					data.MatrixA[ii][jj] = data.MatrixA[ii][jj] + data.MatrixB[jj][ii]
				}
			}
		}
	}
}

func createMatrix(matixLength int, rnd *rand.Rand) [][]int64 {
	matrix := make([][]int64, matixLength)

	for i := 0; i < matixLength; i++ {
		matrix[i] = make([]int64, matixLength)
		for j := 0; j < matixLength; j++ {
			matrix[i][j] = rnd.Int63()
		}
	}

	return matrix
}

func copyMatrix(matixLength int, source [][]int64) [][]int64 {
	matrix := make([][]int64, matixLength)

	for i := 0; i < matixLength; i++ {
		matrix[i] = make([]int64, matixLength)
		_ = copy(matrix[i], source[i])
	}

	return matrix
}

func GetBlockSize() int {
	typeSize := unsafe.Sizeof(0)
	cacheSize := cpuid.CPU.Cache.L1D

	return cacheSize / int(typeSize)
}

type BenchmarkFunctionData struct {
	Name      string
	BenchFunc func(*BenchmarkData) `json:"-"`
}

type BenchmarkData struct {
	RunData        BenchmarkFunctionData
	MatrixSize     int
	BlockSize      int
	MatrixA        [][]int64 `json:"-"`
	MatrixB        [][]int64 `json:"-"`
	TimeElapsed    []time.Duration
	AvgTimeElapsed time.Duration
}

type BenchmarkInfoOutput struct {
	Name           string          `json:"Name"`
	MatrixSize     int             `json:"MatrixSize"`
	BlockSize      int             `json:"BlockSize"`
	TimeElapsed    []time.Duration `json:"TimeElapsed"`
	AvgTimeElapsed time.Duration   `json:"AvgTimeElapsed"`
}

var benchmarkDataCollection []*BenchmarkData

func NewBenchmarkData(runData BenchmarkFunctionData, matrixSize int, blockSize int) *BenchmarkData {
	return &BenchmarkData{
		RunData:    runData,
		MatrixSize: matrixSize,
		BlockSize:  blockSize,
	}
}

func SaveResults(benchmarkDataCollection []*BenchmarkData) {
	var benchmarkInfoCollection []*BenchmarkInfoOutput = make([]*BenchmarkInfoOutput, 0, len(benchmarkDataCollection))
	for _, benchmarkData := range benchmarkDataCollection {
		if benchmarkData == nil {
			continue
		}
		benchmarkInfoCollection = append(benchmarkInfoCollection, &BenchmarkInfoOutput{
			Name:           benchmarkData.RunData.Name,
			MatrixSize:     benchmarkData.MatrixSize,
			BlockSize:      benchmarkData.BlockSize,
			TimeElapsed:    benchmarkData.TimeElapsed,
			AvgTimeElapsed: benchmarkData.AvgTimeElapsed,
		})
	}
	res, err := json.Marshal(benchmarkInfoCollection)
	if err != nil {
		panic(err)
	}

	fileName := fmt.Sprintf("output-%s.json", time.Now().Format("2006-01-13_15-04-05"))
	err = os.WriteFile(fileName, res, 0644)
	if err != nil {
		panic(err)
	}
}

// go test -bench . -benchtime 500ms -cpu 1
func BenchmarkMatrixCombination(b *testing.B) {
	blockSize := GetBlockSize()
	fmt.Printf("Elements per cacheline: %d\n", blockSize)

	amountOfRuns := 4
	fmt.Printf("Runs per test: %d\n", amountOfRuns)

	blockSizes := []int{blockSize / 64}
	matrixSizes := []int{blockSize * 4}

	// amountOfRuns := 10
	// blockSizes := []int{blockSize}
	// matrixSizes := []int{blockSize}
	//var randSeed int64 = 8444154584984

	var randSeed int64 = time.Now().UnixMilli()
	rnd := rand.New(rand.NewSource(randSeed))

	benchmarkFunctionData := []BenchmarkFunctionData{
		{"MatrixCombination", MatrixCombination},
		{"MatrixReversedCombination", MatrixReversedCombination},
	}

	benchmarkFunctionDataWithBlock := []BenchmarkFunctionData{
		{"MatrixCombinationPerBlock", MatrixCombinationPerBlock},
		{"MatrixReversedCombinationPerBlock", MatrixReversedCombinationPerBlock},
	}

	totalAmountOfBenchmarks := len(blockSizes) * len(matrixSizes) * (len(benchmarkFunctionData) + len(benchmarkFunctionDataWithBlock))
	benchmarkDataCollection = make([]*BenchmarkData, 0, totalAmountOfBenchmarks)

	//Horrible things going there. For explanation go to line 198
	var matrixesA [][][]int64 = make([][][]int64, len(matrixSizes))
	matrixesB := make([][][]int64, len(matrixSizes))

	for i, matrixSize := range matrixSizes {
		matrixesA[i] = createMatrix(matrixSize, rnd)
		matrixesB[i] = createMatrix(matrixSize, rnd)
		for _, blockSize := range blockSizes {
			if blockSize <= matrixSize {
				for _, data := range benchmarkFunctionData {
					benchmarkDataCollection = append(benchmarkDataCollection, NewBenchmarkData(data, matrixSize, 0))
				}

				for _, data := range benchmarkFunctionDataWithBlock {
					benchmarkDataCollection = append(benchmarkDataCollection, NewBenchmarkData(data, matrixSize, blockSize))
				}
			}

		}
	}

	//runtime.GOMAXPROCS(1)
	//fmt.Printf("\nCPUS: %d\n", runtime.NumCPU())

	for _, benchmarkData := range benchmarkDataCollection {
		var totalTimeElapsed time.Duration
		for i := 0; i < amountOfRuns; i++ {
			fmt.Printf("Name: %s\n", benchmarkData.RunData.Name)
			fmt.Printf("Current matrix size: %d\n", benchmarkData.MatrixSize)
			fmt.Printf("Current block size: %d\n", benchmarkData.BlockSize)
			for j := 0; j < len(matrixSizes); j++ {
				if len(matrixesA[j]) == benchmarkData.MatrixSize {
					//Only known by me way to stop this test eating 64GB RAM is to create and release copies as soon as possible
					benchmarkData.MatrixA = copyMatrix(benchmarkData.MatrixSize, matrixesA[j])
					benchmarkData.MatrixB = copyMatrix(benchmarkData.MatrixSize, matrixesB[j])
				}
			}

			start := time.Now()
			benchmarkData.RunData.BenchFunc(benchmarkData)
			elapsed := time.Since(start)
			benchmarkData.MatrixA = nil
			benchmarkData.MatrixB = nil
			fmt.Printf("Time Elapsed: %d\n", elapsed)
			benchmarkData.TimeElapsed = append(benchmarkData.TimeElapsed, elapsed)

			totalTimeElapsed += benchmarkData.TimeElapsed[i]
		}
		benchmarkData.AvgTimeElapsed = totalTimeElapsed / time.Duration(amountOfRuns)
		fmt.Printf("Avg Time Elapsed: %dns\n", benchmarkData.AvgTimeElapsed)
	}

	SaveResults(benchmarkDataCollection)
}
