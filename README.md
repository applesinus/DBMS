[Jump to English](#English)

<a name="Russian"></a>
# Русский
<p id="ru"><h3></h3></p>

<p>Набор лабораторных работ, объединенных в единый проект, по предмету "Базы данных II". Оценка в семестре - <b>Отлично</b> (5).</p>

<h2>Web API UI</h2>
<p></p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/SignUpPage.png" title="Sign Up" alt="Sign Up"></p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/LogInPage.png" title="Log In" alt="Log In"></p>
<p>В системе есть регистрация и вход. После регистрации пользователь получает статус user, однако ему не дается прав на просмотр/редактирование ни одной коллекции.</p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/AdminPanel.png" title="Admin Page" alt="Admin Page"></p>
<p>Администратор может сменить пароль, роль и доступы всех пользователей, кроме других администраторов и суперпользователя.</p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/UserScreenWithoutAccess.png" title="Without accesses" alt="Without accesses"></p>
<p>Так выглядит страница пользователя, не имеющего доступов ни к одной коллекции.</p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/UserScreenWithAccess1.png" title="With accesses" alt="With access"></p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/UserScreenWithAccess2.png" title="With access" alt="With access"></p>
<p>Так выглядит страница пользователя, имеющего доступ к нескльким коллекциям. В данном случае это - суперпользователь, так что он имеет доступ ко всей базе данных, а также имеет кнопку для перехода в панель администратора.</p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/UserScreenWithoutAccess.png" title="Select2" alt="Select2"></p>
<p>Для удобства пользователя использован элемент select2 из <a href="https://select2.org/">соответствующего опен-сорс решения</a>. Пользователь видит только коллекции и данные, с которыми может работать.</p>

<h2>ТЗ</h2>
<p>Проект разбит на пакеты с реализацией конкретных частей задания.</p>
<p>Перед описанием задания в квадратных скобках указываются номера пререквизитов (заданий, которые необходимо выполнить для выполнения текущего). Частичная (неполная) реализация дополнительных заданий не допускается (за исключением заданий 6, 7).</p>

<h3>Задание 0</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main">Корневая папка</a>, <a href="https://github.com/applesinus/DBMS/tree/main/task0%2B3">папка task0+3</a></p>
<p>Реализуйте приложение, позволяющее выполнять операции над коллекциями данных заданных типов и контекстами их хранения (коллекциями данных). Коллекция данных описывается набором строковых параметров (набор параметров однозначно идентифицирует коллекцию данных):</p>

* название пула схем данных, хранящего схемы данных;
* название схемы данных, хранящей коллекции данных;
* название коллекции данных.

<p>Коллекция данных представляет собой ассоциативный контейнер (конкретная реализация определяется вариантом), в котором каждый объект данных соответствует некоторому уникальному ключу. Для ассоциативного контейнера необходимо вынести интерфейсную и реализовать этот интерфейс. Взаимодействие с коллекцией объектов происходит посредством выполнения одной из операций над ней:</p>

* добавление новой записи по ключу;
* чтение записи по её ключу;
* чтение набора записей с ключами из диапазона [𝑚𝑖𝑛𝑏𝑜𝑢𝑛𝑑... 𝑚𝑎𝑥𝑏𝑜𝑢𝑛𝑑];
* обновление данных для записи по ключу;
* удаление существующей записи по ключу.

<p>Во время работы приложения возможно выполнение также следующих операций:</p>

* добавление/удаление пулов данных;
* добавление/удаление схем данных для заданного пула данных;
* добавление/удаление коллекций данных для заданной схемы данных заданного пула данных.

<p>Поток команд, выполняемых в рамках работы приложения, поступает из файла, путь к которому подаётся в качестве аргумента командной строки. Формат команд в файле:</p>

```
createPool <name>
deletePool <name>

createScheme <name> in <pool name>
deleteScheme <name> in <pool name>

createCollection <type> <name> in <pool name>.<scheme name>
deleteCollection <name> in <pool name>.<scheme name>

set <key> <secondaryKey> <value> in <pool name>.<scheme name>.<collection name>
update <key> <value> in <pool name>.<scheme name>.<collection name>

get <key> in <pool name>.<scheme name>.<collection name>
getSecondary <secondaryKey> in <pool name>.<scheme name>.<collection name>
getAt <time> <key> in <pool name>.<scheme name>.<collection name>

getRange <key minimum> <key maximum> in <pool name>.<scheme name>.<collection name>
getRangeSecondary <key minimum> <key maximum> in <pool name>.<scheme name>.<collection name>

delete <key> in <pool name>.<scheme name>.<collection name>
```


<h3>Задание 1 [0]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main">Корневая папка</a>
<p>Реализуйте интерактивный диалог с пользователем. Пользователь при этом может вводить конкретные команды (формат ввода определите самостоятельно) и подавать на вход файлы с потоком команд.</p>

<h3>Задание 2 [0]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task2">Папка task2</a></p>
<p>Реализуйте механизм персистентности данных в коллекциях данных, позволяющий выполнять запросы к данным в рамках коллекции данных на заданный момент времени (дата и время, для которых нужно вернуть актуальную версию данных, передаются как параметр). Для реализации механизма персистентности используйте поведенческие паттерны проектирования “Команда” и “Цепочка обязанностей”.</p>

<h3>Задание 3 [0]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task0%2B3">Папка task0+3</a>, <a href="https://github.com/applesinus/DBMS/tree/main/task6%2B3">папка task6+3</a></p>
<p>Реализуйте механизм вторичных индексов, позволяющий выполнять эффективный поиск по различным отношениям порядка на пространстве данных (дублирование объектов данных в коллекциях данных, построенных по различным отношениям порядка на одном и том же наборе объектов данных, при этом запрещается). Обеспечьте выбор индекса для поиска при помощи указания ключа отношения порядка (в виде строки, подаваемой как параметр поиска).</p>

<h3>Задание 4 [0]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task4">Папка task0+3</a></p>
<p>Обеспечьте хранение объектов строк, размещённых в объектах данных, на основе структурного паттерна проектирования “Приспособленец”. Дублирования объектов строк для разных объектов (независимо от контекста хранения) при этом запрещены. Доступ к строковому пулу обеспечьте на основе порождающего паттерна проектирования “Одиночка”.</p>

<h3>Задание 6 [0]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task6%2B3">Папка task6+3</a></p>
<p>Реализуйте возможность кастомизации (при создании) реализаций ассоциативных контейнеров, репрезентирующих коллекции данных:</p>

* АВЛ-дерево
* Красно-чёрное дерево
* B-дерево

<h3>Задание 7 [0]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task789">Папка task789</a>, <a href="https://github.com/applesinus/DBMS/tree/main/web">папка web</a></p>
<p>Реализуйте функционал приложения в виде сервера, запросы на который поступают из клиентских приложений. При этом взаимодействие клиентских приложений с серверным должно быть реализовано посредством средств сетевого взаимодействия HTTP.</p>

<h3>Задание 8 [0,7]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task789">Папка task789</a>, <a href="https://github.com/applesinus/DBMS/tree/main/web">папка web</a></p>
<p>Реализуйте механизмы регистрации и авторизации пользователя в системе (на клиентской стороне) и открытия пользовательской сессии (на серверной стороне) через пару значений &lt;логин, пароль&gt;, с дальнейшим взаимодействием клиентской стороны с серверной на основе передачи вместе с запросом токена аутентификации. Пароль при этом должен храниться и передаваться в виде хеша (используйте хеш-функцию SHA-256). Логины пользователей уникальны в рамках системы, могут содержать только символы букв и цифр в количестве [5..15] (обеспечьте валидацию на стороне сервера). Пароль должен содержать не менее 8 символов (обеспечьте валидацию на стороне клиента). Формат хранения и передачи данных для авторизации пользователей определите самостоятельно.</p>

<h3>Задание 9 [0,7,8]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task789">Папка task789</a>, <a href="https://github.com/applesinus/DBMS/tree/main/web">папка web</a></p>
<p>На основе передаваемого в клиентские запросы токена аутентификации реализуйте различные роли, разграничивающие доступ к выполнению операций в рамках системы:</p>

* суперпользователь - имеет возможности создания новых пользователей; управления (выдача, блокировка) ролями других пользователей (кроме суперпользователя); управления доступом ко всем схемам данных и доступом к операциям на уровне схем данных для заданных коллекций данных (режимы “только для чтения” и “чтение и модификация”)
* администратор - имеет возможности создания новых пользователей; управления (выдача, блокировка) ролями других пользователей (кроме администраторов и суперпользователя); управления доступом к схемам данных и доступом ко всем операциям на уровне схем данных для заданных коллекций данных (режимы “только для чтения” и “чтение и модификация”)
* редактор - имеет возможности управления пулами (добавление/удаление), схемами (добавление/удаление) данных, коллекциями данных (добавление/удаление) в соответствии с предоставленными правами доступа
* пользователь - имеет возможности взаимодействия с коллекциями данных в соответствии с предоставленными правами доступа.


[Перейти к русскому](#Russian)

<a name="English"></a>
# English
<p id="en"><h3></h3></p>
<p>(translated by Google)</p>

<p>A set of laboratory works, combined into a single project, on the subject "Databases II". Semester grade - <b>Excellent</b> (A).</p>

<h2>Web API UI</h2>
<p></p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/SignUpPage.png" title="Sign Up" alt="Sign Up"></p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/LogInPage.png" title="Log In" alt="Log In"></p>
<p>The system has registration and login. After registration, the user receives user status, but is not given rights to view/edit any collection.</p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/AdminPanel.png" title="Admin Page" alt="Admin Page"></p>
<p>The administrator can change the password, role and access of all users except other administrators and the superuser.</p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/UserScreenWithoutAccess.png" title="Without accesses" alt="Without accesses"></p>
<p>This is what the page of a user who does not have access to any collections looks like.</p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/UserScreenWithAccess1.png" title="With accesses" alt="With access"></p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/UserScreenWithAccess2.png" title="With access" alt="With access"></p>
<p>This is what the page of a user who has access to several collections looks like. In this case, it is the superuser, so he has access to the entire database, and also has a button to go to the admin panel.</p>
<p><img src="https://github.com/applesinus/DBMS/blob/main/UIscreens/UserScreenWithoutAccess.png" title="Select2" alt="Select2"></p>
<p>For user convenience, the select2 element from <a href="https://select2.org/">the corresponding open source solution</a> was used. The user sees only collections and data that he can work with.</p>

<h2>Specifications</h2>
<p>The project is divided into packages with the implementation of specific parts of the task.</p>
<p>Before the task description, the numbers of prerequisites (tasks that must be completed to complete the current one) are indicated in square brackets. Partial (incomplete) implementation of additional tasks is not allowed (except for tasks 6, 7).</p>

<h3>Task 0</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main">Root folder</a>, <a href="https://github.com/applesinus/DBMS/tree /main/task0%2B3">task0+3 folder</a></p>
<p>Implement an application that allows you to perform operations on collections of data of specified types and their storage contexts (data collections). A data collection is described by a set of string parameters (the set of parameters uniquely identifies the data collection):</p>

* name of the data schema pool storing data schemas;
* name of the data schema storing data collections;
* name of the data collection.

<p>A data collection is an associative container (the specific implementation is determined by the variant) in which each data object corresponds to some unique key. For an associative container, you need to remove the interface and implement this interface. Interaction with a collection of objects occurs by performing one of the operations on it:</p>

* adding a new entry by key;
* reading a record by its key;
* reading a set of records with keys from the range [𝑚𝑖𝑛𝑏𝑜𝑢𝑛𝑑... 𝑚𝑎𝑥𝑏𝑜𝑢𝑛𝑑];
* updating data for recording by key;
* deleting an existing entry by key.

<p>While the application is running, the following operations can also be performed:</p>

* adding/removing data pools;
* adding/removing data schemas for a given data pool;
* adding/removing data collections for a given data schema of a given data pool.

<p>The stream of commands executed within the application comes from a file, the path to which is supplied as a command line argument. Command format in the file:</p>

```
createPool <name>
deletePool <name>

createScheme <name> in <pool name>
deleteScheme <name> in <pool name>

createCollection <type> <name> in <pool name>.<scheme name>
deleteCollection <name> in <pool name>.<scheme name>

set <key> <secondaryKey> <value> in <pool name>.<scheme name>.<collection name>
update <key> <value> in <pool name>.<scheme name>.<collection name>

get <key> in <pool name>.<scheme name>.<collection name>
getSecondary <secondaryKey> in <pool name>.<scheme name>.<collection name>
getAt <time> <key> in <pool name>.<scheme name>.<collection name>

getRange <key minimum> <key maximum> in <pool name>.<scheme name>.<collection name>
getRangeSecondary <key minimum> <key maximum> in <pool name>.<scheme name>.<collection name>

delete <key> in <pool name>.<scheme name>.<collection name>
```

<h3>Task 1 [0]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main">Root folder</a>
<p>Implement an interactive dialogue with the user. The user can enter specific commands (define the input format yourself) and submit files with a stream of commands as input.</p>

<h3>Task 2 [0]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task2">task2 folder</a></p>
<p>Implement a data persistence mechanism in data collections that allows you to perform queries on data within a data collection at a given point in time (the date and time for which you want to return the current version of the data are passed as a parameter). To implement the persistence mechanism, use the “Team” and “Chain of Responsibility” behavioral design patterns.</p>

<h3>Task 3 [0]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task0%2B3">Task0+3 folder</a>, <a href="https://github.com /applesinus/DBMS/tree/main/task6%2B3">task6+3 folder</a></p>
<p>Implement a secondary index mechanism that allows you to perform efficient searches based on different order relationships in the data space (duplication of data objects in data collections built using different order relationships on the same set of data objects is prohibited). Ensure that the search index is selected by specifying the order relation key (as a string supplied as a search parameter).</p>

<h3>Task 4 [0]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task4">Folder task0+3</a></p>
<p>Provide storage for row objects placed in data objects based on the Opportunist structural design pattern. Duplication of string objects for different objects (regardless of the storage context) is prohibited. Provide access to the string pool based on the generative Singleton design pattern.</p>

<h3>Task 6 [0]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task6%2B3">Task6+3 folder</a></p>
<p>Implement the ability to customize (when creating) implementations of associative containers representing data collections:</p>

* AVL tree
* Red-black tree
* B-tree

<h3>Task 7 [0]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task789">task789 folder</a>, <a href="https://github.com/applesinus/DBMS /tree/main/web">web folder</a></p>
<p>Implement the application functionality as a server, requests for which come from client applications. In this case, the interaction of client applications with the server must be implemented using HTTP network communication.</p>

<h3>Task 8 [0.7]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task789">task789 folder</a>, <a href="https://github.com/applesinus/DBMS /tree/main/web">web folder</a></p>
<p>Implement mechanisms for registering and authorizing a user in the system (on the client side) and opening a user session (on the server side) through a pair of values ​​&lt;login, password&gt;, with further interaction between the client side and the server side based on transmission along with a request for an authentication token . The password must be stored and transmitted as a hash (use the SHA-256 hash function). User logins are unique within the system and can only contain [5..15] alphanumeric characters (ensure server-side validation). The password must be at least 8 characters long (ensure client-side validation). Determine the format for storing and transmitting data for user authorization yourself.</p>

<h3>Task 9 [0,7,8]</h3>
<p><a href="https://github.com/applesinus/DBMS/tree/main/task789">task789 folder</a>, <a href="https://github.com/applesinus/DBMS /tree/main/web">web folder</a></p>
<p>Based on the authentication token passed to client requests, implement various roles that limit access to performing operations within the system:</p>

* superuser - has the ability to create new users; managing (issuing, blocking) the roles of other users (except for the superuser); control access to all data schemas and access to operations at the data schema level for specified data collections (read-only and read-and-modify modes)
* admin - has the ability to create new users; managing (issuing, blocking) the roles of other users (except admins and superuser); control access to data schemas and access to all operations at the data schema level for given data collections (read-only and read-and-modify modes)
* editor - has the ability to manage pools (adding/removing), schemas (adding/removing) data, collections of data (adding/deleting) in accordance with the granted access rights
* user - has the ability to interact with data collections in accordance with the granted access rights.
