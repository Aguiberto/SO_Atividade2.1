package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const bufferSize = 5

var (
	buffer []int
	ctx    = context.Background()

	// Semáforos
	empty = semaphore.NewWeighted(bufferSize) // Espaços vazios disponíveis
	full  = semaphore.NewWeighted(bufferSize) // Itens disponíveis para consumo
	mutex = semaphore.NewWeighted(1)          // Exclusão mútua para acesso ao buffer
)

func main() {
	// Inicialmente, o semáforo 'full' começa zerado (nenhum item produzido ainda)
	// Como NewWeighted cria o semáforo com o valor máximo de recursos, 
	// nós "adquirimos" todos os recursos dele para começar em 0.
	_ = full.Acquire(ctx, bufferSize)

	var wg sync.WaitGroup
	wg.Add(2)

	// Inicia as Goroutines
	go produtor(&wg)
	go consumidor(&wg)

	// Deixa rodar por 10 segundos e encerra o programa
	time.Sleep(10 * time.Second)
	fmt.Println("\n[Fim do tempo] Encerrando simulação.")
}

func produtor(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// Produz um item aleatório
		item := rand.Intn(100)
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond) // Simula tempo de produção

		// 1. Decrementa o número de espaços vazios (bloqueia se o buffer estiver cheio)
		_ = empty.Acquire(ctx, 1)

		// 2. Entra na região crítica
		_ = mutex.Acquire(ctx, 1)

		// Insere o item no buffer
		buffer = append(buffer, item)
		fmt.Printf("[Produtor] Produziu: %d | Buffer: %v\n", item, buffer)

		// Saída da região crítica
		mutex.Release(1)

		// 3. Incrementa o número de itens cheios (libera o consumidor)
		full.Release(1)
	}
}

func consumidor(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		time.Sleep(time.Duration(rand.Intn(800)) * time.Millisecond) // Simula tempo de consumo

		// 1. Decrementa o número de itens cheios (bloqueia se o buffer estiver vazio)
		_ = full.Acquire(ctx, 1)

		// 2. Entra na região crítica
		_ = mutex.Acquire(ctx, 1)

		// Remove o primeiro item do buffer (FIFO)
		item := buffer[0]
		buffer = buffer[1:]
		fmt.Printf("[Consumidor] Consumiu: %d | Buffer: %v\n", item, buffer)

		// Saída da região crítica
		mutex.Release(1)

		// 3. Incrementa o número de espaços vazios (libera o produtor)
		empty.Release(1)
	}
}