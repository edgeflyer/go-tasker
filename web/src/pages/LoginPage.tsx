import { useMemo, useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { login } from "../api/auth";
import { setAuth } from "../auth/auth";
import "../styles/auth.css";

function mapErrorMessage(code?: string, fallback?: string) {
  switch (code) {
    case "INVALID_CREDENTIALS":
      return "用户名或密码错误";
    case "INVALID_JSON":
      return "请求格式不正确";
    case "UNAUTHORIZED":
      return "登录已过期，请重新登录";
    default:
      return fallback || "登录失败，请重试";
  }
}

export default function LoginPage() {
  const nav = useNavigate();
  const wrapRef = useRef<HTMLDivElement>(null);
  const cardRef = useRef<HTMLDivElement>(null);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string>("");

  const canSubmit = useMemo(
    () => username.trim().length > 0 && password.length >= 6,
    [username, password]
  );

  function handlePointerMove(e: React.MouseEvent<HTMLDivElement>) {
    const wrap = wrapRef.current;
    const card = cardRef.current;
    if (!wrap) return;
    const rect = wrap.getBoundingClientRect();
    wrap.style.setProperty("--mx", `${e.clientX - rect.left}px`);
    wrap.style.setProperty("--my", `${e.clientY - rect.top}px`);
    if (card) {
      const cardRect = card.getBoundingClientRect();
      card.style.setProperty("--cx", `${e.clientX - cardRect.left}px`);
      card.style.setProperty("--cy", `${e.clientY - cardRect.top}px`);
    }
  }

  function handleCardClick(e: React.MouseEvent<HTMLDivElement>) {
    const card = cardRef.current;
    if (!card) return;
    const rect = card.getBoundingClientRect();
    card.style.setProperty("--rx", `${e.clientX - rect.left}px`);
    card.style.setProperty("--ry", `${e.clientY - rect.top}px`);
    card.classList.remove("card-ripple");
    // force reflow to restart animation
    void card.offsetWidth;
    card.classList.add("card-ripple");
  }

  function handleButtonRipple(e: React.PointerEvent<HTMLButtonElement>) {
    const btn = e.currentTarget;
    const rect = btn.getBoundingClientRect();
    btn.style.setProperty("--bx", `${e.clientX - rect.left}px`);
    btn.style.setProperty("--by", `${e.clientY - rect.top}px`);
    btn.classList.remove("btn-ripple");
    void btn.offsetWidth;
    btn.classList.add("btn-ripple");
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    if (!canSubmit) {
      setError("请输入用户名，密码至少 6 位");
      return;
    }

    try {
      setSubmitting(true);
      const res = await login(username.trim(), password);
      setAuth(res.token, res.user);
      nav("/tasks");
    } catch (err: any) {
      setError(mapErrorMessage(err?.code, err?.message));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div
      className="auth-wrap"
      ref={wrapRef}
      onMouseMove={handlePointerMove}
      onTouchMove={(e) => {
        const touch = e.touches[0];
        if (!touch) return;
        handlePointerMove({
          ...e,
          clientX: touch.clientX,
          clientY: touch.clientY,
        } as any);
      }}
    >
      <div className="auth-card" ref={cardRef} onClick={handleCardClick}>
        <h1 className="auth-title">登录</h1>
        <p className="auth-subtitle">使用你的账号登录后管理自己的任务</p>

        <form onSubmit={handleSubmit}>
          <div className="auth-field">
            <div className="auth-label">用户名</div>
            <input
              className="auth-input"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="例如：nullix"
              autoComplete="username"
            />
          </div>

          <div className="auth-field">
            <div className="auth-label">密码</div>
            <input
              className="auth-input"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="至少 6 位"
              autoComplete="current-password"
            />
          </div>

          <button
            className="auth-btn"
            disabled={!canSubmit || submitting}
            onPointerDown={handleButtonRipple}
          >
            {submitting ? "登录中..." : "登录"}
          </button>
        </form>

        <div className="auth-row">
          <span style={{ opacity: 0.85, fontSize: 13 }}>还没有账号？</span>
          <Link className="auth-link" to="/register">
            去注册 →
          </Link>
        </div>

        {error && <div className="auth-error">{error}</div>}
      </div>
    </div>
  );
}
