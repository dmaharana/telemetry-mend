import { useState } from 'react';
import { QueryClient, QueryClientProvider, useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { BrowserRouter as Router, Routes, Route, Link, useParams, useNavigate } from 'react-router-dom';
import { 
  LayoutDashboard, 
  Settings, 
  Terminal, 
  PlusCircle, 
  ChevronRight,
  AlertCircle,
  Code2,
  CheckCircle2
} from 'lucide-react';
import ReactDiffViewer from 'react-diff-viewer-continued';

import { 
  fetchApps, 
  fetchClusters, 
  fetchClusterDetail, 
  generateFix, 
  createApp,
} from '@/lib/api';
import type { 
  Application,
  ErrorCluster,
  SuggestedFix
} from '@/lib/api';

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";

const queryClient = new QueryClient();

function Sidebar({ apps }: { apps?: Application[] }) {
  return (
    <div className="w-64 border-r bg-slate-50 flex flex-col h-screen">
      <div className="p-4 border-b">
        <h1 className="text-xl font-bold flex items-center gap-2">
          <Terminal className="w-6 h-6 text-primary" />
          TelemetryMend
        </h1>
      </div>
      <ScrollArea className="flex-1">
        <div className="p-4 space-y-4">
          <div>
            <h2 className="mb-2 px-2 text-xs font-semibold uppercase tracking-wider text-slate-500">
              Overview
            </h2>
            <Link to="/">
              <Button variant="ghost" className="w-full justify-start gap-2">
                <LayoutDashboard className="w-4 h-4" />
                Global Dashboard
              </Button>
            </Link>
          </div>
          <div>
            <div className="flex items-center justify-between mb-2 px-2">
              <h2 className="text-xs font-semibold uppercase tracking-wider text-slate-500">
                Applications
              </h2>
              <Link to="/settings">
                <PlusCircle className="w-4 h-4 text-slate-400 hover:text-primary cursor-pointer" />
              </Link>
            </div>
            <div className="space-y-1">
              {apps?.map((app) => (
                <Link key={app.id} to={`/app/${app.id}`}>
                  <Button variant="ghost" className="w-full justify-start gap-2">
                    <div className="w-2 h-2 rounded-full bg-green-500" />
                    {app.name}
                  </Button>
                </Link>
              ))}
            </div>
          </div>
        </div>
      </ScrollArea>
      <div className="p-4 border-t">
        <Link to="/settings">
          <Button variant="ghost" className="w-full justify-start gap-2 text-slate-500">
            <Settings className="w-4 h-4" />
            Settings
          </Button>
        </Link>
      </div>
    </div>
  );
}

function GlobalDashboard() {
  const { data: clusters, isLoading } = useQuery<ErrorCluster[]>({
    queryKey: ['clusters'],
    queryFn: () => fetchClusters(),
  });

  if (isLoading) return <DashboardSkeleton />;

  return (
    <div className="p-8 space-y-6">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">System Status</h2>
        <p className="text-muted-foreground">Aggregated error clusters across all microservices.</p>
      </div>
      <Separator />
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {clusters?.map(cluster => (
          <Card key={cluster.id} className="hover:shadow-md transition-shadow cursor-pointer">
            <Link to={`/cluster/${cluster.id}`}>
              <CardHeader>
                <div className="flex justify-between items-start">
                  <Badge variant="destructive">Cluster #{cluster.id}</Badge>
                  <span className="text-xs text-muted-foreground">{new Date(cluster.last_seen).toLocaleString()}</span>
                </div>
                <CardTitle className="text-sm mt-2 font-mono line-clamp-2">
                  {cluster.log_template}
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{cluster.count}</div>
                <p className="text-xs text-muted-foreground">Occurrences</p>
              </CardContent>
            </Link>
          </Card>
        ))}
      </div>
    </div>
  );
}

function AppDashboard() {
  const { id } = useParams();
  const { data: clusters, isLoading } = useQuery<ErrorCluster[]>({
    queryKey: ['clusters', id],
    queryFn: () => fetchClusters(Number(id)),
  });

  if (isLoading) return <DashboardSkeleton />;

  return (
    <div className="p-8 space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">App Dashboard</h2>
          <p className="text-muted-foreground">Error clusters for this specific application.</p>
        </div>
      </div>
      <Separator />
      <div className="space-y-4">
        {clusters?.map(cluster => (
          <Link key={cluster.id} to={`/cluster/${cluster.id}`}>
            <Card className="hover:bg-slate-50 transition-colors">
              <CardContent className="p-6 flex items-center justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <Badge variant="secondary">Fingerprint: {cluster.fingerprint.slice(0, 8)}</Badge>
                    <span className="text-xs text-muted-foreground">Last seen: {new Date(cluster.last_seen).toLocaleString()}</span>
                  </div>
                  <code className="text-sm font-mono block bg-slate-100 p-2 rounded mt-2">
                    {cluster.log_template}
                  </code>
                </div>
                <div className="ml-8 text-center">
                  <div className="text-2xl font-bold">{cluster.count}</div>
                  <div className="text-xs text-muted-foreground uppercase">Errors</div>
                </div>
                <ChevronRight className="ml-4 w-5 h-5 text-slate-400" />
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  );
}

function ClusterDetail() {
  const { id } = useParams();
  const queryClient = useQueryClient();
  const { data, isLoading } = useQuery<{ cluster: ErrorCluster; fixes: SuggestedFix[] }>({
    queryKey: ['cluster', id],
    queryFn: () => fetchClusterDetail(Number(id)),
  });

  const mutation = useMutation({
    mutationFn: (id: number) => generateFix(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cluster', id] });
    },
  });

  if (isLoading) return <div className="p-8">Loading cluster details...</div>;
  if (!data) return <div className="p-8 text-destructive">Cluster not found</div>;

  const { cluster, fixes } = data;

  return (
    <div className="p-8 h-screen flex flex-col overflow-hidden">
      <div className="flex justify-between items-start mb-6 shrink-0">
        <div>
          <h2 className="text-2xl font-bold flex items-center gap-2">
            <AlertCircle className="text-destructive w-6 h-6" />
            Error Cluster #{cluster.id}
          </h2>
          <p className="text-muted-foreground">Analyzed from {cluster.count} logs.</p>
        </div>
        <Button onClick={() => mutation.mutate(cluster.id)} disabled={mutation.isPending}>
          {mutation.isPending ? "Generating Fix..." : "Generate AI Fix"}
        </Button>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 flex-1 overflow-hidden">
        <div className="space-y-6 overflow-y-auto pr-2">
          <Card>
            <CardHeader>
              <CardTitle className="text-sm uppercase text-muted-foreground font-semibold">Log Template</CardTitle>
            </CardHeader>
            <CardContent>
              <pre className="bg-slate-900 text-slate-100 p-4 rounded-lg text-xs overflow-x-auto whitespace-pre-wrap">
                {cluster.log_template}
              </pre>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-sm uppercase text-muted-foreground font-semibold">Suggested Fixes</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {fixes.length === 0 && <p className="text-muted-foreground text-sm italic">No fixes generated yet.</p>}
              {fixes.map(fix => (
                <div key={fix.id} className="border rounded-lg p-4 space-y-3">
                  <div className="flex justify-between items-center">
                    <Badge variant="outline" className="flex items-center gap-1">
                      <CheckCircle2 className="w-3 h-3 text-green-500" />
                      Fix #{fix.id}
                    </Badge>
                    <span className="text-xs text-muted-foreground">{new Date(fix.created_at).toLocaleDateString()}</span>
                  </div>
                  <p className="text-sm">{fix.explanation}</p>
                  <Badge className="bg-blue-100 text-blue-700 hover:bg-blue-100 font-mono text-[10px]">
                    {fix.file_path}
                  </Badge>
                </div>
              ))}
            </CardContent>
          </Card>
        </div>

        <Card className="flex flex-col overflow-hidden">
          <CardHeader className="shrink-0 border-b">
            <div className="flex items-center gap-2">
              <Code2 className="w-4 h-4 text-primary" />
              <CardTitle className="text-sm">Diff View</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="flex-1 overflow-y-auto p-0 font-mono text-xs">
            {fixes.length > 0 ? (
              <ReactDiffViewer 
                oldValue={fixes[0].original_code} 
                newValue={fixes[0].fixed_code} 
                splitView={true} 
                useDarkTheme={false}
              />
            ) : (
              <div className="h-full flex items-center justify-center text-muted-foreground italic">
                Select or generate a fix to see the diff.
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function SettingsPage() {
  const [name, setName] = useState('');
  const [repo, setRepo] = useState('');
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const mutation = useMutation({
    mutationFn: createApp,
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['apps'] });
      navigate(`/app/${data.id}`);
    },
  });

  return (
    <div className="p-8 max-w-2xl mx-auto space-y-8">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">System Settings</h2>
        <p className="text-muted-foreground">Configure your microservices and SCM integrations.</p>
      </div>
      <Separator />
      <Card>
        <CardHeader>
          <CardTitle>Register New Application</CardTitle>
          <CardDescription>Add a new microservice for log analysis.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Application Name</Label>
            <Input id="name" placeholder="e.g. auth-service" value={name} onChange={e => setName(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="repo">Git Repository URL</Label>
            <Input id="repo" placeholder="https://github.com/org/repo.git" value={repo} onChange={e => setRepo(e.target.value)} />
          </div>
        </CardContent>
        <CardFooter>
          <Button onClick={() => mutation.mutate({ name, repo_url: repo })} disabled={mutation.isPending}>
            {mutation.isPending ? "Registering..." : "Add Application"}
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}

function DashboardSkeleton() {
  return (
    <div className="p-8 space-y-6">
      <Skeleton className="h-10 w-48" />
      <Skeleton className="h-4 w-full" />
      <Separator />
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {[1, 2, 3].map(i => <Skeleton key={i} className="h-48 w-full" />)}
      </div>
    </div>
  );
}

function AppContent() {
  const { data: apps } = useQuery({
    queryKey: ['apps'],
    queryFn: fetchApps,
  });

  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar apps={apps} />
      <main className="flex-1 overflow-y-auto bg-slate-200/50">
        <Routes>
          <Route path="/" element={<GlobalDashboard />} />
          <Route path="/app/:id" element={<AppDashboard />} />
          <Route path="/cluster/:id" element={<ClusterDetail />} />
          <Route path="/settings" element={<SettingsPage />} />
        </Routes>
      </main>
    </div>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <AppContent />
      </Router>
    </QueryClientProvider>
  );
}

export default App;
