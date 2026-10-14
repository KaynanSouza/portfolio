import { useState, useEffect } from "react";
import { BookOpen, Cpu, CircuitBoard, Download, FileText, FileCode, Archive, File as FileIcon, Loader2 } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { getDocs } from "@/services/projectService";

// Tipagens locais agora que não dependemos mais do mockData
export type DocItem = {
  id: string;
  title: string;
  filename: string;
  type: string;
  size: string;
};

export type DocCategory = {
  id: string;
  label: string;
  icon: string;
  items: DocItem[];
};

const iconMap: Record<string, React.ElementType> = {
  Cpu: Cpu,
  BookOpen: BookOpen,
  Microchip: CircuitBoard,
};

const FileTypeIcon = ({ type }: { type: string }) => {
  switch (type) {
    case "pdf":
      return <FileText className="h-8 w-8 text-red-500" />;
    case "bin":
      return <FileCode className="h-8 w-8 text-blue-500" />;
    case "zip":
    case "rar":
      return <Archive className="h-8 w-8 text-orange-500" />;
    case "docx":
    case "doc":
      return <FileIcon className="h-8 w-8 text-blue-600" />;
    default:
      return <FileIcon className="h-8 w-8 text-gray-500" />;
  }
};

const Documentation = () => {
  const [categories, setCategories] = useState<DocCategory[]>([]);
  const [selectedCat, setSelectedCat] = useState<DocCategory | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchDocuments = async () => {
      try {
        setLoading(true);
        const data: DocItem[] = await getDocs();

        // Criando uma categoria principal dinâmica com os arquivos do servidor
        const serverCategory: DocCategory = {
          id: "server-docs",
          label: "Arquivos do Servidor",
          icon: "BookOpen",
          items: data || [],
        };

        setCategories([serverCategory]);
        setSelectedCat(serverCategory);
      } catch (error) {
        console.error("Erro ao carregar documentos:", error);
        toast.error("Erro ao carregar os arquivos do servidor.");
      } finally {
        setLoading(false);
      }
    };

    fetchDocuments();
  }, []);

  const handleDownload = async (item: DocItem) => {
    try {
      toast.info(`Iniciando download de ${item.filename}...`);

      // Busca o arquivo na rota estática que criamos no Go
      // Lembre-se de ajustar a porta se o seu backend estiver em outra
      const response = await fetch(`http://localhost:8081/api/download?file=${encodeURIComponent(item.filename)}`);

      if (!response.ok) throw new Error("Arquivo não encontrado no servidor");

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", item.filename);
      document.body.appendChild(link);
      link.click();
      link.parentNode?.removeChild(link);
      window.URL.revokeObjectURL(url);

      toast.success("Download concluído com sucesso!");
    } catch (error) {
      console.error("Erro ao baixar arquivo:", error);
      toast.error("Falha ao baixar o arquivo.");
    }
  };

  return (
      <div className="space-y-6 max-w-6xl mx-auto p-4 md:p-0">
        <div className="flex flex-col gap-1">
          <h1 className="text-3xl font-bold tracking-tight text-foreground">Centro de Downloads</h1>
          <p className="text-muted-foreground">
            Acesse e baixe equipamentos, manuais técnicos e firmwares oficiais do servidor.
          </p>
        </div>

        <div className="flex flex-col md:flex-row gap-8">
          {/* Navigation Sidebar */}
          <aside className="w-full md:w-64 space-y-2">
            {loading ? (
                <div className="p-4 flex justify-center">
                  <Loader2 className="h-6 w-6 animate-spin text-primary" />
                </div>
            ) : categories.map((cat) => {
              const Icon = iconMap[cat.icon] || BookOpen;
              const isActive = selectedCat?.id === cat.id;

              return (
                  <button
                      key={cat.id}
                      onClick={() => setSelectedCat(cat)}
                      className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-all duration-200 ${
                          isActive
                              ? "bg-primary text-primary-foreground shadow-lg shadow-primary/20"
                              : "bg-card hover:bg-muted text-muted-foreground border border-transparent hover:border-border"
                      }`}
                  >
                    <Icon className={`h-5 w-5 ${isActive ? "text-primary-foreground" : "text-primary"}`} />
                    {cat.label}
                  </button>
              );
            })}
          </aside>

          {/* Main Content Area */}
          <div className="flex-1">
            {loading ? (
                <div className="flex flex-col items-center justify-center py-20">
                  <Loader2 className="h-8 w-8 animate-spin text-primary mb-4" />
                  <p className="text-muted-foreground">Lendo pasta do servidor...</p>
                </div>
            ) : selectedCat && selectedCat.items.length > 0 ? (
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                  {selectedCat.items.map((item) => (
                      <Card key={item.id} className="glass-card overflow-hidden hover:shadow-xl transition-all duration-300 group border-muted/40">
                        <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-3">
                          <div className="flex items-start gap-4">
                            <div className="p-3 bg-muted rounded-xl group-hover:bg-primary/10 transition-colors duration-300">
                              <FileTypeIcon type={item.type} />
                            </div>
                            <div className="space-y-1">
                              <CardTitle
                                  className="text-base font-bold group-hover:text-primary transition-colors cursor-pointer leading-tight break-all"
                                  onClick={() => handleDownload(item)}
                              >
                                {item.title}
                              </CardTitle>
                              <CardDescription className="text-xs font-mono text-muted-foreground/70">
                                Extensão: .{item.type}
                              </CardDescription>
                            </div>
                          </div>
                          <Button
                              variant="ghost"
                              size="icon"
                              className="rounded-full hover:bg-primary/20 hover:text-primary transition-all shrink-0"
                              onClick={() => handleDownload(item)}
                          >
                            <Download className="h-4 w-4" />
                          </Button>
                        </CardHeader>
                        <CardContent>
                          <div className="flex items-center justify-between mt-2 pt-4 border-t border-muted/30">
                      <span className="text-[10px] uppercase tracking-wider font-bold text-muted-foreground/60 bg-muted/50 px-2 py-1 rounded">
                        {item.size}
                      </span>
                            <Button
                                variant="link"
                                className="h-auto p-0 text-xs font-semibold text-primary hover:text-primary/80 transition-colors"
                                onClick={() => handleDownload(item)}
                            >
                              Baixar Arquivo
                            </Button>
                          </div>
                        </CardContent>
                      </Card>
                  ))}
                </div>
            ) : (
                <div className="flex flex-col items-center justify-center py-20 text-center">
                  <FileIcon className="h-12 w-12 text-muted-foreground/30 mb-4" />
                  <h3 className="text-lg font-medium text-muted-foreground">Nenhum arquivo encontrado</h3>
                  <p className="text-sm text-muted-foreground/60">Coloque arquivos na pasta 'documentos' do seu servidor Go.</p>
                </div>
            )}
          </div>
        </div>
      </div>
  );
};

export default Documentation;