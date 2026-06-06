import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
});

export interface Application {
  id: number;
  name: string;
  language: string;
  repo_url: string;
  default_branch: string;
  created_at: string;
}

export interface ErrorCluster {
  id: number;
  app_id: number;
  fingerprint: string;
  log_template: string;
  last_seen: string;
  count: number;
  created_at: string;
}

export interface SuggestedFix {
  id: number;
  cluster_id: number;
  explanation: string;
  file_path: string;
  original_code: string;
  fixed_code: string;
  git_patch: string;
  status: string;
  created_at: string;
}

export const fetchApps = async () => {
  const { data } = await api.get<Application[]>('/apps');
  return data;
};

export const fetchApp = async (id: number) => {
  const { data } = await api.get<Application>(`/apps/${id}`);
  return data;
};

export const createApp = async (app: Partial<Application>) => {
  const { data } = await api.post<Application>('/apps', app);
  return data;
};

export const updateApp = async ({ id, ...app }: Partial<Application> & { id: number }) => {
  const { data } = await api.put<Application>(`/apps/${id}`, app);
  return data;
};

export const deleteApp = async (id: number) => {
  await api.delete(`/apps/${id}`);
};

export const fetchClusters = async (appId?: number) => {
  const { data } = await api.get<ErrorCluster[]>('/clusters', {
    params: { app_id: appId },
  });
  return data;
};

export const fetchClusterDetail = async (id: number) => {
  const { data } = await api.get<{ cluster: ErrorCluster; fixes: SuggestedFix[] }>(`/clusters/${id}`);
  return data;
};

export const generateFix = async (id: number) => {
  const { data } = await api.post<SuggestedFix>(`/clusters/${id}/fix`);
  return data;
};

export const ingestLog = async (log: { app_id: number; environment: string; commit_hash: string; log_body: string }) => {
  const { data } = await api.post('/logs/ingest', log);
  return data;
};
