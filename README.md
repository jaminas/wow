# Пояснительная записка

## Server
### cmd/server/main.go
Запускалка сервера:
- Определение конфигурации
- Инициализация дерева зависимостей - мы хотим поддержать принцип Dependency Injection SOLID для удобного переиспользования каждого слоя и удобного тестирования
- Запуск сервера

### internal/server - слой взаимодействия с внешней средой
Реализация простого tcp-сервера:
- Слушает порт
- Сервер умеет только принимать/отправлять запросы в формате протокола. Тут есть связанность с протоколом, правильно сделать отвязку протокола, я об этом знаю, но не хотелось зависнуть с ТЗ на много времени.
- При создании сервера ему передается реализация Guard и Handler по описанным интерфейсам.
- Сервер ничего не знает как происходит защита TCP-соединения - этот черный ящик реализуется через интерфейсом Guard где выходным параметром есть success[bool] который говорит что защита прошла и можно отправлять успешный пейлоад.
- Сервер ничего не знает как обрабатывать успешный запрос - это черный ящик реализуется через интерфейс Handler который возвращает результат своей работы
- При подключении клиента создает обрабатывает запросы клиента в отдельной горутине.
- Не реализован пул воркеров, который должен быть дабы не сожрать всю память при большом кол-ве запросов, но наверное выходит за рамки ТЗ
- Не реализованы ограничения буффера на чтение которые должны быть дабы не отправляли километровые запросы, но наверное выходит за рамки ТЗ
- Тестов нет, посчитал что выходит за рамки ТЗ, потому как надо писать качественные моки для хендлера, протектора. Ну и тут много чего можно потестировать
- Сервер читает стоку до delimiter - в реальных условиях так делать не надо, а в нашем допустимо.

### pkg/guard - middleware защиты
Реализация способа защиты
- Фактически это фасад (можно доработать до классического фасада) с реализацией POW защиты, можно вставлять любую другую. 
- Если хочется менять или конфигурацией, то стоит сделать еще одну абстракцию, для данного ТЗ оверхед.
- /pkg/pow - реализация защиты POW скачана откуда-то с гитхаба, внутренности изучены, но писать свою реализацию было бы оверхед.
- на гитхабе есть много готовых реализаций, в принципе можно подключать стандартно

### pkg/protocol - протокол обмена
- просто формат общения
- но при практическом применении лучше разделить на непосредственно транспортный протокол и протокол защиты, сейчас это все в одной куче

### pkg/cache - кеш
- стандартный кеш в памяти
- guard хранит там свои данные для межзапросового взаимодействия
- можно подсунуть любую реализацию, вплоть до внешнего распределенного
- можно сделать другую реализацию, в которой кеш не нужен - передавать данные клиенту и возвращать их. Но тут мы попадаем на сетевой траффик, но уменьшаем затраты памяти. Баланс который нужно решать в конкретном случае.

## Client
### cmd/client/main.go
Запускалка клиента. 
Клиент сделан максимально просто без каких либо абстракций дабы показать что сервер работает. 
В практическом применении нужно реализовывать протокол со стороны клиента отдельным пакетом с его имплементацией.

# Requirements
Go 1.21
Docker


# Install
```
make server
```

# Start native server 
```
make server
```

# Start native client
```
make client
```

# Start on docker:
```
docker-compose up --abort-on-container-exit --force-recreate --build server --build client
```

# Выбор алгоритма
Собственно гугл говорит нам о том что есть несколько алгоритмов POW, наиболее популярные из них:
+ [Hashcash](https://en.wikipedia.org/wiki/Hashcash)
+ [Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree)
+ [Guided tour puzzle](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol)

Сообщество рекомендует Hashcash. Его реализация была скачана из недр гитхаба. Недостатки других алгоритмов в огромном кол-ве гуглятся, тут их копипастить наверное смысла нет.
## Hashcash преимущества (копипаст):
- простота реализации
- много документации и статей с описанием
- простота проверки на стороне сервера
- возможность динамического управления сложностью для клиента путем изменения необходимого количества ведущих нулей
## Hashcash недостатки (копипаст):
- Время вычислений зависит от мощности клиентской машины.
   Например, очень слабые клиенты, возможно, не смогут решить задачу, или слишком мощные компьютеры могут реализовать DDOS-атаки.
   Но сложность задачи может быть решена динамически путем изменения необходимых нулей с сервера.
- Предварительные расчеты перед DDOS-атакой.
   Некоторые клиенты могли анализировать протокол и вычислять множество задач, чтобы применить все это за один момент.
   Эту проблему можно решить путем дополнительной проверки параметров hashcash на сервере.


## Улучшения
- Добавить пул воркеров (для анти-ddos это первое что нужно делать)
- Добавить лимиты буфера чтения сети (для анти-ddos это второе что нужно делать)
- Явно видно что алгоритм на стороне сервера потребялеет какие-то ресурсы, хоть и сильно меньшие чем клиент (это и суть защиты) - данная реализация вычисления имеет конфигурационный параметр сложности вычисления, его можно тюнить в зависимости от ситуации, например ограничением ресурсов.
- Явно видно что алгоритм можно сделать распределенным если топология приложения планируется как one_large_input_many_tiny_outputs
- Вероятно (нужно смотреть на реальных данных в реальной жизни) можно убрать кеш и передавать эти данные от сервер-клиент-сервер в зашифрованном виде. Но это тонкий баланс со многими вводными. Оптимизируем одно, возрастает нагрузка на другое.
- Конечно тесты, их нет не одного, но архитектура приложения построена так чтобы с помощью моков можно было положить каждый пакет на тесты.
- Если конечным клиентом является человек - то есть сильно более совершенные реализации алгоритма. Например гугловые картинки.
- Гуглится еще решение, где алгоритм вычисления передается клиенту в теле запроса, этим самым мы можем динамически передавать(например создавая на лету) алгоритм, тем самым для слабых клиентов расслабляя сложность вычисления.

# PS
Данный сервер исключительно на коленке показательный