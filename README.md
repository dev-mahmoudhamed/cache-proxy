# ⚡ Proxy Cache Lab

**A simple tool to see how caching makes websites faster**

🌐 **Try it here:** [https://cache-proxy.leapcell.app/](https://cache-proxy.leapcell.app/)

---

## What does it do?

Enter any URL and see the difference between:

- **First request** (downloading fresh) → Slow ❌
- **Second request** (served from cache) → Super fast ✅

### Example

```
First time visiting https://github.com
Status: ✗ MISS
Duration: 245ms
```

```
Second time visiting https://github.com
Status: ✓ HIT
Duration: 2ms
```

**That's 122x faster!** 🚀

---

## How it works

1. Enter a URL
2. App downloads the page and saves it to cache
3. Request the same URL again
4. See how much faster it loads from cache

**Green = Fast (cached)**  
**Red = Slow (not cached)**

---

## Features

✅ Real-time performance tracking  
✅ Private cache for each user  
✅ Full request history  
✅ Works on mobile and desktop

---

## Run it locally

```bash
# Clone the repo
git clone https://github.com/dev-mahmoudhamed/cache-proxy.git
cd cache-proxy

# Run the app
go run main.go

# Open browser
http://localhost:8080
```

**Requirements:** Go 1.21+

---

## Why is this useful?

- 🎓 Learn how caching works
- 📊 Show the impact of caching to others
- 🔧 Test your own URLs and APIs
- ⚡ Understand web performance

---

## License

MIT License - Free to use and modify
