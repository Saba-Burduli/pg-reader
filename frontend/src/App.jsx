import { useEffect, useMemo, useState } from 'react'
import { Link, useLocation, useNavigate, useParams } from 'react-router-dom'
import './App.css'

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

export default function App() {
  const { id } = useParams()
  const navigate = useNavigate()
  const location = useLocation()

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
        if (!res.ok) throw new Error('Failed to load articles')
        return res.json()
      })
      .then((data) => {
        if (!canceled) {
          setArticles(data)
          setError('')
        }
      })
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
        if (!res.ok) throw new Error('Article not found')
        return res.json()
      })
      .then((data) => !canceled && setArticle(data))
      .catch((err) => !canceled && setError(err.message))
      .finally(() => !canceled && setLoadingArticle(false))
    return () => {
      canceled = true
    }
  }, [id])

  const selected = useMemo(() => articles.find((a) => a.id === id), [articles, id])

  return (
    <div className="app">
      <aside className="sidebar">
        <div className="brand">
          <h1>Paul Graham</h1>
          <p>Reader</p>
        </div>
        {loadingList ? <p className="muted">Loading essays...</p> : null}
        {error ? <p className="error">{error}</p> : null}
        <nav className="list">
          {articles.map((a) => (
            <Link key={a.id} className={a.id === id ? 'item active' : 'item'} to={`/article/${a.id}`}>
              <span>{a.title}</span>
              <small>{a.wordCount} words</small>
            </Link>
          ))}
        </nav>
      </aside>
      <main className="reader">
        {location.pathname !== '/' ? (
          <button className="back" onClick={() => navigate('/')}>
            All Essays
          </button>
        ) : null}
        {loadingArticle ? <p className="muted">Loading article...</p> : null}
        {article ? (
          <article>
            <h2>{article.title}</h2>
            <p className="meta">{selected?.wordCount ?? article.wordCount} words</p>
            <div className="content">{article.content}</div>
          </article>
        ) : (
          <section className="empty">
            <h2>Pick an essay</h2>
            <p>Dark mode, minimal UI, and fast local loading.</p>
          </section>
        )}
      </main>
    </div>
  )
}
