import { Search, FolderGit2, ClipboardCheck, Server, Activity, Loader2, Filter } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { summaryCards } from "@/data/mockData";
import { getProjects, getKanbanTasks } from "@/services/projectService";
import { useEffect, useState } from "react";
import { Project, Task } from "@/data/mockData";

const Dashboard = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Novo estado para o filtro de criadores
  const [selectedCreator, setSelectedCreator] = useState<string>("Todos os criadores");

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        setError(null);

        // Promise.all com catch individual: Se um falhar, o outro continua funcionando!
        const [projectsData, tasksData] = await Promise.all([
          getProjects().catch((err) => {
            console.error("Erro isolado ao carregar projetos:", err);
            return [];
          }),
          getKanbanTasks().catch((err) => {
            console.error("Erro isolado ao carregar tarefas:", err);
            return null;
          })
        ]);

        // 1. PROCESSA PROJETOS
        const rawProjects = Array.isArray(projectsData) ? projectsData : [];
        const mappedProjects = rawProjects.map((p: any, index: number) => ({
          id: p.id || p.ID || index.toString(),
          name: p.name || p.Name || "Projeto sem nome",
          lastCommit: p.lastCommitDate || p.LastCommitDate || "Sem commits",
          ciStatus: (p.statusCI || p.StatusCI || "unknown").toLowerCase()
        }));
        setProjects(mappedProjects);

        // 2. PROCESSA TAREFAS
        let rawTasks: any[] = [];
        if (tasksData && tasksData.buckets && Array.isArray(tasksData.buckets)) {
          tasksData.buckets.forEach((bucket: any) => {
            if (bucket.tasks && Array.isArray(bucket.tasks)) {
              rawTasks = [...rawTasks, ...bucket.tasks];
            }
          });
        } else if (Array.isArray(tasksData)) {
          rawTasks = tasksData;
        }

        const mappedTasks = rawTasks.map((t: any, index: number) => {
          let status: Task["status"] = "todo";
          const rawStatus = (t.status || t.Status || "").toLowerCase();

          if (rawStatus === "em progresso" || rawStatus === "in_progress") {
            status = "in_progress";
          } else if (rawStatus === "concluído" || rawStatus === "concluido" || rawStatus === "done") {
            status = "done";
          }

          let formattedDate = "Sem data";
          const rawDate = t.dueDateTime || t.DueDateTime || t.dueDate;
          if (rawDate) {
            const date = new Date(rawDate);
            formattedDate = isNaN(date.getTime()) ? rawDate : date.toLocaleDateString();
          }

          return {
            id: t.id || index.toString(),
            title: t.title || t.Title || t.tittle || "Sem título",
            priority: (t.priority || "baixa").toLowerCase(),
            dueDate: formattedDate,
            assignee: t.createdBy || t.CreatedBy || "N/A",
            status: status,
          };
        });

        setTasks(mappedTasks);
      } catch (err: any) {
        console.error("Erro crítico no dashboard:", err);
        setError("Erro ao processar a tela. Verifique a conexão com o backend.");
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  // Cria a lista de opções para o select (removendo duplicatas)
  const uniqueCreators = ["Todos os criadores", ...Array.from(new Set(tasks.map(t => t.assignee)))];

  // Aplica o filtro de criador e pega apenas os que não estão "done"
  const pendingTasks = tasks
      .filter((t) => t.status !== "done")
      .filter((t) => selectedCreator === "Todos os criadores" || t.assignee === selectedCreator)
      .slice(0, 4);

  const recentProjects = projects.slice(0, 3);

  return (
      <div className="space-y-6 max-w-6xl mx-auto">
        {/* Global Search */}
        <div className="relative max-w-xl mx-auto">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
              placeholder="Busca global — projetos, docs, tarefas..."
              className="pl-10 h-11 bg-card"
          />
        </div>

        {/* Summary Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <Card className="glass-card animate-fade-in">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">Projetos Ativos</CardTitle>
              <FolderGit2 className="h-4 w-4 text-primary" />
            </CardHeader>
            <CardContent>
              {loading ? (
                  <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
              ) : (
                  <p className="text-3xl font-bold">{projects.length || summaryCards.activeProjects}</p>
              )}
            </CardContent>
          </Card>

          <Card className="glass-card animate-fade-in" style={{ animationDelay: "0.1s" }}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">Tarefas Pendentes</CardTitle>
              <ClipboardCheck className="h-4 w-4 text-warning" />
            </CardHeader>
            <CardContent>
              {loading ? (
                  <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
              ) : (
                  <p className="text-3xl font-bold">{tasks.filter(t => t.status !== 'done').length}</p>
              )}
            </CardContent>
          </Card>

          <Card className="glass-card animate-fade-in" style={{ animationDelay: "0.2s" }}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">Status dos Servidores</CardTitle>
              <Server className="h-4 w-4 text-success" />
            </CardHeader>
            <CardContent>
              <div className="flex items-baseline gap-1">
                <p className="text-3xl font-bold">{summaryCards.serversOnline}</p>
                <span className="text-sm text-muted-foreground">/ {summaryCards.serversTotal}</span>
              </div>
              <div className="flex gap-1 mt-2">
                {Array.from({ length: summaryCards.serversTotal }).map((_, i) => (
                    <div
                        key={i}
                        className={`h-1.5 flex-1 rounded-full ${
                            i < summaryCards.serversOnline ? "bg-success" : "bg-destructive/50"
                        }`}
                    />
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Recent Activity */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <Card className="glass-card">
            <CardHeader>
              <CardTitle className="text-base flex items-center gap-2">
                <Activity className="h-4 w-4 text-primary" />
                Projetos Recentes
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              {loading && (
                  <div className="flex items-center justify-center py-8">
                    <Loader2 className="h-6 w-6 animate-spin text-primary" />
                  </div>
              )}
              {error && <p className="text-destructive text-sm text-center py-4">{error}</p>}
              {!loading && !error && recentProjects.length === 0 && (
                  <p className="text-sm text-muted-foreground text-center py-4">Nenhum projeto encontrado.</p>
              )}
              {!loading && !error && recentProjects.map((p) => (
                  <div key={p.id} className="flex items-center justify-between p-3 rounded-md bg-muted/50">
                    <div className="min-w-0">
                      <p className="font-medium text-sm truncate">{p.name}</p>
                      <p className="text-xs text-muted-foreground font-mono truncate">{p.lastCommit}</p>
                    </div>
                    {/* O CIBadge já pode ser ativado! */}
                    <CIBadge status={p.ciStatus} />
                  </div>
              ))}
            </CardContent>
          </Card>

          <Card className="glass-card">
            <CardHeader className="flex flex-row items-center justify-between border-b pb-3">
              <CardTitle className="text-base flex items-center gap-2 mt-1">
                <ClipboardCheck className="h-4 w-4 text-warning" />
                Tarefas Pendentes
              </CardTitle>

              {/* Filtro de Criadores Adicionado Aqui */}
              {!loading && tasks.length > 0 && (
                  <div className="flex items-center gap-2 text-sm">
                    <Filter className="h-3 w-3 text-muted-foreground" />
                    <select
                        className="bg-background border rounded-md text-xs px-2 py-1 outline-none focus:ring-1 focus:ring-primary"
                        value={selectedCreator}
                        onChange={(e) => setSelectedCreator(e.target.value)}
                    >
                      {uniqueCreators.map((creator) => (
                          <option key={creator} value={creator}>{creator}</option>
                      ))}
                    </select>
                  </div>
              )}
            </CardHeader>

            <CardContent className="space-y-3 pt-4">
              {loading && (
                  <div className="flex items-center justify-center py-8">
                    <Loader2 className="h-6 w-6 animate-spin text-primary" />
                  </div>
              )}
              {error && <p className="text-destructive text-sm text-center py-4">{error}</p>}
              {!loading && !error && pendingTasks.length === 0 && (
                  <p className="text-sm text-muted-foreground text-center py-4">
                    {selectedCreator === "Todos os criadores"
                        ? "Nenhuma tarefa pendente."
                        : `Nenhuma tarefa para ${selectedCreator}.`}
                  </p>
              )}
              {!loading && !error && pendingTasks.map((t) => (
                  <div key={t.id} className="flex items-center justify-between p-3 rounded-md bg-muted/50">
                    <div className="min-w-0">
                      <p className="font-medium text-sm truncate">{t.title}</p>
                      <p className="text-xs text-muted-foreground">{t.assignee} • {t.dueDate}</p>
                    </div>
                    <PriorityBadge priority={t.priority} />
                  </div>
              ))}
            </CardContent>
          </Card>
        </div>
      </div>
  );
};

// Componentes Auxiliares
function CIBadge({ status }: { status: any }) {
  const map = {
    success: { label: "Success", variant: "default" as const, className: "bg-success text-success-foreground" },
    failed: { label: "Failed", variant: "destructive" as const, className: "" },
    running: { label: "Running", variant: "default" as const, className: "bg-info text-info-foreground" },
    unknown: { label: "N/A", variant: "outline" as const, className: "" }
  };

  const currentStatus = (status && map[status as keyof typeof map]) ? status : 'unknown';
  const cfg = map[currentStatus as keyof typeof map];

  return <Badge variant={cfg.variant} className={cfg.className}>{cfg.label}</Badge>;
}

function PriorityBadge({ priority }: { priority: "alta" | "média" | "baixa" }) {
  const map = {
    alta: { className: "bg-destructive/15 text-destructive border-destructive/30" },
    média: { className: "bg-warning/15 text-warning-foreground border-warning/30" },
    baixa: { className: "bg-muted text-muted-foreground" },
  };
  return <Badge variant="outline" className={map[priority].className}>{priority}</Badge>;
}

export { CIBadge, PriorityBadge };
export default Dashboard;