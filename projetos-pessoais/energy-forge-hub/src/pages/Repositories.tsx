import { ExternalLink, GitBranch } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { CIBadge } from "@/pages/Dashboard";
import { getProjects } from "@/services/projectService";
import { useEffect, useState } from "react";
import { Project } from "@/data/mockData";

const Repositories = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProjects = async () => {
      try {
        const data = await getProjects();
        setProjects(data);
      } catch (err) {
        setError("Failed to fetch projects");
      } finally {
        setLoading(false);
      }
    };

    fetchProjects();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>{error}</div>;
  }

  return (
    <div className="space-y-6 max-w-6xl mx-auto">
      <div>
        <h1 className="text-2xl font-bold">Repositórios</h1>
        <p className="text-sm text-muted-foreground">Projetos GitLab do setor de engenharia</p>
      </div>

      <Card className="glass-card overflow-hidden">
        <CardHeader className="pb-0">
          <CardTitle className="text-base flex items-center gap-2">
            <GitBranch className="h-4 w-4 text-primary" />
            Projetos ({projects.length})
          </CardTitle>
        </CardHeader>
        <CardContent className="p-0 mt-4">
          {/* Desktop table */}
          <div className="hidden md:block">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Projeto</TableHead>
                  <TableHead>Último Commit</TableHead>
                  <TableHead>Branch</TableHead>
                  <TableHead>CI/CD</TableHead>
                  <TableHead className="text-right">Ação</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {projects.map((p) => (
                  <TableRow key={p.id}>
                    <TableCell>
                      <p className="font-medium">{p.name}</p>
                      <p className="text-xs text-muted-foreground">{p.description}</p>
                    </TableCell>
                    <TableCell>
                      <p className="text-sm font-mono truncate max-w-[200px]">{p.lastCommit}</p>
                      <p className="text-xs text-muted-foreground">{p.lastCommitDate}</p>
                    </TableCell>
                    <TableCell>
                      <code className="text-xs bg-muted px-1.5 py-0.5 rounded">{p.branch}</code>
                    </TableCell>
                    <TableCell>
                      <CIBadge status={p.ciStatus} />
                    </TableCell>
                    <TableCell className="text-right">
                      <Button variant="ghost" size="sm" asChild>
                        <a href={p.web_url} target="_blank" rel="noopener noreferrer">
                          <ExternalLink className="h-3.5 w-3.5 mr-1" />
                          GitLab
                        </a>
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          {/* Mobile cards */}
          <div className="md:hidden space-y-3 p-4">
            {projects.map((p) => (
              <div key={p.id} className="rounded-lg border bg-muted/30 p-4 space-y-2">
                <div className="flex items-start justify-between">
                  <div>
                    <p className="font-medium">{p.name}</p>
                    <p className="text-xs text-muted-foreground">{p.description}</p>
                  </div>
                  <CIBadge status={p.ciStatus} />
                </div>
                <p className="text-xs font-mono text-muted-foreground truncate">{p.lastCommit}</p>
                <div className="flex items-center justify-between">
                  <code className="text-xs bg-muted px-1.5 py-0.5 rounded">{p.branch}</code>
                  <Button variant="ghost" size="sm" asChild>
                    <a href={p.web_url} target="_blank" rel="noopener noreferrer">
                      <ExternalLink className="h-3.5 w-3.5 mr-1" />
                      GitLab
                    </a>
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default Repositories;
