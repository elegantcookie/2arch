# Архиватор 2ch.hk
Может скачивать тред в html и json.
Изображения и видео скачиваются без повторов. Также есть функции скачивания только картинок/видео. Видео сохраняются в папку threads

# Установка
 
 <b>Установка с github</b>
 
 1. Установить <a href="https://go.dev/dl/">GO</a> и <a href="https://desktop.github.com/">GitHub Desktop</a>
 
 2. В консоли написать:
 
```
git clone https://github.com/elegantcookie/2arch.git
cd 2arch
go build 2arch
go install 2arch
```

<b>Установка из релизов</b>

1. Скачать, разархивировать папку

# Запуск

<b>Запуск с github</b>
```
2arch --help
```

<b>Запуск из релизов</b>
```
2arch.exe --help
```

# Команды
<b>Скачивание треда в html</b> 

```
2arch.exe --url [ссылка на тред]
2arch.exe --u [ссылка на тред]
```

<b>Скачивание треда в html только с изображениями (видео)</b> 

```
2arch.exe --url [ссылка на тред] --images (--videos)
2arch.exe --u [ссылка на тред] -i (-v)
```

<b>Скачивание треда в json</b>
```
2arch.exe --url [ссылка на тред] --json
2arch.exe -u [ссылка на тред] -j
```

# Примеры использования

Если установлено командой "git install 2arch", то можно запускать, как 2arch вместо 2arch.exe

<b>Скачиваем файл в html</b>
```
2arch.exe -u https://2ch.hk/abu/res/42375.html
```
или 
```
2arch -u https://2ch.hk/abu/res/42375.html
```

<b>Скачиваем только картинки</b>
```
2arch.exe -u https://2ch.hk/abu/res/42375.html -i
```
или
```
2arch -u https://2ch.hk/abu/res/42375.html -i
```

<b>Скачиваем только видео</b>
```
2arch.exe -u https://2ch.hk/abu/res/42375.html -v
```
или
```
2arch -u https://2ch.hk/abu/res/42375.html -v
```

<b>Скачиваем файл в json</b>
```
2arch.exe -u https://2ch.hk/abu/res/42375.html -j
```
или
```
2arch -u https://2ch.hk/abu/res/42375.html -j
```

<b>Треды сохраняются в папку threads</b>

![image](https://user-images.githubusercontent.com/68335351/174891350-782cc811-32db-4f2d-8025-8308f693bb95.png)


# Использованные библиотеки
cobra, json, io, net, os, regexpr, sync, time и другие...
