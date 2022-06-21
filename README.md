# Архиватор 2ch.hk
Может скачивать тред в html и json.

# Установка и запуск
 
```
git clone https://github.com/elegantcookie/2arch.git
cd 2arch
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
2arch -u https://2ch.hk/abu/res/42375.html -j
```

<b>Треды сохраняются в папку threads</b>

![image](https://user-images.githubusercontent.com/68335351/174891017-fd000fae-830d-43da-a6ef-35efe12c25c6.png)


# Использованные библиотеки
cobra, json, io, net, os, regexpr, sync, time и другие...
