import { useState } from 'react'
import './ShortenUrl.css'

export default function ShortenUrl({ onUrlCreated, createdUrl }) {
  const [longUrl, setLongUrl] = useState('')
  const [customAlias, setCustomAlias] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [shortUrl, setShortUrl] = useState('')
  const [copied, setCopied] = useState(false)

  const handleShortenUrl = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    setShortUrl('')
    setCopied(false)

    try {
      const payload = {
        original_url: longUrl,
      }
      if (customAlias) {
        payload.custom_alias = customAlias
      }

      console.log('Sending request to:', `${import.meta.env.VITE_API_URL}/api/shrink`)
      console.log('Payload:', payload)

      const response = await fetch(`${import.meta.env.VITE_API_URL}/api/shrink`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      })

      console.log('Response status:', response.status)
      const data = await response.json()
      console.log('Response data:', data)

      if (!response.ok) {
        const errorData = data
        throw new Error(errorData.message || 'Failed to shorten URL')
      }

      setShortUrl(data.short_url)
      onUrlCreated(data)
      setLongUrl('')
      setCustomAlias('')
    } catch (err) {
      console.error('Error:', err)
      setError(err.message || 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  const copyToClipboard = () => {
    navigator.clipboard.writeText(shortUrl)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="shorten-url-container">
      <form onSubmit={handleShortenUrl} className="shorten-form">
        <div className="form-group">
          <label htmlFor="longUrl">Enter Your Long URL</label>
          <input
            id="longUrl"
            type="url"
            placeholder="https://example.com/very/long/url/that/needs/shortening"
            value={longUrl}
            onChange={(e) => setLongUrl(e.target.value)}
            required
            className="url-input"
            disabled={loading}
          />
        </div>

        <div className="form-group">
          <label htmlFor="customAlias">Custom Alias (Optional)</label>
          <input
            id="customAlias"
            type="text"
            placeholder="my-link (leave empty for auto-generated)"
            value={customAlias}
            onChange={(e) => setCustomAlias(e.target.value)}
            className="url-input"
            disabled={loading}
          />
        </div>

        <button type="submit" disabled={loading} className="shorten-button">
          {loading ? 'Shortening...' : 'Create Short URL'}
        </button>
      </form>

      {error && (
        <div className="alert alert-error">
          {error}
        </div>
      )}

      {shortUrl && (
        <div className="result-card">
          <div className="result-header">
            <h3>Success!</h3>
          </div>

          <div className="result-content">
            <div className="result-item">
              <label>Original URL:</label>
              <div className="url-display">
                <a
                  href={longUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="original-url"
                >
                  {longUrl}
                </a>
              </div>
            </div>

            <div className="result-item">
              <label>Shortened URL:</label>
              <div className="short-url-wrapper">
                <input
                  type="text"
                  value={shortUrl}
                  readOnly
                  className="short-url-display"
                />
                <button
                  type="button"
                  onClick={copyToClipboard}
                  className={`copy-btn ${copied ? 'copied' : ''}`}
                  title="Copy to clipboard"
                >
                  {copied ? 'Copied!' : 'Copy'}
                </button>
              </div>
            </div>

            <div className="result-footer">
              <a
                href={shortUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="visit-btn"
              >
                Visit Short URL
              </a>
            </div>
          </div>
        </div>
      )}

      {createdUrl && !shortUrl && (
        <div className="info-card">
          <p>URLs created in this session will appear here. Create a new one to get started!</p>
        </div>
      )}
    </div>
  )
}
