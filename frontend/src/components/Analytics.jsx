import { useState } from 'react'
import './Analytics.css'

export default function Analytics() {
  const [shortUrl, setShortUrl] = useState('')
  const [stats, setStats] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleGetStats = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    setStats(null)

    try {
      console.log('Fetching stats for:', shortUrl)
      console.log('URL:', `${import.meta.env.VITE_API_URL}/api/stats/${shortUrl}`)

      const response = await fetch(`${import.meta.env.VITE_API_URL}/api/stats/${shortUrl}`)

      console.log('Response status:', response.status)
      console.log('Response ok:', response.ok)

      if (!response.ok) {
        throw new Error('URL not found or no analytics available')
      }

      const data = await response.json()
      console.log('Stats data:', data)
      setStats(data)
    } catch (err) {
      console.error('Error fetching stats:', err)
      setError(err.message || 'Failed to fetch analytics')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="analytics-container">
      <form onSubmit={handleGetStats} className="analytics-form">
        <div className="form-group">
          <label htmlFor="shortUrl">Enter Short URL Code</label>
          <input
            id="shortUrl"
            type="text"
            placeholder="e.g., jA4xbD or my-link"
            value={shortUrl}
            onChange={(e) => setShortUrl(e.target.value)}
            required
            className="url-input"
            disabled={loading}
          />
        </div>

        <button type="submit" disabled={loading} className="search-button">
          {loading ? 'Fetching Analytics...' : 'Get Analytics'}
        </button>
      </form>

      {error && (
        <div className="alert alert-error">
          {error}
        </div>
      )}

      {stats && (
        <div className="stats-container">
          <div className="stats-header">
            <h3>Analytics for {shortUrl}</h3>
          </div>

          <div className="stats-grid">
            <div className="stat-card">
              <div className="stat-content">
                <div className="stat-label">Total Clicks</div>
                <div className="stat-value">{stats.total_clicks || 0}</div>
              </div>
            </div>

            <div className="stat-card">
              <div className="stat-content">
                <div className="stat-label">Countries</div>
                <div className="stat-value">
                  {getUniqueCountries(stats.clicks).length}
                </div>
              </div>
            </div>

            <div className="stat-card">
              <div className="stat-content">
                <div className="stat-label">Cities</div>
                <div className="stat-value">
                  {getUniqueCities(stats.clicks).length}
                </div>
              </div>
            </div>

            <div className="stat-card">
              <div className="stat-content">
                <div className="stat-label">Created</div>
                <div className="stat-value">
                  {new Date(stats.created_at).toLocaleDateString()}
                </div>
              </div>
            </div>
          </div>

          {stats.clicks && stats.clicks.length > 0 && (
            <div className="detailed-stats">
              <h4>Click Details</h4>

              <div className="stats-section">
                <div className="section-title">Geolocation Breakdown</div>
                <div className="geo-list">
                  {getGeoStats(stats.clicks).map((geo, idx) => (
                    <div key={idx} className="geo-item">
                      <div className="geo-info">
                        <span className="geo-location">
                          {geo.city && geo.country
                            ? `${geo.city}, ${geo.country}`
                            : geo.country || 'Unknown'}
                        </span>
                        <span className="geo-count">{geo.count} clicks</span>
                      </div>
                      <div className="geo-bar">
                        <div
                          className="geo-fill"
                          style={{
                            width: `${(geo.count / stats.total_clicks) * 100}%`,
                          }}
                        ></div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <div className="stats-section">
                <div className="section-title">Browser Information</div>
                <div className="browser-list">
                  {getBrowserStats(stats.clicks).map((browser, idx) => (
                    <div key={idx} className="browser-item">
                      <span className="browser-name">{browser.name}</span>
                      <span className="browser-count">{browser.count}</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className="stats-section">
                <div className="section-title">Device Types</div>
                <div className="device-grid">
                  {getDeviceStats(stats.clicks).map((device, idx) => (
                    <div key={idx} className="device-card">
                      <span className="device-label">
                        {device.type.charAt(0).toUpperCase() + device.type.slice(1)}
                      </span>
                      <span className="device-count">{device.count}</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className="stats-section">
                <div className="section-title">Recent Clicks</div>
                <div className="recent-clicks">
                  {stats.clicks.slice(0, 10).map((click, idx) => (
                    <div key={idx} className="click-item">
                      <div className="click-time">
                        {new Date(click.timestamp).toLocaleString()}
                      </div>
                      <div className="click-details">
                        <span className="click-location">
                          {click.city}, {click.country}
                        </span>
                        <span className="click-device">{click.browser}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}

          {(!stats.clicks || stats.clicks.length === 0) && (
            <div className="info-card">
              <span className="info-icon">ℹ️</span>
              <p>No clicks recorded yet for this URL.</p>
            </div>
          )}
        </div>
      )}

      {!stats && !error && (
        <div className="info-card">
          <span className="info-icon">💡</span>
          <p>Enter a short URL code to view its analytics and detailed click statistics.</p>
        </div>
      )}
    </div>
  )
}

// Helper functions
function getUniqueCountries(clicks) {
  if (!clicks) return []
  return [...new Set(clicks.map((c) => c.country))].filter(Boolean)
}

function getUniqueCities(clicks) {
  if (!clicks) return []
  return [...new Set(clicks.map((c) => c.city))].filter(Boolean)
}

function getGeoStats(clicks) {
  if (!clicks) return []
  const geoMap = {}
  clicks.forEach((click) => {
    const key = `${click.city}-${click.country}`
    geoMap[key] = (geoMap[key] || 0) + 1
  })
  return Object.entries(geoMap)
    .map(([key, count]) => {
      const [city, country] = key.split('-')
      return { city, country, count }
    })
    .sort((a, b) => b.count - a.count)
}

function getBrowserStats(clicks) {
  if (!clicks) return []
  const browserMap = {}
  clicks.forEach((click) => {
    const browser = click.browser || 'Unknown'
    browserMap[browser] = (browserMap[browser] || 0) + 1
  })
  return Object.entries(browserMap)
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
}

function getDeviceStats(clicks) {
  if (!clicks) return []
  const deviceMap = {}
  clicks.forEach((click) => {
    const device = click.device_type || 'unknown'
    deviceMap[device] = (deviceMap[device] || 0) + 1
  })
  return Object.entries(deviceMap).map(([type, count]) => ({ type, count }))
}
