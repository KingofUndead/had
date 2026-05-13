# حرفه‌ای Smart Downloader

یک ابزار قدرتمند دانلودر خودکار در گیت‌هاب اکشنز برای دانلود از یوتیوب، تلگرام، bunkr و لینک‌های مستقیم با قابلیت آپلود روی uploadkon.ir و روبیکا

---

**[📥 مشاهده لیست فایل‌های دانلود شده](./Links.md)**

[**🎥 ویدیو آموزش اجرا کردن**](https://up.20script.ir/do.php?filename=ebcf-VivaCut-video-1778618372567.mp4)

---

## ✨ قابلیت‌ها

- 🎬 **دانلود از یوتیوب** با دو روش:
  - **Fast API** (بدون نیاز به کوکی، سریع‌ترین روش)
  - **yt-dlp** (روش قابل اعتماد با پشتیبانی از کوکی)
- 📨 **دانلود از تلگرام** (با استفاده از String Session)
- 📦 **پشتیبانی از bunkr.si** و دامنه‌های مشابه
- 🔗 **دانلود لینک مستقیم** با aria2
- 🗜️ **تقسیم فایل‌های حجیم** به قطعات ۹۵ مگابایتی
- 🔐 **رمزگذاری فایل‌ها** با پسورد دلخواه
- ☁️ **آپلود خودکار روی uploadkon.ir**
- 📱 **آپلود خودکار روی روبیکا** (پیام‌های ذخیره شده)
- 📄 **تولید لینک دانلود** در فایل Links.md

---

## 🛠️ پیش‌نیازها

- یک اکانت **گیت‌هاب** (رایگان)
- مرورگر **کروم** یا **فایرفاکس** (برای اکستنشن کوکی)

---

## 🚀 نصب و راه‌اندازی

### مرحله ۱: فورک کردن پروژه

۱. به [این لینک](https://github.com/Mr-Spect3r/had) برو
۲. روی دکمه **Fork** در بالا سمت راست کلیک کن
۳. ریپازیتوری برای خودت ساخته می‌شه

### مرحله ۲: اضافه کردن فایل workflow

فایل `had.yml` قبلاً در مسیر `.github/workflows/` قرار داره. اگر نیست، خودت بساز:

1. توی ریپازیتوری فورک شده، پوشه `.github/workflows/` رو ایجاد کن
2. فایل `had.yml` رو با محتوای پروژه داخلش بذار

---

## 🍪 راه‌اندازی کوکی یوتیوب (فقط برای روش yt-dlp)

### نصب اکستنشن Cookie Manager

**مشخصات اکستنشن:**
- نام: `extensions-had`
- مسیر: داخل پوشه پروژه
- قابلیت‌ها: 🛡️ پروکسی MITM | 🍪 اکسپورت کوکی Netscape

**نصب در کروم:**
1. برو به `chrome://extensions/`
2. **حالت توسعه‌دهنده** رو روشن کن
3. روی **Load unpacked** کلیک کن
4. پوشه `extensions-had` رو انتخاب کن

**نصب در فایرفاکس:**
1. برو به `about:debugging`
2. روی **This Firefox** کلیک کن
3. روی **Load Temporary Add-on** کلیک کن
4. فایل `manifest.json` داخل پوشه `extensions-had` رو انتخاب کن

### گرفتن کوکی

1. وارد یوتیوب بشو و لاگین کن
2. روی آیکون اکستنشن کلیک کن
3. گزینه **Export Cookies** → **Netscape format** رو بزن
4. فایل `cookies.txt` ذخیره می‌شه

### اضافه کردن به گیت‌هاب سیکرت

1. برو به **Settings** → **Secrets and variables** → **Actions**
2. روی **New repository secret** کلیک کن
3. **Name**: `YOUTUBE_COOKIES`
4. **Value**: محتوای فایل `cookies.txt` رو کامل پیست کن
5. روی **Add secret** کلیک کن

> 💡 **نکته:** اگه از روش **Fast API** استفاده می‌کنی، **نیازی به کوکی نیست**!

---

## 📨 ساخت String Session تلگرام

### نصب کتابخونه

```python
pip install telethon
```

### 📨 ساخت String Session تلگرام

برای دانلود از تلگرام، نیاز به String Session داری.

#### مرحله ۱: گرفتن api_id و api_hash

1. به سایت [my.telegram.org](https://my.telegram.org) برو
2. با شماره تلفن خودت وارد شو
3. روی **API Development Tools** کلیک کن
4. اگر قبلاً اپلیکیشنی نساختی، روی **Create Application** کلیک کن
5. هر اسمی میخوای وارد کن (مثلاً MyDownloader)
6. بعد از ساختن، مقادیر **api_id** و **api_hash** رو کپی کن

> 📸 **api_id** یک عدد است مثل `1234567`
> 📸 **api_hash** یک متن طولانی است مثل `7a8b9c0d1e2f3g4h5i6j7k8l9m0n1o2p`

ساخته نشد یا نتونستی بگیری؟ برو سراغ این بات داخل تلگرام تا برات بسازه و بگیریش: @RailBot

برای دانلود از تلگرام، نیاز به String Session داری. کد زیر رو اجرا کن:

```python
from telethon import TelegramClient
from telethon.sessions import StringSession
import asyncio

api_id = 12345  # عدد API ID خودتون رو وارد کنید
api_hash = 'your_api_hash_here'  # متن API Hash خودتون رو وارد کنید

async def generate():
    client = TelegramClient(StringSession(), api_id, api_hash)
    await client.start()
    string_session = client.session.save()
    print("\n" + "="*50)
    print("🔐 رشته سشن شما:")
    print(string_session)
    print("="*50 + "\n")
    
    # ذخیره در فایل
    with open('string_session.txt', 'w') as f:
        f.write(string_session)
    print("✅ در فایل 'string_session.txt' ذخیره شد.")
    await client.disconnect()

asyncio.run(generate())
```

> ⚠️ **عدد api_id و api_hash رو توی کد جایگزین کن بعد اجرا کن**

---

## 📱 ساخت Auth String روبیکا

### نصب کتابخونه
```python
pip install rubpy
```

برای آپلود روی روبیکا، نیاز به فایل احراز هویت داری.

### مرحله ۱: ساخت فایل rubpy.rp , فشرده‌سازی و تبدیل به متن

اول با کتابخونه `rubpy` وارد حساب روبیکات بشو و فایل سشن رو بساز و کپیش کن:

```python
from rubpy import Client
import lzma
import base64

client = Client(name="rubpy")
client.start()
print("✅ وارد شدید! فایل rubpy.rp ساخته شد.")

with open("rubpy.rp", "rb") as f:
    compressed = lzma.compress(f.read(), preset=9)
    auth_string = base64.b64encode(compressed).decode()
    
print("\n" + "="*50)
print("🔐 رشته احراز هویت روبیکا (این رو ذخیره کن):")
print(auth_string)
print("="*50 + "\n")

# ذخیره در فایل
with open("rubika_auth.txt", "w") as f:
    f.write(auth_string)
print("✅ در فایل 'rubika_auth.txt' ذخیره شد.")
```

> ⚠️ **این رشته رو موقع اجرا توی فیلد `rubika_auth_string` وارد کن**

---

## 🎮 نحوه استفاده

### اجرای دانلودر

1. توی ریپازیتوری خودت برو به تب **Actions**
2. روی workflow **Professional Smart Downloader** کلیک کن
3. روی دکمه **Run workflow** کلیک کن
4. فرم رو پر کن:

| فیلد | توضیح | مثال |
|------|-------|------|
| **URL** | لینک مورد نظر (چندتا با کاما یا فاصله جدا کن) | `https://youtu.be/xxxx, https://t.me/xxx/1` |
| **YouTube Method** | روش دانلود یوتیوب | `Fast API` (توصیه می‌شه) یا `yt-dlp` |
| **Quality** | کیفیت ویدیو | `480p` یا `720p` یا `Audio` |
| **Password** | رمز فایل (اختیاری) | `mysecret123` |
| **Split Size** | حداکثر حجم هر قطعه (مگابایت) | `95` |
| **Upload to uploadkon** | آپلود روی سرور خارجی | `true` یا `false` |
| **Upload to Rubika** | آپلود روی روبیکا | `true` یا `false` |
| **Rubika Auth String** | رشته احراز هویت روبیکا | (از مرحله قبل) |
| **Telegram String Session** | رشته سشن تلگرام | (از مرحله قبل) |

> 💡 **توصیه:** 
> - برای یوتیوب از **Fast API** استفاده کن (نیاز به کوکی نداره)
> - فقط یکی از گزینه‌های آپلود رو فعال کن

---

## 📥 خروجی

بعد از اتمام اجرا:

- فایل‌ها توی پوشه `dl/` ذخیره می‌شن
- لینک‌های دانلود توی فایل `Links.md` قرار می‌گیره
- اگه آپلود روی uploadkon رو فعال کرده باشی، لینک دانلود + لینک حذف
- اگه آپلود روی روبیکا رو فعال کرده باشی، فایل توی Saved Messages ذخیره می‌شه

---

## 🔧 عیب‌یابی

### خطای `Sign in to confirm you’re not a bot`

**راه حل:** از روش **Fast API** استفاده کن.

### Fast API کار نمی‌کنه

**راه حل:** به روش **yt-dlp** سوییچ کن و کوکی رو اضافه کن.

### خطای lzma در روبیکا

**راه حل:** نیازی به نصب lzma نیست، کتابخانه داخلی پایتونه. فقط `rubpy` رو نصب کن.

### تلگرام دانلود نمی‌کنه

**راه حل:** 
1. String Session رو دوباره بگیر
2. مطمئن شو api_id و api_hash درست باشن
3. لینک تلگرام درست باشه

---

## 📝 نکات مهم

- ✅ روش **Fast API** برای ۹۰٪ ویدیوهای یوتیوب جواب میده و نیازی به کوکی نداره
- ✅ حجم هر فایل آپلودی روی گیت‌هاب نباید بیشتر از **۹۵ مگابایت** باشه (خودکار تقسیم می‌شه)
- ✅ می‌تونی چندین لینک رو همزمان وارد کنی (با فاصله یا کاما جدا کن)
- ✅ آپلود روی uploadkon.ir محدودیت حجم نداره
- ✅ آپلود روی روبیکا محدودیت حجم نداره (به Saved Messages می‌ره)

---

## 🤝 مشارکت

اگه باگی پیدا کردی یا ایده‌ای برای بهبود داری، یه **Issue** یا **Pull Request** باز کن.

---

## 🔗 لینک‌های مفید

- [پروژه اصلی had](https://github.com/Mr-Spect3r/had)
- [گرفتن api_id از my.telegram.org](https://my.telegram.org)
- [آموزش ویدئویی](https://youtu.be/vqHtozeY7iA)

---

⭐ اگر از این ابزار استفاده می‌کنی، بهش **ستاره** بده!
