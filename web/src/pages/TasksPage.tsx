import { useEffect, useMemo, useRef, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { clearAuth, getUser } from "../auth/auth";
import type { Task, TaskStatus, TaskPage } from "../api/tasks";
import { createTask, deleteTask, listTasks, updateTask } from "../api/tasks";
import "../styles/tasks.css";

type StatusFilter = "all" | TaskStatus;
type SortOption = "created_desc" | "created_asc" | "status";

function parseStatus(val: string | null): StatusFilter {
  return val === "pending" || val === "completed" ? val : "all";
}

export default function TasksPage() {
  const nav = useNavigate();
  const user = getUser();
  const [searchParams, setSearchParams] = useSearchParams();

  const wrapRef = useRef<HTMLDivElement>(null);
  const heroRef = useRef<HTMLDivElement>(null);

  // 查询状态 & URL 同步
  const [status, setStatus] = useState<StatusFilter>(() =>
    parseStatus(searchParams.get("status"))
  );
  const [page, setPage] = useState<number>(() => {
    const p = Number(searchParams.get("page") || 1);
    return Number.isNaN(p) || p <= 0 ? 1 : p;
  });
  const [pageSize, setPageSize] = useState<number>(() => {
    const ps = Number(searchParams.get("page_size") || 10);
    return Number.isNaN(ps) || ps <= 0 ? 10 : ps;
  });
  const [sort, setSort] = useState<SortOption>(() => {
    const s = searchParams.get("sort");
    return s === "created_asc" || s === "status" ? s : "created_desc";
  });
  const [keyword, setKeyword] = useState<string>(() => searchParams.get("q") || "");

  // 数据
  const [data, setData] = useState<TaskPage>({
    items: [],
    total: 0,
    page: 1,
    page_size: 10,
  });

  // UI 状态
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // 新建任务表单
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const canCreate = useMemo(() => title.trim().length > 0, [title]);

  const totalPages = useMemo(() => {
    const ps = data.page_size || pageSize;
    return Math.max(1, Math.ceil((data.total || 0) / ps));
  }, [data.total, data.page_size, pageSize]);

  // 将状态写入 URL
  useEffect(() => {
    const sp = new URLSearchParams();
    sp.set("status", status);
    sp.set("page", String(page));
    sp.set("page_size", String(pageSize));
    if (keyword.trim()) sp.set("q", keyword.trim());
    sp.set("sort", sort);
    setSearchParams(sp, { replace: true });
  }, [status, page, pageSize, keyword, sort, setSearchParams]);

  // 监听浏览器回退（searchParams变化）同步 state
  useEffect(() => {
    const s = parseStatus(searchParams.get("status"));
    const p = Number(searchParams.get("page") || 1);
    const ps = Number(searchParams.get("page_size") || 10);
    const q = searchParams.get("q") || "";
    const so = searchParams.get("sort");
    const parsedSort: SortOption =
      so === "created_asc" || so === "status" ? so : "created_desc";
    setStatus(s);
    setPage(!Number.isNaN(p) && p > 0 ? p : 1);
    setPageSize(!Number.isNaN(ps) && ps > 0 ? ps : 10);
    setKeyword(q);
    setSort(parsedSort);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams.toString()]);

  async function refresh(next?: {
    status?: StatusFilter;
    page?: number;
    pageSize?: number;
    sort?: SortOption;
    keyword?: string;
  }) {
    setError("");
    const s = next?.status ?? status;
    const p = next?.page ?? page;
    const ps = next?.pageSize ?? pageSize;
    const so = next?.sort ?? sort;
    const q = next?.keyword ?? keyword;

    try {
      setLoading(true);
      const res = await listTasks({
        status: s,
        page: p,
        page_size: ps,
        sort: so,
        q: q.trim() || undefined,
      });
      setData(res);

      // 防御：如果后端返回的 page 超界，前端修正一下
      const tp = Math.max(1, Math.ceil(res.total / res.page_size));
      if (p > tp && tp > 0) {
        setPage(tp);
      }
    } catch (e: any) {
      setError(e?.message ?? "加载任务失败");
      if (e?.code === "UNAUTHORIZED") {
        clearAuth();
        nav("/login");
      }
    } finally {
      setLoading(false);
    }
  }

  // 条件变化刷新
  useEffect(() => {
    refresh({ page: Math.max(1, page) });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [status, page, pageSize, sort, keyword]);

  async function onCreate() {
    setError("");
    if (!canCreate) return;

    try {
      await createTask(title.trim(), description);

      setTitle("");
      setDescription("");

      setPage(1);
      await refresh({ page: 1 });
    } catch (e: any) {
      setError(e?.message ?? "创建失败");
    }
  }

  async function toggleStatus(t: Task) {
    setError("");
    const nextStatus: TaskStatus =
      t.status === "pending" ? "completed" : "pending";

    try {
      await updateTask(t.id, {
        title: t.title,
        description: t.description,
        status: nextStatus,
      });

      await refresh();
    } catch (e: any) {
      setError(e?.message ?? "更新失败");
    }
  }

  async function onDelete(t: Task) {
    setError("");
    const ok = window.confirm(`删除任务「${t.title}」？`);
    if (!ok) return;

    try {
      await deleteTask(t.id);

      const afterTotal = Math.max(0, data.total - 1);
      const tp = Math.max(1, Math.ceil(afterTotal / data.page_size));
      const nextPage = Math.min(page, tp);

      setPage(nextPage);
      await refresh({ page: nextPage });
    } catch (e: any) {
      setError(e?.message ?? "删除失败");
    }
  }

  function logout() {
    clearAuth();
    nav("/login");
  }

  function handlePointerMove(e: React.MouseEvent<HTMLDivElement>) {
    const wrap = wrapRef.current;
    const hero = heroRef.current;
    if (!wrap) return;
    const rect = wrap.getBoundingClientRect();
    wrap.style.setProperty("--mx", `${e.clientX - rect.left}px`);
    wrap.style.setProperty("--my", `${e.clientY - rect.top}px`);
    if (hero) {
      const hRect = hero.getBoundingClientRect();
      hero.style.setProperty("--cx", `${e.clientX - hRect.left}px`);
      hero.style.setProperty("--cy", `${e.clientY - hRect.top}px`);
    }
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

  function handleCardRipple(e: React.MouseEvent<HTMLDivElement>) {
    const card = e.currentTarget as HTMLDivElement;
    const rect = card.getBoundingClientRect();
    card.style.setProperty("--rx", `${e.clientX - rect.left}px`);
    card.style.setProperty("--ry", `${e.clientY - rect.top}px`);
    card.classList.remove("card-ripple");
    void card.offsetWidth;
    card.classList.add("card-ripple");
  }

  return (
    <div
      className="tasks-wrap"
      ref={wrapRef}
      onMouseMove={handlePointerMove}
      onTouchMove={(e) => {
        const t = e.touches[0];
        if (!t) return;
        handlePointerMove({
          ...e,
          clientX: t.clientX,
          clientY: t.clientY,
        } as any);
      }}
    >
      <div className="tasks-shell">
        <header className="tasks-top">
          <div>
            <div className="pill pulse">TASKER</div>
            <h2 className="title">任务控制台</h2>
            <div className="subtitle">
              欢迎，{user?.username ?? "-"} · {loading ? "同步中…" : `共 ${data.total} 项`}
            </div>
          </div>
          <button className="ghost-btn" onClick={logout} onPointerDown={handleButtonRipple}>
            退出
          </button>
        </header>

        <section className="hero" ref={heroRef} onClick={handleCardRipple}>
          <div className="hero-head">
            <div className="pill neon">闪电创建</div>
            <h3>起草一个新任务</h3>
            <p>输入标题与描述，点击创建即可加入任务流。</p>
          </div>
          <div className="hero-form">
            <input
              className="field"
              placeholder="任务标题（必填）"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />
            <textarea
              className="field area"
              placeholder="任务描述（可选）"
              rows={3}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
            <div className="hero-actions">
              <button
                className="cta-btn"
                disabled={!canCreate}
                onClick={onCreate}
                onPointerDown={handleButtonRipple}
              >
                创建
              </button>
              <button
                className="ghost-btn"
                disabled={loading}
                onClick={() => refresh()}
                onPointerDown={handleButtonRipple}
              >
                {loading ? "刷新中…" : "刷新列表"}
              </button>
            </div>
          </div>
        </section>

        <div className="controls">
          <div className="filters">
            {(["all", "pending", "completed"] as StatusFilter[]).map((s) => (
              <button
                key={s}
                className={`chip ${status === s ? "active" : ""}`}
                onClick={() => {
                  setStatus(s);
                  setPage(1);
                }}
                onPointerDown={handleButtonRipple}
              >
                {s === "all" ? "全部" : s === "pending" ? "未完成" : "已完成"}
              </button>
            ))}
          </div>

          <div className="search-sort">
            <div className="search-box">
              <input
                value={keyword}
                onChange={(e) => setKeyword(e.target.value)}
                placeholder="搜索标题 / 描述"
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    setPage(1);
                    refresh({ page: 1, keyword });
                  }
                }}
              />
              <button
                className="mini-btn"
                onClick={() => {
                  setPage(1);
                  refresh({ page: 1, keyword });
                }}
                onPointerDown={handleButtonRipple}
              >
                搜索
              </button>
            </div>

            <label className="select-wrap">
              <span>排序</span>
              <select
                value={sort}
                onChange={(e) => {
                  setSort(e.target.value as SortOption);
                  setPage(1);
                }}
                disabled={loading}
              >
                <option value="created_desc">创建时间 新→旧</option>
                <option value="created_asc">创建时间 旧→新</option>
                <option value="status">状态</option>
              </select>
            </label>
          </div>

          <div className="controls-right">
            <label className="select-wrap">
              <span>每页</span>
              <select
                value={pageSize}
                onChange={(e) => {
                  setPageSize(Number(e.target.value));
                  setPage(1);
                }}
                disabled={loading}
              >
                {[5, 10, 20, 50].map((v) => (
                  <option key={v} value={v}>
                    {v}
                  </option>
                ))}
              </select>
            </label>

            <div className="pager">
              <button
                className="mini-btn"
                disabled={loading || page <= 1}
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                onPointerDown={handleButtonRipple}
              >
                上一页
              </button>
              <span className="pager-text">
                第 {page} / {totalPages} 页 · 本页 {data.items.length} 条
              </span>
              <button
                className="mini-btn"
                disabled={loading || page >= totalPages}
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                onPointerDown={handleButtonRipple}
              >
                下一页
              </button>
            </div>
          </div>
        </div>

        {error && <div className="toast error">{error}</div>}

        <div className="tasks-grid">
          {data.items.map((t, idx) => (
            <div
              key={t.id}
              className="task-card"
              style={{ animationDelay: `${idx * 60}ms` }}
              onClick={handleCardRipple}
            >
              <div className="task-head">
                <div>
                  <div className="task-title">{t.title}</div>
                  {t.description && <div className="task-desc">{t.description}</div>}
                  <div className="task-meta">
                    <span className={`status-pill ${t.status}`}>{t.status}</span>
                    <span className="task-id">#{t.id}</span>
                  </div>
                </div>
                <div className="task-actions">
                  <button
                    className="mini-btn"
                    onClick={(e) => {
                      e.stopPropagation();
                      toggleStatus(t);
                    }}
                    onPointerDown={handleButtonRipple}
                  >
                    {t.status === "pending" ? "标记完成" : "改回未完成"}
                  </button>
                  <button
                    className="mini-btn danger"
                    onClick={(e) => {
                      e.stopPropagation();
                      onDelete(t);
                    }}
                    onPointerDown={handleButtonRipple}
                  >
                    删除
                  </button>
                </div>
              </div>
            </div>
          ))}

          {!loading && data.items.length === 0 && (
            <div className="empty">暂无任务，先创建一个吧。</div>
          )}
        </div>
      </div>
    </div>
  );
}
