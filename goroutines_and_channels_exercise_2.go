package main

import (
	"fmt"
	"io"
	"net/http"
)

func fetch(url string, getResults chan<- string) {
	// verifica se o host está acessível.
	response, err := http.Get(url)
	if err != nil {
		getResults <- fmt.Sprintf("\nFailed fetching %s: %v\n", url, err)
		return
	}

	// o 'defer' adia a execução do código para quando o retorno da funcção (ou o final dela) é chamado.
	defer response.Body.Close()

	// tenta fazer uma leitura do corpo da Response.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		getResults <- fmt.Sprintf("Failed to read data. Lost %d bytes.", len(body))
		return
	}

	// Envia os serultados para o canal.
	getResults <- fmt.Sprintf("Fetched %s: %d bytes.", url, len(body))
}

func main() {
	urls := []string{
		"https://www.google.com",
		"https://www.golang.org",
		"https://www.github.com",
		"https://www.non-existent-site.xyz",
		"https://www.wikipedia.org",
		"https://youtube.com",
		"https://www.sistemaq.com",
		"https://www.murissocas-onlie.br",
	}

	fmt.Println("Starting fetching the selected URLs...")

	// executa uma goroutine para cada url AO MESMO TEMPO!
	getResults := make(chan string)
	for _, url := range urls {
		go fetch(url, getResults)
	}

	// os resultados são impressos de acordo com a finaliação de cada goroutine
	// a ordem de impressão pode ser aleatória, uma vez que as tarefas estão sendo exexutadas os mesmo tempo.
	for i := 0; i < len(urls); i++ {
		result := <-getResults // result recebe o sinal 'send' do canal getResults
		fmt.Println(result)
	}

	fmt.Println("Operation has finished.")
}
