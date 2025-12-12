import { useMemo, useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { register } from "../api/auth";
import "../styles/auth.css";

function mapErrorMessage(code?: string, fallback?: string) {
  switch (code) {
    case "USERNAME_EXISTS":
      return "用户名已存在，请换一个";
    case "INVALID_USERNAME":
      return "用户名不能为空";
    case "INVALID_PASSWORD":
      return "密码至少 6 位";
    default:
      return fallback || "注册失败，请重试";
  }
}

export default function RegisterPage() {
  const nav = useNavigate();
  const wrapRef = useRef<HTMLDivElement>(null);
  const cardRef = useRef<HTMLDivElement>(null);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");

  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const canSubmit = useMemo(() => {
    if (username.trim().length === 0) return false;
    if (password.length < 6) return false;
    if (password !== confirm) return false;
    return true;
  }, [username, password, confirm]);

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
    setSuccess("");

    if (!canSubmit) {
      if (password !== confirm) setError("两次密码不一致");
      else setError("用户名不能为空，密码至少 6 位");
      return;
    }

    try {
      setSubmitting(true);
      await register(username.trim(), password);
      setSuccess("注册成功！即将跳转到登录页…");
      setTimeout(() => nav("/login"), 800);
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
        <h1 className="auth-title">注册</h1>
        <p className="auth-subtitle">创建账号后即可登录管理任务</p>

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
              autoComplete="new-password"
            />
          </div>

          <div className="auth-field">
            <div className="auth-label">确认密码</div>
            <input
              className="auth-input"
              type="password"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              placeholder="再次输入密码"
              autoComplete="new-password"
            />
          </div>

          <button
            className="auth-btn"
            disabled={!canSubmit || submitting}
            onPointerDown={handleButtonRipple}
          >
            {submitting ? "注册中..." : "注册"}
          </button>
        </form>

        <div className="auth-row">
          <span style={{ opacity: 0.85, fontSize: 13 }}>已有账号？</span>
          <Link className="auth-link" to="/login">
            去登录 →
          </Link>
        </div>

        {error && <div className="auth-error">{error}</div>}
        {success && <div className="auth-success">{success}</div>}
      </div>
    </div>
  );
}
