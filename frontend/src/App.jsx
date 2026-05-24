import { useEffect, useMemo, useRef, useState } from 'react'
import { Link, useLocation, useNavigate, useParams } from 'react-router-dom'
import './App.css'

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

export default function App() {
  const { id } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const readerRef = useRef(null)

  const [articles, setArticles] = useState([])
  const [article, setArticle] = useState(null)
  const [loadingList, setLoadingList] = useState(true)
  const [loadingArticle, setLoadingArticle] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    let canceled = false
    setLoadingList(true)
    fetch(`${API_BASE}/articles`)
      .then((res) => {
        if (!res.ok) throw new Error('Failed to load essays')
        return res.json()
      })
      .then((data) => !canceled && setArticles(data))
      .catch((err) => !canceled && setError(err.message))
      .finally(() => !canceled && setLoadingList(false))

    return () => {
      canceled = true
    }
  }, [])

  useEffect(() => {
    if (!id) {
      setArticle(null)
      return
    }
    let canceled = false
    setLoadingArticle(true)
    fetch(`${API_BASE}/articles/${id}`)
      .then((res) => {
        if (!res.ok) throw new Error('Essay not found')
        return res.json()
      })
      .then((data) => !canceled && setArticle(data))
      .catch((err) => !canceled && setError(err.message))
      .finally(() => !canceled && setLoadingArticle(false))
    return () => {
      canceled = true
    }
  }, [id])

  useEffect(() => {
    const node = readerRef.current
    if (!node || !article || article.isRead) return

    const onScroll = () => {
      const reachedBottom = node.scrollTop + node.clientHeight >= node.scrollHeight - 24
      if (!reachedBottom) return
      fetch(`${API_BASE}/articles/${article.id}/read`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ isRead: true }),
      })
        .then((res) => {
          if (!res.ok) throw new Error('Failed to update read status')
          return res.json()
        })
        .then((updated) => {
          setArticle(updated)
          setArticles((prev) => prev.map((item) => (item.id === updated.id ? { ...item, isRead: true } : item)))
        })
        .catch(() => {})
    }

    node.addEventListener('scroll', onScroll, { passive: true })
    return () => node.removeEventListener('scroll', onScroll)
  }, [article])

  const stats = useMemo(() => {
    const total = articles.length
    const read = articles.filter((a) => a.isRead).length
    const remaining = total - read
    const percent = total === 0 ? 0 : Math.round((read / total) * 100)
    return { total, read, remaining, percent }
  }, [articles])

  const unreadEssays = useMemo(() => articles.filter((a) => !a.isRead), [articles])
  const completedEssays = useMemo(() => articles.filter((a) => a.isRead), [articles])

  return (
    <div className="app">
      <aside className="sidebar">
        <div className="brand">
          <h1>Paul Graham</h1>
          <p>Personal Reading Library</p>
        </div>

        <section className="stats">
          <div className="stats-row">
            <span>Total</span>
            <strong>{stats.total}</strong>
          </div>
          <div className="stats-row">
            <span>Read</span>
            <strong>{stats.read}</strong>
          </div>
          <div className="stats-row">
            <span>Remaining</span>
            <strong>{stats.remaining}</strong>
          </div>
          <div className="progress">
            <div className="progress-fill" style={{ width: `${stats.percent}%` }} />
          </div>
          <small>{stats.percent}% complete</small>
        </section>

        {loadingList ? <p className="muted">Loading essays...</p> : null}
        {error ? <p className="error">{error}</p> : null}

        <section className="group">
          <h3>Unread</h3>
          <nav className="list">
            {unreadEssays.map((essay) => (
              <EssayRow key={essay.id} essay={essay} activeId={id} />
            ))}
          </nav>
        </section>

        <section className="group">
          <h3>Completed</h3>
          <nav className="list">
            {completedEssays.map((essay) => (
              <EssayRow key={essay.id} essay={essay} activeId={id} />
            ))}
          </nav>
        </section>
      </aside>

      <main ref={readerRef} className="reader">
        {location.pathname !== '/' ? (
          <button className="back" onClick={() => navigate('/')}>
            All Essays
          </button>
        ) : null}

        {loadingArticle ? <p className="muted">Loading essay...</p> : null}
        {article ? (
          <article>
            <h2>{article.title}</h2>
            <p className="meta">
              {formatDate(article.publishedAt)} · {article.wordCount} words · {article.isRead ? 'Read' : 'Unread'}
            </p>
            <div className="content">{article.content}</div>
          </article>
        ) : (
          <section className="empty">
            <h2>Pick an essay</h2>
            <p>Read to the end to move it to completed.</p>
          </section>
        )}
      </main>
    </div>
  )
}

function EssayRow({ essay, activeId }) {
  return (
    <Link className={`item ${essay.id === activeId ? 'active' : ''} ${essay.isRead ? 'done' : ''}`} to={`/article/${essay.id}`}>
      <div className="item-title">{essay.isRead ? '✓ ' : '○ '}{essay.title}</div>
      <small>{formatDate(essay.publishedAt)} · {essay.wordCount} words</small>
    </Link>
  )
}

function formatDate(raw) {
  const date = new Date(raw)
  if (Number.isNaN(date.getTime())) return 'Unknown date'
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: '2-digit',
  }).format(date)
}
