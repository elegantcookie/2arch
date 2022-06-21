# Архиватор 2ch.hk
Может скачивать тред в html и json

# Установка

```
git clone https://github.com/elegantcookie/2arch.git
cd 2arch
```

# Запуск

```
2arch --help
```

# Команды
<b>Скачивание треда в html</b> 

```
2arch --url [ссылка на тред]
2arch --u [ссылка на тред]
```

<b>Скачивание треда в json</b>
```
2arch --url [ссылка на тред] --json
2arch -u [ссылка на тред] -j
```
# Примеры использования
<b>Скачиваем файл в html</b>
```
2arch -u https://2ch.hk/abu/res/42375.html
```

<b>Скачиваем файл в json</b>
```
2arch -u https://2ch.hk/abu/res/42375.html
```

# Использованные библиотеки
cobra, json, io, net, os, regexpr, sync, time и другие...