package main

func main() {
	handler, err := NewHandler()
	if err != nil {
		panic(err)
	}
	handler.Run()
}
