# 🚀 HAD - دانلودر فوق پیشرفته

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS%20%7C%20ARM64-lightgrey)]()

یک دانلودر قدرتمند و چندنخی با پشتیبانی از **HTTP/HTTPS**، **FTP/FTPS**، **SFTP**، **Metalink**، **پروکسی ضبط MITM**، **اتصالات پروکسی**، **قابلیت از سرگیری هوشمند**، **رابط JSON-RPC** و **ابزار کامل پشتیبان‌گیری از وب‌سایت**. ساخته شده برای سرعت و قابلیت اطمینان.

[**اموزش دانلود فایل و ویدیو در مواقع نت ملی با گیت هاب*](https://github.com/Mr-Spect3r/had/blob/main/DownloaderWithAction.md)


## 🎉 جدید در نسخه v3.2.5

HAD اکنون شامل **پروکسی ضبط MITM** است که ترافیک HTTPS را رهگیری می‌کند تا به طور خودکار فایل‌های قابل دانلود را شناسایی و ضبط کند، **دانلود دسته‌ای از JSON ضبط شده**، و **نصب خودکار گواهی CA** برای رهگیری یکپارچه HTTPS.

<a href="https://github.com/Mr-Spect3r/had/releases/tag/3.2.5">نسخه v3.2.5</a>

## ✨ ویژگی‌ها

### دانلودر اصلی
- 🧵 **دانلود چندنخی** - حداکثر استفاده از پهنای باند
- 📡 **پروتکل‌های متعدد** - HTTP، HTTPS، FTP، FTPS، SFTP
- 🔄 **پشتیبانی از ادامه دانلود** - قطع و ادامه دانلودها بدون مشکل
- 🕸️ **پشتیبانی از پروکسی** - پروکسی‌های SOCKS4، SOCKS5 و HTTP
- 📁 **دانلود دسته‌ای** - دانلود از لیست فایل‌ها
- 🕷️ **وب اسکرپینگ** - استخراج و دانلود لینک‌ها از صفحات وب
- 🔍 **فیلتر پسوند** - فیلتر دانلودها بر اساس پسوند فایل
- ⚡ **بافر تطبیقی** - بهینه‌سازی خودکار اندازه بافر (16KB-1MB)
- 💾 **ذخیره جلسه** - ذخیره پیشرفت و ادامه بعداً
- 🎨 **نوار پیشرفت زیبا** - بازخورد بصری لحظه‌ای
- 🌍 **چندسکویی** - Windows، Linux، macOS، ARM64
- 🔐 **بررسی یکپارچگی** - تأیید SHA256، SHA1، MD5

### پروکسی ضبط MITM (جدید در v3.2.5)
- 🔒 **رهگیری HTTPS** - قابلیت کامل پروکسی مرد میانی
- 🎯 **تشخیص خودکار** - تشخیص خودکار ویدیوها، موسیقی، تصاویر، اسناد، آرشیوها
- 🔍 **استخراج URL پنهان** - اسکرپینگ HTML، JSON و جاوااسکریپت برای لینک‌های دفن شده
- 📊 **امتیازدهی اطمینان** - سیستم اطمینان ۰-۱۰۰٪ برای فیلتر مثبت‌های کاذب
- 📝 **خروجی چندفرمتی** - ذخیره لینک‌های ضبط شده در هر دو فرمت TXT و JSON
- 🔧 **پسوندهای سفارشی** - اضافه کردن پسوندهای فایل دلخواه برای ضبط
- 🎯 **فیلتر دامنه** - تمرکز روی دامنه‌های خاص
- 📐 **فیلتر اندازه** - محدودیت حداقل/حداکثر اندازه فایل
- 🔄 **دانلود خودکار** - گزینه دانلود خودکار فایل‌های ضبط شده
- 🍪 **پشتیبانی از کوکی** - عبور کوکی‌ها از طریق پروکسی
- 🔐 **گواهی خودکار** - نصب خودکار گواهی CA

### دانلود از JSON ضبط شده (جدید در v3.2.5)
- 📦 **دانلود دسته‌ای** - دانلود همه فایل‌های ضبط شده با یک دستور
- ⚡ **نخی هوشمند** - تنظیم خودکار تعداد رشته‌ها بر اساس اندازه فایل
- 🔄 **دانلود همزمان** - دانلود چندین فایل به طور همزمان
- 📊 **پیشرفت یکپارچه** - پیگیری همه دانلودها در یک رابط زیبا
- 💾 **پشتیبانی از ادامه** - دانلودهای قطع شده قابل ادامه هستند
- 🧹 **پاکسازی نام فایل** - پاکسازی خودکار عنوان‌ها برای نام‌های فایل امن

### دانلودر وب‌سایت
- 🌐 **آینه‌سازی کامل سایت** - خزش و پشتیبان‌گیری از کل وب‌سایت‌ها
- 📄 **پشتیبان‌گیری صفحات تکی** - ذخیره صفحه با همه وابستگی‌ها
- 🎯 **پشتیبانی SPA** - مدیریت مسیریابی مبتنی بر هش (#!/ و مسیرهای #!)
- 🖼️ **بازنویسی دارایی** - بازنویسی خودکار URL برای مرور آفلاین
- 🚀 **خزش همزمان** - کارگرهای قابل تنظیم برای صفحات و دارایی‌ها
- 💾 **خزش قابل ادامه** - ذخیره و ادامه پشتیبان‌گیری‌های قطع شده
- 🎨 **پردازش CSS/JS** - بازنویسی URLها در استایل‌شیت‌ها و اسکریپت‌ها
- 🖼️ **پشتیبانی Iframe** - دانلود محتوای iframe به صورت بازگشتی
- ⚡ **محدودیت نرخ** - کنترل درخواست به ازای هر دامنه

### پشتیبانی Metalink (RFC 5854)
- 📦 **پشتیبانی نسخه ۳ و ۴** - انطباق کامل با مشخصات Metalink
- 🔄 **انتخاب آینه** - انتخاب خودکار آینه بر اساس اولویت
- ✅ **تأیید یکپارچگی** - اعتبارسنجی جمع‌بازی داخلی
- 📊 **اطلاعات تکه‌ها** - تقسیم‌بندی فایل با تأیید هش

### رابط RPC
- 🔌 **API JSON-RPC** - قابلیت کنترل کامل از راه دور
- 🌐 **نقطه‌های پایانی REST** - نقاط پایانی HTTP برای وضعیت و کنترل
- 📊 **مانیتورینگ بی‌درنگ** - پیشرفت و آمار دانلود
- 🎮 **کنترل پویا** - توقف، ادامه، محدودیت سرعت از طریق API

### پشتیبانی SFTP
- 🔑 **احراز هویت کلید SSH** - پشتیبانی RSA، ECDSA، Ed25519
- 🔐 **احراز هویت رمز عبور** - با بازگشت به کلیدها
- 📁 **قابلیت کامل ادامه** - قطع و ادامه انتقال‌های SFTP
- ⏱️ **تایم‌اوت‌های قابل تنظیم** - تایم‌اوت اتصال و عملیات

## 📸 تصاویر

| ویژگی | پیش‌نمایش |
|---------|---------|
| **دانلود چندنخی** | ![چندنخی](https://github.com/user-attachments/assets/633999ce-c3da-4db5-b4be-b4714164a504) |
| **پیشرفت بی‌درنگ** | ![پیشرفت](https://github.com/user-attachments/assets/2e3f4fde-4e6c-4e79-bdae-fa9094bf2993) |
| **ادامه جلسه (JSON)** | ![ادامه جلسه](https://github.com/user-attachments/assets/370683d3-0d54-4b9c-8e40-f1ce2f515667) |
| **پروکسی ضبط MITM** | ![پروکسی ضبط](https://github.com/user-attachments/assets/6cba171d-632c-4aef-a654-ae33c9b21b4f) |
| **خروجی JSON از پروکسی** | ![خروجی JSON](https://github.com/user-attachments/assets/7f098862-e7b4-4baa-9663-b28489e9b5e3) |
| **پشتیبان وب (کلون)** | ![پشتیبان وب](https://github.com/user-attachments/assets/09883fe7-14d7-4045-8269-ea3c5bd5b1ae) |
| **وب اسکرپینگ** | ![اسکرپینگ](https://github.com/user-attachments/assets/fa49cd59-418d-4690-a8dc-b7a8ab2f043b) |
| **افزونه HAD** | ![افزونه](https://github.com/user-attachments/assets/ab0ea00b-7d44-45cf-8824-eb998a6c453d) |

## 📦 نصب

### نصب با Go

```bash
go install github.com/Mr-Spect3r/had@3.2.5
```

### نصب از روی سورس

```bash
git clone https://github.com/Mr-Spect3r/had.git
cd had
go build -o had .
```

## 🚀 شروع سریع

### پروکسی ضبط MITM (جدید)

ضبط لینک‌های دانلود در حین مرور:

```bash
# شروع پروکسی ضبط روی پورت 8085
./had -capture-proxy :8085 -capture-types video,music

# ضبط با پسوندهای سفارشی
./had -capture-proxy :9090 -capture-types video,archive -capture-exts .webm,.mka

# دانلود خودکار فایل‌های ضبط شده
./had -capture-proxy :8085 -capture-auto -capture-output ./downloads

# فیلتر بر اساس دامنه و میزان اطمینان
./had -capture-proxy :8085 -filter-domain example.com -capture-confidence 50

# فقط نصب گواهی
./had -install-cert
```

### دانلود از JSON ضبط شده (جدید)

بعد از ضبط لینک‌ها، همه چیز را با یک دستور دانلود کنید:

```bash
# دانلود همه فایل‌های ضبط شده
./had -download-json captured_links.json

# خروجی سفارشی با دانلود همزمان
./had -download-json captured_links.json -o ./videos -u 5

# کارایی بالا (8 رشته برای هر فایل، 4 دانلود همزمان)
./had -download-json captured_links.json -t 8 -u 4
```

### دانلود‌های پایه

```bash
# دانلود یک فایل
./had https://example.com/file.zip

# دانلود با 16 رشته
./had -t 16 https://example.com/large-file.zip

# دانلود چند فایل
./had https://example.com/file1.zip https://example.com/file2.zip

# دانلود از لیست فایل
./had -f urls.txt

# دانلود با محدودیت سرعت (1 مگابایت بر ثانیه)
./had -max-speed 1048576 https://example.com/file.zip

# دانلود با بررسی جمع‌بازی
./had -checksum-sha256 abc123... https://example.com/file.zip
```

### پشتیبان‌گیری از وب‌سایت

```bash
# پشتیبان کامل پایه از سایت
./had web -url https://example.com -mode full

# پشتیبان کامل از کل وب‌سایت در دایرکتوری خاص
./had web -url https://example.com -mode full -output ./backup

# پشتیبان‌گیری از صفحه تکی با همه دارایی‌ها
./had web -url https://example.com/about -mode single -download-external

# پشتیبان‌گیری با دارایی‌های CDN خارجی
./had web -url https://example.com -mode full -download-external -external-domains cdn.example.com,images.example.com

# خزش با کارایی بالا (10 کارگر همزمان)
./had web -url https://example.com -mode full -concurrency 10 -max-pages 500

# ادامه پشتیبان‌گیری قطع شده
./had web -url https://example.com -mode full -resume -output ./backup

# SPA با پشتیبانی از مسیریابی هش
./had web -url https://app.example.com/#!/home -mode full -crawl-hash-routes

# محدودیت اندازه دارایی و نرخ
./had web -url https://example.com -mode full -max-asset-size 20 -rate-limit 5
```

### دانلود Metalink

```bash
# دانلود از آدرس Metalink
./had -metalink https://example.com/file.metalink

# دانلود از فایل Metalink محلی
./had -metalink ./downloads/ubuntu.metalink4

# Metalink با دایرکتوری خروجی سفارشی
./had -metalink https://example.com/file.metalink -o ./downloads
```

### حالت سرور RPC

```bash
# شروع سرور RPC روی پورت پیش‌فرض
./had -rpc

# شروع سرور RPC روی آدرس سفارشی
./had -rpc -rpc-addr 0.0.0.0:6800

# فعال کردن WebSocket RPC (آزمایشی)
./had -rpc -rpc-websocket -rpc-addr :6800

# RPC با دایرکتوری دانلودها
./had -rpc -rpc-addr localhost:6800 -o /downloads
```

**نمونه درخواست‌های RPC:**

```bash
# دریافت اطلاعات نسخه
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"had.version","id":1}'

# دریافت آمار کلی
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"had.getGlobalStat","id":2}'

# اضافه کردن آدرس دانلود
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"had.addUri","params":{"uris":["https://example.com/file.zip"]},"id":3}'

# دریافت وضعیت همه فایل‌ها
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"had.tellAllStatus","id":4}'

# لیست همه متدهای موجود
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"system.listMethods","id":5}'

# توقف همه دانلودها
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"had.pauseAll","id":6}'

# تنظیم محدودیت سرعت به 5 مگابایت بر ثانیه
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"had.setSpeedLimit","params":{"speed":5242880},"id":7}'

# خاموش کردن had
curl -X POST http://localhost:6800/jsonrpc -d '{"method":"had.shutdown","id":8}'
```

**نقاط پایانی REST API:**

```bash
# دریافت وضعیت کلی
curl http://localhost:6800/api/status

# دریافت همه فایل‌ها
curl http://localhost:6800/api/files

# توقف همه دانلودها
curl http://localhost:6800/api/pause

# ادامه همه دانلودها
curl http://localhost:6800/api/resume

# اطلاعات نسخه
curl http://localhost:6800/api/version
```

### 🔐 دانلود از طریق پروکسی

```bash
# پروکسی SOCKS5
./had -proxy socks5://127.0.0.1:1080 https://example.com/file.zip

# پروکسی SOCKS5 با احراز هویت
./had -proxy socks5://user:pass@127.0.0.1:1080 https://example.com/file.zip

# پروکسی SOCKS4 با رشته‌های سفارشی
./had -proxy socks4://192.168.1.1:9050 -t 16 https://example.com/file.zip

# پروکسی HTTP
./had -proxy http://proxy.company.com:8080 https://example.com/file.zip

# پروکسی HTTPS با احراز هویت
./had -proxy https://user:pass@proxy.company.com:8080 https://example.com/file.zip
```

### 📡 دانلود FTP/SFTP

```bash
# FTP استاندارد
./had -protocol ftp ftp://example.com/file.zip

# FTP با مشخصات سفارشی
./had -protocol ftp -ftp-user myuser -ftp-pass mypass ftp://example.com/file.zip

# FTPS (FTP روی TLS)
./had -protocol ftps ftps://example.com/secure-file.zip

# SFTP با رمز عبور
./had -protocol sftp -sftp-user myuser -sftp-pass mypass sftp://example.com/file.zip

# SFTP با کلید SSH
./had -protocol sftp -ssh-key ~/.ssh/id_rsa sftp://example.com/file.zip

# SFTP با کلید SSH رمزگذاری شده
./had -protocol sftp -ssh-key ~/.ssh/id_rsa -ssh-key-pass mypassphrase sftp://example.com/file.zip
```

### 🕷️ وب اسکرپینگ

```bash
# استخراج و دانلود همه فایل‌ها از یک صفحه
./had -scrape https://example.com/downloads/

# فیلتر بر اساس پسوندها
./had -scrape https://example.com/downloads/ -ex .mp4,.mp3,.zip

# اسکرپینگ با رشته‌های سفارشی
./had -scrape https://example.com/files/ -t 16 -ex .pdf,.doc,.xls

# اسکرپینگ و دانلود با پیشرفت
./had -scrape https://example.com/media/ -ex .jpg,.png,.gif -v
```

### 🔄 آدرس‌های پارامتری

```bash
# جایگذاری عددی ساده
./had -parameterized-url 'https://example.com/file{}.zip' -start 1 -end 50

# جایگذاری با صفر تا ابتدا
./had -parameterized-url 'https://example.com/image{0}.jpg' -start 1 -end 100

# سه صفر تا ابتدا
./had -parameterized-url 'https://example.com/page{00}.html' -start 1 -end 500 -step 2

# اندازه گام سفارشی
./had -parameterized-url 'https://example.com/chunk{}.bin' -start 0 -end 200 -step 10
```

### 🔄 ادامه دانلودها

```bash
# ادامه از جلسه ذخیره شده
./had session_20231215_143022.json

# جلسه به طور خودکار با قطع شدن (Ctrl+C) ذخیره می‌شود
# پیشرفت هر 10 ثانیه به طور خودکار ذخیره می‌شود
```

### 🍪 پشتیبانی از کوکی

```bash
# بارگذاری کوکی از فایل فرمت Netscape (خروجی Firefox/Chrome)
./had -load-cookies cookies.txt https://example.com/private-file.zip

# ذخیره کوکی بعد از دانلود
./had -save-cookies output.txt https://example.com/file.zip

# رشته کوکی مستقیم
./had -c "sessionid=abc123; user=test" https://example.com/file.zip

# بارگذاری و ذخیره کوکی
./had -load-cookies cookies.txt -save-cookies newcookies.txt https://example.com/file.zip
```

### 🔐 احراز هویت NetRC

```bash
# استفاده از فایل .netrc برای احراز هویت
./had -netrc ~/.netrc https://example.com/private/file.zip

# فرمت فایل .netrc:
# machine example.com login myuser password mypass
# default login anonymous password user@example.com
```
