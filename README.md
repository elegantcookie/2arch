# Архиватор 2ch.hk
Может скачивать тред в html и json.

# Установка и запуск
 
 <b>Установка с github</b>
```
git clone https://github.com/elegantcookie/2arch.git
cd 2arch
go build 2arch
```

<b>Установка из релизов</b>

Скачать, разархивировать папку

<b>Запуск</b>
```
2arch.exe --help
```

# Команды
<b>Скачивание треда в html</b> 

```
2arch.exe --url [ссылка на тред]
2arch.exe --u [ссылка на тред]
```

<b>Скачивание треда в json</b>
```
2arch.exe --url [ссылка на тред] --json
2arch.exe -u [ссылка на тред] -j
```
# Примеры использования
<b>Скачиваем файл в html</b>
```
2arch.exe -u https://2ch.hk/abu/res/42375.html
```

<b>Скачиваем файл в json</b>
```
2arch.exe -u https://2ch.hk/abu/res/42375.html -j
```

<b>Треды сохраняются в папку threads</b>

![image](https://user-images.githubusercontent.com/68335351/174891350-782cc811-32db-4f2d-8025-8308f693bb95.png)


# Использованные библиотеки
cobra, json, io, net, os, regexpr, sync, time и другие...
