package main

import (
	"fmt"
	"time"
)

func doWork(done chan<- bool) {
	fmt.Println("Trabalhando...")
	// finjindo que está trabalhando...
	time.Sleep(10 * time.Second)

	fmt.Println("Trabalho concluído.")
	done <- true // envia um sinal de 'true' para o canal de booleanos.
}

func main() {
	fmt.Println("Início da função 'main'.")

	fmt.Println("Iniciando a Goroutine...")
	// cria um canal de booleanos para gerenciar o sinal 'send' da goroutine.
	done := make(chan bool)
	go doWork(done) // chamada da goroutine.

	// a main estará bloqueada até o canal 'done' receber um sinal 'send' da goroutine.

	<-done

	// assim que o canal de booleanos receber o sinal de 'doWork', a main é liberada.
	fmt.Println("Final da 'main'.")
}
