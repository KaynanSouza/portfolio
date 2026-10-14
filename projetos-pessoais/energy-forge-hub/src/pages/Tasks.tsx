import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { getKanbanTasks } from "@/services/projectService";
import { CalendarDays, User, Loader2, AlertCircle, Filter } from "lucide-react";
import { type Task } from "@/data/mockData";

const columns: { key: Task["status"]; label: string; accent: string }[] = [
  { key: "todo", label: "A Fazer", accent: "border-t-primary" },
  { key: "in_progress", label: "Em Progresso", accent: "border-t-warning" },
  { key: "done", label: "Concluído", accent: "border-t-success" },
];

const priorityStyles: Record<string, string> = {
  alta: "bg-destructive/15 text-destructive border-destructive/30",
  média: "bg-warning/15 text-warning-foreground border-warning/30",
  baixa: "bg-muted text-muted-foreground",
};

const TasksPage = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Estado para o filtro de criadores
  const [selectedCreator, setSelectedCreator] = useState<string>("Todos os criadores");

  useEffect(() => {
    const fetchTasks = async () => {
      try {
        setLoading(true);
        const data = await getKanbanTasks();

        console.log("Dados puros do backend (Tasks.tsx):", data);

        let rawList: any[] = [];

        // Busca ESTRITA: Olha apenas para os buckets e extrai as tarefas
        if (data && data.buckets && Array.isArray(data.buckets)) {
          data.buckets.forEach((bucket: any) => {
            // Cobre a possibilidade da propriedade vir minúscula (tasks) ou maiúscula (Tasks) do Go
            const bucketTasks = bucket.tasks || bucket.Tasks;
            if (bucketTasks && Array.isArray(bucketTasks)) {
              rawList = [...rawList, ...bucketTasks];
            }
          });
        }

        const mappedTasks = rawList.map((t: any, index: number) => {
          // Normalizando a data
          let formattedDate = "Sem data";
          const rawDate = t.dueDateTime || t.DueDateTime || t.dueDate;
          if (rawDate) {
            const date = new Date(rawDate);
            formattedDate = isNaN(date.getTime()) ? rawDate : date.toLocaleDateString();
          }

          // Normalizando o status
          let status: Task["status"] = "todo";
          const rawStatus = (t.status || t.Status || "").toLowerCase().trim();

          if (rawStatus === "em progresso" || rawStatus === "in_progress" || rawStatus === "em andamento") {
            status = "in_progress";
          } else if (rawStatus === "concluído" || rawStatus === "concluido" || rawStatus === "done") {
            status = "done";
          }

          return {
            id: t.id || t.ID || index.toString(),
            title: t.title || t.Title || t.tittle || "Sem título",
            priority: (t.priority || t.Priority || "baixa").toLowerCase(),
            dueDate: formattedDate,
            assignee: t.createdBy || t.CreatedBy || "N/A",
            status: status,
          };
        });

        setTasks(mappedTasks);
        setError(null);
      } catch (err: any) {
        console.error("Erro ao processar tarefas:", err);
        setError(err.message || "Erro ao processar os dados do servidor.");
      } finally {
        setLoading(false);
      }
    };

    fetchTasks();
  }, []);

  // Cria a lista única de criadores
  const uniqueCreators = ["Todos os criadores", ...Array.from(new Set(tasks.map(t => t.assignee)))];

  // Filtra as tarefas de acordo com o dropdown
  const filteredTasks = tasks.filter(
      (t) => selectedCreator === "Todos os criadores" || t.assignee === selectedCreator
  );

  if (loading) {
    return (
        <div className="flex flex-col items-center justify-center h-[60vh] space-y-4">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <p className="text-muted-foreground animate-pulse">Carregando quadro Kanban...</p>
        </div>
    );
  }

  if (error) {
    return (
        <div className="flex flex-col items-center justify-center h-[60vh] space-y-4 max-w-lg mx-auto text-center px-6">
          <div className="bg-destructive/10 text-destructive p-6 rounded-xl border border-destructive/20 w-full">
            <AlertCircle className="h-10 w-10 mx-auto mb-3 opacity-80" />
            <h2 className="text-lg font-bold mb-1">Erro nos Dados</h2>
            <p className="text-sm opacity-90 mb-4">{error}</p>
            <button
                onClick={() => window.location.reload()}
                className="text-xs underline font-medium"
            >
              Tentar novamente
            </button>
          </div>
        </div>
    );
  }

  return (
      <div className="space-y-6 max-w-6xl mx-auto">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-end gap-4">
          <div>
            <h1 className="text-2xl font-bold">Tarefas</h1>
            <p className="text-sm text-muted-foreground">Quadro Kanban (Dados do Backend)</p>
          </div>

          {/* Filtro de Criadores e Contador */}
          <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
            <div className="flex items-center gap-2 text-sm bg-background border px-3 py-1.5 rounded-md shadow-sm">
              <Filter className="h-4 w-4 text-muted-foreground" />
              <span className="text-muted-foreground">Filtrar por:</span>
              <select
                  className="bg-transparent border-none outline-none font-medium focus:ring-0 text-sm cursor-pointer"
                  value={selectedCreator}
                  onChange={(e) => setSelectedCreator(e.target.value)}
              >
                {uniqueCreators.map((creator) => (
                    <option key={creator} value={creator}>{creator}</option>
                ))}
              </select>
            </div>

            <div className="text-xs text-muted-foreground bg-muted px-2 py-1 rounded hidden sm:block">
              Total: {filteredTasks.length} tarefas
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {columns.map((col) => {
            // Usando as tarefas filtradas (filteredTasks) em vez de todas (tasks)
            const colTasks = filteredTasks.filter((t) => t.status === col.key);
            return (
                <Card key={col.key} className={`glass-card border-t-2 ${col.accent}`}>
                  <CardHeader className="pb-3">
                    <CardTitle className="text-sm font-medium flex items-center justify-between">
                      {col.label}
                      <span className="text-xs bg-muted px-2 py-0.5 rounded-full text-muted-foreground">
                    {colTasks.length}
                  </span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-3">
                    {colTasks.length === 0 ? (
                        <div className="text-center py-8 border-2 border-dashed rounded-lg opacity-40">
                          <p className="text-xs">Vazio</p>
                        </div>
                    ) : (
                        colTasks.map((task) => (
                            <div
                                key={task.id}
                                className="rounded-lg border bg-card p-3 space-y-2 hover:shadow-sm transition-shadow"
                            >
                              <div className="flex items-start justify-between gap-2">
                                <p className="text-sm font-medium leading-snug">{task.title}</p>
                                <Badge variant="outline" className={priorityStyles[task.priority] || priorityStyles.baixa}>
                                  {task.priority}
                                </Badge>
                              </div>
                              <div className="flex items-center gap-3 text-xs text-muted-foreground">
                        <span className="flex items-center gap-1">
                          <User className="h-3 w-3" />
                          {task.assignee}
                        </span>
                                <span className="flex items-center gap-1">
                          <CalendarDays className="h-3 w-3" />
                                  {task.dueDate}
                        </span>
                              </div>
                            </div>
                        ))
                    )}
                  </CardContent>
                </Card>
            );
          })}
        </div>
      </div>
  );
};

export default TasksPage;