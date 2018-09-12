package main

func ConvertByteArrayToIntArray(in []byte) (out []int) {
	outArray := []int{}

	for _, v := range in {
		outArray = append(outArray, int(v))
	}

	return outArray
}