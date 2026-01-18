# Shrink - URL Shortener API

A Go-based URL shortener API built with Gin and MongoDB that tracks user clicks with geolocation data.

## Features

- ✅ Create shortened URLs with custom aliases
- ✅ Automatic redirect and click tracking
- ✅ Browser geolocation support
- ✅ IP-based geolocation fallback
- ✅ Complete click analytics (IP, location, browser, device, referrer)
- ✅ View all URL details and click history

## API Endpoints

### Create Shortened URL
```bash
POST /api/shrink
Content-Type: application/json

{
  "original_url": "https://example.com/very/long/url",
  "custom_alias": "my-link"  # optional
}
```

### Redirect to Original URL
```bash
GET /:shortUrl
# Automatically redirects and tracks the click
```

With browser geolocation:
```bash
GET /:shortUrl?lat=40.7128&lng=-74.0060&country=United%20States&city=New%20York&region=NY
```

### Get URL Details & Analytics
```bash
GET /info/:shortUrl
```

### Get URL Statistics
```bash
GET /api/stats/:shortUrl
```

### Health Check
```bash
GET /health
```

## Deployment to Render

### Prerequisites
- GitHub account with this repository
- Render account (https://render.com)
- MongoDB Atlas URI

### Steps

1. **Push to GitHub**
```bash
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/YOUR_USERNAME/shrink.git
git push -u origin main
```

2. **Create New Web Service on Render**
   - Go to https://dashboard.render.com
   - Click "New +" → "Web Service"
   - Connect your GitHub repository
   - Configure:
     - **Name**: `shrink-api`
     - **Runtime**: `Go`
     - **Build Command**: `cd backend && go build -o app .`
     - **Start Command**: `cd backend && ./app`
     - **Plan**: Free (or your preferred plan)

3. **Add Environment Variables**
   - In Render dashboard, go to your service
   - Click "Environment"
   - Add variable:
     - **Key**: `MONGODB_URI`
     - **Value**: Your MongoDB Atlas connection string
     - **Scope**: Build and Runtime

4. **Deploy**
   - Click "Deploy"
   - Wait for build to complete
   - Your API will be live at: `https://YOUR_SERVICE_NAME.onrender.com`

## Environment Variables

- `MONGODB_URI` - MongoDB connection string (Atlas or local)
- `PORT` - Server port (default: 8080)

## Local Development

### Install Dependencies
```bash
cd backend
go mod download
```

### Set Environment Variables
```bash
export MONGODB_URI="your_mongodb_uri"
export PORT=8080
```

### Run
```bash
go run main.go
```

### Build
```bash
go build -o app .
./app
```

## Project Structure
```
backend/
├── main.go           # Server and MongoDB connection
├── models/
│   └── url.go        # Data models
├── routes/
│   └── url.go        # API handlers
├── go.mod
├── go.sum
└── Dockerfile        # Container configuration
```

## Data Models

### URL
- `id` - MongoDB ObjectID
- `original_url` - Long URL
- `short_url` - Shortened code
- `custom_alias` - Custom short code (optional)
- `total_clicks` - Click counter
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp
- `is_active` - Active status
- `clicks` - Array of click records

### Click
- `id` - Click record ID
- `timestamp` - When clicked
- `ip` - User IP address
- `country` - Country name
- `city` - City name
- `region` - State/region
- `latitude` - Geo coordinate
- `longitude` - Geo coordinate
- `browser` - Browser name
- `browser_version` - Browser version
- `os` - Operating system
- `os_version` - OS version
- `device_type` - mobile/tablet/desktop
- `user_agent` - Full user agent string
- `referrer` - HTTP referrer

## Example Workflow

1. **Create a short URL**
```bash
curl -X POST https://YOUR_API.onrender.com/api/shrink \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'
```

Response:
```json
{
  "id": "696d23f4c17842a42503b214",
  "original_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "short_url": "jA4xbD",
  "total_clicks": 0,
  "created_at": "2026-01-18T18:18:28Z",
  "is_active": true
}
```

2. **Share the short URL**
```
https://YOUR_API.onrender.com/jA4xbD
```

3. **View analytics**
```bash
curl https://YOUR_API.onrender.com/info/jA4xbD
```

## Geolocation Data

The API captures geolocation in two ways:

1. **Browser Geolocation API** (Most accurate)
   - Pass coordinates via query parameters
   - Example: `/?lat=40.7128&lng=-74.0060&country=USA&city=NY&region=NY`

2. **IP Geolocation** (Fallback)
   - Automatically resolves country, city, region from IP
   - Uses ipapi.co service

## Notes

- Free Render plans have limitations (15-min auto-sleep after inactivity)
- For production, consider a paid plan
- MongoDB Atlas free tier includes 512MB storage
- Local testing shows "Local/Testing" for localhost IP (127.0.0.1)

## License

MIT
