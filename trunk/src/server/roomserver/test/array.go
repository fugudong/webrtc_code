package main

func main()  {
	var data1 [] byte
	data1 = append(data1, 1)
	data1 = append(data1, 2)
	data1 = append(data1, 3)
	data1 = append(data1, 4)
	data1 = append(data1, 5)
	println("orig:",len(data1))
	println("half:",len(data1[2:]))
}
