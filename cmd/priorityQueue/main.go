package main

import (
	priorityQueue "algorithms/internal/priorityQueue/priorityQueueByLinkedList"
	"fmt"
)

func main() {
	queue := priorityQueue.NewPriorityQueue[string](3)

	queue.Insert("Незашифрованный пароль в базе данных", 10)
	queue.Insert("Пользовательский интерфейс не работает в браузере X", 9)
	queue.Insert("Стиль CSS нарушает выравнивание", 8)
	queue.Insert("Загрузка страницы занимает более 2 сек", 7)
	fmt.Println(queue.Top())
	queue.Insert("Стиль CSS вызывает смещение 1 пикселя", 5)
	queue.Insert("Переработать CSS используя SASS", 3)
	queue.Insert("Необязательное поле формы заблокировано", 8)

	queue.Insert("Утечка памяти", 9)
	queue.Insert("Добавить исключение для суперкуба", 9.5)

	fmt.Println(queue.AsciiTree())
}
