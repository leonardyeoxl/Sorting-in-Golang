package main

import(
	"fmt"
	"sync"
)
         
func Merge(left, right [] int) [] int{
	merged := make([] int, 0, len(left) + len(right))

	for len(left) > 0 || len(right) > 0{
		if len(left) == 0 {
			return append(merged,right...)

		} else if len(right) == 0 {
			return append(merged,left...)

		} else if left[0] < right[0] {
			merged = append(merged, left[0])
			left = left[1:]

		} else{
			merged = append(merged, right[0])
			right = right[1:]
		}
	}

	return merged
}

func SingleMergeSort(data [] int) [] int {
	if len(data) <= 1 {
	  return data
	}

	mid := len(data)/2
	left := SingleMergeSort(data[:mid])
	right := SingleMergeSort(data[mid:])
	
	return Merge(left,right)
}

func ConcurrentMergeSort(data [] int, c chan struct{}) [] int {
	if len(data) <= 1 {
		return data
	}
	
	mid := len(data)/2

	var wg sync.WaitGroup
	wg.Add(2)

	var leftData []int
	var rightData []int

	select {
	case c <- struct{}{}:
		go func() {
			leftData = ConcurrentMergeSort(data[:mid], c)
			wg.Done()
		}()
	default:
		leftData = SingleMergeSort(data[:mid])
		wg.Done()
	}

	select {
	case c <- struct{}{}:
		go func() {
			rightData = ConcurrentMergeSort(data[mid:], c)
			wg.Done()
		}()
	default:
		rightData = SingleMergeSort(data[mid:])
		wg.Done()
	}
		
	wg.Wait()
	return Merge(leftData,rightData)
}

func RunMergeSort(data []int) []int {
	c := make(chan struct{}, 4)
	return ConcurrentMergeSort(data, c)
}

func main(){
	data := [] int{9,4,3,6,1,2,10,5,7,8}
	fmt.Printf("%v\n%v\n", data, RunMergeSort(data))
}