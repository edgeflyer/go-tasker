import { api } from "./client";

export type TaskStatus = "pending" | "completed";

export interface Task {
  id: number;
  user_id: number;
  title: string;
  description: string;
  status: TaskStatus;
  created_at: string;
  updated_at: string;
}

export interface TaskPage {
  items: Task[];
  total: number;
  page: number;
  page_size: number;
}

export interface ListTaskParams {
  status?: TaskStatus | "all";
  page?: number;
  page_size?: number;
  q?: string;
  sort?: "created_desc" | "created_asc" | "status";
}

// 把 query 拼到 URL 上
function toQuery(params: ListTaskParams) {
  const sp = new URLSearchParams();

  if (params.status && params.status !== "all") sp.set("status", params.status);
  if (params.page && params.page > 0) sp.set("page", String(params.page));
  if (params.page_size && params.page_size > 0)
    sp.set("page_size", String(params.page_size));
  if (params.q && params.q.trim()) sp.set("q", params.q.trim());
  if (params.sort) sp.set("sort", params.sort);

  const qs = sp.toString();
  return qs ? `?${qs}` : "";
}

export interface UpsertTaskPayload {
  title: string;
  description?: string;
  status?: TaskStatus;
}

/**
 * 兼容两种后端：
 * A) 分页：data = { items, total, page, page_size }
 * B) 非分页：data = Task[]
 */
export async function listTasks(
  params: ListTaskParams = {}
): Promise<TaskPage> {
  const url = `/tasks${toQuery(params)}`;

  // 先尝试按分页结构解析
  const dataAny = await api<any>(url);

  // A: 分页结构
  if (
    dataAny &&
    Array.isArray(dataAny.items) &&
    typeof dataAny.total === "number"
  ) {
    return dataAny as TaskPage;
  }

  // B: 旧结构，data 是 Task[]
  if (Array.isArray(dataAny)) {
    const all = dataAny as Task[];

    // 前端降级：做本地筛选+分页（临时）
    const status =
      params.status && params.status !== "all" ? params.status : undefined;
    const filtered = status ? all.filter((t) => t.status === status) : all;

    const page_size = params.page_size ?? 10;
    const page = params.page ?? 1;

    const total = filtered.length;
    const start = (page - 1) * page_size;
    const items = filtered.slice(start, start + page_size);

    return { items, total, page, page_size };
  }

  // 兜底
  return {
    items: [],
    total: 0,
    page: params.page ?? 1,
    page_size: params.page_size ?? 10,
  };
}

export async function createTask(
  title: string,
  description: string
): Promise<Task> {
  return api<Task>("/tasks", {
    method: "POST",
    body: JSON.stringify({ title, description }),
  });
}

export async function updateTask(
  id: number,
  payload: UpsertTaskPayload
): Promise<Task> {
  return api<Task>(`/tasks/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload),
  });
}

export async function deleteTask(id: number): Promise<{ ok: boolean } | null> {
  // 你的后端可能返回 data: {} 或 data: null，都没关系
  return api<any>(`/tasks/${id}`, { method: "DELETE" });
}
